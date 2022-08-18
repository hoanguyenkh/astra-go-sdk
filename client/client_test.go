package client

import (
	"encoding/json"
	"fmt"
	"github.com/AstraProtocol/astra-go-sdk/bank"
	"github.com/AstraProtocol/astra-go-sdk/common"
	"github.com/AstraProtocol/astra-go-sdk/config"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/types"
	signingTypes "github.com/cosmos/cosmos-sdk/types/tx/signing"
	ethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"log"
	"math"
	"math/big"
	"os"
	"sync"
	"testing"
)

type AstraSdkTestSuite struct {
	suite.Suite
	Client *Client
}

func (suite *AstraSdkTestSuite) SetupTest() {
	err := godotenv.Load("./../dev.env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

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
	bankClient := suite.Client.NewBankClient()

	amount := big.NewInt(0).Mul(big.NewInt(10), big.NewInt(0).SetUint64(uint64(math.Pow10(18))))
	fmt.Println("amount", amount.String())

	request := &bank.TransferRequest{
		PrivateKey: "valve season sauce knife burden benefit zone field ask carpet fury vital action donate trade street ability artwork ball uniform garbage sugar warm differ",
		Receiver:   "astra1p6sscujfpygmrrxqlwqeqqw6r5lxk2x9gz9glh",
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

	fmt.Println(ethCommon.BytesToHash(txByte).String())

	res, err := suite.Client.rpcClient.BroadcastTxCommit(txByte)
	if err != nil {
		panic(err)
	}

	fmt.Println(res)
}

func (suite *AstraSdkTestSuite) TestTransferMultiSign() {
	//main address
	/*
		addr astra1ha0vgh05zzlwdeejxq9aq7gqr6jzs7stdhlfra
		pucKey {"@type":"/cosmos.crypto.multisig.LegacyAminoPubKey","threshold":2,"public_keys":[{"@type":"/ethermint.crypto.v1.ethsecp256k1.PubKey","key":"A0ATAOfWQM6XXCA5po9DBsKVGmWudnIN55arHhDYhR89"},{"@type":"/ethermint.crypto.v1.ethsecp256k1.PubKey","key":"A0ks8ww7AVKYQRsKgZSQi9wTfoQzKNt30gLOMpOJNSPn"},{"@type":"/ethermint.crypto.v1.ethsecp256k1.PubKey","key":"A9Q4nSS73SG+Tclghh1JEtfng5vd41dgmG7HJrYW4/Ml"}]}
	*/

	//child address
	/*
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

	listPrivate := []string{
		"ignore risk morning strike school street radar silk recipe health december system inflict gold foster item end twenty magic shine oppose island loop impact",
		"seven mean snap illness couch excite item topic tobacco erosion tourist blue van possible wolf gadget combine excess brush goddess glory subway few mind",
	}

	thread := 2
	listRawdata := make([][]byte, 0)

	for i := 0; i < thread; i++ {
		amount := big.NewInt(0).Mul(big.NewInt(10+int64(i)), big.NewInt(0).SetUint64(uint64(math.Pow10(18))))
		fmt.Println("amount", amount.String())

		fmt.Println("start signer")
		signList := make([][]signingTypes.SignatureV2, 0)
		for i, s := range listPrivate {
			fmt.Println("index", i)
			request := &bank.SignTxWithSignerAddressRequest{
				SignerPrivateKey: s,
				SignerPublicKey:  pk,
				Receiver:         "astra156dh69y8j39eynue4jahrezg32rgl8eck5rhsl",
				Amount:           amount,
				GasLimit:         200000,
				GasPrice:         "0.001aastra",
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

		listRawdata = append(listRawdata, txByte)
	}

	var wg sync.WaitGroup
	wg.Add(thread)

	go func(item []byte, client client.Context) {
		_, err := client.BroadcastTxCommit(item)
		if err != nil {
			panic(err)
		}

		//fmt.Println("BroadcastTxCommit", res)
		defer wg.Done()
	}(listRawdata[0], suite.Client.rpcClient)

	go func(item []byte, client client.Context) {
		_, err := client.BroadcastTxCommit(item)
		if err != nil {
			panic(err)
		}

		//fmt.Println("BroadcastTxCommit", res)
		defer wg.Done()
	}(listRawdata[1], suite.Client.rpcClient)

	wg.Wait()
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
	eth := "0xAfaF3500571b703c3bc55190a257C4881e9823c3"
	cosmos := "astra147hn2qzhrdcrcw792xg2y47y3q0fsg7rg8wfh9"

	rs, _ := common.EthAddressToCosmosAddress(eth)
	fmt.Println(rs)
	assert.Equal(suite.T(), cosmos, rs)

	rs1, _ := common.CosmosAddressToEthAddress(cosmos)
	fmt.Println(rs1)
	assert.Equal(suite.T(), eth, rs1)
}

func (suite *AstraSdkTestSuite) TestCheckTx() {
	bankClient := suite.Client.NewBankClient()
	//rs, err := bankClient.CheckTx("646F944DCDB201F674C109E6EF9A594ADBCC33B8F0FA054D7B3F4ABE4CCA2AEB")
	rs, err := bankClient.CheckTx("25D2704C3ABDFE3DBBC1A8202A15D43A2E86D7F1F24AD2E704A7F50FCB75FB94")
	if err != nil {
		panic(err)
	}

	fmt.Println(rs.Code)
	if rs != nil && common.BlockedStatus(rs.Code) {
		fmt.Println("blocked")
	}
}

func (suite *AstraSdkTestSuite) TestImportAccountViaHdPath() {
	accClient := suite.Client.NewAccountClient()

	_, err := common.VerifyHdPath("m/44'/60'/0'/0/0")
	if err != nil {
		panic(err)
	}

	nmemonic := "secret immense amount trial polar security mother scare useful hen squeeze confirm right size best trash team clock matter grow copy quiz capital ill"

	for i := 100083357; i <= (100083357 + 20); i++ {
		s := fmt.Sprintf("m/44'/60'/%v'/1/0", i)
		wallet, err := accClient.ImportHdPath(
			nmemonic,
			s,
		)

		if err != nil {
			panic(err)
		}

		fmt.Println("index ", i, s)
		fmt.Println(wallet.String())
	}
}

func (suite *AstraSdkTestSuite) TestImportByNmemonic() {
	accClient := suite.Client.NewAccountClient()
	key, err := accClient.ImportAccount("effort behave trash gaze youth food north brain poverty drive armed split kind script fox frog breeze cliff bright raise napkin question payment upset")
	if err != nil {
		panic(err)
	}

	fmt.Println(key.String())
}

func (suite *AstraSdkTestSuite) TestImportByPrivatekey() {
	accClient := suite.Client.NewAccountClient()
	key, err := accClient.ImportPrivateKey("b8f7f2e5bab9c0b08df50cb5aa93ca8d1f5fe4aa11677ebf05232930d28349a9")
	if err != nil {
		panic(err)
	}

	fmt.Println(key.String())
}

func (suite *AstraSdkTestSuite) TestScanner() {
	client := suite.Client.NewScanner()
	//listTx, err := client.ScanByBlockHeight(1032719) //cosmos
	//listTx, err := client.ScanByBlockHeight(450666) //eth
	//listTx, err := client.ScanByBlockHeight(596365) //error
	listTx, err := client.ScanByBlockHeight(1036356) //error

	if err != nil {
		panic(err)
	}

	rs, _ := json.MarshalIndent(listTx, " ", " ")

	fmt.Println(string(rs))

}

func (suite *AstraSdkTestSuite) TestSequenceNumberFromPk() {
	mulSignAccPubKey := "{\"@type\":\"/cosmos.crypto.multisig.LegacyAminoPubKey\",\"threshold\":2,\"public_keys\":[{\"@type\":\"/ethermint.crypto.v1.ethsecp256k1.PubKey\",\"key\":\"A0UjEVXxXA7JY2oou5HPH7FuPSyJ2hAfDMc4XThXiopM\"},{\"@type\":\"/ethermint.crypto.v1.ethsecp256k1.PubKey\",\"key\":\"A6DFr74kQmk/k88fCTPCxmf9kyFJMhFUF21IPFY7XoV2\"},{\"@type\":\"/ethermint.crypto.v1.ethsecp256k1.PubKey\",\"key\":\"AgPQELGzKmlAaSb01OKbmuL1f17MHJshkh9s9xAWxMa3\"}]}"

	walletMultiPub, err := common.DecodePublicKey(suite.Client.RpcClient(), mulSignAccPubKey)
	if err != nil {
		panic(err)
	}

	masterHexAddr := ethCommon.BytesToAddress(walletMultiPub.Address().Bytes())
	fmt.Println(masterHexAddr)

	bankClient := suite.Client.NewBankClient()
	accNum, accSeq, err := bankClient.AccountRetriever(masterHexAddr.String())
	if err != nil {
		panic(err)
	}

	fmt.Println(accNum)
	fmt.Println(accSeq)
}
