package app

const (
	DEFAULT_PORT = "8887"
)

type Addr struct {
	Net  string `yaml:"net"`
	Addr string `yaml:"addr"`
}

/**
 * Server config.
 */
type ServerConfig struct {
	Version string `yaml:"version"`
	Config  Addr   `yaml:"config"`
	Logs    Logs   `yaml:"logs"`
	Users   []User `yaml:"users"`
}

type Logs struct {
	Allow []string `yaml:"allow"`
	Outs  []LogOut `yaml:"outs"`
}

type LogOut struct {
	To    string `yaml:"to"`
	Color bool   `yaml:"color,omitempty"` // only relevant for stdout
	Path  string `yaml:"path,omitempty"`  // only relevant for file
}

type User struct {
	ID          string      `yaml:"id"`
	Enabled     bool        `yaml:"enabled"`
	Key         string      `yaml:"key"`
	MaxPipes    int         `yaml:"max_pipes"`
	MaxServices MaxServices `yaml:"max_services"`
}

type MaxServices struct {
	Unix int `yaml:"unix"`
	TCP  int `yaml:"tcp"`
	UDP  int `yaml:"udp"`
}

/**
 * Client config.
 */
type ClientConfig struct {
	Version  string    `yaml:"version"`
	Auth     Auth      `yaml:"auth"`
	Server   Addr      `yaml:"server"`
	Pipes    []Pipe    `yaml:"pipes"`
	Services []Service `yaml:"services"`
}

type Auth struct {
	ID  string `yaml:"id"`
	Key string `yaml:"key"`
}

type Pipe struct {
	ID      string `yaml:"id"`
	Enabled bool   `yaml:"enabled"`
	Encrypt bool   `yaml:"encrypt"`
}

type Service struct {
	ID      uint16     `yaml:"id"`
	Enabled bool       `yaml:"enabled"`
	Local   Addr       `yaml:"local"`
	Remote  RemoteConf `yaml:"remote"`
}

type RemoteConf struct {
	Net   string   `yaml:"net"`
	Port  int      `yaml:"port"`
	Pipes []string `yaml:"pipes"`
}
