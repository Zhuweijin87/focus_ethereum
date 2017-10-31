package main

import (
	"fmt"
	mrand "math/rand"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	wsp "github.com/ethereum/go-ethereum/whisper/whisperv5"
)

var seed int64
var aesKeyLength = 32

func NewRandSeed() {
	seed = time.Now().Unix()
	mrand.Seed(seed)
}

// 初始生产一个Filter
func newFilter(symmetric bool) (*wsp.Filter, error) {
	var f wsp.Filter
	f.Messages = make(map[common.Hash]*wsp.ReceivedMessage) // 接收的消息缓存的初始化

	const topicNum = 8
	f.Topics = make([][]byte, topicNum) //
	for i := 0; i < topicNum; i++ {
		f.Topics[i] = make([]byte, 4)
		mrand.Read(f.Topics[i][:])
		f.Topics[i][0] = 0x01
	}

	key, err := crypto.GenerateKey()
	if err != nil {
		fmt.Printf("generateFilter 1 failed with seed %d.", seed)
		return nil, err
	}

	f.Src = &key.PublicKey // 消息发送者地址 (公钥)

	if symmetric { // 使用对称密钥
		f.KeySym = make([]byte, aesKeyLength) // 与Topic相关的密钥
		mrand.Read(f.KeySym)
		f.SymKeyHash = crypto.Keccak256Hash(f.KeySym) // Keccak256Hash对称密钥，
	} else {
		f.KeyAsym, err = crypto.GenerateKey() // 收件人私钥
		if err != nil {
			fmt.Printf("generateFilter 2 failed with seed %d.", seed)
			return nil, err
		}
	}

	return &f, nil
}

// 初始化消息
func generateMessageParams() (*wsp.MessageParams, error) {
	buf := make([]byte, 4)
	mrand.Read(buf)
	sz := mrand.Intn(400)

	//Options specifies the exact way a message should be wrapped into an Envelope
	var p wsp.MessageParams

	p.PoW = 0.01
	p.WorkTime = 1
	p.TTL = uint32(mrand.Intn(1024))      // 有效时间
	p.Payload = make([]byte, sz)          // 传输的数据
	p.KeySym = make([]byte, aesKeyLength) // 32位， 对称密钥
	mrand.Read(p.Payload)                 // 随机产生数据
	mrand.Read(p.KeySym)                  // 随机产生对称密钥
	p.Topic = wsp.BytesToTopic(buf)

	var err error
	p.Src, err = crypto.GenerateKey() // 私钥
	if err != nil {
		return nil, err
	}

	return &p, nil
}

// 简单的Filter(添加)注册
func InstallFilter() {
	NewRandSeed()
	const SizeFilter = 256
	w := wsp.New(&wsp.Config{})  // 新建一个whisper(基于P2P网络应用),使用默认的Config
	filters := wsp.NewFilters(w) // 创建一个Filters监视平台

	f, err := newFilter(true)
	if err != nil {
		fmt.Println(err)
		return
	}

	fid, err := filters.Install(f) // 主要是给filter生产一个随机ID，并注册到监视容器中(id:Filter)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(fid)

	getFilter := filters.Get(fid)
	if getFilter == nil {
		fmt.Printf("Filter %s not exist\n", fid)
		return
	}

	fmt.Println("Filter Get OK!")

	// 删除Filter
	filters.Uninstall(fid)
	getFilter2 := filters.Get(fid)
	if getFilter2 == nil {
		fmt.Printf("Filter %s not exist\n", fid)
		return
	}
}

// 匹配相应的Envelope
func MatchEnvelope() {
	NewRandSeed()

	fsym, err := newFilter(true) //
	if err != nil {
		fmt.Println(err)
		return
	}

	fasym, err := newFilter(false) //
	if err != nil {
		fmt.Println(err)
		return
	}

	params, err := generateMessageParams()
	if err != nil {
		fmt.Println(err)
		return
	}

	params.Topic[0] = 0xFF

	msg, err := wsp.NewSentMessage(params)
	if err != nil {
		fmt.Println(err)
		return
	}

	env, err := msg.Wrap(params)
	if err != nil {
		fmt.Println(err)
		return
	}

	// 匹配envelope: 检验Topic是否匹配
	// 如果Topic 为空,则所有Envelope都过
	match := fsym.MatchEnvelope(env)
	if !match {
		fmt.Printf("failed MatchEnvelope symmetric with seed %d.\n", seed)
	}

	// 匹配envelope
	match = fasym.MatchEnvelope(env)
	if !match {
		fmt.Printf("failed MatchEnvelope symmetric with seed %d.\n", seed)
	}
}

// 消息体
func MatchMessageSym() {
	NewRandSeed()

	params, err := generateMessageParams()
	if err != nil {
		fmt.Println(err)
		return
	}

	f, err := newFilter(true)
	if err != nil {
		fmt.Println(err)
		return
	}

	const index = 1
	params.KeySym = f.KeySym
	params.Topic = wsp.BytesToTopic(f.Topics[index])

	sentmsg, err := wsp.NewSentMessage(params)
	if err != nil {
		fmt.Println(err)
		return
	}

	env, err := sentmsg.Wrap(params)
	if err != nil {
		fmt.Println(err)
		return
	}

	// 打开envelope, 并返回的是一个RecievedMessage
	// 将Envelope ==> RecievedMessage
	msg := env.Open(f)
	if msg == nil {
		fmt.Println("fail to open envelope")
		return
	}

	// MatchMeesage: 需要判断公钥是否一致即 Filter公钥 与 RecievedMessage公钥
	// 发送者消息不匹配
	if !f.MatchMessage(msg) {
		fmt.Println("fail to match message with src mismatch")
	}

	// 发送者消息匹配: 即公钥信息匹配
	*f.Src.X = *params.Src.PublicKey.X
	*f.Src.Y = *params.Src.PublicKey.Y
	if !f.MatchMessage(msg) {
		fmt.Println("fail to match message with src match")
	}

	// PoW 不足
	f.PoW = msg.PoW + 1.0
	if !f.MatchMessage(msg) {
		fmt.Println("fail to match message with insufficient PoW")
	}

	// PoW 足够
	f.PoW = msg.PoW / 2
	if !f.MatchMessage(msg) {
		fmt.Println("fail to match message with sufficient PoW")
	}

	// Topic 不匹配
	f.Topics[index][0]++
	if f.MatchMessage(msg) {
		fmt.Println("fail to match message with topic mismatch")
	}

	f.Topics[index][0]--
}

func MatchMessageAsym() {
	NewRandSeed()

	params, err := generateMessageParams()
	if err != nil {
		fmt.Println(err)
		return
	}

	f, err := newFilter(false)
	if err != nil {
		fmt.Println(err)
		return
	}

	const index = 1
	params.Topic = wsp.BytesToTopic(f.Topics[index])
	params.Dst = &f.KeyAsym.PublicKey
	keySymOrig := params.KeySym
	params.KeySym = nil

	sentmsg, err := wsp.NewSentMessage(params)
	if err != nil {
		fmt.Println(err)
		return
	}

	env, err := sentmsg.Wrap(params)
	if err != nil {
		fmt.Println(err)
		return
	}

	msg := env.Open(f)
	if msg == nil {
		fmt.Println(err)
		return
	}

	if !f.MatchMessage(msg) {
		fmt.Println("fail to match message with mismatch src")
	}

	*f.Src.X = *params.Src.PublicKey.X
	*f.Src.Y = *params.Src.PublicKey.Y
	if !f.MatchMessage(msg) {
		fmt.Println("fail to match message with match src")
	}

	// 加密方式不匹配
	f.KeySym = keySymOrig
	f.KeyAsym = nil
	if f.MatchMessage(msg) {
		fmt.Println("fail to match message with encrypto method mismatch")
	}
}

func main() {
	//InstallFilter()

	//MatchEnvelope()

	//MatchMessageSym()

	MatchMessageAsym()
}
