package main

import (
	"fmt"
	"log"
	"net"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/p2p"
	"github.com/ethereum/go-ethereum/p2p/discover"
)

func randomID() (id discover.NodeID) {
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

	srv := newServer(randId, func(p *p2p.Peer) {
		if p.ID() != remid {
			log.Fatal("peer func called with wrong node id")
		}
		if p == nil {
			log.Fatal("peer func called with nil conn")
		}
		conn <- p
	})

	defer close(conn)
	defer srv.Stop()

	conn, err := net.DialTimeout("tcp", srv.ListenAddr, 5*time.Second)
	if err != nil {
		log.Fatal("could not dial:", err)
	}
	defer conn.Close()

	select {
	case peer <- conn:

	}
}
