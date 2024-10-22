package main

import (
	"io/ioutil"

	yaml "gopkg.in/yaml.v3"
)

type Config struct {
	Iface  string     `yaml:"iface"`
	IPAddr string     `yaml:"ipaddr"`
	DHCP   DHCPConfig `yaml:"dhcp"`
	TFTP   TFTPConfig `yaml:"tftp"`
}

type DHCPConfig struct {
	Port int `yaml:"port"`
}

type TFTPConfig struct {
	Port int `yaml:"port"`
}

func (c *Config) ParseConfig(filepath string) {
	file, err := ioutil.ReadFile(filepath)
	if err != nil {
		Error("failed to open config file: ", err)
	}
	err = yaml.Unmarshal(file, c)
	if err != nil {
		Error("failed to unmarshal yaml file: ", err)
	}
}
