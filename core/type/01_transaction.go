package main

import (
	"fmt"
	"math/big"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
)

// type HomesteadSigner struct{ FrontierSigner }
// type FrontierSigner struct{} 

func main() {
	addr := common.HexToAddress("b94f5374fce5edbc8e2a8697c15331677e6ebf0b")
	// 创建一个交易 
	trans := types.NewTransaction(2, addr, big.NewInt(10), big.NewInt(2000), big.NewInt(1), common.FromHex("4455"))
	if trans == nil {
		fmt.Println("fail to new transaction")
		return 
	}
	fmt.Println(trans)

	// 签名
	hash := trans.SigHash(types.HomesteadSigner{})
	fmt.Printf("Hash: %x\n", hash)

	// 编码 
	encode, err := rlp.EncodeToBytes(trans)
	if err != nil {
		fmt.Println("fail to encode:", encode)
		return 
	}

	fmt.Printf("encode: %x\n", encode)

	from, err := types.Sender(types.HomesteadSigner{}, trans)
	if err != nil {
		fmt.Println("fail to sender: ", err)
		return 
	}

	fmt.Printf("From: %v\n", from)
}