package common

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestECBEncrypt(t *testing.T) {
	data := "mammal initial effort joke public daring fish puppy risk famous cream occur else busy cable cruel vacant brick used patient choose object teach special"
	result, err := ECBEncrypt([]byte(data), []byte("VbKerprP*=wSql9X"))
	if err != nil {
		panic(err)
	}

	fmt.Println(result)

	result = "B4GLoBV4acCzFJ9jcdxM6tW9uzdi3E1qzkHpy+b2iuJRoIBxYZEfz/qEtHQh1j14p+NP7JsZymThPb7CdHqjzgO2yG1x6DRidu+6nKD4ONJt2CKbct8QtLUm9ut6g1EZOPvhkQX2Qx8gagesMdZM5n8uf78OsSZIS+ar+49P9t2dUEaqnr0elylIa48hPD8AIqHtuX/70fXxiIbDZASaQg=="
	result1, err := ECBDecrypt(result, []byte("VbKerprP*=wSql9X"))

	if err != nil {
		panic(err)
	}

	fmt.Println(result1)
	assert.Equal(t, data, result1)
}

func TestGenerateKey(t *testing.T) {
	token1, _ := GenerateRandomBytes(24)
	token2, _ := GenerateRandomBytes(16)

	key := encodeBase64(token1)
	iv := encodeBase64(token2)

	token3 := GenerateRandomString(16)
	fmt.Println(key)
	fmt.Println(iv)
	fmt.Println(token3)

	token4, _ := GenerateSecretKeyRandomString(16)
	fmt.Println(token4)
}
