package channel

import (
	"context"
	"fmt"
	"github.com/AstraProtocol/astra-go-sdk/account"
	"github.com/AstraProtocol/astra-go-sdk/common"
	channelTypes "github.com/AstraProtocol/channel/x/channel/types"
	"github.com/cosmos/cosmos-sdk/client"
	cryptoTypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/pkg/errors"
)

type Channel struct {
	rpcClient client.Context
}

func NewChannel(rpcClient client.Context) *Channel {
	return &Channel{rpcClient}
}

func (cn *Channel) SignOpenChannel(request OpenChannelRequest,
	account *account.PrivateKeySerialized,
	multiSigPubkey cryptoTypes.PubKey) (string, error) {
	msg := channelTypes.NewMsgOpenChannel(request.Creator, request.PartA, request.PartB,
		request.CoinA, request.CoinB, request.MultisigAddr, request.Sequence)
	err := msg.ValidateBasic()
	if err != nil {
		return "", err
	}
	newTx := common.NewTxMulSign(cn.rpcClient, account, request.GasLimit, request.GasPrice, 0, 2)
	txBuilder, err := newTx.BuildUnsignedTx(msg)

	err = newTx.SignTxWithSignerAddress(txBuilder, multiSigPubkey)
	if err != nil {
		return "", errors.Wrap(err, "SignTx")
	}

	sign, err := common.TxBuilderSignatureJsonEncoder(cn.rpcClient.TxConfig, txBuilder)
	if err != nil {
		return "", err
	}
	fmt.Println("sign-data", string(sign))
	return sign, nil
}

func (cn *Channel) ListChannel() (*channelTypes.QueryAllChannelResponse, error) {
	channelClient := channelTypes.NewQueryClient(cn.rpcClient)
	return channelClient.ChannelAll(context.Background(), &channelTypes.QueryAllChannelRequest{})
}
