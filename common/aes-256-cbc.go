package common

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"
)

func CBCEncrypt(plaintext string, key []byte, iv []byte, blockSize int) (string, error) {
	bKey := key
	bIV := iv

	if !validKey(bIV) {
		return "", fmt.Errorf("the length of the secret key is wrong, the current incoming length is %d", len(iv))
	}

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

func CBCDecrypt(cipherText string, encKey []byte, iv []byte) (string, error) {
	bKey := encKey
	bIV := iv
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

	dst, err := PKCS5UnPadding(cipherTextDecoded)
	if err != nil {
		return "", err
	}

	return string(dst), nil
}
