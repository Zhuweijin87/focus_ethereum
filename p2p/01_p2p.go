package main

import (
	"fmt"
	"log"
	"net"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/p2p"
	"github.com/ethereum/go-ethereum/p2p/discover"
)

func newProtocol() []p2p.Protocol {
	protocol := p2p.Protocol {
		Name:    "MyProtocol",
		Version: 1,
		Length:  10,
		Run:     whisper.HandlePeer,
		NodeInfo: func() interface{} {
			return map[string]interface{}{
				"version":        "Protocol_v0.5",
				"maxMessageSize": 200,
				"minimumPoW":     0.2,
			}
		},
	}

	return []p2p.Protocol{protocol}
}

func main() {
	uniqueKey, _ := crypto.GenerateKey()
	var peers []*discover.Node 

	peer := discover.MustParseNode("enode://b89172e36cb79202dd0c0822d4238b7a7ddbefe8aa97489049c9afe68f71b10c5c9ce588ef9b5df58939f982c718c59243cc5add6cebf3321b88d752eac02626@182.254.155.208:33333")
	peers = append(peers, peer)

	server := p2p.Server {
		Config: p2p.Config {
			PrivateKey: uniqueKey, // 必须(私钥)
			Name: "MyP2PTest",
			MaxPeers: 10,
			Protocols: newProtocol(), // P2P 协议信息
			// NAT: nat.Any(),  // NAT 协议
			BootstrapNodes: peers,
			StaticNodes: peers,
			TrustedNodes: peers,
		},
	}

	err := server.Start()
	if err != nil {
		fmt.Println("fail to start")
		return 
	}

	defer server.Stop()
	select{}
}