package main

import (
	"fmt"
	"time"
	"errors"
	"github.com/ethereum/go-ethereum/p2p"
)

var errProtocolReturned = errors.New("protocol returned")

func main() {
	baseProtocolLength := uint64(16)
	proto := p2p.Protocol{
		Name:   "a",
		Length: 5,
		Run: func(peer *p2p.Peer, rw p2p.MsgReadWriter) error {
			if err := p2p.ExpectMsg(rw, 2, []uint{1}); err != nil {
				fmt.Println(err)
				return err 
			}
			fmt.Println("Expect Msg 2")
			if err := p2p.ExpectMsg(rw, 3, []uint{2}); err != nil {
				fmt.Println(err)
				return err 
			}
			fmt.Println("Expect Msg 3")
			if err := p2p.ExpectMsg(rw, 4, []uint{3}); err != nil {
				fmt.Println(err)
				return err 
			}
			fmt.Println("Expect Msg 4")
			return nil
		},
	}

	closer, rw, _, errc := p2p.MyTestPeer([]p2p.Protocol{proto})
	defer closer()

	p2p.Send(rw, baseProtocolLength+2, []uint{1})
	p2p.Send(rw, baseProtocolLength+3, []uint{2})
	p2p.Send(rw, baseProtocolLength+4, []uint{3})

	select {
	case err := <-errc:
		if err != errProtocolReturned {
			fmt.Printf("peer returned error: %v\n", err)
		}
	case <-time.After(2 * time.Second):
		fmt.Printf("receive timeout\n")
	}
}