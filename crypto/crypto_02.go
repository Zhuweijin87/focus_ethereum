package main

import (
	"fmt"

	"github.com/ethereum/go-ethereum/crypto" 
	"github.com/ethereum/go-ethereum/common" 
)

func main() {
	s := []byte("123456789abcd")

	hash256 := crypto.Keccak256(s)  // 32位
	fmt.Println(len(hash256), hash256)

	hash512 := crypto.Keccak512(s) // 64位
	fmt.Println(len(hash512), hash512)

	hash := crypto.Keccak256Hash(s)
	fmt.Println(hash)

	// ### 签名 ###
	address := "970e8128ab834e8eac17ab8e3812f010678cf791"
	privKey := "606c90b9ce05ab4bb2535957334eba473f9ff197e98d546e71f8556097031f0e"
	// 解析secp256k1密钥
	key, _ := crypto.HexToECDSA(privKey) 
	addr := common.HexToAddress(address)
	fmt.Println(addr)
	msg := crypto.Keccak256([]byte("foo")) 
	 // 获取签名
	sig, err := crypto.Sign(msg, key)
	if err != nil {
		fmt.Println(err)
		return 
	}
	fmt.Printf("sign: %x\n", sig)

	// 检验签名
	recoverPub, err := crypto.Ecrecover(msg, sig)
	if err != nil {
		fmt.Println(err)
		return 
	}
	fmt.Printf("recover: %x\n", recoverPub)

	pubKey := crypto.ToECDSAPub(recoverPub) // 转化为ECDSA公钥
	recoverAddr := crypto.PubkeyToAddress(*pubKey)
	if addr != recoverAddr {
		fmt.Println("Address mismatch")
		fmt.Println("\t", addr, "\n\t", recoverAddr)
		return 
	}
}