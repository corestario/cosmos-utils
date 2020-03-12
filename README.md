# cosmos_utils
repo for common cosmos utils and wrappers 

## Client
Some code from cosmos-sdk for querying and sending transactions to a Cosmos application but with explicit configuration.

It allows you to conveniently and easily send transactions or make queries to a Cosmos-application directly from a code.

Example:
```go
cliCtx, err := cliCTX.NewCLIContext(chainID, nodeEndpoint, validatorName, genOnly, broadcastMode, vfrHome, height, trustNode, cliHome, "")
if err != nil {
    return nil, nil, err
}
cliCtx = cliCtx.WithCodec(cdc).WithAccountDecoder(cdc)
accNumber, err := cliCtx.GetAccountNumber(cliCtx.FromAddress)
if err != nil {
    return nil, nil, err
}
kb, err := keys.NewKeyBaseFromDir(cliCtx.Home)
if err != nil {
    return nil, nil, err
}
// TxBuilder implements tx generation
txBldr := authtxb.NewTxBuilder(utils.GetTxEncoder(cdc), accNumber, 0, 0, 0.0, false, cliCtx.Verifier.ChainID(), "", nil, nil).WithKeybase(kb)
if err := cliCtx.EnsureAccountExists(); err != nil {
    return nil, nil, fmt.Errorf("failed to find account: %v", err)
}

//Query some data from an app
res, _, err := cli.QueryWithData("custom/app/SOME_ENDPOINT", nil)

//Send transaction to an app
msg := msgs.NewSomeMsg(item, cli.GetFromAddress())
err = utils.GenerateOrBroadcastMsgs(*cli, *txBldr, msg, false)
```

## StoreWrapper
The Cosmos KVStore has limit on size of the value, so the wrapper divide large value on little pieces and stores them separately.

And then it combines them when getting a value

Example:
```go
ctx := sdk.NewContext(ms, abci.Header{}, false, log.NewNopLogger())
stw := storewrapper.NewKVStore(ctx.KVStore(testKey), 0)
key := "key"


_, err := stw.SetW(key, []byte(LARGE_STRING))

isValueExist, _, err := stw.HasW(key)

res, err := stw.GetW(key) // res contains large value

```