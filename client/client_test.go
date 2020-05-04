package client

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/corestario/cosmos-utils/client/context"
	"testing"
)

const (
	nodeEndpoint  = "tcp://localhost:26657"
	chainID       = "mpchain"
	trustNode     = false
	validatorName = "user1"

	cliHome = "/Users/pr0n00gler/.mpcli"
)

func TestWow(t *testing.T) {
	cdc := codec.New()

	ctx, err := context.NewContext(chainID, nodeEndpoint, cliHome)
	if err != nil {
		t.Fatal(err.Error())
		return
	}
	ctx = ctx.WithCodec(cdc)
}