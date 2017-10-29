### P2P 代码解析
-----------------

#### server.go  
p2p服务，整个p2p初始的结构  

首先是它的结构部分 
``` go
type Server struct {
    Config      //配置文件

    newTransport func(net.Conn) transport // 传输接口
	newPeerHook  func(*Peer)

	lock    sync.Mutex // 保护运行的节点
	running bool

	ntab         discoverTable  // 节点动态存储的表，用于快速搜索 其关联的结构discovery.Table
	listener     net.Listener   // 作为服务的监听接口
	ourHandshake *protoHandshake  // 握手协议
	lastLookup   time.Time  // 最后一次节点探测
	DiscV5       *discv5.Network  // discovery V5的网络

	// These are for Peers, PeerCount (and nothing else).
	peerOp     chan peerOpFunc  // peer 处理的函数
	peerOpDone chan struct{}

	quit          chan struct{}  // 服务退出通知
	addstatic     chan *discover.Node  // 添加节点通知
	removestatic  chan *discover.Node  // 删除节点通知
	posthandshake chan *conn
	addpeer       chan *conn   // 添加peer的通知
	delpeer       chan peerDrop
	loopWG        sync.WaitGroup //  listen Loop
}
```
配置信息(很重要) 
```go 
type Config struct {
    PrivateKey          *ecdsa.PrivateKey // 私钥信息
    MaxPeers            int   // 最大节点数
    MaxPendingPeers     int 
    NoDiscovery         bool  // 设置节点是否曝露
    DiscoveryV5         bool  // V5版探测节点协议
    DiscoveryV5Addr     string
    Name                string     // 当前节点的服务名称
    BootstrapNodes      []*discover.Node // 启动节点, 可以选择默认的
    BootstrapNodesV5    []*discv5.Node   // V5 启动节点
    StaticNodes         []*discover.Node  // 静态节点 
    TrustedNodes        []*discover.Node  // 信任节点
    NetRestrict         *netutil.Netlist  // 限制IP, 只有主节点才能匹配到这些IP
    NodeDatabase        string      // 节点数据库
    Protocols           []Protocol  // 
    ListenAddr          string      // server监听的地址, 为空的话，由系统指定实际的IP地址
    NAT                 nat.Interface // 本地的端口映射到对用的对外接入Internet的端口
    Dialer              *net.Dialer  // 连接对外的对等节点，返回的是个连接
    NoDial              bool    // true: server不会连接任意节点,也不会去监听端口
}
```

主要函数说明:  
* Start() error: 启动p2p服务  
这里会检测一侧配置参数，初始化一些信号通知，对整个网络的通知的处理  

* AddPeer(node *discover.Node) : 添加节点，主要通过RPC调用  
peer 将会添加到 dialstatic.static 存储中，然后在创建newTasks时，会检验这些节点的连接状态，如果连接没问题则连接保存起来，否则删除这个节点。  

* Peers() : 获取当前对等节点，通过RPC调用  

* RemovePeer() : 删除节点，通过RPC调用  

* Self() *discover.Node : 获取节点本身的信息  

* NodeInfo() : 节点信息, RPC调用  

! 内部函数  
* run(dialstate dialer)  不断的扫描节点，更新节点的状态（是否联通）  
	+ scheduleTasks: 任务调度器，将可能的节点添加到任务队列中   
	+ startTasks(): 启动并执行任务  
	+ delTask() : 删除任务  


``` go
// P2P 握手协议
type protoHandshake struct {
    Version     uint64      // 协议版本
	Name        string      // 协议名称
	Caps        []Cap       // 
	ListenPort  uint64      // 监听端口
	ID          discover.NodeID // 节点标识

	Rest []rlp.RawValue     // 兼容性数据
}

// 读写协议(对读写的封装)
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

// 基于net.Conn协议扩展
type conn struct {
	fd          net.Conn // 套接字
	transport
	flags       connFlag  // 链接标识
	cont        chan error      // 
	id          discover.NodeID // 节点ID
	caps        []Cap           //
	name        string          // 
}
```

```go
type Node struct {
    IP          net.IP // IPv4 IP 地址
    UDP         uint16  // UDP 端口
    TCP         uint16  // TCP 端口
    ID          NodeID   // 每个节点的唯一标识(64byte)， 使用椭圆曲线算出来的公钥
    sha         common.Hash
    contested   bool 
}
```

节点连接状态 p2p/dail.go
```go
type dialstate struct {
	maxDynDials     int            // 配置中的最大节点 + 1 / 2
	ntab            discoverTable
	netrestrict     *netutil.Netlist

	lookupRunning   bool
	dialing         map[discover.NodeID]connFlag
	lookupBuf       []*discover.Node    // 当前发现的结果
	randomNodes     []*discover.Node    // filled from Table
	static          map[discover.NodeID]*dialTask
	hist            *dialHistory

	start     time.Time        // time when the dialer was first used
	bootnodes []*discover.Node // default dials when there are no peers
}
```

``` go
discover.MustParseNode(node)
```
+ 默认的节点IP  
    + 182.254.155.208:33333  (不可PING)
    + 52.16.188.185:30303   (不可PING)
    + 13.93.211.84:30303  (不可PING)
    + 52.74.57.123:30303  (可PING)
    + 191.235.84.50:30303 (不可PING)
    + 5.1.83.226:30303  (可PING)

另外还有相关的测试节点 <br>
略

P2P 协议信息
```go
type Protocol struct {
    Name    string  //协议名称
    Version string  // 版本号
    Length  uint64  // 协议中消息代码的个数
    Run func(peer *Peer, rw MsgReadWriter) error // 当Protocol与peer建立协议时，新启动线程执行
    NodeInfo func() interface{}
    PeerInfo func(id discover.NodeID) interface{} // 节点信息
}
```

### message.go
消息读写  
```go
type MsgReadWriter interface {
	MsgReader
	MsgWriter
}

type MsgWriter interface {
    // 发送消息
	WriteMsg(Msg) error
}

type MsgReader interface {
	ReadMsg() (Msg, error)
}
```

### dail.go 
```go
type dialstate struct {
	maxDynDials int   // 动态链接的最大数
	ntab        discoverTable  // 节点探测表 
	netrestrict *netutil.Netlist

	lookupRunning bool  
	dialing       map[discover.NodeID]connFlag // 当前链接服务
	lookupBuf     []*discover.Node // 当前查找的结果
	randomNodes   []*discover.Node // 
	static        map[discover.NodeID]*dialTask
	hist          *dialHistory  // 链接历史揭露

	start     time.Time        // 
	bootnodes []*discover.Node // 如果没有要链接的节点时，默认的启动节点
}

// 探测节点存储的表操作
type discoverTable interface {
	Self() *discover.Node
	Close()
	Resolve(target discover.NodeID) *discover.Node
	Lookup(target discover.NodeID) []*discover.Node
	ReadRandomNodes([]*discover.Node) int
}

// 服务处理操作，p2p.Server启动后，调度任务是，会按照介个借口处理
type task interface {
	Do(*Server)
}

// 链接多的节点
type pastDial struct {
	id  discover.NodeID  // 节点ID
	exp time.Time  // 过期时间
}

type connFlag int  // 链接标识 

// 节点链接 
type dialTask struct {
	flags        connFlag  // 默认为 staticDialedConn
	dest         *discover.Node // 链接的节点
	lastResolved time.Time  // 
	resolveDelay time.Duration
}
```

### discover/table.go

维护节点的表结构
```go
type Table struct {
	mutex   sync.Mutex        // protects buckets, their content, and nursery
	buckets [nBuckets]*bucket // index of known nodes by distance
	nursery []*Node           // 托管节点
	db      *nodeDB           // 已知节点的数据存储

	refreshReq chan chan struct{}  // 刷新请求信号
	closeReq   chan struct{}  // 关闭请求信号
	closed     chan struct{}

	bondmu    sync.Mutex
	bonding   map[NodeID]*bondproc
	bondslots chan struct{} // limits total number of active bonding processes

	nodeAddedHook func(*Node) // for testing

	net  transport
	self *Node // 本地节点的信息元
}
```