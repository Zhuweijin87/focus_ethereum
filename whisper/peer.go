package main

import (
	"fmt"
	mrand "math/rand"
	"net"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/p2p"
	"github.com/ethereum/go-ethereum/p2p/discover"
	"github.com/ethereum/go-ethereum/p2p/nat"
)

// 初始化消息
func generateMessageParams() (*wsp.MessageParams, error) {
	buf := make([]byte, 4)
	mrand.Read(buf)
	sz := mrand.Intn(400)

	//Options specifies the exact way a message should be wrapped into an Envelope
	var p wsp.MessageParams

	p.PoW = 0.01
	p.WorkTime = 1
	p.TTL = uint32(mrand.Intn(1024))      // 有效时间
	p.Payload = make([]byte, sz)          // 传输的数据
	p.KeySym = make([]byte, aesKeyLength) // 32位， 对称密钥
	mrand.Read(p.Payload)                 // 随机产生数据
	mrand.Read(p.KeySym)                  // 随机产生对称密钥
	p.Topic = wsp.BytesToTopic(buf)

	var err error
	p.Src, err = crypto.GenerateKey() // 私钥
	if err != nil {
		return nil, err
	}

	return &p, nil
}

func BasicPeer() {
	seed := time.Now().Unix()
	mrand.Seed(seed)

	params, err := generateMessageParams()
	if err != nil {
		fmt.Println(err)
		return
	}

	params.PoW = 0.001
	msg, err := wsp.NewSentMessage(params)
	if err != nil {
		fmt.Println(err)
		return
	}

	env, err := msg.Wrap(params)
	if err != nil {
		fmt.Println(err)
		return
	}

}
