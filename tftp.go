package main

import (
	"fmt"
	"time"

	tftp "github.com/pin/tftp/v3"
)

type TFTPServer struct {
	Handler *TFTPHandler
	Port    int
}

func (s *TFTPServer) Run() {
	address := fmt.Sprintf(":%d", s.Port)
	server := tftp.NewServer(s.Handler.Read, nil)
	server.SetTimeout(5 * time.Second)
	err := server.ListenAndServe(address)
	if err != nil {
		Panic("tftp server: ", err)
	}
}
