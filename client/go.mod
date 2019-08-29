module github.com/dgamingfoundation/cosmos-utils/client

go 1.12

require (
	github.com/bgentry/speakeasy v0.1.0
	github.com/cosmos/cosmos-sdk v0.36.0
	github.com/golang/mock v1.3.1-0.20190508161146-9fa652df1129
	github.com/mattn/go-isatty v0.0.7
	github.com/mattn/go-runewidth v0.0.4 // indirect
	github.com/olekukonko/tablewriter v0.0.1
	github.com/pkg/errors v0.8.1
	github.com/spf13/cobra v0.0.5
	github.com/spf13/viper v1.4.0
	github.com/stretchr/testify v1.3.0
	github.com/tendermint/go-amino v0.15.0
	github.com/tendermint/tendermint v0.31.7
	gopkg.in/yaml.v2 v2.2.2
)

replace golang.org/x/crypto => github.com/tendermint/crypto v0.0.0-20180820045704-3764759f34a5

replace github.com/cosmos/cosmos-sdk => github.com/cosmos/cosmos-sdk v0.0.0-20190625145210-5fb86c661ea1
