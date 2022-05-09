package client

import (
	"fmt"
	"github.com/AstraProtocol/astra-go-sdk/bank"
	"github.com/AstraProtocol/astra-go-sdk/common"
	"github.com/AstraProtocol/astra-go-sdk/config"
	"github.com/cosmos/cosmos-sdk/types"
	signingTypes "github.com/cosmos/cosmos-sdk/types/tx/signing"
	"github.com/stretchr/testify/assert"
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
	accClient := suite.Client.NewAccountClient()
	acc, err := accClient.CreateAccount()
	if err != nil {
		panic(err)
	}

	data, _ := acc.String()

	fmt.Println(data)
}

func (suite *AstraSdkTestSuite) TestGenMulSignAccount() {
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

	amount := big.NewInt(0).Mul(big.NewInt(10), big.NewInt(0).SetUint64(uint64(math.Pow10(18))))
	fmt.Println("amount", amount.String())

	request := &bank.TransferRequest{
		PrivateKey: "saddle click spawn install mutual visa usage eyebrow awesome inherit rifle moon giraffe deposit reduce east gossip ice salute hill fire require knife traffic",
		Receiver:   "astra156dh69y8j39eynue4jahrezg32rgl8eck5rhsl",
		//Receiver: "astra1p394e9sg72fq48z8k2skdh8s270n5zm04kzffu",
		//Receiver: "astra147hn2qzhrdcrcw792xg2y47y3q0fsg7rg8wfh9",
		//Receiver: "astra1pu3yrnyg9mq3lj43r6t3u08mdwq20fqj2gahlc",
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
}

func (suite *AstraSdkTestSuite) TestTransferMultiSign() {
	/*
			 addr astra18dgn6vxsyk69xglsp8z0r6ltc5q2slzc2nglwd
		pucKey {"@type":"/cosmos.crypto.multisig.LegacyAminoPubKey","threshold":2,"public_keys":[{"@type":"/ethermint.crypto.v1.ethsecp256k1.PubKey","key":"A0UjEVXxXA7JY2oou5HPH7FuPSyJ2hAfDMc4XThXiopM"},{"@type":"/ethermint.crypto.v1.ethsecp256k1.PubKey","key":"A6DFr74kQmk/k88fCTPCxmf9kyFJMhFUF21IPFY7XoV2"},{"@type":"/ethermint.crypto.v1.ethsecp256k1.PubKey","key":"AgPQELGzKmlAaSb01OKbmuL1f17MHJshkh9s9xAWxMa3"}]}
		list key
		index 0
		{
		 "address": "astra1p394e9sg72fq48z8k2skdh8s270n5zm04kzffu",
		 "hexAddress": "0x0c4b5c9608f2920a9C47B2A166dcf0579F3A0B6f",
		 "mnemonic": "project fat comfort regular strong dream crack palace boost act reform minor rack where vicious photo cat pass bounce dune fuel movie tennis sausage",
		 "privateKey": "e7c769631951cf96909cefcf1352e50f15d4d9644b184d0500b5004edd686a6a",
		 "publicKey": "{\"@type\":\"/ethermint.crypto.v1.ethsecp256k1.PubKey\",\"key\":\"A0UjEVXxXA7JY2oou5HPH7FuPSyJ2hAfDMc4XThXiopM\"}",
		 "type": "eth_secp256k1",
		 "validatorKey": "astravaloper1p394e9sg72fq48z8k2skdh8s270n5zm0s0rcjj"
		} <nil>
		index 1
		{
		 "address": "astra147hn2qzhrdcrcw792xg2y47y3q0fsg7rg8wfh9",
		 "hexAddress": "0xAfaF3500571b703c3bc55190a257C4881e9823c3",
		 "mnemonic": "embody thrive world there siren afraid sport utility dove rural few guess grid own strategy orbit vacuum soft gold muffin wrestle shoulder detect record",
		 "privateKey": "6be1f98d3fb3421bd965b4aecaabecce1035d92999cc12d64d7af2ccb9f99c68",
		 "publicKey": "{\"@type\":\"/ethermint.crypto.v1.ethsecp256k1.PubKey\",\"key\":\"A6DFr74kQmk/k88fCTPCxmf9kyFJMhFUF21IPFY7XoV2\"}",
		 "type": "eth_secp256k1",
		 "validatorKey": "astravaloper147hn2qzhrdcrcw792xg2y47y3q0fsg7rd70cvt"
		} <nil>
		index 2
		{
		 "address": "astra1pu3yrnyg9mq3lj43r6t3u08mdwq20fqj2gahlc",
		 "hexAddress": "0x0F2241CC882EC11fcAb11E971e3Cfb6B80a7A412",
		 "mnemonic": "property cactus cannon talent priority silk ice nurse such arctic dove wonder blue stumble chalk engine start know unable tool arctic tone sugar grass",
		 "privateKey": "d85f14f8c268f516f6cfcbc293a687f6918db94a29a6cacd59a28e10ea9908a1",
		 "publicKey": "{\"@type\":\"/ethermint.crypto.v1.ethsecp256k1.PubKey\",\"key\":\"AgPQELGzKmlAaSb01OKbmuL1f17MHJshkh9s9xAWxMa3\"}",
		 "type": "eth_secp256k1",
		 "validatorKey": "astravaloper1pu3yrnyg9mq3lj43r6t3u08mdwq20fqj03uxyk"
		} <nil>

	*/

	pk, err := common.DecodePublicKey(suite.Client.rpcClient, "{\"@type\":\"/cosmos.crypto.multisig.LegacyAminoPubKey\",\"threshold\":2,\"public_keys\":[{\"@type\":\"/ethermint.crypto.v1.ethsecp256k1.PubKey\",\"key\":\"A0UjEVXxXA7JY2oou5HPH7FuPSyJ2hAfDMc4XThXiopM\"},{\"@type\":\"/ethermint.crypto.v1.ethsecp256k1.PubKey\",\"key\":\"A6DFr74kQmk/k88fCTPCxmf9kyFJMhFUF21IPFY7XoV2\"},{\"@type\":\"/ethermint.crypto.v1.ethsecp256k1.PubKey\",\"key\":\"AgPQELGzKmlAaSb01OKbmuL1f17MHJshkh9s9xAWxMa3\"}]}")
	if err != nil {
		panic(err)
	}

	//astra18dgn6vxsyk69xglsp8z0r6ltc5q2slzc2nglwd
	from := types.AccAddress(pk.Address())
	fmt.Println("from", from.String())

	bankClient := suite.Client.NewBankClient()

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

	res, err := suite.Client.rpcClient.BroadcastTxCommit(txByte)
	if err != nil {
		panic(err)
	}

	fmt.Println(res)
}

func (suite *AstraSdkTestSuite) TestAddresValid() {
	addressCheck := "astra1hann2zj3sx3ympd40ptxdmpd4nd4eypm45zhhr"
	receiver, err := types.AccAddressFromBech32(addressCheck)
	if err != nil {
		panic(err)
	}

	fmt.Println(receiver.String())
	assert.Equal(suite.T(), addressCheck, receiver.String(), "they should be equal")

	rs, _ := common.IsAddressValid(addressCheck)
	assert.Equal(suite.T(), rs, true)
}

func (suite *AstraSdkTestSuite) TestConvertHexToCosmosAddress() {
	//"address": "astra147hn2qzhrdcrcw792xg2y47y3q0fsg7rg8wfh9",
	//"hexAddress": "0xAfaF3500571b703c3bc55190a257C4881e9823c3",
	rs, _ := common.EthAddressToCosmosAddress("AfaF3500571b703c3bc55190a257C4881e9823c3")
	fmt.Println(rs)
	assert.Equal(suite.T(), rs, "astra147hn2qzhrdcrcw792xg2y47y3q0fsg7rg8wfh9")
}
