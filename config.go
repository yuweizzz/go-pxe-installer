package main

import (
	"io/ioutil"

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
	Port int `yaml:"port"`
}

type Entry struct {
	Label   string `yaml:"label"`
	Config  string `yaml:"config"`
	Display string `yaml:"display"`
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
