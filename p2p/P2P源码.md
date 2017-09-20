### P2P 解析
-----------------

#### P2P 目录结构
+ discover/
+ discv5/
+ nat/
+ netutil/

#### 数据结构 
服务信息 p2p/server.go
```
type Server struct {
    Config      //配置文件

    // 一些初始化的内部数据
    newTransport    func(net.Conn) transport  //传输协议
    newPeerHook     func(*Peer)

    lock            sync.Mutex
    ntab            discoverTable
    ourHandshake    *protoHandshake  // 握手协议
}

// 配置信息
type Config struct {
    PrivateKey          *ecdsa.PrivateKey // 私钥信息
    MaxPeers            int   // 最大节点数
    MaxPendingPeers     int 
    NoDiscovery         bool  // 设置节点是否曝露
    DiscoveryV5         bool  // V5版探测节点协议
    DiscoveryV5Addr     string
    Name                string     // 当前节点的服务名称
    BootstrapNodes      []*discover.Node // 启动节点, 可以选择默认的
    BootstrapNodesV5    []*discv5.Node   // V5 启动节点
    StaticNodes         []*discover.Node
    TrustedNodes        []*discover.Node
    NetRestrict         *netutil.Netlist  // 限制IP, 只有主节点才能匹配到这些IP
    NodeDatabase        string      // 节点数据库
    Protocols           []Protocol  // 
    ListenAddr          string      // server监听的地址, 为空的话，由系统指定实际的IP地址
    NAT                 nat.Interface // 本地的端口映射到对用的对外接入Internet的端口
    Dialer              *net.Dialer  // 连接对外的对等节点，返回的是个连接
    NoDial              bool    // true: server不会连接任意节点
}
```

P2P 协议相关 p2p/peer.go
```
// P2P 握手协议
type protoHandshake struct {
    Version     uint64      // 协议版本
	Name        string      // 协议名称
	Caps        []Cap       // 
	ListenPort  uint64      // 监听端口
	ID          discover.NodeID // 节点标识

	Rest []rlp.RawValue     // 兼容性数据
}

// 读写协议
type protoRW struct {
	Protocol
	in          chan Msg        // 接收读进消息
	closed      <-chan struct{} // 节点关闭
	wstart      <-chan struct{} // 开始写数据
	werr        chan<- error    // 执行写的结果
	offset      uint64
	w           MsgWriter       // 
}

// 
type Cap struct {
	Name    string
	Version uint
}

// 网络对等节点的连接信息
type Peer struct {
    rw          *conn       // 网络连接的信息 （两次握手后的信息)
    running     map[string]*protoRW  // 节点读写协议相关
    log         log.Logger      // 日志存储
    created     mclock.AbsTime
    wg          sync.WaitGroup
    protoErr    chan error
    closed      chan struct{}
    disc        chan DiscReason
}

// 基于net.Conn 协议扩展
type conn struct {
	fd          net.Conn
	transport
	flags       connFlag
	cont        chan error      // 
	id          discover.NodeID // 
	caps        []Cap           //
	name        string          // 
}
```


关于节点信息 p2p/discover/node.go
```
type Node struct {
    IP          net.IP // IPv4 IP 地址
    UDP         uint16  // UDP 端口
    TCP         uint16  // TCP 端口
    ID          NodeID   // 每个节点的唯一标识(64 byte)， 使用椭圆曲线算出来的公钥
    sha         common.Hash
    contested   bool 
}
```

以太坊默认的几个主节点 params/bootnodes.go <br>
```
var MainnetBootnodes = []string{
	// Ethereum Foundation Go Bootnodes
	"enode://a979fb575495b8d6db44f750317d0f4622bf4c2aa3365d6af7c284339968eef29b69ad0dce72a4d8db5ebb4968de0e3bec910127f134779fbcb0cb6d3331163c@52.16.188.185:30303", // IE
	"enode://3f1d12044546b76342d59d4a05532c14b85aa669704bfe1f864fe079415aa2c02d743e03218e57a33fb94523adb54032871a6c51b2cc5514cb7c7e35b3ed0a99@13.93.211.84:30303",  // US-WEST
	"enode://78de8a0916848093c73790ead81d1928bec737d565119932b98c6b100d944b7a95e94f847f689fc723399d2e31129d182f7ef3863f2b4c820abbf3ab2722344d@191.235.84.50:30303", // BR
	"enode://158f8aab45f6d19c6cbf4a089c2670541a8da11978a2f90dbf6a502a4a3bab80d288afdbeb7ec0ef6d92de563767f3b1ea9e8e334ca711e9f8e2df5a0385e8e6@13.75.154.138:30303", // AU
	"enode://1118980bf48b0a3640bdba04e0fe78b1add18e1cd99bf22d53daac1fd9972ad650df52176e7c7d89d1114cfef2bc23a2959aa54998a46afcf7d91809f0855082@52.74.57.123:30303",  // SG

	// Ethereum Foundation Cpp Bootnodes
	"enode://979b7fa28feeb35a4741660a16076f1943202cb72b6af70d327f053e248bab9ba81760f39d0701ef1d8f89cc1fbd2cacba0710a12cd5314d5e0c9021aa3637f9@5.1.83.226:30303", // DE
}
```
可以通过下面函数解析上面地址信息
```
discover.MustParseNode(node)
```
另外还有相关的测试节点 <br>
略

P2P 协议信息
```
type Protocol struct {
    Name    string  //协议名称
    Version string  // 版本号
    Length  uint64  // 协议中消息代码的个数
    Run func(peer *Peer, rw MsgReadWriter) error // 当Protocol与peer建立协议时，新启动线程执行
    NodeInfo func() interface{}
    PeerInfo func(id discover.NodeID) interface{} // 节点信息
}
```

P2P 通信层相关的信息 whisper/whisper.go  
```
type Whisper struct {
    protocol        p2p.Protocol   // 协议信息(描述，参数)
    filters         *Filters        // 订阅消息相关
    privateKeys     map[string]*ecdsa.PrivateKey  // 存储私钥相关 （[公钥]|[私钥]）
    symKeys         map[string][]byte   // 对称密钥存储
    keyMu           sync.RWMutex        // 秘钥存储锁
    poolMu          sync.RWMutex        // 同步消息与有效期池锁
    expirations     map[uint32]*set.SetNonTS  // 消息有效期池
    peerMu          sync.RWMutex    // 同步活跃节点的锁
    peers           map[*Peer]struct{}   // 当前活跃节点
    messageQueue    chan *Envelope  // 普通消息队列
    p2pMsgQueue     chan *Envelope  // 针对P2P通讯的消息队列
    quit            chan struct{}  // 优雅退出
    settings        syncmap.Map // 
    statsMu         sync.Mutex
    stats           Statistics  // 通讯节点统计: 消息个数，内存使用
    mailServer      MailServer  // 邮件服务
}
```

消息打包的参数信息 whisper/message.go
```
type MessageParams struct {
	TTL             uint32     // 允许消息alive的最长时间,默认50秒
	Src             *ecdsa.PrivateKey   // 私钥 (个人)
	Dst             *ecdsa.PublicKey    // 公约 (个人)
	KeySym          []byte              // 
	Topic           TopicType
	WorkTime        uint32
	PoW             float64     // 工作量证明, 默认 0.2
	Payload         []byte
	Padding         []byte
}
```
这里 TTL 可以认定为消息有效期，如果有效期过了，这会从网络中删除

文本包 （传输过程中的数据包)
```
type Envelope struct {
    Version         []byte  // 版本
	Expiry          uint32  // 有效期
	TTL             uint32
	Topic           TopicType
	AESNonce        []byte
	Data            []byte  // 数据
	EnvNonce        uint64  // Nonce 值

	pow             float64  // 工作量   
	hash            common.Hash 
}
```

+ 公钥与私钥 crypto/ecdsa <br>
使用go自带的库：椭圆曲线数字签名算法 <br>
ecdsa.PrivateKey <br>
ecdsa.PublicKey

```
// 消息选择
type Filter struct {
    Src             *ecdsa.PublicKey  // 发送消息者公钥
	KeyAsym         *ecdsa.PrivateKey // 私钥
	KeySym          []byte            // topic相关的密钥
	Topics          [][]byte          // filter topic
	PoW             float64           // 工作量证明
	AllowP2P        bool              // 
	SymKeyHash      common.Hash       // 对称密钥Hash, 

	Messages        map[common.Hash]*ReceivedMessage
	mutex           sync.RWMutex
}
```

MessageParams : 用于发送消息
Filter : 用于接收消息


#### P2P实现

+ p2p/server.go
    + Start() 

+ p2p/discover/udp.go <br>
    + ListenUDP(): 启动本地监听服务
    + newUDP(): 实现本地IP到外网IP的转化(环回地址)接收peer返回
    + makeEndpoint() : 创建RPC接口
    + loop():
    + readLoop(): 处理的接收的UDP数据包

+ whisper/whisper.go
    + add(): 发送消息的处理
        + 时间: 
        + wh.mailServer.Archive(envelope) : 数据包存档

+ whisper/doc.go
    + MailServer: 数据发送处理
