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

func (cn *Channel) ListChannel() (*channelTypes.QueryAllChannelResponse, error) {
	channelClient := channelTypes.NewQueryClient(cn.rpcClient)
	return channelClient.ChannelAll(context.Background(), &channelTypes.QueryAllChannelRequest{})
}

func (cn *Channel) OpenChannel(req OpenChannelRequest) (string, error) {
	msg := channelTypes.NewMsgOpenChannel(req.Creator, req.PartA, req.PartB,
		req.CoinA, req.CoinB, req.MultisigAddr, req.Sequence)
	return cn.GetMsgJSONEncoder(msg)
}

func (cn *Channel) CloseChannel(req CloseChannelRequest) (string, error) {
	msg := channelTypes.NewMsgCloseChannel(req.Creator, req.From, req.ToA,
		req.CoinA, req.ToB, req.CoinB, req.Channelid)
	return cn.GetMsgJSONEncoder(msg)
}

func (cn *Channel) Commitment(req CommitmentRequest) (string, error) {
	msg := channelTypes.NewMsgCommitment(req.Creator, req.From, req.Cointocreator, req.ToTimelock,
		req.Blockheight, req.ToHashlock, req.Hashcode, req.Coinhtlc, req.Channelid)
	return cn.GetMsgJSONEncoder(msg)
}

func (cn *Channel) WithdrawTimelock(req WithdrawTimelockRequest) (string, error) {
	msg := channelTypes.NewMsgWithdrawTimelock(req.Creator, req.To, req.Index)
	return cn.GetMsgJSONEncoder(msg)
}

func (cn *Channel) WithdrawHashlock(req WithdrawHashlockRequest) (string, error) {
	msg := channelTypes.NewMsgWithdrawHashlock(req.Creator, req.To, req.Index, req.Secret)
	return cn.GetMsgJSONEncoder(msg)
}

func (cn *Channel) Fund(req FundRequest) (string, error) {
	msg := channelTypes.NewMsgFund(req.Creator, req.From, req.Channelid, req.Coinlock, req.Hashcode, req.Timelock, req.Multisig)
	return cn.GetMsgJSONEncoder(msg)
}

func (cn *Channel) AcceptFund(req AcceptFundRequest) (string, error) {
	msg := channelTypes.NewMsgAcceptfund(req.Creator, req.From, req.Channelid, req.Coin, req.Hashcode, req.Timelock, req.Multisig)
	return cn.GetMsgJSONEncoder(msg)
}

func (cn *Channel) SenderCommit(req SenderCommitRequest) (string, error) {
	msg := channelTypes.NewMsgSendercommit(req.Creator, req.From, req.Channelid, req.Cointosender, req.Cointohtlc, req.Hashcodehtlc, req.Timelockhtlc, req.Cointransfer, req.Hashcodedest, req.Timelockreceiver, req.Timelocksender, req.Multisig)
	return cn.GetMsgJSONEncoder(msg)
}
func (cn *Channel) SenderWithdrawTimelock(req SenderWithdrawTimelockRequest) (string, error) {
	msg := channelTypes.NewMsgSenderwithdrawtimelock(req.Creator, req.Transferindex, req.To)
	return cn.GetMsgJSONEncoder(msg)
}
func (cn *Channel) SenderWithdrawHashlock(req SenderWithdrawHashlockRequest) (string, error) {
	msg := channelTypes.NewMsgSenderwithdrawhashlock(req.Creator, req.Transferindex, req.To, req.Secret)
	return cn.GetMsgJSONEncoder(msg)
}

func (cn *Channel) ReceiverCommit(req ReceiverCommitRequest) (string, error) {
	msg := channelTypes.NewMsgReceivercommit(req.Creator, req.From, req.Channelid, req.Cointoreceiver, req.Cointohtlc, req.Hashcodehtlc, req.Timelockhtlc, req.Cointransfer, req.Hashcodedest, req.Timelocksender, req.Multisig)
	return cn.GetMsgJSONEncoder(msg)
}
func (cn *Channel) ReceiverWithdraw(req ReceiverWithdrawRequest) (string, error) {
	msg := channelTypes.NewMsgReceiverwithdraw(req.Creator, req.Transferindex, req.To, req.Secret)
	return cn.GetMsgJSONEncoder(msg)
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
	if err != nil {
		return "", err
	}

	json, err := cn.rpcClient.TxConfig.TxJSONEncoder()(txBuilder.GetTx())
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s\n", json), nil
}
