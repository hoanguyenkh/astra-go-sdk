package common

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestECBEncrypt(t *testing.T) {
	data := "property cactus cannon talent priority silk ice nurse such arctic dove wonder blue stumble chalk engine start know unable tool arctic tone sugar grass"
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
