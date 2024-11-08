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

const OPTIONS = `DISPLAY message
PROMPT 1
`

type Entry struct {
	Label  string
	Config string
	Kernel string
	Initrd string
	Append string
}

type BootConfig struct {
	DefaultEntry string
	Entries      map[string]Entry
	Options      string
}

type TFTPHandler struct {
	Root       embed.FS
	TftpAddr   net.IP
	BootConfig BootConfig
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
		help := Entry{
			Label:  "help",
			Config: "pxelinux.cfg/default",
		}
		debian12 := Entry{
			Label:  "1",
			Kernel: "images/debian-bookworm-amd64/linux",
			Initrd: "images/debian-bookworm-amd64/initrd.gz",
			Append: fmt.Sprintf("vga=normal fb=false auto=true priority=critical preseed/url=tftp://%s/images/debian-bookworm-amd64/preseed.cfg", h.TftpAddr.String()),
		}
		entries := make(map[string]Entry, 10)
		entries["help"] = help
		entries["1"] = debian12
		boot := &BootConfig{
			DefaultEntry: "help",
			Entries:      entries,
			Options:      OPTIONS,
		}
		var reader *bytes.Buffer
		if filename == "pxelinux.cfg/default" {
			reader = boot.String()
		}
		if filename == "message" {
			reader = boot.DisplayMessage()
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

func (b *BootConfig) String() *bytes.Buffer {
	const BootConfigTpl = `
{{define "entryTpl"}}LABEL {{.Label}}
    {{if .Config}}CONFIG {{.Config}}{{end}}
    {{if .Kernel}}KERNEL {{.Kernel}}{{end}}
    {{if .Initrd}}INITRD {{.Initrd}}{{end}}
    {{if .Append}}APPEND {{.Append}}{{end}}
{{end}}
DEFAULT {{.DefaultEntry}}
{{.Options}}

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

func (b *BootConfig) DisplayMessage() *bytes.Buffer {
	const MessageTpl = `Select the boot option and Press the corresponding number:
{{ range $key, $value := .Entries }}{{ $value.Label }}	{{ $value.Label }}
{{ end }}
`
	buf := &bytes.Buffer{}
	tpl := template.Must(template.New("MessageTpl").Parse(MessageTpl))
	if err := tpl.ExecuteTemplate(buf, "MessageTpl", b); err != nil {
		Error(err)
	}
	return buf
}
