package context

import (
	"fmt"
	"io"
	"path/filepath"
	"sync"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	cryptokeys "github.com/cosmos/cosmos-sdk/crypto/keys"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dgamingfoundation/cosmos-utils/client/keys"
	"github.com/pkg/errors"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/libs/log"
	tmlite "github.com/tendermint/tendermint/lite"
	tmliteProxy "github.com/tendermint/tendermint/lite/proxy"
	rpcclient "github.com/tendermint/tendermint/rpc/client"
)

// Context implements a typical CLI context created in SDK modules for
// transaction handling and queries.
type Context struct {
	Codec         *codec.Codec
	Client        rpcclient.Client
	Keybase       cryptokeys.Keybase
	Output        io.Writer
	OutputFormat  string
	Height        int64
	NodeURI       string
	From          string
	AccountStore  string
	TrustNode     bool
	UseLedger     bool
	BroadcastMode string
	PrintResponse bool
	verifier      tmlite.Verifier
	VerifierHome  string
	Simulate      bool
	GenerateOnly  bool
	FromAddress   sdk.AccAddress
	FromName      string
	Indent        bool
	SkipConfirm   bool
	Home          string
	Passphrase    string
	PrivKey       crypto.PrivKey
	mtx           sync.RWMutex
}

// NewContext returns a new initialized Context
func NewContext(chainID string, nodeURI string, home string) (*Context, error) {
	var (
		rpc rpcclient.Client
		err error
	)

	if nodeURI != "" {
		rpc = rpcclient.NewHTTP(nodeURI, "/websocket")
	}

	verifier, err := createVerifier(chainID, home, nodeURI)
	if err != nil {
		return nil, err
	}

	return &Context{
		Client:        rpc,
		NodeURI:       nodeURI,
		AccountStore:  AccountStoreKey,
		verifier:      verifier,
		Home:          home,
		BroadcastMode: BroadcastSync,
	}, nil
}

func NewContextWithDelay(chainID string, nodeURI string, home string) (*Context, error) {
	var (
		rpc rpcclient.Client
		err error
		ctx *Context
	)

	if nodeURI != "" {
		rpc = rpcclient.NewHTTP(nodeURI, "/websocket")
	} else {
		return nil, fmt.Errorf("no nodeURI specified")
	}

	var verifier tmlite.Verifier
	ctx = &Context{
		Client:        rpc,
		NodeURI:       nodeURI,
		AccountStore:  AccountStoreKey,
		verifier:      verifier,
		Home:          home,
		BroadcastMode: BroadcastSync,
	}
	go func() {
		for {
			node := rpcclient.NewHTTP(nodeURI, "/websocket")
			st, err := node.Status()
			if err != nil || !node.IsRunning() {
				fmt.Printf("node is not running, status: %#+v", st)
				time.Sleep(time.Second * 4)
			} else {
				break
			}
		}

		verifier, err = createVerifier(chainID, home, nodeURI)
		ctx.WithVerifier(verifier)
	}()

	return ctx, nil
}

func createVerifier(chainID string, home string, nodeURI string) (tmlite.Verifier, error) {
	if chainID == "" {
		return nil, errors.New("Invalid chainID")
	}
	if home == "" {
		return nil, errors.New("Invalid home")
	}
	if nodeURI == "" {
		return nil, errors.New("Invalid nodeURI")
	}

	node := rpcclient.NewHTTP(nodeURI, "/websocket")
	cacheSize := 10 // TODO: determine appropriate cache size
	verifier, err := tmliteProxy.NewVerifier(
		chainID, filepath.Join(home, ".gaialite"),
		node, log.NewNopLogger(), cacheSize,
	)

	if err != nil {
		return nil, err
	}

	return verifier, nil
}

// WithCodec returns a copy of the context with an updated codec.
func (ctx *Context) WithCodec(cdc *codec.Codec) *Context {
	ctx.Codec = cdc
	return ctx
}

// WithHeight returns a copy of the context with an updated height.
func (ctx *Context) WithHeight(height int64) *Context {
	ctx.Height = height
	return ctx
}

// WithHome returns a copy of the context with an updated home.
func (ctx *Context) WithHome(home string) *Context {
	ctx.Home = home
	return ctx
}

// WithOutput returns a copy of the context with an updated output writer (e.g. stdout).
func (ctx *Context) WithOutput(w io.Writer) *Context {
	ctx.Output = w
	return ctx
}

// WithAccountStore returns a copy of the context with an updated AccountStore.
func (ctx *Context) WithAccountStore(accountStore string) *Context {
	ctx.AccountStore = accountStore
	return ctx
}

// WithFrom returns a copy of the context with an updated from address or name.
func (ctx *Context) WithFrom(from string) *Context {
	ctx.From = from
	return ctx
}

// WithTrustNode returns a copy of the context with an updated TrustNode flag.
func (ctx *Context) WithTrustNode(trustNode bool) *Context {
	ctx.TrustNode = trustNode
	return ctx
}

// WithNodeURI returns a copy of the context with an updated node URI.
func (ctx *Context) WithNodeURI(nodeURI string) *Context {
	ctx.NodeURI = nodeURI
	ctx.Client = rpcclient.NewHTTP(nodeURI, "/websocket")
	return ctx
}

// WithClient returns a copy of the context with an updated RPC client
// instance.
func (ctx *Context) WithClient(client rpcclient.Client) *Context {
	ctx.Client = client
	return ctx
}

// WithUseLedger returns a copy of the context with an updated UseLedger flag.
func (ctx *Context) WithUseLedger(useLedger bool) *Context {
	ctx.UseLedger = useLedger
	return ctx
}

// WithPassphrase returns a copy of the context with an passphrase for signing tx.
func (ctx *Context) WithPassphrase(passphrase string) *Context {
	ctx.Passphrase = passphrase
	return ctx
}

// WithPassphrase returns a copy of the context with an private key for signing tx.
func (ctx *Context) WithPrivKey(privKey crypto.PrivKey) *Context {
	ctx.PrivKey = privKey
	return ctx
}

// WithVerifier - return a copy of the context with an updated Verifier
func (ctx *Context) WithVerifier(verifier tmlite.Verifier) *Context {
	ctx.mtx.Lock()
	ctx.mtx.Unlock()
	ctx.verifier = verifier
	return ctx
}

func (ctx *Context) GetVerifier() tmlite.Verifier {
	ctx.mtx.RLock()
	defer ctx.mtx.RUnlock()
	return ctx.verifier
}

// WithGenerateOnly returns a copy of the context with updated GenerateOnly value
func (ctx *Context) WithGenerateOnly(generateOnly bool) *Context {
	ctx.GenerateOnly = generateOnly
	return ctx
}

// WithSimulation returns a copy of the context with updated Simulate value
func (ctx *Context) WithSimulation(simulate bool) *Context {
	ctx.Simulate = simulate
	return ctx
}

// WithFromName returns a copy of the context with an updated from account name.
func (ctx *Context) WithFromName(name string) *Context {
	ctx.FromName = name
	return ctx
}

// WithFromAddress returns a copy of the context with an updated from account
// address.
func (ctx *Context) WithFromAddress(addr sdk.AccAddress) *Context {
	ctx.FromAddress = addr
	return ctx
}

// WithBroadcastMode returns a copy of the context with an updated broadcast
// mode.
func (ctx *Context) WithBroadcastMode(mode string) *Context {
	ctx.BroadcastMode = mode
	return ctx
}

// PrintOutput prints output while respecting output and indent flags
// NOTE: pass in marshalled structs that have been unmarshaled
// because this function will panic on marshaling errors
func (ctx *Context) PrintOutput(toPrint fmt.Stringer) (err error) {
	var out []byte

	switch ctx.OutputFormat {
	case "text":
		out = []byte(toPrint.String())

	case "json":
		if ctx.Indent {
			out, err = ctx.Codec.MarshalJSONIndent(toPrint, "", "  ")
		} else {
			out, err = ctx.Codec.MarshalJSON(toPrint)
		}
	}

	if err != nil {
		return
	}

	fmt.Println(string(out))
	return
}

// GetFromFields returns a from account address and Keybase name given either
// an address or key name. If genOnly is true, only a valid Bech32 cosmos
// address is returned.
func GetFromFields(from string, home string) (sdk.AccAddress, string, error) {
	if from == "" {
		return nil, "", nil
	}

	keybase, err := keys.NewKeyBaseFromDir(home)
	if err != nil {
		return nil, "", err
	}

	var info cryptokeys.Info
	if addr, err := sdk.AccAddressFromBech32(from); err == nil {
		info, err = keybase.GetByAddress(addr)
		if err != nil {
			return nil, "", err
		}
	} else {
		info, err = keybase.Get(from)
		if err != nil {
			return nil, "", err
		}
	}

	return info.GetAddress(), info.GetName(), nil
}
