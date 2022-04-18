package config

import "sync"

var Version string = "v0.0.1"

var Wg sync.WaitGroup

type SSHServerConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	KeyFile  string `yaml:"keyFile"`
}

type RemoteConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type LocalConfig struct {
	Bind string `yaml:"bind"`
	Port int    `yaml:"port"`
}

type TunnelConfig struct {
	Name   string       `yaml:"name"`
	Remote RemoteConfig `yaml:"remote"`
	Local  LocalConfig  `yaml:"local"`
}

type CulvertConfig struct {
	Server  SSHServerConfig `yaml:"server"`
	Tunnels []TunnelConfig  `yaml:"tunnels"`
}
