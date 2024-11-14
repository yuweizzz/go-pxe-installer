package main

import (
	"bytes"
	"io/ioutil"
	"text/template"

	yaml "gopkg.in/yaml.v3"
)

type Config struct {
	Iface  string     `yaml:"iface"`
	IPAddr string     `yaml:"ipaddr"`
	Logger Logger     `yaml:"logger"`
	DHCP   DHCPConfig `yaml:"dhcp"`
	TFTP   TFTPConfig `yaml:"tftp"`
	PXE    PXEConfig  `yaml:"pxe"`
}

type Logger struct {
	Level string `yaml:"level"`
	File  string `yaml:"file"`
}

type DHCPConfig struct {
	Port int `yaml:"port"`
}

type TFTPConfig struct {
	Port     int    `yaml:"port"`
	External string `yaml:"external"`
}

type Entry struct {
	Label   string `yaml:"label"`
	Config  string `yaml:"config"`
	Display string `yaml:"display"`
	Prefix  string `yaml:"prefix`
	Kernel  string `yaml:"kernel"`
	Initrd  string `yaml:"initrd"`
	Append  string `yaml:"append"`
}

type PXEConfig struct {
	DefaultEntry string  `yaml:"default"`
	Entries      []Entry `yaml:"entries"`
}

func (c *Config) ParseConfig(filepath string) {
	file, err := ioutil.ReadFile(filepath)
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(file, c)
	if err != nil {
		panic(err)
	}
}

func (c *PXEConfig) ConfigReader() (*bytes.Buffer, error) {
	const BootConfigTpl = `
{{define "entryTpl"}}LABEL {{.Label}}
    {{if .Config}}CONFIG {{.Config}}{{end}}
    {{if .Kernel}}KERNEL {{.Kernel}}{{end}}
    {{if .Initrd}}INITRD {{.Initrd}}{{end}}
    {{if .Append}}APPEND {{.Append}}{{end}}
{{end}}
DEFAULT {{.DefaultEntry}}
DISPLAY prompt
PROMPT 1

{{ range $value := .Entries }}{{ template "entryTpl" $value }}{{ end }}
`
	buf := &bytes.Buffer{}
	tpl := template.Must(template.New("BootConfigTpl").Parse(BootConfigTpl))
	if err := tpl.ExecuteTemplate(buf, "BootConfigTpl", c); err != nil {
		return nil, err
	}
	Debug("Parse 'pxelinux.cfg/default':\n", buf.String())
	return buf, nil
}

func (c *PXEConfig) PromptReader() (*bytes.Buffer, error) {
	const PromptTpl = `Select the boot option and Press the corresponding number:
{{ range $key, $value := .Entries }}{{ $value.Label }}	{{ $value.Display }}
{{ end }}
`
	buf := &bytes.Buffer{}
	tpl := template.Must(template.New("PromptTpl").Parse(PromptTpl))
	if err := tpl.ExecuteTemplate(buf, "PromptTpl", c); err != nil {
		return nil, err
	}
	Debug("Parse 'prompt':\n", buf.String())
	return buf, nil
}
