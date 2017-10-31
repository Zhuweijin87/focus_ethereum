package main

import (
	"fmt"
	"math/rand"

	"github.com/ethereum/go-ethereum/p2p"
	"github.com/ethereum/go-ethereum/p2p/discover"
)

func randomID() (id discover.NodeID) {
	for i := range id {
		id[i] = byte(rand.Intn(255))
	}
	return id
}

// 简单Peer创建
func SimplePeer() {
	var id discover.NodeID
	caps := []p2p.Cap{{"foo", 2}, {"bar", 3}}
	id = randomID()
	peer := p2p.NewPeer(id, "MyPeer", caps)
	fmt.Println(peer.String())
}

func NormalPeer() {
	protocol := p2p.Protocol{
		Name:   "Normal Peer",
		Length: 5,
		Run: func(peer *p2p.Peer, rw *p2p.MsgReadWriter) error {

			return nil
		},
	}

	caps := []p2p.Cap{{"foo", 2}, {"bar", 3}}
	peer := p2p.NewPeer(randomID(), "My Peer", caps)
}

func main() {
	SimplePeer()

}
