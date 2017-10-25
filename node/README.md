## Node 节点

### node.go

```go
type Node struct {
	eventmux *event.TypeMux   // Event multiplexer used between the services of a stack
	config   *Config
	accman   *accounts.Manager

	ephemeralKeystore string     // if non-empty, the key directory that will be removed by Stop
	instanceDirLock   flock.Releaser // prevents concurrent use of instance directory

	serverConfig p2p.Config
	server       *p2p.Server    // p2p网络服务层

	serviceFuncs []ServiceConstructor     // 服务接口
	services     map[reflect.Type]Service // 运行的服务

	rpcAPIs       []rpc.API   // 节点提供的API
	inprocHandler *rpc.Server // In-process RPC request handler to process the API requests

	ipcEndpoint string       // IPC 断点
	ipcListener net.Listener 
	ipcHandler  *rpc.Server 

	httpEndpoint  string       // HTTP 断点
	httpWhitelist []string     // HTTP 白名单
	httpListener  net.Listener 
	httpHandler   *rpc.Server 

	wsEndpoint string       // WEbsocket 断点
	wsListener net.Listener 
	wsHandler  *rpc.Server 

	stop chan struct{} // 结束信号
	lock sync.RWMutex
}
```

+ Start(): 节点启动  
	+ 加载P2P节点信息  
	+ 启动注册的服务，如果有失败，就返回失败 
	+ 启动RPC进行交互处理

+ Stop(): 节点停止  

+ startRPC(): 启动RPC服务

