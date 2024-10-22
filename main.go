package main

import (
	"os"
	"os/signal"
	"syscall"
)

var serverIPAddr = "10.0.2.5"

func main() {
	Initial("debug", os.Stdout)
	go Rundhcp("enp0s3", 67)
	go Runtftp()

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
