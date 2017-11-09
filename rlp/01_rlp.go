package main

import (
	"fmt"
	_ "encoding/binary"
	"github.com/ethereum/go-ethereum/rlp"
)

func main() {
	var b [][]byte
	b = append(b, []byte("1000"))
	b = append(b, []byte("Hello"))
	b = append(b, []byte("234"))

	fmt.Println(b)
	o, err := rlp.EncodeToBytes(b)
	if err != nil {
		fmt.Println("err:", err)
		return 
	}

	fmt.Println(o)

	var c [][]byte
	err = rlp.DecodeBytes(o, &c)
	if err != nil {
		fmt.Println("err:", err)
	}

	fmt.Println(c)
}