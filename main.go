package main

import (
	"os"
	"os/signal"
	"syscall"
)

var serverIPAddr string

func main() {
	Initial("debug", os.Stdout)
	Conf := &Config{}
	Conf.ParseConfig("config.yaml")
	serverIPAddr = Conf.IPAddr
	go Rundhcp(Conf.Iface, Conf.DHCP.Port)
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
