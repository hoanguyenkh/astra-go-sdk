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

func TestECBEncrypt(t *testing.T) {
	data := "testing seeding"
	result, err := ECBEncrypt([]byte(data), []byte("WA5Nyx4ODj1Y%9j3"))
	if err != nil {
		panic(err)
	}

	fmt.Println(result)

	result1, err := ECBDecrypt(result, []byte("WA5Nyx4ODj1Y%9j3"))

	if err != nil {
		panic(err)
	}

	fmt.Println(result1)
	assert.Equal(t, data, result1)
}

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
	privateKey, publickKey := GenPrivateKeySign()
	fmt.Println("privateKey", privateKey)
	fmt.Println("publickKey", publickKey)

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
