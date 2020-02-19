module github.com/corestario/cosmos-utils/client

go 1.12

require (
	github.com/bgentry/speakeasy v0.1.0
	github.com/cosmos/cosmos-sdk v0.28.2-0.20190827131926-5aacf454e1b6
	github.com/mattn/go-isatty v0.0.11
	github.com/pkg/errors v0.8.1
	github.com/stretchr/testify v1.4.0
	github.com/tendermint/go-amino v0.15.1
	github.com/tendermint/tendermint v0.32.8
)

replace (
	github.com/corestario/dkglib => github.com/corestario/dkglib v0.1.0
	github.com/cosmos/cosmos-sdk => github.com/corestario/cosmos-sdk v0.1.0
	github.com/tendermint/tendermint => github.com/corestario/tendermint v0.1.0
	go.dedis.ch/kyber/v3 => github.com/corestario/kyber/v3 v3.0.0-20200218082721-8ed10c357c05
	golang.org/x/crypto => github.com/tendermint/crypto v0.0.0-20180820045704-3764759f34a5
)
