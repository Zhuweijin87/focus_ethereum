package main

import (
	"fmt"
	_ "time"

	wsp "github.com/ethereum/go-ethereum/whisper/whisperv5"
)

func BasicWhisper() {
	w := wsp.New(&wsp.DefaultConfig)  // 默认的配置
	p := w.Protocols()
	shh := p[0]

	if shh.Name != wsp.ProtocolName {
		fmt.Printf("failed Protocol Name: %v\n", shh.Name)
	}

	if uint64(shh.Version) != wsp.ProtocolVersion {
		fmt.Printf("failed Protocol Version: %v\n", shh.Version)
	}
	if shh.Length != wsp.NumberOfMessageCodes {
		fmt.Printf("failed Protocol Length: %v\n", shh.Length)
	}
	if shh.Run == nil {
		fmt.Printf("failed shh.Run.")
	}

	if uint64(w.Version()) != wsp.ProtocolVersion {
		fmt.Printf("failed whisper Version: %v\n", shh.Version)
	}

	mail := w.Envelopes() 
	if len(mail) == 0 {
		fmt.Println("no envelope")
	}

	id, err := w.NewKeyPair() // 生产新的加密身份，并保存相应的密钥信息
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(id) 

	if !w.HasKeyPair(id) {
		fmt.Printf("%v has no keypair\n", id)
	}
	
	// 获取密钥
	pk, err := w.GetPrivateKey(id) 
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(pk.PublicKey)
}

func main() {
	BasicWhisper()
}