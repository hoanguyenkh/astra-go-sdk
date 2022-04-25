package bank

import (
	"math/big"
)

type TransferRequest struct {
	PrivateKey    string
	From          string
	FromPublicKey string
	Receiver      string
	Amount        *big.Int
	GasLimit      uint64
	GasPrice      string
}
