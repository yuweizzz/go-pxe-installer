package main

import (
	"embed"
	"io"
	"io/fs"
	"net"
	"path/filepath"
	"strings"
)

type TFTPHandler struct {
	Root      embed.FS
	TftpAddr  net.IP
	PXEConfig PXEConfig
}

func (h *TFTPHandler) Read(filename string, rf io.ReaderFrom) error {
	path := filename
	if filepath.IsAbs(path) {
		// use relative path to access file
		path = strings.Replace(path, "/", "", 1)
	}
	var reader any
	var err error
	switch path {
	case "pxelinux.cfg/default":
		reader, err = h.PXEConfig.ConfigReader()
		if err != nil {
			return err
		}
	case "message":
		reader, err = h.PXEConfig.MessageReader()
		if err != nil {
			return err
		}
	default:
		// enter root filesystem
		root, _ := fs.Sub(h.Root, "tftpboot")
		reader, err = root.Open(path)
		if err != nil {
			return err
		}
	}
	n, err := rf.ReadFrom(reader.(io.Reader))
	if err != nil {
		return err
	}
	Info("Raw filename: ", filename, ", overWrited filename: ", path, ", Sent ", n, " bytes")
	return nil
}
