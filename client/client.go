package client

import (
	"fmt"
	"github.com/AstraProtocol/astra-go-sdk/account"
	"github.com/AstraProtocol/astra-go-sdk/bank"
	"github.com/AstraProtocol/astra-go-sdk/config"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/types"
	authTypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/tharsis/ethermint/app"
	"github.com/tharsis/ethermint/encoding"
	ethermintTypes "github.com/tharsis/ethermint/types"
)

type Client struct {
	coinType      uint32
	prefixAddress string
	tokenSymbol   string
	rpcClient     client.Context
}

func (c *Client) RpcClient() client.Context {
	return c.rpcClient
}

func NewClient(cfg *config.Config) *Client {
	cli := new(Client)
	cli.Init(cfg)
	return cli
}

func (c *Client) Init(cfg *config.Config) {
	c.coinType = cfg.CoinType
	c.prefixAddress = cfg.PrefixAddress
	c.tokenSymbol = cfg.TokenSymbol

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

	ar := authTypes.AccountRetriever{}

	//github.com/cosmos/cosmos-sdk/simapp/app.go
	//github.com/tharsis/ethermint@v0.14.0/app/app.go -> selected
	encodingConfig := encoding.MakeConfig(app.ModuleBasics)

	rpcHttp, err := client.NewClientFromNode(cfg.Endpoint)
	if err != nil {
		panic(err)
	}

	rpcClient := client.Context{}
	rpcClient = rpcClient.
		WithClient(rpcHttp).
		//WithNodeURI(c.endpoint).
		WithTxConfig(encodingConfig.TxConfig).
		WithAccountRetriever(ar).
		WithChainID(cfg.ChainId).
		WithInterfaceRegistry(encodingConfig.InterfaceRegistry). //codec.NewProtoCodec
		WithCodec(encodingConfig.Marshaler).
		WithLegacyAmino(encodingConfig.Amino).
		WithBroadcastMode(flags.BroadcastBlock)

	c.rpcClient = rpcClient
}

func (c *Client) NewAccountClient() *account.Account {
	return account.NewAccount(c.coinType)
}

func (c *Client) NewBankClient() *bank.Bank {
	return bank.NewBank(c.rpcClient, c.tokenSymbol, c.coinType)
}
