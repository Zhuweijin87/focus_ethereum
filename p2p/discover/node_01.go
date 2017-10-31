package main

import (
	"fmt"
	"net"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/p2p/discover"
	"github.com/ethereum/go-ethereum/rlp"
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
	// 将公钥转化节点ID
	nodeId = discover.PubkeyID(&key.PublicKey)
	//fmt.Println("NodeID: ",nodeId)
	fmt.Println(nodeId.String())
	fmt.Println(nodeId.GoString())

	// 创建节点Node
	// 全节点: IP不为nil，且不为组播地址或者未指定地址, UDP,TCP端口不为0，
	node := discover.NewNode(nodeId, net.ParseIP("127.0.0.1"), 0, 3001)

	fmt.Println(node.Incomplete()) // 是否为全节点
	fmt.Println(node.String())     // 几点转化string

	// Node 编码
	blob, err := rlp.EncodeToBytes(node)
	if err != nil {
		fmt.Println("fail to rlp encode bytes: ", err)
		return
	}
	fmt.Println("Node编码:", string(blob))

	//Node 解码
	var new discover.Node
	err = rlp.DecodeBytes(blob, &new)
	if err != nil {
		fmt.Println("fail to rlp decode byte: ", err)
		return
	}

	fmt.Println("Node解码:", new)
}
