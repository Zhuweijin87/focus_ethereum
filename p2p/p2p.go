package main

import (
	"fmt"
	"os"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/p2p"
)

const messageID = 0

type Message string

func MyProtocal() p2p.Protocol {
	return p2p.Protocol {
		Name : "MyProtocal",
		Version: 1,
		Length: 1,
		Run: msgHandle,
	}
}

func msgHandler(p2p *p2p.Peer, ws p2p.MsgReaderWriter) error {
	for {
		msg, err := ws.ReadMsg()
		if err != nil {
			return err
		}
		
	}
}