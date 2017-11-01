## 数据结构说明 

### whisper.go 
p2p层之上的协议，用于p2p网络数据的传输。  

首先看下whisper的数据结构
```go
type Whisper struct {
	protocol p2p.Protocol   // 用于定义whisper的p2p协议
	filters  *Filters       // 消息过滤器 

	privateKeys map[string]*ecdsa.PrivateKey // 存储的私钥
	symKeys     map[string][]byte            // 
	keyMu       sync.RWMutex                 // 

	poolMu      sync.RWMutex              // 读写锁
	envelopes   map[common.Hash]*Envelope // 装载消息的信封
	expirations map[uint32]*set.SetNonTS  // 

	peerMu sync.RWMutex       // 节点锁
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
基本的配置  
```go
type Config struct {
	MaxMessageSize     uint32    消息大小 (默认 1024 * 1024)
	MinimumAcceptedPOW float64   接受最小的PoW (默认 0.2)
}
```
主要函数说明:  
* New(cfg Config) *Whisper : 创建一个用于p2p通信的whisper客户端  
whisper中定义的协议  
```go
	whisper.protocol = p2p.Protocol{
		Name:    ProtocolName,
		Version: uint(ProtocolVersion),
		Length:  NumberOfMessageCodes, (默认64)
		Run:     whisper.HandlePeer,  // 主要处理p2p通信  
		NodeInfo: func() interface{} {
			return map[string]interface{}{
				"version":        ProtocolVersionStr,
				"maxMessageSize": whisper.MaxMessageSize(),
				"minimumPoW":     whisper.MinPow(),
			}
		},
	}
```
* Start(*p2p.Server) error : 启动whisper  
这个函数自带的参数是p2p服务，也就是说它是适用于p2p服务的。主要处理消息队列中的数据。 
也就是whisper结构体中的messageQueue，p2pMsgQueue。  

* Stop() : 关闭whisper 

+ HandlePeer(peer *p2p.Peer, rw p2p.MsgReadWriter) : whisper定义的协议  
创建对等节点， 
具体处理逻辑参考peer.go  

+ RegisterServer(server MailServer) : 注册邮件服务  
归档，备份数据处理  

### peer.go 
基于whisper协议对等节点(与之对应的p2p的对等节点)  

基本数据结构  
```go
type Peer struct {
	host    *Whisper  // whisper协议
	peer    *p2p.Peer 
	ws      p2p.MsgReadWriter // p2p读写
	trusted bool

	known *set.Set // Messages already known by the peer to avoid wasting bandwidth

	quit chan struct{}
}
```

```go
// 消息体数据(不含加密，签名)
type MessageParams struct {
	TTL      uint32             // 消息体有效
	Src      *ecdsa.PrivateKey  // 公钥 : 做为显示消息的来源
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
	Src     *ecdsa.PublicKey // 消息接受者公钥 : 可以判断接收的消息是不是自己的，是自己的可以不处理
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

