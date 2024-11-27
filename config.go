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
	Display string `yaml:"display"`
	Kernel  string `yaml:"kernel"`
	Initrd  string `yaml:"initrd"`
	Append  string `yaml:"append"`
}

type PXEConfig struct {
	DefaultEntry string  `yaml:"default"`
	Timeout      int64   `yaml:"timeout"`
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

func (c *PXEConfig) ScriptRender() (*bytes.Buffer, error) {
	const BootConfigTpl = `#!ipxe
{{define "entryTpl"}}:{{.Label}}
kernel {{.Kernel}} {{.Append}} || goto failed
initrd {{.Initrd}} || goto failed
boot || goto failed
goto start
{{end}}

set menu-timeout {{.Timeout}}
set submenu-timeout ${menu-timeout}
set protocol tftp
isset ${menu-default} || set menu-default {{.DefaultEntry}}

:start
menu iPXE Boot Menu -- ${buildarch}-${platform}
item --gap -- --------------------------------- Images -------------------------------
{{ with .Entries }}{{ range . }}item {{ .Label }} {{ .Display }}
{{ end }}{{ end }}
item --gap -- -------------------------------- Advanced ------------------------------
item --key c config [C] Configure settings
item --key s shell [S] Drop to iPXE Shell
item --key r reboot [R] Reboot the Computer
item --key x exit [X] Exit iPXE and Continue BIOS Booting

choose --timeout ${menu-timeout} --default ${menu-default} selected
goto ${selected}

{{ range $value := .Entries }}{{ template "entryTpl" $value }}
{{ end }}

:failed
echo Booting failed, dropping to shell
goto shell

:config
config
goto start

:shell
echo Type 'exit' to get the back to the menu
shell
set menu-timeout 0
goto start

:reboot
reboot

:exit
exit
`
	buf := &bytes.Buffer{}
	tpl := template.Must(template.New("BootConfigTpl").Parse(BootConfigTpl))
	if err := tpl.ExecuteTemplate(buf, "BootConfigTpl", c); err != nil {
		return nil, err
	}
	Debug("Parse 'ipxe.script':\n", buf.String())
	return buf, nil
}
