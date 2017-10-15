package main

import (
	"fmt"
	"net"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/p2p/discover"
)

// 关于NodeID的算法
func main() {
	// 创建一个NodeID(节点ID)
	var id string 
	id = "0x000000000000000000000000000000000000000000000000000000000000000000000000000000806ad9b61fa5ae014307ebdc964253adcd9f2c0a392aa11abc"
	nodeId, err := discover.HexID(id) // 128位，以0x开头，节点的公钥
	if err != nil {
		fmt.Println(err)
		return 
	}
	fmt.Println(nodeId)

	
	key, _ := crypto.GenerateKey()
	nodeId = discover.PubkeyID(&key.PublicKey)
	//fmt.Println("NodeID: ",nodeId)
	fmt.Println(nodeId.String())
	fmt.Println(nodeId.GoString())

	// 创建节点Node
	// 全节点: IP不为nil，且不为组播地址或者未指定地址, UDP,TCP端口不为0，
	node := discover.NewNode(nodeId, net.ParseIP("127.0.0.1"), 0, 3001)

	fmt.Println(node.Incomplete())
	fmt.Println(node.String())  // 几点转化string

}