package channel

import (
	"context"
	"fmt"

	"github.com/AstraProtocol/astra-go-sdk/account"
	"github.com/AstraProtocol/astra-go-sdk/common"
	channelTypes "github.com/AstraProtocol/channel/x/channel/types"
	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Channel struct {
	rpcClient client.Context
}

func NewChannel(rpcClient client.Context) *Channel {
	return &Channel{rpcClient}
}

func (cn *Channel) OpenChannel(request OpenChannelRequest) (string, error) {
	msg := channelTypes.NewMsgOpenChannel(request.Creator, request.PartA, request.PartB,
		request.CoinA, request.CoinB, request.MultisigAddr, request.Sequence)
	err := msg.ValidateBasic()
	if err != nil {
		return "", err
	}
	tmpPrivKey := account.PrivateKeySerialized{}
	newTx := common.NewTx(cn.rpcClient,
		account.NewPrivateKeySerialized("", tmpPrivKey.PrivateKey()), 0, "")
	txBuilder, err := newTx.BuildUnsignedTx(msg)

	json, err := cn.rpcClient.TxConfig.TxJSONEncoder()(txBuilder.GetTx())
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s\n", json), nil
}

func (cn *Channel) ListChannel() (*channelTypes.QueryAllChannelResponse, error) {
	channelClient := channelTypes.NewQueryClient(cn.rpcClient)
	return channelClient.ChannelAll(context.Background(), &channelTypes.QueryAllChannelRequest{})
}

func (cn *Channel) CloseChannel(request CloseChannelRequest) (string, error) {
}

func (cn *Channel) GetMsgJSONEncoder(msg sdk.Msg) (string, error) {
	err := msg.ValidateBasic()
	if err != nil {
		return "", err
	}
	tmpPrivKey := account.PrivateKeySerialized{}
	newTx := common.NewTx(cn.rpcClient,
		account.NewPrivateKeySerialized("", tmpPrivKey.PrivateKey()), 0, "")
	txBuilder, err := newTx.BuildUnsignedTx(msg)

	json, err := cn.rpcClient.TxConfig.TxJSONEncoder()(txBuilder.GetTx())
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s\n", json), nil
}
