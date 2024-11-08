package main

import (
	"bytes"
	"embed"
	"fmt"
	"io"
	"io/fs"
	"net"
	"path/filepath"
	"strings"
	"text/template"
)

type TFTPHandler struct {
	Root      embed.FS
	TftpAddr  net.IP
	PXEConfig PXEConfig
}

func (h *TFTPHandler) Read(filename string, rf io.ReaderFrom) error {
	// enter root filesystem
	root, _ := fs.Sub(h.Root, "tftpboot")
	// use relative path to access file
	Debug("Raw filename: ", filename)
	if filepath.IsAbs(filename) {
		filename = strings.Replace(filename, "/", "", 1)
	}
	Debug("overWrited filename: ", filename)
	if filename == "pxelinux.cfg/default" || filename == "message" {
		var reader *bytes.Buffer
		if filename == "pxelinux.cfg/default" {
			reader = h.PXEConfig.String()
		}
		if filename == "message" {
			reader = h.PXEConfig.DisplayMessage()
		}
		n, err := rf.ReadFrom(reader)
		if err != nil {
			Error(err)
			return err
		}
		Info("sent ", n, " bytes")
	} else {
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
	}
	return nil
}

func (b *PXEConfig) String() *bytes.Buffer {
	const BootConfigTpl = `
{{define "entryTpl"}}LABEL {{.Label}}
    {{if .Config}}CONFIG {{.Config}}{{end}}
    {{if .Kernel}}KERNEL {{.Kernel}}{{end}}
    {{if .Initrd}}INITRD {{.Initrd}}{{end}}
    {{if .Append}}APPEND {{.Append}}{{end}}
{{end}}
DEFAULT {{.DefaultEntry}}
DISPLAY message
PROMPT 1

{{ range $value := .Entries }}{{ template "entryTpl" $value }}{{ end }}
`
	buf := &bytes.Buffer{}
	tpl := template.Must(template.New("BootConfigTpl").Parse(BootConfigTpl))
	if err := tpl.ExecuteTemplate(buf, "BootConfigTpl", b); err != nil {
		Error(err)
	}
	fmt.Println(buf.String())
	return buf
}

func (b *PXEConfig) DisplayMessage() *bytes.Buffer {
	const MessageTpl = `Select the boot option and Press the corresponding number:
{{ range $key, $value := .Entries }}{{ $value.Label }}	{{ $value.Display }}
{{ end }}
`
	buf := &bytes.Buffer{}
	tpl := template.Must(template.New("MessageTpl").Parse(MessageTpl))
	if err := tpl.ExecuteTemplate(buf, "MessageTpl", b); err != nil {
		Error(err)
	}
	return buf
}
