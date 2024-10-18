package main

import (
	"net"
	"os"

	"github.com/insomniacslk/dhcp/dhcpv4/server4"
)

func main() {
	Initial("debug", os.Stdout)
	laddr := net.UDPAddr{
		IP:   net.ParseIP("0.0.0.0"),
		Port: 67,
	}
	server, err := server4.NewServer("enp0s3", &laddr, DHCPHandler)
	if err != nil {
		Fatal(err)
	}

	// This never returns. If you want to do other stuff, dump it into a
	// goroutine.
	go server.Serve()
	Runtftp()
}
