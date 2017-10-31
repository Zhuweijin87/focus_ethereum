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

func newCompatibeEnvelope(f *wsp.Filter) *wsp.Envelope {
	params, err := generateMessageParams()
	if err != nil {
		fmt.Println(err)
		return nil
	}

	params.KeySym = f.KeySym
	params.Topic = wsp.BytesToTopic(f.Topics[2])
	sentmsg, err := wsp.NewSentMessage(params)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	env, err := sentmsg.Wrap(params)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return env
}

// Filter测试
type FilterCase struct {
	f      *wsp.Filter
	id     string
	alive  bool
	msgCnt int
}

func newFilterCases(num int) []FilterCase {
	cases := make([]FilterCase, num)
	for i := 0; i < num; i++ {
		f, _ := newFilter(true)
		cases[i].f = f
		cases[i].alive = (mrand.Int()&int(1) == 0)
	}

	return cases
}

func FilterWatcher() {
	NewRandSeed()

	const NumMsg = 256
	const NumFilter = 16

	var (
		firstId string
		e       *wsp.Envelope
		j       uint32
	)
	w := wsp.New(&wsp.Config{})
	filters := wsp.NewFilters(w)

	cases := newFilterCases(NumFilter)
	for i := 0; i < NumFilter; i++ {
		cases[i].f.Src = nil
		fid, err := filters.Install(cases[i].f)
		if err != nil {
			fmt.Println(err)
			return
		}

		cases[i].id = fid
		if len(firstId) == 0 {
			firstId = fid
		}
	}

	var envelopes [NumMsg]*wsp.Envelope
	for i := 0; i < NumMsg; i++ {
		j = mrand.Uint32() % NumFilter
		e = newCompatibeEnvelope(cases[j].f)
		envelopes[i] = e
		cases[j].msgCnt++
	}

	for i := 0; i < NumMsg; i++ {
		filters.NotifyWatchers(envelopes[i], false)
	}

	var total int
	var mail []*wsp.ReceivedMessage
	var count [NumFilter]int

	for i := 0; i < NumFilter; i++ {
		mail = cases[i].f.Retrieve()
		count[i] = len(mail)
		total += len(mail)
	}

	if total != NumMsg {
		fmt.Printf("failed with: total = %d, want: %d\n", total, NumMsg)
	}

	for i := 0; i < NumFilter; i++ {
		mail = cases[i].f.Retrieve()
		if len(mail) != 0 {
			fmt.Printf("failed with: i = %d.", i)
		}

		if cases[i].msgCnt != count[i] {
			fmt.Printf("failed with: count[%d]: get %d, want %d.", i, cases[i].msgCnt, count[i])
		}
	}
}

func main() {
	FilterWatcher()
}
