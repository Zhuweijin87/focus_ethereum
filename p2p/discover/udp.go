package main

import (
	"fmt"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/p2p/discover"
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
