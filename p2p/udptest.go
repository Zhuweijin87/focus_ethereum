package main

import (
	"fmt"
	"net"
	"os"
)

type conn interface {
	ReadFromUDP(b []byte) (n int, addr *net.UDPAddr, err error)
	WriteToUDP(b []byte, addr *net.UDPAddr) (n int, err error)
	Close() error
	LocalAddr() net.Addr
}

type UDP struct {
	conn
}

func main() {
	addr, err := net.ResolveUDPAddr("udp", "127.0.0.1:8200")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(addr)

	lsn, err := net.ListenUDP("udp4", addr)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer lsn.Close()

	realAddr := lsn.LocalAddr().(*net.UDPAddr)
	fmt.Println(realAddr)

	fmt.Println(realAddr.IP.IsLoopback())

	addrs, _ := net.InterfaceAddrs()
	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				os.Stdout.WriteString(ipnet.IP.String() + "\n")
			}
		}
	}

	var udp UDP

	udp.LocalAddr()
}
