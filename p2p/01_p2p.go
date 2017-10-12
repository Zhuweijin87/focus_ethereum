package main

import (
	"fmt"
	"os"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/p2p"
)

const messageId = 0

type Message string

// 自定义协议
func MyProtocol() p2p.Protocol {
	return p2p.Protocol{
		Name:    "MyProtocol",
		Version: 1,
		Length:  1,
		Run:     MyHandler,
	}
}

func main() {
	nodekey, _ := crypto.GenerateKey() // 生成密钥
	config := p2p.Config{
		MaxPeers:   2,
		PrivateKey: nodekey,
		Name:       "MyP2PServer",
		ListenAddr: ":3001",
		Protocols:  []p2p.Protocol{MyProtocol()},
		//NoDial: true,  
	}

	srv := &p2p.Server{
		Config: config,
	}

	fmt.Println("P2P Server Start...")
	err := srv.Start()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	select {}
}

func MyHandler(peer *p2p.Peer, ws p2p.MsgReadWriter) error {
	for {
		msg, err := ws.ReadMsg()
		if err != nil {
			return err
		}

		var myMessage Message
		err = msg.Decode(&myMessage)
		if err != nil {
			continue
		}

		switch myMessage {
		case "foo":
			err = p2p.SendItems(ws, messageId, "bar")
			if err != nil {
				return err
			}
		default:
			fmt.Println("recv:", myMessage)
		}
	}

	return nil
}