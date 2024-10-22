package main

import (
	"embed"
	"fmt"
	"io"
	"io/fs"
	"path/filepath"
	"strings"
	"time"

	tftp "github.com/pin/tftp/v3"
)

//go:embed tftpboot
var tftpRoot embed.FS

func overWrite(filename string) string {
	Debug("Raw filename: ", filename)
	if filepath.IsAbs(filename) {
		filename = strings.Replace(filename, "/", "", 1)
	}
	Debug("overWrited filename: ", filename)
	return filename
}

// readHandler is called when client starts file download from server
func readHandler(filename string, rf io.ReaderFrom) error {
	// enter root filesystem
	root, _ := fs.Sub(tftpRoot, "tftpboot")
	// use relative path to access file
	filename = overWrite(filename)
	file, err := root.Open(filename)
	if err != nil {
		Error(err)
		return err
	}
	n, err := rf.ReadFrom(file)
	if err != nil {
		Error(err)
		return err
	}
	Info("sent ", n, " bytes")
	return nil
}

func Runtftp(port int) {
	// use nil in place of handler to disable read or write operations
	address := fmt.Sprintf(":%d", port)
	s := tftp.NewServer(readHandler, nil)
	s.SetTimeout(5 * time.Second)
	err := s.ListenAndServe(address)
	if err != nil {
		Panic("tftp server: ", err)
	}
}
