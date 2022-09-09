package bank

import (
	"context"
	"fmt"
	"github.com/AstraProtocol/astra-go-sdk/account"
	"github.com/AstraProtocol/astra-go-sdk/common"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/types"
	authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	bankTypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	emvTypes "github.com/evmos/ethermint/x/evm/types"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"math/big"
	"strings"
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
	if b.coinType == 60 {
		queryClient := emvTypes.NewQueryClient(b.rpcClient)
		cosmosAccount, err := queryClient.CosmosAccount(context.Background(), &emvTypes.QueryCosmosAccountRequest{Address: addr})
		if err != nil {
			return 0, 0, errors.Wrap(err, "CosmosAccount")
		}

		accNum := cosmosAccount.AccountNumber
		accSeq := cosmosAccount.Sequence

		return accNum, accSeq, nil

	}

	addrAcc, err := types.AccAddressFromBech32(addr)
	if err != nil {
		return 0, 0, errors.Wrap(err, "AccAddressFromBech32")
	}

	accNum, accSeq, err := b.rpcClient.AccountRetriever.GetAccountNumberSequence(b.rpcClient, addrAcc)
	if err != nil {
		return 0, 0, errors.Wrap(err, "GetAccountNumberSequence")
	}

	return accNum, accSeq, nil
}

func (b *Bank) CheckTx(txHash string) (*types.TxResponse, error) {
	output, err := authtx.QueryTx(b.rpcClient, txHash)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return nil, fmt.Errorf("no transaction found with hash %s. err = %v", txHash, err.Error())
		}

		return nil, err
	}

	if output.Empty() {
		return nil, fmt.Errorf("no transaction found with hash %s", txHash)
	}

	return output, nil
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

	newTx := common.NewTx(b.rpcClient, acc, param.GasLimit, param.GasPrice)

	txBuilder, err := newTx.BuildUnsignedTx(msg)
	if err != nil {
		return nil, errors.Wrap(err, "BuildUnsignedTx")
	}

	err = newTx.SignTx(txBuilder)
	if err != nil {
		return nil, errors.Wrap(err, "SignTx")
	}

	return txBuilder, nil
}

func (b *Bank) TransferRawDataWithPrivateKey(param *TransferRequest) (client.TxBuilder, error) {
	auth := account.NewAccount(b.coinType)
	acc, err := auth.ImportPrivateKey(param.PrivateKey)
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

	newTx := common.NewTx(b.rpcClient, acc, param.GasLimit, param.GasPrice)

	txBuilder, err := newTx.BuildUnsignedTx(msg)
	if err != nil {
		return nil, errors.Wrap(err, "BuildUnsignedTx")
	}

	err = newTx.SignTx(txBuilder)
	if err != nil {
		return nil, errors.Wrap(err, "SignTx")
	}

	return txBuilder, nil
}

func (b *Bank) SignTxWithSignerAddress(param *SignTxWithSignerAddressRequest) (client.TxBuilder, error) {
	auth := account.NewAccount(b.coinType)
	acc, err := auth.ImportAccount(param.SignerPrivateKey)
	if err != nil {
		return nil, errors.Wrap(err, "ImportAccount")
	}

	from := types.AccAddress(param.SignerPublicKey.Address())
	receiver, err := types.AccAddressFromBech32(param.Receiver)
	if err != nil {
		return nil, errors.Wrap(err, "AccAddressFromBech32")
	}

	amount := types.NewCoin(b.tokenSymbol, types.NewIntFromBigInt(param.Amount))
	msg := bankTypes.NewMsgSend(
		from,
		receiver,
		types.NewCoins(amount),
	)

	newTx := common.NewTxMulSign(b.rpcClient, acc, param.GasLimit, param.GasPrice, param.SequeNum, param.AccNum)

	txBuilder, err := newTx.BuildUnsignedTx(msg)
	if err != nil {
		return nil, errors.Wrap(err, "BuildUnsignedTx")
	}

	err = newTx.SignTxWithSignerAddress(txBuilder, param.SignerPublicKey)
	if err != nil {
		return nil, errors.Wrap(err, "SignTxWithSignerAddress")
	}

	return txBuilder, nil
}

func (b *Bank) TransferMultiSignRawData(param *TransferMultiSignRequest) (client.TxBuilder, error) {
	mulSignAccPublicKey := param.MulSignAccPublicKey

	from := types.AccAddress(mulSignAccPublicKey.Address())
	receiver, err := types.AccAddressFromBech32(param.Receiver)
	if err != nil {
		return nil, errors.Wrap(err, "AccAddressFromBech32")
	}

	amount := types.NewCoin(b.tokenSymbol, types.NewIntFromBigInt(param.Amount))
	msg := bankTypes.NewMsgSend(
		from,
		receiver,
		types.NewCoins(amount),
	)

	newTx := common.NewTxMulSign(b.rpcClient,
		nil,
		param.GasLimit,
		param.GasPrice,
		param.SequeNum,
		param.AccNum)

	txBuilder, err := newTx.BuildUnsignedTx(msg)
	if err != nil {
		return nil, errors.Wrap(err, "BuildUnsignedTx")
	}

	err = newTx.CreateTxMulSign(txBuilder, mulSignAccPublicKey, b.coinType, param.Sigs)
	if err != nil {
		return nil, errors.Wrap(err, "CreateTxMulSign")
	}

	return txBuilder, nil
}

func (b *Bank) ParserEthMsg(txs *Txs, msgEth *emvTypes.MsgEthereumTx) error {
	ethTx := msgEth.AsTransaction()

	var txDataType string
	switch ethTx.Type() {
	case ethTypes.LegacyTxType:
		txDataType = "legacy_tx_type"
	case ethTypes.AccessListTxType:
		txDataType = "access_list_tx_type"
	case ethTypes.DynamicFeeTxType:
		txDataType = "dynamic_fee_tx_type"
	default:
		return errors.New("type invalid")
	}

	txs.Type = msgEth.Type()
	txs.TxDataType = txDataType
	txs.EthTxHash = ethTx.Hash().String()

	from, err := msgEth.GetSender(ethTx.ChainId())
	if err != nil {
		return errors.Wrap(err, "GetSender")
	}

	sender, err := common.EthAddressToCosmosAddress(from.String())
	if err != nil {
		return errors.Wrap(err, "CosmosAddressToEthAddress")
	}

	receiver, err := common.EthAddressToCosmosAddress(ethTx.To().String())
	if err != nil {
		return errors.Wrap(err, "EthAddressToCosmosAddress")
	}

	txs.Sender = sender
	txs.EthSender = from.String()

	txs.Receiver = receiver
	txs.EthReceiver = ethTx.To().String()

	dataExternal := ethTx.Data()
	if dataExternal != nil {
		txs.IsUnNativeCoin = true
	}

	amount, ok := big.NewFloat(0).SetString(ethTx.Value().String())
	if !ok {
		return errors.New("Parser amount invalid")
	}

	txs.AmountDecimal = big.NewFloat(0).Quo(amount, big.NewFloat(1e18)).String()
	txs.Amount = ethTx.Value().String()

	txs.TokenSymbol = ""

	return nil
}

func (b *Bank) ParserCosmosMsg(txs *Txs, msgSend *bankTypes.MsgSend) error {
	txs.TxDataType = "cosmos"
	txs.Type = msgSend.Type()
	ethSender, err := common.CosmosAddressToEthAddress(msgSend.FromAddress)
	if err != nil {
		return errors.Wrap(err, "CosmosAddressToEthAddress")
	}

	txs.Sender = msgSend.FromAddress
	txs.EthSender = ethSender

	receiver, err := common.CosmosAddressToEthAddress(msgSend.ToAddress)
	if err != nil {
		return errors.Wrap(err, "CosmosAddressToEthAddress")
	}

	txs.Receiver = msgSend.ToAddress
	txs.EthReceiver = receiver

	txs.Amount = msgSend.Amount[0].Amount.String()

	amount, ok := big.NewInt(0).SetString(msgSend.Amount[0].Amount.String(), 10)
	if !ok {
		return errors.New("Parser amount invalid")
	}

	txs.AmountDecimal = big.NewInt(0).Div(amount, big.NewInt(1e18)).String()

	txs.TokenSymbol = msgSend.Amount[0].Denom

	return nil
}

func (b *Bank) TxDetail(txHash string) (*Txs, error) {
	rs, err := b.CheckTx(txHash)
	if err != nil {
		return nil, err
	}

	tx, err := b.rpcClient.TxConfig.TxDecoder()(rs.Tx.GetValue())
	if err != nil {
		return nil, err
	}

	txBytes, err := b.rpcClient.TxConfig.TxJSONEncoder()(tx)
	if err != nil {
		return nil, err
	}

	msg := tx.GetMsgs()[0]
	txs := &Txs{
		Code:        rs.Code,
		IsOk:        common.BlockedStatus(rs.Code),
		Time:        rs.Timestamp,
		BlockHeight: rs.Height,
		TxHash:      rs.TxHash,
		RawData:     string(txBytes),
	}

	msgEth, ok := msg.(*emvTypes.MsgEthereumTx)
	if ok {
		err := b.ParserEthMsg(txs, msgEth)
		if err != nil {
			panic(err)
		}
	}

	msgBankSend, ok := msg.(*bankTypes.MsgSend)
	if ok {
		err := b.ParserCosmosMsg(txs, msgBankSend)
		if err != nil {
			panic(err)
		}
	}

	return txs, nil
}
