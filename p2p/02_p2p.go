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

func main() {

	lsn, err := net.Listen("tcp", "127.0.0.1:8200")
	if err != nil {
		log.Fatalf("Could not setup listener: %v\n", lsn)
	}
	defer lsn.Close()

	accepted := make(chan net.Conn)

	go func() {
		conn, err := lsn.Accept()
		if err != nil {
			log.Fatal("accept error:", err)
		}
		accepted <- conn
	}()

	conned := make(chan *p2p.Peer)
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

	defer close(conned)
	defer srv.Stop()

	tcpAddr := lsn.Addr().(*net.TCPAddr)
	srv.AddPeer(&discover.Node{ID: randId, IP: tcpAddr.IP, TCP: uint16(tcpAddr.Port)})

	select {
	case conn := <- accepted:
		defer conn.CLose()

		select {
		case peer := <- conned:
				if peer.ID() != randId {
					log.Printf("peer has wrong id\n")
				}

				if peer.Name() != "test" {
					log.Printf("peer has wrong name\n")
				}

				if peer.RemoteAddr().String() != conn.LocalAddr().String() {
					log.Printf("peer started with wrong conn: got %v, want %v", peer.RemoteAddr(), conn.LocalAddr())
				}

				peers := srv.Peers()

				if !reflect.DeepEqual(peers, []*Peer{peer}) {
					log.Printf("Peers mismatch: got %v, want %v", peers, []*Peer{peer})
				}

		case <- time.After(1 * time.Second):
			log.Printf("server did not launch peer within one second")
		}
	case <- time.After(1 * time.Second):
			log.Printf("server did not launch peer within one second")
	}
}