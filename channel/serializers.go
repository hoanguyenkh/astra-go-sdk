package channel

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type OpenChannelRequest struct {
	Creator      string
	PartA        string
	PartB        string
	CoinA        *sdk.Coin
	CoinB        *sdk.Coin
	MultisigAddr string
	Sequence     string
}
