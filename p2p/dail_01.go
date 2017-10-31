package main

import (
	"fmt"
	_ "github.com/ethereum/go-ethereum/p2p"
	"github.com/ethereum/go-ethereum/p2p/discover"
	"net"
)

func main() {

	id := discover.MustHexID("1dd9d65c4552b5eb43d5ad55a2ee3f56c6cbc1c64a5c8d659f51fcd51bace24351232b8d7821617d2b29b54b81cdefb9b3e9c37d7fd5f63270bcc9e1a6f6a439")
	// 新建节点
	node := discover.NewNode(id, net.IP{127, 0, 0, 1}, 3333, 0)

	fmt.Println(node)
}
