package bank

import (
	"context"
	"github.com/AstraProtocol/astra-go-sdk/account"
	"github.com/AstraProtocol/astra-go-sdk/common"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/types"
	bankTypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/pkg/errors"
	emvTypes "github.com/tharsis/ethermint/x/evm/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"math/big"
)

type Bank struct {
	rpcClient   client.Context
	tokenSymbol string
	coinType    uint32
}

func NewBank(rpcClient client.Context, tokenSymbol string, coinType uint32) *Bank {
	return &Bank{rpcClient: rpcClient, tokenSymbol: tokenSymbol, coinType: coinType}
}

func (b *Bank) Balance(addr string) (*big.Int, error) {
	var header metadata.MD

	bankClient := bankTypes.NewQueryClient(b.rpcClient)
	bankRes, err := bankClient.Balance(
		context.Background(),
		&bankTypes.QueryBalanceRequest{Address: addr, Denom: b.tokenSymbol},
		grpc.Header(&header),
	)

	if err != nil {
		return nil, errors.Wrap(err, "Balance")
	}

	return bankRes.Balance.Amount.BigInt(), nil
}

func (b *Bank) AccountRetriever(addr string) (uint64, uint64, error) {
	addrAcc, err := types.AccAddressFromBech32(addr)
	if err != nil {
		return 0, 0, errors.Wrap(err, "AccAddressFromBech32")
	}

	if b.coinType == 60 {
		queryClient := emvTypes.NewQueryClient(b.rpcClient)
		cosmosAccount, err := queryClient.CosmosAccount(context.Background(), &emvTypes.QueryCosmosAccountRequest{Address: addrAcc.String()})
		if err != nil {
			return 0, 0, errors.Wrap(err, "CosmosAccount")
		}

		accNum := cosmosAccount.AccountNumber
		accSeq := cosmosAccount.Sequence

		return accNum, accSeq, nil

	}

	accNum, accSeq, err := b.rpcClient.AccountRetriever.GetAccountNumberSequence(b.rpcClient, addrAcc)
	if err != nil {
		return 0, 0, errors.Wrap(err, "GetAccountNumberSequence")
	}

	return accNum, accSeq, nil
}

func (b *Bank) TransferRawData(param *TransferRequest) (client.TxBuilder, error) {
	auth := account.NewAccount(b.coinType)
	acc, err := auth.ImportAccount(param.PrivateKey)
	if err != nil {
		return nil, errors.Wrap(err, "ImportAccount")
	}

	receiver, err := types.AccAddressFromBech32(param.Receiver)
	if err != nil {
		return nil, errors.Wrap(err, "AccAddressFromBech32")
	}

	coin := types.NewCoin(b.tokenSymbol, types.NewIntFromBigInt(param.Amount))
	msg := bankTypes.NewMsgSend(
		acc.AccAddress(),
		receiver,
		types.NewCoins(coin),
	)

	tx := common.NewTx(b.rpcClient, acc, param.GasLimit, param.GasPrice)

	txBuilder, err := tx.BuildUnsignedTx(msg)
	if err != nil {
		return nil, errors.Wrap(err, "BuildUnsignedTx")
	}

	err = tx.SignTx(txBuilder)
	if err != nil {
		return nil, errors.Wrap(err, "SignTx")
	}

	return txBuilder, nil
}

func (b *Bank) TransferMulSignRawData(param *TransferMultiSignRequest) (client.TxBuilder, error) {
	pk, err := common.DecodePublicKey(param.FromPublicKey)
	from, err := types.AccAddressFromBech32(param.From)
	if err != nil {
		return nil, errors.Wrap(err, "AccAddressFromBech32")
	}

	receiver, err := types.AccAddressFromBech32(param.Receiver)
	if err != nil {
		return nil, errors.Wrap(err, "AccAddressFromBech32")
	}

	msg := bankTypes.NewMsgSend(
		from,
		receiver,
		types.NewCoins(types.NewInt64Coin(b.tokenSymbol, param.Amount.Int64())),
	)

	tx := common.NewTx(b.rpcClient, nil, param.GasLimit, param.GasPrice)

	txBuilder, err := tx.BuildUnsignedTx(msg)
	if err != nil {
		return nil, errors.Wrap(err, "BuildUnsignedTx")
	}

	err = tx.MulSignTx(txBuilder, pk, b.coinType, nil)
	if err != nil {
		return nil, errors.Wrap(err, "MulSignTx")
	}

	return txBuilder, nil
}
