package main

import (
	"fmt"

	"github.com/ethereum/go-ethereum/p2p/discover"
	"github.com/ethereum/go-ethereum/crypto"
)

func main() {
	priv, _ := crypto.GenerateKey()

	table, err := discover.ListenUDP(priv, ":8001", nil, "", nil)
	if err != nil {
		fmt.Println(err)
		return 
	}

	fmt.Println(table.Self())
}