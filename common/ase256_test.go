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

	plaintext := "abcdefghijklmnopqrstuvwxyzABCDEF"

	fmt.Println("Data to encode: ", plaintext)

	cipherText := fmt.Sprintf("%v", Ase256Encode(plaintext, key, iv, aes.BlockSize))
	fmt.Println("Encode Result:\t", cipherText)
	fmt.Println("Decode Result:\t", Ase256Decode(cipherText, key, iv))
}
