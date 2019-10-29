module github.com/dgamingfoundation/cosmos-utils/client

go 1.12

require (
	github.com/bgentry/speakeasy v0.1.0
	github.com/cosmos/cosmos-sdk v0.28.2-0.20190827131926-5aacf454e1b6
	github.com/mattn/go-isatty v0.0.7
	github.com/pkg/errors v0.8.1
	github.com/stretchr/testify v1.4.0
	github.com/tendermint/go-amino v0.15.0
	github.com/tendermint/tendermint v0.32.4
)

replace golang.org/x/crypto => github.com/tendermint/crypto v0.0.0-20180820045704-3764759f34a5
