package common

import (
	"crypto/elliptic"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestVerifySignature(t *testing.T) {
	publicKey := "A/tdneaDL83fm2BNjYRKacJvRo81iDaYSiybfaDUSM3I"
	publicKey = "A/tdneaDL83fm2BNjYRKacJvRo81iDaYSiybfaDUSM3I"
	signature := "MEQCIFKsQUbx0dzVLSqtfz8CGKGeY0/p9xEwED/76X1EdznaAiBk2XZTkOEiJpBoiKXbw3bQklw+8M3AffqGwBNJlj+xYQ=="
	msg := "ECDSA is cool."

	isValid, err := VerifySignature(publicKey, signature, msg)
	if err != nil {
		panic(err)
	}

	fmt.Println(isValid)
}

func TestGenKeyAndSignData(t *testing.T) {
	privateKey := "c24de61f4915339130e97ec556c41d1d23c5a83a9f45340c45621ba4e5a60a99"
	publickKey := "AxkzPiZ03U9HdKqmKPzQz3URzeZ2pFaPd/4HSuzh017n"

	msg := "ECDSA is cool."
	sign, err := SignatureData(privateKey, msg)

	if err != nil {
		panic(err)
	}

	fmt.Println("sign", sign)

	isValid, err := VerifySignature(publickKey, sign, msg)
	if err != nil {
		panic(err)
	}

	fmt.Println(isValid)
}

func TestGenKey(t *testing.T) {
	privateKey, publickKey := GenPrivateKeySign()
	fmt.Println("privateKey", privateKey)
	fmt.Println("publickKey", publickKey)
}

func TestImportPrivateKey(t *testing.T) {
	privateKey := "1214c33e0ea1815464124ca3566aa406cc59f59d272b183ff33bd4cbea7d8dba"
	key, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		panic(err)
	}

	pubkey := elliptic.MarshalCompressed(crypto.S256(), key.X, key.Y)
	privkey := make([]byte, 32)
	blob := key.D.Bytes()
	copy(privkey[32-len(blob):], blob)

	privkeyStr := hex.EncodeToString(privkey)
	pubkeyStr := base64.StdEncoding.EncodeToString(pubkey)

	assert.Equal(t, privateKey, privkeyStr)

	fmt.Println(privkeyStr)
	fmt.Println(pubkeyStr)
}
