package main

import (
	"net"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// config
	Conf := &Config{}
	Conf.ParseConfig("config.yaml")
	// logger
	Initial(Conf.Logger.Level, InitialFD(Conf.Logger.File))
	// dhcp
	ip := net.ParseIP(Conf.IPAddr)
	dhcpServer := &DHCPServer{
		Handler: &DHCPHandler{
			DHCPAddr: ip,
			TFTPAddr: ip,
		},
		Iface: Conf.Iface,
		Port:  Conf.DHCP.Port,
	}
	go dhcpServer.Run()
	// tftp
	go Runtftp(Conf.TFTP.Port)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	done := make(chan bool, 1)
	go func() {
		<-sigs
		done <- true
	}()
	Info("Awaiting signal ......")
	<-done
	Info("Except signal, exiting ......")
}
