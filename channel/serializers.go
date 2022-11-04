package channel

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type OpenChannelRequest struct {
	Creator      string    `protobuf:"bytes,1,opt,name=creator,proto3" json:"creator,omitempty"`
	PartA        string    `protobuf:"bytes,2,opt,name=partA,proto3" json:"partA,omitempty"`
	PartB        string    `protobuf:"bytes,3,opt,name=partB,proto3" json:"partB,omitempty"`
	CoinA        *sdk.Coin `protobuf:"bytes,4,opt,name=coinA,proto3" json:"coinA,omitempty"`
	CoinB        *sdk.Coin `protobuf:"bytes,5,opt,name=coinB,proto3" json:"coinB,omitempty"`
	MultisigAddr string    `protobuf:"bytes,6,opt,name=multisigAddr,proto3" json:"multisigAddr,omitempty"`
	Sequence     string    `protobuf:"bytes,7,opt,name=sequence,proto3" json:"sequence,omitempty"`
}

type CloseChannelRequest struct {
	Creator   string    `protobuf:"bytes,1,opt,name=creator,proto3" json:"creator,omitempty"`
	From      string    `protobuf:"bytes,2,opt,name=from,proto3" json:"from,omitempty"`
	ToA       string    `protobuf:"bytes,3,opt,name=toA,proto3" json:"toA,omitempty"`
	CoinA     *sdk.Coin `protobuf:"bytes,4,opt,name=coinA,proto3" json:"coinA,omitempty"`
	ToB       string    `protobuf:"bytes,5,opt,name=toB,proto3" json:"toB,omitempty"`
	CoinB     *sdk.Coin `protobuf:"bytes,6,opt,name=coinB,proto3" json:"coinB,omitempty"`
	Channelid string    `protobuf:"bytes,7,opt,name=channelid,proto3" json:"channelid,omitempty"`
}
type CommitmentRequest struct {
	Creator       string    `protobuf:"bytes,1,opt,name=creator,proto3" json:"creator,omitempty"`
	From          string    `protobuf:"bytes,2,opt,name=from,proto3" json:"from,omitempty"`
	Cointocreator *sdk.Coin `protobuf:"bytes,3,opt,name=cointocreator,proto3" json:"cointocreator,omitempty"`
	ToTimelock    string    `protobuf:"bytes,4,opt,name=toTimelock,proto3" json:"toTimelock,omitempty"`
	Blockheight   uint64    `protobuf:"varint,5,opt,name=blockheight,proto3" json:"blockheight,omitempty"`
	ToHashlock    string    `protobuf:"bytes,6,opt,name=toHashlock,proto3" json:"toHashlock,omitempty"`
	Hashcode      string    `protobuf:"bytes,7,opt,name=hashcode,proto3" json:"hashcode,omitempty"`
	Coinhtlc      *sdk.Coin `protobuf:"bytes,8,opt,name=coinhtlc,proto3" json:"coinhtlc,omitempty"`
	Channelid     string    `protobuf:"bytes,9,opt,name=channelid,proto3" json:"channelid,omitempty"`
}

type WithdrawTimelockRequest struct {
	Creator string `protobuf:"bytes,1,opt,name=creator,proto3" json:"creator,omitempty"`
	To      string `protobuf:"bytes,2,opt,name=to,proto3" json:"to,omitempty"`
	Index   string `protobuf:"bytes,3,opt,name=index,proto3" json:"index,omitempty"`
}

type WithdrawHashlockRequest struct {
	Creator string `protobuf:"bytes,1,opt,name=creator,proto3" json:"creator,omitempty"`
	To      string `protobuf:"bytes,2,opt,name=to,proto3" json:"to,omitempty"`
	Index   string `protobuf:"bytes,3,opt,name=index,proto3" json:"index,omitempty"`
	Secret  string `protobuf:"bytes,4,opt,name=secret,proto3" json:"secret,omitempty"`
}

type FundRequest struct {
	Creator   string    `protobuf:"bytes,1,opt,name=creator,proto3" json:"creator,omitempty"`
	From      string    `protobuf:"bytes,2,opt,name=from,proto3" json:"from,omitempty"`
	Channelid string    `protobuf:"bytes,3,opt,name=channelid,proto3" json:"channelid,omitempty"`
	Coinlock  *sdk.Coin `protobuf:"bytes,4,opt,name=coinlock,proto3" json:"coinlock,omitempty"`
	Hashcode  string    `protobuf:"bytes,5,opt,name=hashcode,proto3" json:"hashcode,omitempty"`
	Timelock  string    `protobuf:"bytes,6,opt,name=timelock,proto3" json:"timelock,omitempty"`
	Multisig  string    `protobuf:"bytes,7,opt,name=multisig,proto3" json:"multisig,omitempty"`
}

type AcceptFundRequest struct {
	Creator   string    `protobuf:"bytes,1,opt,name=creator,proto3" json:"creator,omitempty"`
	From      string    `protobuf:"bytes,2,opt,name=from,proto3" json:"from,omitempty"`
	Channelid string    `protobuf:"bytes,3,opt,name=channelid,proto3" json:"channelid,omitempty"`
	Coin      *sdk.Coin `protobuf:"bytes,4,opt,name=coin,proto3" json:"coin,omitempty"`
	Hashcode  string    `protobuf:"bytes,5,opt,name=hashcode,proto3" json:"hashcode,omitempty"`
	Timelock  string    `protobuf:"bytes,6,opt,name=timelock,proto3" json:"timelock,omitempty"`
	Multisig  string    `protobuf:"bytes,7,opt,name=multisig,proto3" json:"multisig,omitempty"`
}
