package main

import (
	"fmt"
	"os"
	"time"

	tftp "github.com/pin/tftp/v3"
)

type TFTPServer struct {
	Handler  *TFTPHandler
	Port     int
	External string
}

func (s *TFTPServer) Run() {
	address := fmt.Sprintf(":%d", s.Port)
	if len(s.External) > 0 {
		s.Handler.ExternalRoot = os.DirFS(s.External)
	}
	server := tftp.NewServer(s.Handler.Read, nil)
	server.SetTimeout(5 * time.Second)
	err := server.ListenAndServe(address)
	if err != nil {
		Panic("tftp server: ", err)
	}
}
