package account

import (
	"fmt"
	"testing"
)

func TestCreateAccount(t *testing.T) {
	n := NewAccount(60)
	key, err := n.CreateAccount()
	if err != nil {
		panic(err)
	}
	fmt.Println("key type 60")
	fmt.Println(key.String())

	n1 := NewAccount(118)
	key1, err := n1.CreateAccount()
	if err != nil {
		panic(err)
	}

	fmt.Println("key type 118")
	fmt.Println(key1.String())
}

func TestCreateMulAccount(t *testing.T) {
	n := NewAccount(60)
	key, addr, pucKey, err := n.CreateMulSignAccount(3, 2)
	if err != nil {
		panic(err)
	}

	fmt.Println("key type 60")

	fmt.Println("addr", addr)
	fmt.Println("pucKey", pucKey)
	fmt.Println("list key")
	for i, serialized := range key {
		fmt.Println("index", i)
		fmt.Println(serialized.String())
	}

}