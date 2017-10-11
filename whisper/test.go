package main

import (
	"fmt"
	"time"
	"math/rand"
)

func main() {
	fmt.Println(time.Now().Add(time.Second * time.Duration(1)).Unix(), time.Now().Unix())

	seed := time.Now().Unix()
	rand.Seed(seed)
	
	var a [][]byte
	a = make([][]byte, 8)
	for i:=0; i<8; i++ {
		a[i] = make([]byte, 4)
		rand.Read(a[i][:])
		a[i][0] = 0x01
	}

	fmt.Println(a)
}