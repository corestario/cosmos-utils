module github.com/corestario/cosmos-utils/client

go 1.12

require (
	github.com/bgentry/speakeasy v0.1.0
	github.com/cosmos/cosmos-sdk v0.28.2-0.20190827131926-5aacf454e1b6
	github.com/mattn/go-isatty v0.0.12
	github.com/pkg/errors v0.9.1
	github.com/spf13/viper v1.6.3
	github.com/stretchr/testify v1.5.1
	github.com/tendermint/go-amino v0.15.1
	github.com/tendermint/iavl v0.13.3
	github.com/tendermint/tendermint v0.33.4
)

replace golang.org/x/crypto => github.com/tendermint/crypto v0.0.0-20180820045704-3764759f34a5

replace github.com/cosmos/cosmos-sdk => github.com/cosmos/cosmos-sdk v0.34.4-0.20200507135526-b3cada10017d
