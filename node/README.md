## Node 节点

### node.go

```go
type Node struct {
	eventmux *event.TypeMux   // Event multiplexer used between the services of a stack
	config   *Config
	accman   *accounts.Manager

	ephemeralKeystore string         // if non-empty, the key directory that will be removed by Stop
	instanceDirLock   flock.Releaser // prevents concurrent use of instance directory

	serverConfig p2p.Config
	server       *p2p.Server    // p2p网络服务层

	serviceFuncs []ServiceConstructor     // Service constructors (in dependency order)
	services     map[reflect.Type]Service // 运行的服务

	rpcAPIs       []rpc.API   // 节点提供的API
	inprocHandler *rpc.Server // In-process RPC request handler to process the API requests

	ipcEndpoint string       // IPC 断点
	ipcListener net.Listener // IPC RPC listener socket to serve API requests
	ipcHandler  *rpc.Server  // IPC RPC request handler to process the API requests

	httpEndpoint  string       // HTTP 断点
	httpWhitelist []string     // HTTP RPC modules to allow through this endpoint
	httpListener  net.Listener // HTTP RPC listener socket to server API requests
	httpHandler   *rpc.Server  // HTTP RPC request handler to process the API requests

	wsEndpoint string       // WEbsocket 断点
	wsListener net.Listener // Websocket RPC listener socket to server API requests
	wsHandler  *rpc.Server  // Websocket RPC request handler to process the API requests

	stop chan struct{} // Channel to wait for termination notifications
	lock sync.RWMutex
}
```

