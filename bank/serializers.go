package bank

import (
	"math/big"
)

type TransferRequest struct {
	PrivateKey string
	Receiver   string
	Amount     *big.Int
	GasLimit   uint64
	GasPrice   string
}

type TransferMultiSignRequest struct {
	PrivateKey    string
	From          string
	FromPublicKey string
	Receiver      string
	Amount        *big.Int
	GasLimit      uint64
	GasPrice      string
}
