package common

import (
	"crypto/aes"
	"fmt"
	"testing"
)

func TestAce256(t *testing.T) {
	key := "A3xNe7sEB9HixkmBhVrYaB0NhtHpHgAWeTnLZpTSxCI="
	iv := "rNIIoAcpOUh/aZnrnRikRw=="

	keyByte, _ := decodeBase64(key)
	ivByte, _ := decodeBase64(iv)

	key = string(keyByte)
	iv = string(ivByte)

	plaintext := "1214c33e0ea1815464124ca3566aa406cc59f59d272b183ff33bd4cbea7d8dba"

	fmt.Println("Data to encode: ", plaintext)

	cipherText, _ := CBCEncrypt(plaintext, key, iv, aes.BlockSize)
	fmt.Println("Encode Result:\t", cipherText)
	rs, _ := CBCDecrypt(cipherText, key, iv)
	fmt.Println("Decode Result:\t", rs)
}
