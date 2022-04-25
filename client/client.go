package client

import (
	"fmt"
	"github.com/AstraProtocol/astra-go-sdk/account"
	"github.com/AstraProtocol/astra-go-sdk/bank"
	"github.com/AstraProtocol/astra-go-sdk/config"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/simapp/params"
	"github.com/cosmos/cosmos-sdk/types"
	authTypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/tharsis/ethermint/app"
	"github.com/tharsis/ethermint/encoding"
	ethermintTypes "github.com/tharsis/ethermint/types"
)

type Client struct {
	endpoint       string
	chainId        string
	coinType       uint32
	prefixAddress  string
	tokenSymbol    string
	encodingConfig params.EncodingConfig
	sdkConfig      *types.Config
	rpcClient      client.Context
}

func NewClient(cfg *config.Config) (*Client, error) {
	cli := new(Client)
	err := cli.Init(cfg)
	return cli, err
}

func (c *Client) Init(cfg *config.Config) (err error) {
	c.chainId = cfg.ChainId
	c.endpoint = cfg.Endpoint
	c.coinType = cfg.CoinType
	c.prefixAddress = cfg.PrefixAddress
	c.tokenSymbol = cfg.TokenSymbol

	c.encodingConfig = encoding.MakeConfig(app.ModuleBasics)

	sdkConfig := types.GetConfig()
	sdkConfig.SetPurpose(44)

	switch cfg.CoinType {
	case 60:
		sdkConfig.SetCoinType(ethermintTypes.Bip44CoinType)
	case 118:
		sdkConfig.SetCoinType(types.CoinType)
	default:
		panic("Coin type invalid!")
	}

	bech32PrefixAccAddr := fmt.Sprintf("%v", c.prefixAddress)
	bech32PrefixAccPub := fmt.Sprintf("%vpub", c.prefixAddress)
	bech32PrefixValAddr := fmt.Sprintf("%vvaloper", c.prefixAddress)
	bech32PrefixValPub := fmt.Sprintf("%vvaloperpub", c.prefixAddress)
	bech32PrefixConsAddr := fmt.Sprintf("%vvalcons", c.prefixAddress)
	bech32PrefixConsPub := fmt.Sprintf("%vvalconspub", c.prefixAddress)

	sdkConfig.SetBech32PrefixForAccount(bech32PrefixAccAddr, bech32PrefixAccPub)
	sdkConfig.SetBech32PrefixForValidator(bech32PrefixValAddr, bech32PrefixValPub)
	sdkConfig.SetBech32PrefixForConsensusNode(bech32PrefixConsAddr, bech32PrefixConsPub)

	c.sdkConfig = sdkConfig

	rpcHttp, err := client.NewClientFromNode(c.endpoint)
	if err != nil {
		panic(err)
	}

	ar := authTypes.AccountRetriever{}

	rpcClient := client.Context{}
	rpcClient = rpcClient.
		//WithNodeURI(c.endpoint).
		WithTxConfig(c.encodingConfig.TxConfig).
		WithAccountRetriever(ar).
		WithChainID(c.chainId).
		WithClient(rpcHttp).
		WithInterfaceRegistry(c.encodingConfig.InterfaceRegistry) //codec.NewProtoCodec

	c.rpcClient = rpcClient

	return nil
}

func (c *Client) NewAccountClient() *account.Account {
	return account.NewAccount(c.coinType)
}

func (c *Client) NewBankClient() *bank.Bank {
	return bank.NewBank(c.rpcClient, c.tokenSymbol, c.coinType)
}
