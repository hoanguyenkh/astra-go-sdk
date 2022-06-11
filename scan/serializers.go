package scan

type Txs struct {
	Code          uint32
	Type          string
	TxDataType    string
	Time          string
	BlockHeight   int64
	TxHash        string
	EthTxHash     string
	Sender        string
	EthSender     string
	Receiver      string
	EthReceiver   string
	Amount        string
	AmountDecimal string
	TokenSymbol   string
	RawData       string
}
