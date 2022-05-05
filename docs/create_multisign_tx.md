### Transfer Token With MultiSign Account

```
    pk, err := common.DecodePublicKey(suite.Client.rpcClient, "{\"@type\":\"/cosmos.crypto.multisig.LegacyAminoPubKey\",\"threshold\":2,\"public_keys\":[{\"@type\":\"/ethermint.crypto.v1.ethsecp256k1.PubKey\",\"key\":\"A0UjEVXxXA7JY2oou5HPH7FuPSyJ2hAfDMc4XThXiopM\"},{\"@type\":\"/ethermint.crypto.v1.ethsecp256k1.PubKey\",\"key\":\"A6DFr74kQmk/k88fCTPCxmf9kyFJMhFUF21IPFY7XoV2\"},{\"@type\":\"/ethermint.crypto.v1.ethsecp256k1.PubKey\",\"key\":\"AgPQELGzKmlAaSb01OKbmuL1f17MHJshkh9s9xAWxMa3\"}]}"
    
    if err != nil {
        panic(err)
    }
	
    from := types.AccAddress(pk.Address())
    
    bankClient := client.NewBankClient()
    
    amount := big.NewInt(0).Mul(big.NewInt(10), big.NewInt(0).SetUint64(uint64(math.Pow10(18))))
    fmt.Println("amount", amount.String())
    
    listPrivate := []string{
        "project fat comfort regular strong dream crack palace boost act reform minor rack where vicious photo cat pass bounce dune fuel movie tennis sausage",
        "embody thrive world there siren afraid sport utility dove rural few guess grid own strategy orbit vacuum soft gold muffin wrestle shoulder detect record",
        "property cactus cannon talent priority silk ice nurse such arctic dove wonder blue stumble chalk engine start know unable tool arctic tone sugar grass",
    }
    
    signList := make([][]signingTypes.SignatureV2, 0)
    for i, s := range listPrivate {
        fmt.Println("index", i)
        request := &bank.SignTxWithSignerAddressRequest{
            SignerPrivateKey:    s,
            MulSignAccPublicKey: pk,
            Receiver:            "astra156dh69y8j39eynue4jahrezg32rgl8eck5rhsl",
            Amount:              amount,
            GasLimit:            200000,
            GasPrice:            "0.001aastra",
        }
    
        txBuilder, err := bankClient.SignTxWithSignerAddress(request)
        if err != nil {
            panic(err)
        }
    
        sign, err := common.TxBuilderSignatureJsonEncoder(suite.Client.rpcClient.TxConfig, txBuilder)
        if err != nil {
            panic(err)
        }
    
        fmt.Println("sign-data", string(sign))
    
        signByte, err := common.TxBuilderSignatureJsonDecoder(suite.Client.rpcClient.TxConfig, sign)
        if err != nil {
            panic(err)
        }
    
        signList = append(signList, signByte)
    }
    
    fmt.Println("start transfer")
    
    request := &bank.TransferMultiSignRequest{
        MulSignAccPublicKey: pk,
        Receiver:            "astra156dh69y8j39eynue4jahrezg32rgl8eck5rhsl",
        Amount:              amount,
        GasLimit:            200000,
        GasPrice:            "0.001aastra",
        Sigs:                signList,
    }
    
    txBuilder, err := bankClient.TransferMultiSignRawData(request)
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