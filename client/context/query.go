package context

import (
	"fmt"
	"strings"

	"github.com/cosmos/cosmos-sdk/store/rootmulti"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/pkg/errors"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto/merkle"
	cmn "github.com/tendermint/tendermint/libs/common"
	tmliteErr "github.com/tendermint/tendermint/lite/errors"
	tmliteProxy "github.com/tendermint/tendermint/lite/proxy"
	rpcclient "github.com/tendermint/tendermint/rpc/client"
	tmtypes "github.com/tendermint/tendermint/types"
)

// GetNode returns an RPC client. If the context's client is not defined, an
// error is returned.
func (ctx Context) GetNode() (rpcclient.Client, error) {
	if ctx.Client == nil {
		return nil, errors.New("no RPC client defined")
	}

	return ctx.Client, nil
}

// Query performs a query for information about the connected node.
func (ctx Context) Query(path string, data cmn.HexBytes) ([]byte, int64, error) {
	return ctx.query(path, data)
}

// Query information about the connected node with a data payload
func (ctx Context) QueryWithData(path string, data []byte) ([]byte, int64, error) {
	return ctx.query(path, data)
}

// QueryStore performs a query from a Tendermint node with the provided key and
// store name.
func (ctx Context) QueryStore(key cmn.HexBytes, storeName string) ([]byte, int64, error) {
	return ctx.queryStore(key, storeName, "key")
}

// QuerySubspace performs a query from a Tendermint node with the provided
// store name and subspace.
func (ctx Context) QuerySubspace(subspace []byte, storeName string) (res []sdk.KVPair, err error) {
	resRaw, _, err := ctx.queryStore(subspace, storeName, "subspace")
	if err != nil {
		return res, err
	}

	ctx.Codec.MustUnmarshalBinaryLengthPrefixed(resRaw, &res)
	return
}

// GetFromAddress returns the from address from the context's name.
func (ctx Context) GetFromAddress() sdk.AccAddress {
	return ctx.FromAddress
}

// GetFromName returns the key name for the current context.
func (ctx Context) GetFromName() string {
	return ctx.FromName
}

// EnsureAccountExists ensures that an account exists for a given context. An
// error is returned if it does not.
func (ctx Context) EnsureAccountExists() error {
	addr := ctx.GetFromAddress()
	return ctx.EnsureAccountExistsFromAddr(addr)
}

// EnsureAccountExistsFromAddr ensures that an account exists for a given
// address. Instead of using the context's from name, a direct address is
// given. An error is returned if it does not.
func (ctx Context) EnsureAccountExistsFromAddr(addr sdk.AccAddress) error {
	_, err := ctx.queryAccount(addr)
	return err
}

type QueryAccountParams struct {
	Address sdk.AccAddress
}

// queryAccount queries an account using custom query endpoint of auth module
// returns an error if result is `null` otherwise account data
func (ctx Context) queryAccount(addr sdk.AccAddress) ([]byte, error) {
	bz, err := ctx.Codec.MarshalJSON(QueryAccountParams{Address: addr})
	if err != nil {
		return nil, err
	}

	route := fmt.Sprintf("custom/%s/%s", ctx.AccountStore, "account")

	res, _, err := ctx.QueryWithData(route, bz)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// query performs a query from a Tendermint node with the provided store name
// and path.
func (ctx *Context) query(path string, key cmn.HexBytes) (res []byte, height int64, err error) {
	node, err := ctx.GetNode()
	if err != nil {
		return res, height, err
	}

	// When a client did not provide a query height, manually query for it so it can
	// be injected downstream into responses.
	if ctx.Height == 0 {
		status, err := node.Status()
		if err != nil {
			return res, height, err
		}
		ctx.WithHeight(status.SyncInfo.LatestBlockHeight)
	}

	opts := rpcclient.ABCIQueryOptions{
		Height: ctx.Height,
		Prove:  !ctx.TrustNode,
	}

	result, err := node.ABCIQueryWithOptions(path, key, opts)
	if err != nil {
		return res, height, err
	}

	resp := result.Response
	if !resp.IsOK() {
		return res, height, errors.New(resp.Log)
	}

	// data from trusted node or subspace query doesn't need verification
	if ctx.TrustNode || !isQueryStoreWithProof(path) {
		return resp.Value, resp.Height, nil
	}

	err = ctx.verifyProof(path, resp)
	if err != nil {
		return res, height, err
	}

	return resp.Value, resp.Height, nil
}

// Verify verifies the consensus proof at given height.
func (ctx *Context) Verify(height int64) (tmtypes.SignedHeader, error) {
	check, err := tmliteProxy.GetCertifiedCommit(height, ctx.Client, ctx.GetVerifier())
	switch {
	case tmliteErr.IsErrCommitNotFound(err):
		return tmtypes.SignedHeader{}, ErrVerifyCommit(height)
	case err != nil:
		return tmtypes.SignedHeader{}, err
	}

	return check, nil
}

// verifyProof perform response proof verification.
func (ctx *Context) verifyProof(queryPath string, resp abci.ResponseQuery) error {
	if ctx.GetVerifier() == nil {
		return fmt.Errorf("missing valid certifier to verify data from distrusted node")
	}

	// the AppHash for height H is in header H+1
	commit, err := ctx.Verify(resp.Height + 1)
	if err != nil {
		return err
	}

	// TODO: Instead of reconstructing, stash on Context field?
	prt := rootmulti.DefaultProofRuntime()

	// TODO: Better convention for path?
	storeName, err := parseQueryStorePath(queryPath)
	if err != nil {
		return err
	}

	kp := merkle.KeyPath{}
	kp = kp.AppendKey([]byte(storeName), merkle.KeyEncodingURL)
	kp = kp.AppendKey(resp.Key, merkle.KeyEncodingURL)

	if resp.Value == nil {
		err = prt.VerifyAbsence(resp.Proof, commit.Header.AppHash, kp.String())
		if err != nil {
			return errors.Wrap(err, "failed to prove merkle proof")
		}
		return nil
	}
	err = prt.VerifyValue(resp.Proof, commit.Header.AppHash, kp.String(), resp.Value)
	if err != nil {
		return errors.Wrap(err, "failed to prove merkle proof")
	}

	return nil
}

// queryStore performs a query from a Tendermint node with the provided a store
// name and path.
func (ctx *Context) queryStore(key cmn.HexBytes, storeName, endPath string) ([]byte, int64, error) {
	path := fmt.Sprintf("/store/%s/%s", storeName, endPath)
	return ctx.query(path, key)
}

// isQueryStoreWithProof expects a format like /<queryType>/<storeName>/<subpath>
// queryType must be "store" and subpath must be "key" to require a proof.
func isQueryStoreWithProof(path string) bool {
	if !strings.HasPrefix(path, "/") {
		return false
	}

	paths := strings.SplitN(path[1:], "/", 3)
	switch {
	case len(paths) != 3:
		return false
	case paths[0] != "store":
		return false
	case rootmulti.RequireProof("/" + paths[2]):
		return true
	}

	return false
}

// parseQueryStorePath expects a format like /store/<storeName>/key.
func parseQueryStorePath(path string) (storeName string, err error) {
	if !strings.HasPrefix(path, "/") {
		return "", errors.New("expected path to start with /")
	}

	paths := strings.SplitN(path[1:], "/", 3)
	switch {
	case len(paths) != 3:
		return "", errors.New("expected format like /store/<storeName>/key")
	case paths[0] != "store":
		return "", errors.New("expected format like /store/<storeName>/key")
	case paths[2] != "key":
		return "", errors.New("expected format like /store/<storeName>/key")
	}

	return paths[1], nil
}
