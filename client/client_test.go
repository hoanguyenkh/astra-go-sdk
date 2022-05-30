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
	"os"
	"testing"
)

type AstraSdkTestSuite struct {
	suite.Suite
	Client *Client
}

func (suite *AstraSdkTestSuite) SetupTest() {
	cfg := &config.Config{
		ChainId:       os.Getenv("CHAIN_ID"),
		Endpoint:      os.Getenv("END_POINT"),
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
	balance, err := bankClient.Balance("astra1gvwtjcv36nfqe8w3qyem6h600n750jqg6a576j")
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

func (suite *AstraSdkTestSuite) TestTransferMultiSign() {
	/*
		    addr astra1ha0vgh05zzlwdeejxq9aq7gqr6jzs7stdhlfra
			pucKey {"@type":"/cosmos.crypto.multisig.LegacyAminoPubKey","threshold":2,"public_keys":[{"@type":"/ethermint.crypto.v1.ethsecp256k1.PubKey","key":"A0ATAOfWQM6XXCA5po9DBsKVGmWudnIN55arHhDYhR89"},{"@type":"/ethermint.crypto.v1.ethsecp256k1.PubKey","key":"A0ks8ww7AVKYQRsKgZSQi9wTfoQzKNt30gLOMpOJNSPn"},{"@type":"/ethermint.crypto.v1.ethsecp256k1.PubKey","key":"A9Q4nSS73SG+Tclghh1JEtfng5vd41dgmG7HJrYW4/Ml"}]}
			list key
			index 0
			{
			 "address": "astra1dmdsy082730stdletm7z6zulfxuez4lsx3tztx",
			 "hexAddress": "0x6Edb023ceAF45F05b7f95efC2d0B9f49B99157F0",
			 "mnemonic": "ignore risk morning strike school street radar silk recipe health december system inflict gold foster item end twenty magic shine oppose island loop impact",
			 "privateKey": "7f1d3df4044f09b1edfab34c7e3fee92396ea23861e96a8ac7429efcf158d794",
			 "publicKey": "{\"@type\":\"/ethermint.crypto.v1.ethsecp256k1.PubKey\",\"key\":\"A0ATAOfWQM6XXCA5po9DBsKVGmWudnIN55arHhDYhR89\"}",
			 "type": "eth_secp256k1",
			 "validatorKey": "astravaloper1dmdsy082730stdletm7z6zulfxuez4lsrg2nsg"
			} <nil>
			index 1
			{
			 "address": "astra1fd39nlc4hsl7ma9knpjwlhcrnunz66dnvf5agx",
			 "hexAddress": "0x4b6259ff15Bc3FEdf4B69864EfdF039F262d69B3",
			 "mnemonic": "seven mean snap illness couch excite item topic tobacco erosion tourist blue van possible wolf gadget combine excess brush goddess glory subway few mind",
			 "privateKey": "8dca20a27b0bfdcf1dacc9b2f71d4b7e7d269a4b87949707c12ef2ba328fd0e9",
			 "publicKey": "{\"@type\":\"/ethermint.crypto.v1.ethsecp256k1.PubKey\",\"key\":\"A0ks8ww7AVKYQRsKgZSQi9wTfoQzKNt30gLOMpOJNSPn\"}",
			 "type": "eth_secp256k1",
			 "validatorKey": "astravaloper1fd39nlc4hsl7ma9knpjwlhcrnunz66dnfs4vng"
			} <nil>
			index 2
			{
			 "address": "astra1gc0v03kjrg9uv7duvzqsndv3nhkhehvkwuhkdr",
			 "hexAddress": "0x461EC7C6D21a0BC679bC608109b5919DEd7Cdd96",
			 "mnemonic": "swap exhaust letter left light trust diet piano pride rifle trust orbit clip suggest achieve unaware please guess lawsuit doctor use bargain jealous weekend",
			 "privateKey": "e3f46776e933129611b3cb6418176dcd2a9badd8188fb4804d5b822548200bac",
			 "publicKey": "{\"@type\":\"/ethermint.crypto.v1.ethsecp256k1.PubKey\",\"key\":\"A9Q4nSS73SG+Tclghh1JEtfng5vd41dgmG7HJrYW4/Ml\"}",
			 "type": "eth_secp256k1",
			 "validatorKey": "astravaloper1gc0v03kjrg9uv7duvzqsndv3nhkhehvkt9k8kd"
			}
	*/

	pk, err := common.DecodePublicKey(
		suite.Client.rpcClient,
		"{\"@type\":\"/cosmos.crypto.multisig.LegacyAminoPubKey\",\"threshold\":2,\"public_keys\":[{\"@type\":\"/ethermint.crypto.v1.ethsecp256k1.PubKey\",\"key\":\"A0ATAOfWQM6XXCA5po9DBsKVGmWudnIN55arHhDYhR89\"},{\"@type\":\"/ethermint.crypto.v1.ethsecp256k1.PubKey\",\"key\":\"A0ks8ww7AVKYQRsKgZSQi9wTfoQzKNt30gLOMpOJNSPn\"},{\"@type\":\"/ethermint.crypto.v1.ethsecp256k1.PubKey\",\"key\":\"A9Q4nSS73SG+Tclghh1JEtfng5vd41dgmG7HJrYW4/Ml\"}]}",
	)
	if err != nil {
		panic(err)
	}

	from := types.AccAddress(pk.Address())
	fmt.Println("from", from.String())

	bankClient := suite.Client.NewBankClient()

	amount := big.NewInt(0).Mul(big.NewInt(10), big.NewInt(0).SetUint64(uint64(math.Pow10(18))))
	fmt.Println("amount", amount.String())

	listPrivate := []string{
		"ignore risk morning strike school street radar silk recipe health december system inflict gold foster item end twenty magic shine oppose island loop impact",
		"seven mean snap illness couch excite item topic tobacco erosion tourist blue van possible wolf gadget combine excess brush goddess glory subway few mind",
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

	//200
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

func (suite *AstraSdkTestSuite) TestAddressValid() {
	addressCheck := "astra1hann2zj3sx3ympd40ptxdmpd4nd4eypm45zhhr"
	addressCheck = "astra19a3mu6k0y326mcny60m3x70qfxtkms20sn5j8p"
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
	rs, _ := common.EthAddressToCosmosAddress("AfaF3500571b703c3bc55190a257C4881e9823c3")
	fmt.Println(rs)
	assert.Equal(suite.T(), rs, "astra147hn2qzhrdcrcw792xg2y47y3q0fsg7rg8wfh9")
}

func (suite *AstraSdkTestSuite) TestCheckTx() {
	bankClient := suite.Client.NewBankClient()
	rs, err := bankClient.CheckTx("B593ED27B335EAB7DC664C39B56D0BBE36A4977F25B66FAD75503F5EECE21321")
	if err != nil {
		panic(err)
	}

	if rs != nil && common.IsBlocked(rs.Code) {
		fmt.Println("blocked")
	}
}
