package main

import (
	"fmt"
	"time"
	"github.com/ethereum/go-ethereum/p2p"
)

func readMsg(p *p2p.MsgPipeRW) {
	for {
		msg, err := p.ReadMsg()
		if err != nil {
			fmt.Println("Read Msg Err:", err)
			break
		}

		var data []byte
		msg.Decode(&data)
		fmt.Println("MsgCode:", msg.Code, "Data: ", string(data))
	}
	fmt.Println("End Of ReadMsg")
}

func main() {
	p1, p2 := p2p.MsgPipe() 
	defer p1.Close()
	defer p2.Close()

	go func() {
		for {
			p2p.Send(p1, 10, []byte("Hello Peter"))
			time.Sleep(1 * time.Second)
			p2p.Send(p2, 20, []byte("Hello Sofia"))
		}
	}()

	go readMsg(p1)
	readMsg(p2)
}