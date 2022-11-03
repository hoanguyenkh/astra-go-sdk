package channel

import (
	"context"
	"fmt"
	"github.com/AstraProtocol/astra-go-sdk/account"
	"github.com/AstraProtocol/astra-go-sdk/common"
	channelTypes "github.com/AstraProtocol/channel/x/channel/types"
	"github.com/cosmos/cosmos-sdk/client"
)

type Channel struct {
	rpcClient client.Context
}

func NewChannel(rpcClient client.Context) *Channel {
	return &Channel{rpcClient}
}

func (channel *Channel) OpenChannel(request OpenChannelRequest) error {
	msg := channelTypes.NewMsgOpenChannel(request.Creator, request.PartA, request.PartB,
		request.CoinA, request.CoinB, request.MultisigAddr, request.Sequence)
	err := msg.ValidateBasic()
	if err != nil {
		return err
	}
	tmpPrik := account.PrivateKeySerialized{}
	newTx := common.NewTx(channel.rpcClient,
		account.NewPrivateKeySerialized("", tmpPrik.PrivateKey()), 0, "")
	txBuilder, err := newTx.BuildUnsignedTx(msg)

	json, err := channel.rpcClient.TxConfig.TxJSONEncoder()(txBuilder.GetTx())
	if err != nil {
		return err
	}

	return channel.rpcClient.PrintString(fmt.Sprintf("%s\n", json))
}

func (channel *Channel) ListChannel() (*channelTypes.QueryAllChannelResponse, error) {
	channelClient := channelTypes.NewQueryClient(channel.rpcClient)
	return channelClient.ChannelAll(context.Background(), &channelTypes.QueryAllChannelRequest{})
}
