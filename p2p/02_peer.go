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

func main() {
	var id discover.NodeID 
	caps := []p2p.Cap{{"foo", 2}, {"bar", 3}}
	id = randomID()

	peer := p2p.NewPeer(id, "MyPeer", caps)

	fmt.Println(peer.ID(), peer.Name())

	fmt.Println(peer.String())
}