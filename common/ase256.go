package common

import (
	"crypto/aes"
	"crypto/cipher"
)

func Ase256Encode(plaintext string, key string, iv string, blockSize int) (string, error) {
	bKey := []byte(key)
	bIV := []byte(iv)
	bPlaintext := PKCS5Padding([]byte(plaintext), blockSize)
	block, err := aes.NewCipher(bKey)
	if err != nil {
		return "", err
	}
	ciphertext := make([]byte, len(bPlaintext))
	mode := cipher.NewCBCEncrypter(block, bIV)
	mode.CryptBlocks(ciphertext, bPlaintext)
	return encodeBase64(ciphertext), nil
}

func Ase256Decode(cipherText string, encKey string, iv string) (string, error) {
	bKey := []byte(encKey)
	bIV := []byte(iv)
	cipherTextDecoded, err := decodeBase64(cipherText)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(bKey)
	if err != nil {
		return "", err
	}

	mode := cipher.NewCBCDecrypter(block, bIV)
	mode.CryptBlocks([]byte(cipherTextDecoded), []byte(cipherTextDecoded))
	return string(cipherTextDecoded), nil
}
