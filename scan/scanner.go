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

	queryStr := fmt.Sprintf("tm.event='NewBlock' AND block.height='450816'")
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
	/*chainHeight, err := b.GetChainHeight()
	if err != nil {
		return nil, errors.Wrap(err, "GetChainHeight")
	}

	if height > chainHeight {
		return nil, errors.New(fmt.Sprintf("block request large than current block %v > %v", height, chainHeight))
	}*/

	lisTx := make([]*Txs, 0)

	output, err := b.getBlock(&height)
	if err != nil {
		return nil, errors.Wrap(err, "getBlock")
	}

	blkHeight := output.Block.Height
	blockTime := output.Block.Time
	layout := "2006-01-02T15:04:05.000Z"

	for _, rawData := range output.Block.Txs {
		tx, err := b.rpcClient.TxConfig.TxDecoder()(rawData)
		if err != nil {
			return nil, errors.Wrap(err, "TxDecoder")
		}

		txBytes, err := b.rpcClient.TxConfig.TxJSONEncoder()(tx)
		if err != nil {
			return nil, errors.Wrap(err, "TxJSONEncoder")
		}

		ts := blockTime.Format(layout)

		txs := &Txs{
			//Code        uint32
			Time:        ts,
			BlockHeight: blkHeight,
			TxHash:      fmt.Sprintf("%X", rawData.Hash()),
			RawData:     string(txBytes),
		}

		h := fmt.Sprintf("%X", rawData.Hash())

		if h == "6DAA4CAA168236B738E221241FECFEFF11422B72E0CE7AD4CDDA2980896E1BCF" {
			fmt.Println(fmt.Sprintf("%X", rawData.Hash()))

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

	return lisTx, nil
}

func (b *Scanner) getEthMsg(txs *Txs, msgEth *evmtypes.MsgEthereumTx) error {
	data, err := evmtypes.UnpackTxData(msgEth.Data)
	if err != nil {
		return errors.Wrap(err, "UnpackTxData")
	}

	var amountStr string
	var to string

	switch data.(type) {
	case *evmtypes.AccessListTx:
		//nothing
	case *evmtypes.LegacyTx:
		var legacyTx *evmtypes.LegacyTx
		legacyTx = data.(*evmtypes.LegacyTx)
		amountStr = legacyTx.Amount.String()
		to = legacyTx.To
	case *evmtypes.DynamicFeeTx:
		amountStr = data.GetValue().String()
		to = data.GetTo().String()
	default:
		return errors.Wrap(err, "UnpackTxData")
	}

	sig := msgEth.GetSigners()
	from := sig[0].String()

	txs.Type = msgEth.Type()
	txs.EthTxHash = msgEth.Hash

	ethSender, err := common.CosmosAddressToEthAddress(from)
	if err != nil {
		return errors.Wrap(err, "CosmosAddressToEthAddress")
	}

	txs.Sender = from
	txs.EthSender = ethSender

	receiver, err := common.EthAddressToCosmosAddress(to)
	if err != nil {
		return errors.Wrap(err, "EthAddressToCosmosAddress")
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

func (b *Scanner) getBankSendMsg(txs *Txs, msgEth *banktypes.MsgSend) error {
	//txs.EthTxHash = msgEth.Hash
	txs.Type = msgEth.Type()
	ethSender, err := common.CosmosAddressToEthAddress(msgEth.FromAddress)
	if err != nil {
		return errors.Wrap(err, "CosmosAddressToEthAddress")
	}

	txs.Sender = msgEth.FromAddress
	txs.EthSender = ethSender

	receiver, err := common.CosmosAddressToEthAddress(msgEth.ToAddress)
	if err != nil {
		return errors.Wrap(err, "CosmosAddressToEthAddress")
	}

	txs.Receiver = msgEth.ToAddress
	txs.EthReceiver = receiver

	txs.Amount = msgEth.Amount[0].Amount.String()

	amount, ok := big.NewInt(0).SetString(msgEth.Amount[0].Amount.String(), 10)
	if !ok {
		return errors.New("Parser amount invalid")
	}

	txs.AmountDecimal = big.NewInt(0).Div(amount, big.NewInt(1e18)).String()

	txs.TokenSymbol = msgEth.Amount[0].Denom

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

func (b *Scanner) getBlock(height *int64) (*ctypes.ResultBlock, error) {
	// get the node
	node, err := b.rpcClient.GetNode()
	if err != nil {
		return nil, err
	}

	res, err := node.Block(context.Background(), height)
	if err != nil {
		return nil, err
	}

	return res, nil
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
