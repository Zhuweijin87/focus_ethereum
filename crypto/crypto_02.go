package main

import (
	"fmt"

	"github.com/ethereum/go-ethereum/crypto" 
)

func main() {
	s := []byte("123456789abcd")

	hash256 := crypto.Keccak256(s)  // 32位
	fmt.Println(len(hash256), hash256)

	hash512 := crypto.Keccak512(s) // 64位
	fmt.Println(len(hash512), hash512)

	hash := crypto.Keccak256Hash(s)
	fmt.Println(hash)
}