package main

import (
	"fmt"
	"github.com/AstraProtocol/astra-go-sdk/client"
	"github.com/AstraProtocol/astra-go-sdk/config"
)

func main() {
	cfg := &config.Config{
		ChainId:       "astra_11110-1",
		Endpoint:      "http://206.189.43.55:26657",
		CoinType:      60,
		PrefixAddress: "astra",
		TokenSymbol:   "aastra",
	}

	astraClient := client.NewClient(cfg)

	keyWallet := astraClient.NewAccountClient()
	acc, err := keyWallet.CreateAccount()
	if err != nil {
		panic(err)
	}

	data, _ := acc.String()
	fmt.Println(data)
}
