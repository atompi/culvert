package tunnel

var Version string = "v1.0.0 for specific tunnel only"

var ConfigYaml = `---
tunnels:
  - name: specific_tunnel
    mode: L
    host:
      ip: 123.123.123.123
      port: 2222
      username: tunneluser
      password: "123456"
      keyFile: "./id_ed25519"
      keyPassword: "123456"
      knownHost: "[123.123.123.123]:2222 ecdsa-sha2-nistp256 AAAAE2VjZHNhLXNoYTItbmlzdHAyNTYAAAAIbmlzdHAyNTYAAABBBOUDn8tF9i1XwSnYYKnoyR9z4g+pgdMR16vFFVH1UpskxgpAjgjBubdqTmIs1JQ8OJyWBomqandNM2WtIgQqAPc="
    keepalive:
      interval: 30
      countMax: 2
    remote:
      bind: 192.168.15.128
      port: 3306
    local:
      bind: 0.0.0.0
      port: 12333
    retryInterval: 5

log:
  path: "./culvert.log"
  level: "INFO"
`

var KeyStr string = `-----BEGIN OPENSSH PRIVATE KEY-----
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=
-----END OPENSSH PRIVATE KEY-----
`

type HostConfig struct {
	IP          string `yaml:"ip"`
	Port        int    `yaml:"port"`
	Username    string `yaml:"username"`
	Password    string `yaml:"password"`
	KeyFile     string `yaml:"keyFile"`
	KeyPassword string `yaml:"keyPassword"`
	KnownHost   string `yaml:"knownHost"`
}

type KeepaliveConfig struct {
	Interval int `yaml:"interval"`
	CountMax int `yaml:"countMax"`
}

type RemoteConfig struct {
	Bind string `yaml:"bind"`
	Port int    `yaml:"port"`
}

type LocalConfig struct {
	Bind string `yaml:"bind"`
	Port int    `yaml:"port"`
}

type TunnelConfig struct {
	Name          string          `yaml:"name"`
	Mode          string          `yaml:"mode"`
	Host          HostConfig      `yaml:"host"`
	Keepalive     KeepaliveConfig `yaml:"keepalive"`
	Remote        RemoteConfig    `yaml:"remote"`
	Local         LocalConfig     `yaml:"local"`
	RetryInterval int             `yaml:"retryInterval"`
}

type LogConfig struct {
	Path  string `yaml:"path"`
	Level string `yaml:"level"`
}

type CulvertConfig struct {
	Tunnels []TunnelConfig `yaml:"tunnels"`
	Log     LogConfig      `yaml:"log"`
}
