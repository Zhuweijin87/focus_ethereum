package main

import (
	"fmt"

	"github.com/ethereum/go-ethereum/node"
	"github.com/ethereum/go-ethereum/p2p"
	"github.com/ethereum/go-ethereum/rpc"
)

type MyService struct{}

func (ms *MyService) Protocols() []p2p.Protocol {return nil}
func (ms *MyService) APIs() []rpc.API { return nil }
func (ms *MyService) Start(*p2p.Server) error { fmt.Println("Service starting..."); return nil }
func (ms *MyService) Stop() error { fmt.Println("Service stopping"); return nil }

func main() {
	// 创建一个节点 
	stack, err := node.New(&node.Config{})
	if err != nil {
		fmt.Println("new node:", err)
		return 
	}

	constructor := func(context *node.ServiceContext) (node.Service, error) {
		return new(MyService), nil
	}

	// 注册一个服务
	if err := stack.Register(constructor); err != nil {
		fmt.Println("Fail to register: ", err)
		return 
	}

	if err := stack.Start(); err != nil {
		fmt.Println("fail to start: ", err)
		return 
	}

	if err := stack.Restart(); err != nil {
		fmt.Println("fail to restart: ", err)
		return 
	}

	if err := stack.Stop(); err != nil {
		fmt.Println("Fail to stop: ", err)
		return 
	}
}