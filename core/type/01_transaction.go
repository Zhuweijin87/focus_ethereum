package main

import (
	"fmt"
	"math/big"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/common"
)

func main() {
	addr := common.HexToAddress("b94f5374fce5edbc8e2a8697c15331677e6ebf0b")
	trans := types.NewTransaction(1, addr, big.NewInt(10), big.NewInt(2000), big.NewInt(1), common.FromHex("4455"))
	if trans != nil {
		fmt.Println("fail to new transaction")
		return 
	}
	fmt.Println(trans)
}