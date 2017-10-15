package main

import (
	"bufio"
	"crypto/ecdsa"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/p2p"
	"github.com/ethereum/go-ethereum/p2p/discover"
	_ "github.com/ethereum/go-ethereum/p2p/nat"
	whisper "github.com/ethereum/go-ethereum/whisper/whisperv5"
)

var server *p2p.Server
var quit chan struct{} 

func main() {
	// 创建一个whisper
	wsp := whisper.New(nil)

	asymKeyId, err := wsp.NewKeyPair()
	if err != nil {
		fmt.Println("fail to New Key Pair")
		return
	}

	asymKey, err := wsp.GetPrivateKey(asymKeyId)
	if err != nil {
		fmt.Println("fail to get private key")
		return
	}

	// 设置中心节点
	var peers []*discover.Node
	peer := discover.MustParseNode("enode://b89172e36cb79202dd0c0822d4238b7a7ddbefe8aa97489049c9afe68f71b10c5c9ce588ef9b5df58939f982c718c59243cc5add6cebf3321b88d752eac02626@182.254.155.208:33333")
	peers = append(peers, peer)

	// 定义P2P Server
	server = &p2p.Server{
		Config: p2p.Config{
			PrivateKey: asymKey,
			MaxPeers:   10,
			Name:       common.MakeName("My Chat", "5.0"),
			Protocols:  wsp.Protocols(), //使用whisper协议 
			NAT:            nat.Any(),
			BootstrapNodes: peers,
			StaticNodes:    peers,
			TrustedNodes:   peers,
		},
	}

	// 配置时，需要将信任节点添加到数据库
	err = server.Start() // 开启服务
	if err != nil {
		fmt.Println(err)
	}
	defer server.Stop()


	// 等待节点连接
	waitCount := 100
	var count int
	var connected bool
	for !connected {
		time.Sleep(time.Millisecond * 500)
		connected = server.PeerCount() > 0

		if count > waitCount {
			fmt.Println("Connect timeout, exit!")
			return
		}
		count++
	}
	fmt.Println("Connect to peer OK")


	symkeyId, err := wsp.AddSymKeyFromPassword("123456")
	if err != nil {
		fmt.Println("fail to add symkey: ", err)
		return
	}

	symKey, err := wsp.GetSymKey(symkeyId)
	if err != nil {
		fmt.Println("fail to get symkey: ", err)
		return
	}

	// 创建TOPIC
	var topic whisper.TopicType
	copy(topic[:], common.FromHex("3ea456f2"))
	fmt.Printf("Topic %v Create\n", topic)

	// 创建Filter
	filter := whisper.Filter{
		KeySym:   symKey,
		KeyAsym:  asymKey,
		Topics:   [][]byte{topic[:]},
		AllowP2P: true, // 设置false,也可以通信
	}
	filterId, err := wsp.Subscribe(&filter)
	if err != nil {
		fmt.Println("fail to subscribe: ", err)
		return
	}
	fmt.Println("subscribe filter")

	// 启动whisper
	wsp.Start(nil)
	defer wsp.Stop()

	// 接收消息
	go RecvMessage(wsp, filterId, asymKey)
	// 发送消息
	SendMessage(wsp, asymKey, symKey, topic)
}

func RecvMessage(wsp *whisper.Whisper, filterId string, asymKey *ecdsa.PrivateKey) {
	filter := wsp.GetFilter(filterId)
	if filter == nil {
		fmt.Println("fail to Get filter:", filterId)
		return
	}

	ticker := time.NewTicker(time.Millisecond * 50)
	for {
		select {
		case <-ticker.C:
			messages := filter.Retrieve()
			for _, msg := range messages {
				printMessage(msg, asymKey)
			}
		}
		case <-quit:
			return 
	}
}

func printMessage(msg *whisper.ReceivedMessage, asymKey *ecdsa.PrivateKey) {
	text := string(msg.Payload)
	timestamp := time.Unix(int64(msg.Sent), 0).Format("2006-01-02 15:04:05")
	var address common.Address
	if msg.Src != nil {
		address = crypto.PubkeyToAddress(*msg.Src)
	}
	if whisper.IsPubKeyEqual(msg.Src, &asymKey.PublicKey) {
		fmt.Printf("%s <self> => %s\n", timestamp, text) // 自己发出信息
	} else {
		fmt.Printf("%s from [%x]=> %s\n", timestamp, address, text) // 其他节点的信息
	}
}

func SendMessage(wsp *whisper.Whisper, asymKey *ecdsa.PrivateKey, symKey []byte, topic whisper.TopicType) {
	for {
		buf := scanLine()
		if buf == "quit" {
			quit <- struct{}{}
			return 
		}
		sendMsg(wsp, []byte(buf), asymKey, symKey, topic)
	}
}

func waitConnect(times int) bool {
	for i := 0; i < times; i++ {
		if server.PeerCount() > 0 {
			//fmt.Println("Peer Connected.")
			return true
		}
		time.Sleep(time.Millisecond * 500)
	}
	return false
}

func scanLine() string {
	txt, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		fmt.Println("input error: %s", err)
		return ""
	}
	txt = strings.TrimRight(txt, "\n\r")
	return txt
}

func sendMsg(wsp *whisper.Whisper, payload []byte, asymKey *ecdsa.PrivateKey, symKey []byte, topic whisper.TopicType) common.Hash {
	params := whisper.MessageParams{
		Src:      asymKey,
		KeySym:   symKey,
		Payload:  payload,
		Topic:    topic,
		TTL:      whisper.DefaultTTL,
		PoW:      whisper.DefaultMinimumPoW,
		WorkTime: 5,
	}

	// 生成一个没有签名，没有加密的消息
	msg, err := whisper.NewSentMessage(&params)
	if err != nil {
		fmt.Printf("failed to create new message: %s\n", err)
		return common.Hash{}
	}

	// 打包消息
	envelope, err := msg.Wrap(&params)
	if err != nil {
		fmt.Printf("failed to seal message: %v\n", err)
		return common.Hash{}
	}

	// 将envelope 添加到whisper发送队列中, 并没有实际发送出去
	err = wsp.Send(envelope)
	if err != nil {
		fmt.Printf("failed to send message: %v\n", err)
		return common.Hash{}
	}

	// Hash
	return envelope.Hash()
}
