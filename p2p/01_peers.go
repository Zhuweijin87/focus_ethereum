package main

import (
	"fmt"
	"net"
	"crypto/ecdsa"
	"github.com/ethereum/go-ethereum/p2p"
	"github.com/ethereum/go-ethereum/rpl"
	"github.com/ethereum/go-ethereum/crypto/sha3"
	"github.com/ethereum/go-ethereum/p2p/discover"
)

type connFlag int

// server.go
type conn struct {
	fd 		net.Conn
	transport
	flags 	connFlag
	cont  	chan error      	// The run loop uses cont to signal errors to setupConn.
	id    	discover.NodeID 	// valid after the encryption handshake
	caps  	[]p2p.Cap           // valid after the protocol handshake
	name  	string          	// valid after the protocol handshake
}

// server.go
type transport interface {
	// The two handshakes.
	doEncHandshake(prv *ecdsa.PrivateKey, dialDest *discover.Node) (discover.NodeID, error)
	doProtoHandshake(our *protoHandshake) (*protoHandshake, error)
	
	p2p.MsgReadWriter
	close(err error)
}

// peer.go
type protoHandshake struct {
	Version    uint64
	Name       string
	Caps       []p2p.Cap
	ListenPort uint64
	ID         discover.NodeID
	Rest 		[]rlp.RawValue `rlp:"tail"`
}

func newTestTransport(id discover.NodeID, fd net.Conn) transport {
	wrapped := newRLPX(fd).(*rlpx)
	wrapped.rw = newRLPXFrameRW(fd, secrets{
		MAC:        zero16,
		AES:        zero16,
		IngressMAC: sha3.NewKeccak256(),
		EgressMAC:  sha3.NewKeccak256(),
	})
	return &testTransport{id: id, rlpx: wrapped}
}

func newPeer(conn *conn, protocols []p2p.Protocol) *Peer {
	protomap := matchProtocols(protocols, conn.caps, conn)
	p := &Peer{
		rw:       conn,
		running:  protomap,
		created:  mclock.Now(),
		disc:     make(chan DiscReason),
		protoErr: make(chan error, len(protomap)+1), // protocols + pingLoop
		closed:   make(chan struct{}),
		log:      log.New("id", conn.id, "conn", conn.flags),
	}
	return p
}

func testPeer(proto []p2p.Protocol) (func(), *conn, *p2p.Peer, <-chan error) {
	fd1, fd2 := net.net.Pipe()
	c1 := &conn{fd: fd1, transport: newTestTransport(randomID(), fd1)}
	c2 := &conn{fd: fd2, transport: newTestTransport(randomID(), fd2)}
	for _, p := range protos {
		c1.caps = append(c1.caps, p.cap())
		c2.caps = append(c2.caps, p.cap())
	}

	peer := newPeer(c1, protos)
	errc := make(chan error, 1)
	go func() {
		_, err := peer.run()
		errc <- err
	}()

	closer := func() { c2.close(errors.New("close func called")) }
	return closer, c2, peer, errc
}

func main() {
	proto := p2p.Protocol {
		Name: "Test",
		Length: 5,
		Run : func(peer *p2p.Peer, rw p2p.MsgReadWriter) error {
			if err := p2p.ExpectMsg(rw, 2, []uint{1}); err != nil {
				fmt.Println("ExpectMsg 2:", err)
				return err
			}
			if err := p2p.ExpectMsg(rw, 3, []uint{2}); err != nil {
				fmt.Println("ExpectMsg 3:", err)
				return err
			}
			if err := p2p.ExpectMsg(rw, 4, []uint{3}); err != nil {
				fmt.Println("ExpectMsg 3:", err)
				return err
			}
			return nil
		}
	}


}