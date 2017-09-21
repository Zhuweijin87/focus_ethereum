// NAT-PMP 测试
// go get github.com/jackpal/gateway
// go get github.com/jackpal/go-nat-pmp
package main

import (
	"fmt"
    "github.com/jackpal/gateway"
    natpmp "github.com/jackpal/go-nat-pmp"
)

func main() {
	gatewayIP, err := gateway.DiscoverGateway()
	if err != nil {
		fmt.Println(err)
    	return
	}

	client := natpmp.NewClient(gatewayIP)
	response, err := client.GetExternalAddress()
	if err != nil {
		fmt.Println(err)
    	return
	}
	fmt.Println("External IP address:", response.ExternalIPAddress)
}