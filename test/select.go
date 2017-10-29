package main

import (
	"fmt"
	"time"
)

type Server struct {
	peer    []string
	addpeer chan string
	quit    chan struct{}
}

func newServer() *Server {
	return &Server{
		addpeer: make(chan string),
		quit:    make(chan struct{}),
	}
}

func (s *Server) addPeer(add string) {
	select {
	case s.addpeer <- add:
	case <-s.quit:
	}
}

func (s *Server) handle() {
	for i := 0; i < 10; i++ {
		s.addPeer("Hello")
		time.Sleep(1 * time.Second)
	}
	s.quit <- struct{}{}
}

func (s *Server) stop() {
	fmt.Println("server exit")
	close(s.addpeer)
	close(s.quit)
}

func main() {
	server := newServer()
	defer server.stop()
	go server.handle()


loop:
	for {
		select {
		case newadd := <-server.addpeer:
			fmt.Println("new add :", newadd)
		case <-server.quit:
			fmt.Println("quit handle select")
			break loop
		}
	}
}
