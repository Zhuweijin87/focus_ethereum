## 探测节点模块

发现与管理节点 

### udp.go 
```go
// 待处理的回复
type pending struct {
	from  NodeID
	ptype byte

	deadline time.Time  // 请求完成的时间

    // 接受到匹配的回复时，调用的回调函数，当返回为true 删除待处理回复队列的数据
	callback func(resp interface{}) (done bool) 

	errc chan<- error
}

// 回复信息
type reply struct {
	from   NodeID
	ptype  byte
	data   interface{}  // 数据
	matched chan<- bool // 是否匹配
}
```

```go
// udp协议处理
type udp struct {
	conn        conn
	netrestrict *netutil.Netlist
	priv        *ecdsa.PrivateKey  /// 节点私钥
	ourEndpoint rpcEndpoint  // rpc端点信息

	addpending chan *pending  // 添加到待处理回复的队列
	gotreply   chan reply  // 当接受回复时，处理回复

	closing    chan struct{}  // 关闭信号
	nat     nat.Interface

	*Table   // 节点存储表
}
```

### table.go 

存放节点的表
```go
type Table struct {
	mutex   sync.Mutex        // protects buckets, their content, and nursery
	buckets [nBuckets]*bucket // index of known nodes by distance
	nursery []*Node           // bootstrap nodes
	db      *nodeDB           // database of known nodes

	refreshReq chan chan struct{}  // 
	closeReq   chan struct{}
	closed     chan struct{}

	bondmu    sync.Mutex
	bonding   map[NodeID]*bondproc
	bondslots chan struct{} // limits total number of active bonding processes

	nodeAddedHook func(*Node) // for testing

	net  transport
	self *Node // metadata of the local node
}
```

主要接口实现：  
+ Self() *discover.Node  自身节点  

+ Close()  

+ Resolve(target discover.NodeID) *discover.Node  

+ Lookup(target discover.NodeID) []*discover.Node  搜寻节点  

+ ReadRandomNodes([]*discover.Node) int  读取随机节点  


### node.go 
节点信息  

+ NewnNode() : 顾名思义，创建一个节点  

+ Incomplete() : 无IP的节点 

+ String() 