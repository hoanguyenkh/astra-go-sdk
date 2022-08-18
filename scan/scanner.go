package scan

import (
	"context"
	"fmt"
	"github.com/AstraProtocol/astra-go-sdk/common"
	"github.com/cosmos/cosmos-sdk/client"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/pkg/errors"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
	"github.com/tendermint/tendermint/types"
	evmtypes "github.com/tharsis/ethermint/x/evm/types"
	"math/big"
	"time"
)

type Scanner struct {
	rpcClient client.Context
}

func NewScanner(rpcClient client.Context) *Scanner {
	return &Scanner{rpcClient: rpcClient}
}

func (b *Scanner) ScanViaWebsocket() {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	subscription := b.rpcClient.Client

	err := subscription.Start()
	if err != nil {
		panic(err)
	}
	defer subscription.Stop()

	queryStr := fmt.Sprintf("tm.event='NewBlock' AND block.height='1038312'")
	fmt.Println(queryStr)
	blockHeadersSub, err := subscription.Subscribe(
		ctx,
		"test-client",
		queryStr,
	)

	if err != nil {
		panic(err)
	}

	go func() {
		for e := range blockHeadersSub {
			eventDataHeader := e.Data.(types.EventDataNewBlock)
			height := eventDataHeader.Block.Height
			data := eventDataHeader.Block.Data

			fmt.Println(height)
			for _, rawData := range data.Txs {
				tx, err := b.rpcClient.TxConfig.TxDecoder()(rawData)
				if err != nil {
					panic(err)
				}

				_, err = b.rpcClient.TxConfig.TxJSONEncoder()(tx)
				if err != nil {
					panic(err)
				}

				fmt.Printf("%X\n", rawData.Hash())
			}
		}
	}()

	select {}
}

func (b *Scanner) ScanByBlockHeight(height int64) ([]*Txs, error) {
	startTime := time.Now()

	lisTx := make([]*Txs, 0)

	blockInfo, _, err := b.getBlock(&height)
	if err != nil {
		return nil, errors.Wrap(err, "getBlock")
	}

	blkHeight := blockInfo.Block.Height
	blockTime := blockInfo.Block.Time
	layout := "2006-01-02T15:04:05.000Z"

	fmt.Printf("scan block = %v total = %v\n", height, len(blockInfo.Block.Txs))
	for _, rawData := range blockInfo.Block.Txs {
		tx, err := b.rpcClient.TxConfig.TxDecoder()(rawData)
		if err != nil {
			return nil, errors.Wrap(err, "TxDecoder")
		}

		txBytes, err := b.rpcClient.TxConfig.TxJSONEncoder()(tx)
		if err != nil {
			return nil, errors.Wrap(err, "TxJSONEncoder")
		}

		/*txResult := blockResults.TxsResults[i]
		if !txResult.IsOK() {
			fmt.Printf("Tx = %X at block = %v is failed\n", rawData.Hash(), height)
			continue
		}*/

		ts := blockTime.Format(layout)
		txs := &Txs{
			//Code:        txResult.Code,
			//IsOk:        txResult.IsOK(),
			Time:        ts,
			BlockHeight: blkHeight,
			TxHash:      fmt.Sprintf("%X", rawData.Hash()),
			RawData:     string(txBytes),
		}

		msg := tx.GetMsgs()[0]

		msgEth, ok := msg.(*evmtypes.MsgEthereumTx)
		if ok {
			err := b.getEthMsg(txs, msgEth)
			if err != nil {
				return nil, errors.Wrap(err, "getEthMsg")
			}
		}

		msgBankSend, ok := msg.(*banktypes.MsgSend)
		if ok {
			err := b.getBankSendMsg(txs, msgBankSend)
			if err != nil {
				return nil, errors.Wrap(err, "getBankSendMsg")
			}
		}

		lisTx = append(lisTx, txs)
	}

	fmt.Printf("end scan block = %v in = %v\n", height, time.Since(startTime).String())
	return lisTx, nil
}

func (b *Scanner) getEthMsg(txs *Txs, msgEth *evmtypes.MsgEthereumTx) error {
	data, err := evmtypes.UnpackTxData(msgEth.Data)
	if err != nil {
		return errors.Wrap(err, "UnpackTxData")
	}

	var txDataType string
	switch data.(type) {
	case *evmtypes.AccessListTx:
		txDataType = "access_list_tx"
	case *evmtypes.LegacyTx:
		txDataType = "legacy_tx"
	case *evmtypes.DynamicFeeTx:
		txDataType = "dynamic_fee_tx"
	default:
		return errors.Wrap(err, "UnpackTxData")
	}

	txType := msgEth.Type()

	sig := msgEth.GetSigners()
	from := sig[0].String()

	amountStr := "0"
	if data.GetValue() != nil {
		amountStr = data.GetValue().String()
	}

	txs.Type = txType
	txs.TxDataType = txDataType
	txs.EthTxHash = msgEth.Hash

	ethSender, err := common.CosmosAddressToEthAddress(from)
	if err != nil {
		return errors.Wrap(err, "CosmosAddressToEthAddress")
	}

	txs.Sender = from
	txs.EthSender = ethSender

	to := ""
	receiver := ""
	if data.GetTo() != nil {
		to = data.GetTo().String()

		receiver, err = common.EthAddressToCosmosAddress(to)
		if err != nil {
			return errors.Wrap(err, "EthAddressToCosmosAddress")
		}
	}

	txs.Receiver = receiver
	txs.EthReceiver = to

	amount, ok := big.NewInt(0).SetString(amountStr, 10)
	if !ok {
		return errors.New("Parser amount invalid")
	}

	txs.AmountDecimal = big.NewInt(0).Div(amount, big.NewInt(1e18)).String()

	txs.Amount = amountStr
	txs.TokenSymbol = ""

	return nil
}

func (b *Scanner) getBankSendMsg(txs *Txs, msgSend *banktypes.MsgSend) error {
	//txs.EthTxHash = msgEth.Hash

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

func (b *Scanner) getBlockResults(height *int64) (*ctypes.ResultBlockResults, error) {
	node, err := b.rpcClient.GetNode()
	if err != nil {
		return nil, err
	}

	res, err := node.BlockResults(context.Background(), height)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (b *Scanner) getBlock(height *int64) (*ctypes.ResultBlock, *ctypes.ResultBlockResults, error) {
	// get the node
	node, err := b.rpcClient.GetNode()
	if err != nil {
		return nil, nil, err
	}

	res, err := node.Block(context.Background(), height)
	if err != nil {
		return nil, nil, err
	}

	/*res1, err := node.BlockResults(context.Background(), height)
	if err != nil {
		return nil, nil, err
	}*/

	return res, nil, nil
}

// get the current blockchain height
func (b *Scanner) GetChainHeight() (int64, error) {
	node, err := b.rpcClient.GetNode()
	if err != nil {
		return -1, err
	}

	status, err := node.Status(context.Background())
	if err != nil {
		return -1, err
	}

	height := status.SyncInfo.LatestBlockHeight
	return height, nil
}
