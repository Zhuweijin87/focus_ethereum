package main

import (
	"fmt"
	"log"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/p2p"
)

func randomID() (id p2p.discover.NodeID) {
	for i := range id {
		id[i] = byte(rand.Intn(255))
	}
	return id
}

func newServer(id discover.NodeID, pf func(*Peer)) *p2p.Server {
	uniqueKey, _ := crypto.GenerateKey()
	config := p2p.Config{
		Name:       "test",
		MaxPeers:   10,
		ListenAddr: "192.168.1.195:8200",
		PrivateKey: uniqueKey,
	}

	serv := p2p.Server{
		Config:       config,
		newPeerHook:  pf,
		newTransport: func(fd net.Conn) transport { return newTestTransport(id, fd) },
	}

	if err := serv.Start(); err != nil {
		log.Println("fail to start server: ", err)
	}

	return serv
}

func servListen() {
	conn := make(chan *p2p.Peer)
	randId := randomID()

	srv := newServer()
	
}
