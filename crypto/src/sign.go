package main

import (
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
)

func main() {
	privKey, _ := crypto.GenerateKey()
	pubKey := privKey.PublicKey

	fmt.Println("pubkey: ", pubKey)

	hash := "23456789abcdef324ad2def4566a23d0"
	sign, err := crypto.Sign([]byte(hash), privKey)
	if err != nil {
		fmt.Println("sign err:", err)
	}

	fmt.Println("sign: ", sign)

	// 转回PublicKey
	toPubkey, err := crypto.SigToPub([]byte(hash), sign)
	if err != nil {
		fmt.Println("sig to pub err:", toPubkey)
	}

	fmt.Println("sigtoPub:", *toPubkey)
}