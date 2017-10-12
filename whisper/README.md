## 数据结构说明 

```go
// 消息体数据(不含加密，签名)
type MessageParams struct {
	TTL      uint32             // 消息体有效
	Src      *ecdsa.PrivateKey  // 公钥
	Dst      *ecdsa.PublicKey   // 私钥
	KeySym   []byte             // 对称加密密钥
	Topic    TopicType          // Topic: 网络传输相关
	WorkTime uint32             // 
	PoW      float64            // 
	Payload  []byte             // 数据
	Padding  []byte
} 
```

```go
// 接收消息包
type ReceivedMessage struct {
	Raw []byte          // 接收到包的原数据

	Payload   []byte    // 实际数据
	Padding   []byte 
	Signature []byte    // 签名

	PoW     float64       // 工作量
	Sent    uint32        // 消息发送时间
	TTL     uint32        // 消息的生命期
	Src     *ecdsa.PublicKey // 消息接受者公钥
	Dst     *ecdsa.PublicKey // 消息接受者公钥 Filter公钥
	Topic   TopicType   // 

	SymKeyHash      common.Hash // 与Topic相关的对称密钥
	EnvelopeHash    common.Hash // Message envelope唯一Hash ID
	EnvelopeVersion uint64
}
```

```go
// 数据包过滤器
type Filter struct {
	Src        *ecdsa.PublicKey     // 消息发送者公钥
	KeyAsym    *ecdsa.PrivateKey    // 消息接受者私钥 
	KeySym     []byte            // Topic相关的密钥 
	Topics     [][]byte          // 过滤消息包的Topic
	PoW        float64           // 
	AllowP2P   bool              // 是否基于P2P网络
	SymKeyHash common.Hash       // 

	Messages map[common.Hash]*ReceivedMessage  // 接收到所有消息包
	mutex    sync.RWMutex
}

```
+ 关于KeyAsym与KeySym说明   
KeyAsym: 使用非对称密钥时，需要创建该私钥  
KeySym: 使用对称密钥时，KeySym需要创建(一般都是随机采用生成)，通过crypto.Keccak256Hash(KeySym)算出SymKeyHash


```go
type Filters struct {
	watchers map[string]*Filter   // 过滤器容器, 可以添加多种过滤器
	whisper  *Whisper           // whisper 网络
	mutex    sync.RWMutex
}
```

```go
type Whisper struct {
	protocol p2p.Protocol   // p2p协议
	filters  *Filters     // 消息过滤器

	privateKeys map[string]*ecdsa.PrivateKey // 存储的私钥
	symKeys     map[string][]byte            // 
	keyMu       sync.RWMutex                 // 

	poolMu      sync.RWMutex              // 
	envelopes   map[common.Hash]*Envelope // 装载消息的信封
	expirations map[uint32]*set.SetNonTS  // 

	peerMu sync.RWMutex       // 
	peers  map[*Peer]struct{} // 

	messageQueue chan *Envelope // Message queue for normal whisper messages
	p2pMsgQueue  chan *Envelope // Message queue for peer-to-peer messages 
	quit         chan struct{}  // 
	settings     syncmap.Map // 保存动态更改的配置设置: 如节点要求最小PoW, 最大消息大小，消息队列溢出提示等

	statsMu sync.Mutex // guard stats
	stats   Statistics // whisper节点统计信息

	mailServer MailServer // 邮件服务, 基于邮件服务发送数据
}
```

```go
type Envelope struct {
	Version  []byte
	Expiry   uint32
	TTL      uint32   // 有效时长 
	Topic    TopicType 
	AESNonce []byte
	Data     []byte   // 传输数据 
	EnvNonce uint64   // Envelope 封装时，所计算的Nonce

	pow  float64     // 消息所包含的PoW
	hash common.Hash // 
}
```

### Filter Match规则
+ Envelope匹配
1. Filter.PoW > 0, 且Envelope.pow < Filter.PoW 
3. Envelope.Topic存在于Filter.Topic

+ RecievedMessage匹配  
1. Filter.PoW > 0, 且Envelope.pow < RecievedMessage.PoW 
2. Filter.Src != nil，且 Filter.Src == RecievedMessage.Src  
3. Filter.KeyAsym.PublicKey 与 RecievedMessage.Dst 匹配问题  
4. Filter.SymKeyHash 与 RecievedMessage 匹配问题
5. 确认 Envelope.Topic存在于Filter.Topic 


### whisper
+ Start() : 启动多个线程，并处理p2p消息队列和普通whisper消息队列  
+ Stop() : 关闭消息队列的处理  
+ Subscribe() : 订阅，主要是添加filter，filter.Install() 
+ Unsubscribe(): 取消订阅，删除filter，filter.Uninstall()