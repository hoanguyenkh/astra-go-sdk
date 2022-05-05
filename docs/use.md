### Initilized Client

```
    cfg := &config.Config{
        ChainId:       "astra_11110-1",
        Endpoint:      "http://206.189.43.55:26657",
        CoinType:      60,
        PrefixAddress: "astra",
        TokenSymbol:   "aastra",
    }
    
    client := NewClient(cfg)
```

### Create account

```
    accClient := client.NewAccountClient()
    acc, err := accClient.CreateAccount()
    if err != nil {
        panic(err)
    }
    
    data, _ := acc.String()
```

### Create MultiSign Account

```
    accClient := suite.Client.NewAccountClient()
    acc, addr, pubKey, err := accClient.CreateMulSignAccount(3, 2)
    if err != nil {
        panic(err)
    }
    
    fmt.Println("addr", addr)
    fmt.Println("pucKey", pubKey)
    fmt.Println("list key")
    for i, item := range acc {
        fmt.Println("index", i)
        fmt.Println(item.String())
    }
```

### Transfer Token

```
    bankClient := suite.Client.NewBankClient()
    
    amount := big.NewInt(0).Mul(big.NewInt(10), big.NewInt(0).SetUint64(uint64(math.Pow10(18))))
     
    request := &bank.TransferRequest{
        PrivateKey: "saddle click spawn install mutual visa usage eyebrow awesome inherit rifle moon giraffe deposit reduce east gossip ice salute hill fire require knife traffic",
        Receiver:   "astra156dh69y8j39eynue4jahrezg32rgl8eck5rhsl",
        Amount:   amount,
        GasLimit: 200000,
        GasPrice: "0.001aastra",
    }
    
    txBuilder, err := bankClient.TransferRawData(request)
    if err != nil {
        panic(err)
    }
    
    txJson, err := common.TxBuilderJsonEncoder(suite.Client.rpcClient.TxConfig, txBuilder)
    if err != nil {
        panic(err)
    }
    
    fmt.Println("rawData", string(txJson))
    
    txByte, err := common.TxBuilderJsonDecoder(suite.Client.rpcClient.TxConfig, txJson)
    if err != nil {
        panic(err)
    }
    
    txHash := common.TxHash(txByte)
    fmt.Println("txHash", txHash)
    
    res, err := suite.Client.rpcClient.BroadcastTxCommit(txByte)
    if err != nil {
        panic(err)
    }
    fmt.Println(res)
```