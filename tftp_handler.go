package main

import (
	"embed"
	"io"
	"io/fs"
	"net"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"

	tftp "github.com/pin/tftp/v3"
)

const DefaultPxeCfgPath = "pxelinux.cfg/default"
const DefaultPxePrompt = "prompt"

type TFTPHandler struct {
	EmbedRoot    embed.FS
	ExternalRoot fs.FS
	TftpAddr     net.IP
	PXEConfig    PXEConfig
}

func HttpReader(path string) (io.ReadCloser, int64, error) {
	resp, err := http.Get(path)
	if err != nil {
		return nil, 0, err
	}
	return resp.Body, resp.ContentLength, nil
}

func (h *TFTPHandler) PatchfilePath(path string) string {
	for _, entry := range h.PXEConfig.Entries {
		if path == entry.Kernel || path == entry.Initrd {
			path = entry.Prefix + path
			return path
		}
	}
	if filepath.IsAbs(path) {
		// use relative path to access file
		path = strings.Replace(path, "/", "", 1)
	}
	return path
}

func (h *TFTPHandler) Read(filename string, rf io.ReaderFrom) error {
	path := h.PatchfilePath(filename)
	u, err := url.Parse(path)
	if err != nil {
		return err
	}
	var reader any
	// http
	if u.Scheme == "http" || u.Scheme == "https" {
		var size int64
		reader, size, err = HttpReader(path)
		if err != nil {
			return err
		}
		rf.(tftp.OutgoingTransfer).SetSize(size)
	} else if h.ExternalRoot != nil {
		// tftp external
		reader, err = h.ExternalRoot.Open(path)
		if err != nil {
			Error("lookup external fs failed, will fallback to embed fs, err: ", err)
		}
	}
	// fallback to tftp embed
	if reader == nil {
		switch path {
		case DefaultPxeCfgPath:
			reader, err = h.PXEConfig.ConfigReader()
		case DefaultPxePrompt:
			reader, err = h.PXEConfig.PromptReader()
		default:
			// enter root filesystem
			root, _ := fs.Sub(h.EmbedRoot, "tftpboot")
			reader, err = root.Open(path)
		}
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
