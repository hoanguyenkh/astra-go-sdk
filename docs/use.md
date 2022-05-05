### Initilized Client

```
import (
    "fmt"
    "github.com/AstraProtocol/astra-go-sdk/client"
)

func main() {
    cfg := &config.Config{
        ChainId:       "astra_11110-1",
        Endpoint:      "http://206.189.43.55:26657",
        CoinType:      60,
        PrefixAddress: "astra",
        TokenSymbol:   "aastra",
    }
    
    client := NewClient(cfg)
    
    //todo
}
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