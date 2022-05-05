### Transfer Token

```
    bankClient := client.NewBankClient()
    
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
    
    res, err := client.rpcClient.BroadcastTxCommit(txByte)
    if err != nil {
        panic(err)
    }
    fmt.Println(res)
```