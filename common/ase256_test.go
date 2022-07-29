package common

import (
	"crypto/aes"
	"fmt"
	"testing"
)

func TestAce256(t *testing.T) {
	key := "12345678901234567890123456789012"
	iv := "1234567890123456"

	key = "tF+EgPEnMsAGUunOWkNzIV8SzD3fYatlHrCndy1y4QY="
	iv = "o5L+qPKtMjcgnFXyERHuoQ=="

	keyByte, _ := decodeBase64(key)
	ivByte, _ := decodeBase64(iv)

	key = string(keyByte)
	iv = string(ivByte)

	plaintext := "abcdefghijklmnopqrstuvwxyzABCDEF"

	fmt.Println("Data to encode: ", plaintext)

	cipherText := fmt.Sprintf("%v", Ase256Encode(plaintext, key, iv, aes.BlockSize))
	fmt.Println("Encode Result:\t", cipherText)
	fmt.Println("Decode Result:\t", Ase256Decode(cipherText, key, iv))
}
