package common

import (
	"bytes"
	"crypto/aes"
	"encoding/base64"
	"fmt"
)

func encodeBase64(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}

func decodeBase64(s string) ([]byte, error) {
	data, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return nil, err
	}
	return data, nil
}

//ECC mode decryption
func ECBDecrypt(cryptedStr string, key []byte) (string, error) {
	crypted, err := decodeBase64(cryptedStr)
	if err != nil {
		return "", err
	}

	if !validKey(key) {
		return "", fmt.Errorf("the length of the secret key is wrong, the current incoming length is %d", len(key))
	}

	if len(crypted) < 1 {
		return "", fmt.Errorf("source data length cannot be 0")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	if len(crypted)%block.BlockSize() != 0 {
		return "", fmt.Errorf("the source data length must be an integer multiple of %d, the current length is %d", block.BlockSize(), len(crypted))
	}

	var dst []byte
	tmpData := make([]byte, block.BlockSize())

	for index := 0; index < len(crypted); index += block.BlockSize() {
		block.Decrypt(tmpData, crypted[index:index+block.BlockSize()])
		dst = append(dst, tmpData...)
	}

	dst, err = PKCS5UnPadding(dst)
	if err != nil {
		return "", err
	}

	return string(dst[:]), nil
}

//ECC mode encryption
func ECBEncrypt(src, key []byte) (string, error) {
	if !validKey(key) {
		return "", fmt.Errorf("the length of the secret key is wrong, the current incoming length is %d", len(key))
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	if len(src) < 1 {
		return "", fmt.Errorf("source data length cannot be 0")
	}

	src = PKCS5Padding(src, block.BlockSize())
	if len(src)%block.BlockSize() != 0 {
		return "", fmt.Errorf("the source data length must be an integer multiple of %d, the current length is %d", block.BlockSize(), len(src))
	}

	var dst []byte
	tmpData := make([]byte, block.BlockSize())
	for index := 0; index < len(src); index += block.BlockSize() {
		block.Encrypt(tmpData, src[index:index+block.BlockSize()])
		dst = append(dst, tmpData...)
	}

	encodeB64 := encodeBase64(dst)

	return encodeB64, nil
}

//Pkcs5 filling
func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

//Remove pkcs5 filling
func PKCS5UnPadding(origData []byte) ([]byte, error) {
	length := len(origData)
	unpadding := int(origData[length-1])

	if length < unpadding {
		return nil, fmt.Errorf("invalid unpadding length")
	}
	return origData[:(length - unpadding)], nil
}

//Key length verification
func validKey(key []byte) bool {
	k := len(key)
	switch k {
	default:
		return false
	case 16, 24, 32:
		return true
	}
}
