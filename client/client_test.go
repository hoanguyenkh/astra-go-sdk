package client

import (
	"fmt"
	"github.com/AstraProtocol/astra-go-sdk/bank"
	"github.com/AstraProtocol/astra-go-sdk/common"
	"github.com/AstraProtocol/astra-go-sdk/config"
	"github.com/stretchr/testify/suite"
	"math"
	"math/big"
	"testing"
)

type AstraSdkTestSuite struct {
	suite.Suite
	Client *Client
}

func (suite *AstraSdkTestSuite) SetupTest() {
	cfg := &config.Config{
		ChainId:       "astra_11110-1",
		Endpoint:      "http://206.189.43.55:26657",
		CoinType:      60,
		PrefixAddress: "astra",
		TokenSymbol:   "aastra",
	}

	client := NewClient(cfg)
	suite.Client = client
}

func TestAstraSdkTestSuite(t *testing.T) {
	suite.Run(t, new(AstraSdkTestSuite))
}

func (suite *AstraSdkTestSuite) TestInitBank() {
	bankClient := suite.Client.NewBankClient()
	balance, err := bankClient.Balance("astra12nnueg3904ukfjel4u695ma6tvrkqvqmrqstx6")
	if err != nil {
		panic(err)
	}

	fmt.Println(balance.String())
}

func (suite *AstraSdkTestSuite) TestGenAccount() {
	client := suite.Client.NewAccountClient()
	acc, err := client.CreateAccount()
	if err != nil {
		panic(err)
	}

	data, _ := acc.String()

	fmt.Println(data)
}

func (suite *AstraSdkTestSuite) TestTransfer() {
	/*
			{
		 "address": "astra156dh69y8j39eynue4jahrezg32rgl8eck5rhsl",
		 "hexAddress": "0xa69b7d1487944b924f99ACBb71e4488a868F9f38",
		 "mnemonic": "strike enhance tray always bulb pioneer message hair swim alone soda possible print vivid delay winner guilt test skirt liar slab fog balance fence",
		 "privateKey": "5efa3b891b24832c185f3133f4eb978345d3f1563346ff26465dac8f1832dcc1",
		 "publicKey": "{\"@type\":\"/ethermint.crypto.v1.ethsecp256k1.PubKey\",\"key\":\"Ag/YFj8pCpDgLkz+B4D2IrRplrqZHyCf4oDy5As9AFbI\"}",
		 "type": "eth_secp256k1",
		 "validatorKey": "astravaloper156dh69y8j39eynue4jahrezg32rgl8ecndzxt3"
		}
	*/

	bankClient := suite.Client.NewBankClient()

	amount := big.NewInt(0).Mul(big.NewInt(100), big.NewInt(0).SetUint64(uint64(math.Pow10(18))))
	fmt.Println("amount", amount.String())

	request := &bank.TransferRequest{
		PrivateKey: "saddle click spawn install mutual visa usage eyebrow awesome inherit rifle moon giraffe deposit reduce east gossip ice salute hill fire require knife traffic",
		Receiver:   "astra156dh69y8j39eynue4jahrezg32rgl8eck5rhsl",
		Amount:     amount,
		GasLimit:   200000,
		GasPrice:   "0.001aastra",
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
}
