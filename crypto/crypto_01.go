package main

import (
	"fmt"

	// 底层使用 crypto/ecdsa 
	"github.com/ethereum/go-ethereum/crypto" 
)

func main() {
	// 生成私钥 PrivateKey
 	privKey, _ := crypto.GenerateKey() 
	
	// 将私钥转化为[]byte 
	dump_priv := crypto.FromECDSA(privKey)
	fmt.Printf("private key: %x\n", dump_priv)

	// 获取公钥 
	pubKey := privKey.PublicKey  

	// 将公钥转化为 []byte 
	dump_pub := crypto.FromECDSAPub(&pubKey)
	fmt.Printf("public key: %x\n", dump_pub)

	// 转化为地址(取后20位)
	addr := crypto.PubkeyToAddress(pubKey) 

	fmt.Printf("address str: %x\n", addr) 

	// 哈希
	/*
	hash := crypto.Keccak256Hash([]byte("zhuweijin"))
	fmt.Println("hash: ", hash)
	*/

	hash := "ab12345ba2e43567ab12345ba2e43567"
	// 用私钥签名 
	sign, _ := crypto.Sign([]byte(hash), privKey)
	fmt.Printf("签名之后: %x\n", sign)
}
