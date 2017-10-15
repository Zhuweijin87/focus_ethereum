package main

import (
	"fmt"

	_ "github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/p2p/discover"
)

// 关于NodeID的算法
func main() {
	// 创建一个NodeID(节点ID)
	var id string 
	id = "0x000000000000000000000000000000000000000000000000000000000000000000000000000000806ad9b61fa5ae014307ebdc964253adcd9f2c0a392aa11abc"
	nodeId, err := discover.HexID(id) // 128位，以0x开头
	if err != nil {
		fmt.Println(err)
		return 
	}
	fmt.Println(nodeId)

	/*
	key, _ := crypto.GenerateKey()
	nodeId := discover.PubkeyID(key)
	fmt.Println(key)
	fmt.Println(nodeId)
	*/
}