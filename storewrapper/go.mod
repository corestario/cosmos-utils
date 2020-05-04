module github.com/corestario/cosmos-utils/storewrapper

go 1.12

require (
	github.com/cosmos/cosmos-sdk v0.28.2-0.20190827131926-5aacf454e1b6
	github.com/stretchr/testify v1.5.1
	github.com/tendermint/tendermint v0.32.8
)

replace github.com/tendermint/tendermint => ./../../tendermint

replace github.com/cosmos/cosmos-sdk => ./../../cosmos-sdk
