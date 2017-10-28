package main

import (
	"fmt"
	"sync"
	"time"
)

type Server struct {
	lock sync.Mutex
	name string
}

func newServer() *Server {
	return &Server{
		name: "server",
	}
}

func (s *Server) Run() {
	s.lock.Lock()
	defer s.lock.Unlock()

	fmt.Printf("server:%s is running\n", s.name)
	time.Sleep(5 * time.Second)
}

func main() {
	serv := newServer()
	go serv.Run()
	go serv.Run()
	serv.Run()
}