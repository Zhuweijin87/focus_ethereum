package main

import (
	"fmt"
	"net"
	"context"
	"encoding/json"
	"github.com/ethereum/go-ethereum/rpc"
)

type Service struct{}

type Args struct {
	S string
}

type Result struct {
	String string
	Int    int
	Args   *Args
}

func (s *Service) Echo(str string, i int, args *Args) Result {
	return Result{str, i, args}
}

func (s *Service) EchoWithCtx(ctx context.Context, str string, i int, args *Args) Result {
	return Result{str, i, args}
}

func main() {
	serv := rpc.NewServer()
	service := new(Service)
	
	err := serv.RegisterName("calc", service)
	if err != nil {
		fmt.Println("fail to register name")
		return 
	}
/*
	svc, ok := serv.services["calc"]
	if !ok {
		fmt.Println("Expected service calc to be registered")
	}
	*/

	stringArg := "string arg"
	intArg := 1122
	argsArg := &Args{"abcde"}
	params := []interface{}{stringArg, intArg, argsArg}

	request := map[string]interface{}{
		"id":      12345,
		"method":  "test_" + method,
		"version": "2.0",
		"params":  params,
	}

	clientConn, serverConn := net.Pipe()
	defer clientConn.Close()

	go server.ServeCodec(rpc.NewJSONCodec(serverConn), OptionMethodInvocation)

	out := json.NewEncoder(clientConn)
	in := json.NewDecoder(clientConn)

	if err := out.Encode(request); err != nil {
		t.Fatal(err)
	}
}