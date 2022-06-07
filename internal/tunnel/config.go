package tunnel

var Version string = "v1.0.0"

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
