package main

import (
	"fmt"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/p2p"
)

// p2p 简单启动
func main() {
	privKey, _ := crypto.GenerateKey()
	config := p2p.Config{
		PrivateKey: privKey, 
	}

	server := &p2p.Server{
		Config: config,
	}

	if err := server.Start(); err != nil {
		fmt.Println("Fail to start server:", err)
		return 
	}
	defer server.Stop()

	fmt.Println("Server Start ...")
	select{}
}