package common

import (
	"crypto/aes"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAce256(t *testing.T) {
	key, _ := GenerateSecretKeyRandomString(32)
	iv, _ := GenerateSecretKeyRandomString(16)

	fmt.Println(key, iv)

	keyByte, _ := DecodeBase64(key)
	ivByte, _ := DecodeBase64(iv)

	fmt.Println(len(keyByte), len(ivByte))

	plaintext := "1214c33e0ea1815464124ca3566aa406cc59f59d272b183ff33bd4cbea7d8dba"

	fmt.Println("Data to encode: ", plaintext)

	cipherText, _ := CBCEncrypt(plaintext, keyByte, ivByte, aes.BlockSize)
	fmt.Println("Encode Result:\t", cipherText)

	rs, _ := CBCDecrypt(cipherText, keyByte, ivByte)
	fmt.Println("Decode Result:\t", rs)

	assert.Equal(t, plaintext, rs)
}
