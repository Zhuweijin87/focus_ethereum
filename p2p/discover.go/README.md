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

	// 继承table
	*Table
}
```
