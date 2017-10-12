package main


import (
	"fmt"
	mrand "math/rand"
	"github.com/ethereum/go-ethereum/crypto"
	wsp "github.com/ethereum/go-ethereum/whisper/whisperv5"
	"github.com/ethereum/go-ethereum/rlp"
)

const (
	aesKeyLength = 32
)

// 初始化消息
func generateMessageParams() (*wsp.MessageParams, error) {
	buf := make([]byte, 4)
	mrand.Read(buf)
	sz := mrand.Intn(400)

	//Options specifies the exact way a message should be wrapped into an Envelope
	var p wsp.MessageParams 

	p.PoW = 0.01 
	p.WorkTime = 1
	p.TTL = uint32(mrand.Intn(1024))  // 有效时间 
	p.Payload = make([]byte, sz) // 传输的数据
	p.KeySym = make([]byte, aesKeyLength) // 32位， 对称密钥
	mrand.Read(p.Payload)  // 随机产生数据
	mrand.Read(p.KeySym)   // 随机产生对称密钥
	p.Topic = wsp.BytesToTopic(buf)

	var err error
	p.Src, err = crypto.GenerateKey()  // 私钥
	if err != nil {
		return nil, err
	}

	return &p, nil
}

// Envelope 信息
func ShowEnvelop(env *wsp.Envelope) {
	fmt.Println("Version:", env.Version)
	fmt.Println("Expiry:", env.Expiry)  // 当前时间 + TTL(有效期)
	fmt.Println("TTL:", env.TTL)
	fmt.Println("Topic:", env.Topic)
	fmt.Printf("AESNonce: %x\n", env.AESNonce)  // 对称加密 通过params.KeySym 对称密钥
	fmt.Printf("Data: %x\n", env.Data)
	fmt.Println("EnvNonce: ", env.EnvNonce)
}

// 消息打包
func MsgWrap() {
	seed := int64(1777444222)
	mrand.Seed(seed)

	target := 128.0

	params, err := generateMessageParams()
	if err != nil {
		fmt.Println(err)
		return 
	}

	// 创建一个初始的没有签名，没有加密的消息(参数)
	msg, err := wsp.NewSentMessage(params) // msg 为[]byte
	if err != nil {
		fmt.Println(err)
		return 
	}

	params.TTL = 1  // 最长时间
	params.WorkTime = 12
	params.PoW = target  // 工作量 

	// 将消息体打包如信封，用于在p2p网络传输 
	env, err := msg.Wrap(params) 
	if err != nil { 
		fmt.Println(err)
		return 
	}

	ShowEnvelop(env)

	pow := env.PoW() // 将env数据通过RPL编码,取前32位字节 ...
	if pow < target {
		fmt.Println("failed Wrap with seed %d: pow < target (%f vs. %f)", seed, pow, target)
	}
	fmt.Println("PoW:", pow)

	// 算Hash
	hash := env.Hash() 
	fmt.Println("Hash:", hash)

	msg2, err := wsp.NewSentMessage(params)
	if err != nil {
		fmt.Println(err)
		return 
	}

	// 设置参数过大
	params.TTL = 1000000
	params.WorkTime = 1
	params.PoW = 10000000.0

	_, err = msg2.Wrap(params)
	if err == nil {
		fmt.Println("unexpectedly reached the PoW target with seed %d", seed)
		return 
	}
	
}

func MsgSeal() {
	// 生产随机数的随机种子
	seed := int64(1976726903)
	mrand.Seed(seed)

	params, err := generateMessageParams()
	if err != nil {
		fmt.Println(err)
		return 
	}

	msg, err := wsp.NewSentMessage(params)
	if err != nil {
		fmt.Println(err)
		return 
	}

	params.TTL = 1
	aesnonce := make([]byte, 12)
	mrand.Read(aesnonce)

	// 创建一个信封
	env := wsp.NewEnvelope(params.TTL, params.Topic, aesnonce, msg)
	if err != nil {
		fmt.Println(err)
		return 
	}

	ShowEnvelop(env)

	env.Expiry = uint32(seed)
	target := 32.0
	params.WorkTime = 4
	params.PoW = target 
	env.Seal(params)

	pow := env.PoW()
	if pow < target {
		fmt.Println("failed Wrap with seed %d: pow < target (%f vs. %f).", seed, pow, target)
		return 
	}

	fmt.Println("PoW: ", pow)
}

// 消息填充
func MsgPadding() {
	
}

func RplEncode() {
	// 随机种子
	seed := time.Now().Unix()
	mrand.Seed(seed)

	params, err := generateMessageParams()
	if err != nil {
		fmt.Println(err)
		return 
	}

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

	// RLP 编码
	raw, err := rlp.EncodeToBytes(env)
	if err != nil {
		fmt.Println(err)
		return 
	}

	//fmt.Println("Encode Raw:", raw)

	// RLP 解码
	var  decoded wsp.Envelope 
	err = rlp.DecodeBytes(raw, &decoded)
	if err != nil {
		fmt.Println("解码失败")
		return 
	}

	he := env.Hash()
	hd := decoded.Hash()

	if he != hd {
		fmt.Println("编码-解码校验失败")
		return 
	}

	fmt.Println("编码-解码成功")
}

func main() {

	//MsgWrap()

	//MsgSeal()

	RplEncode()

}
