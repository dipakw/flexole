package server

import (
	"flexole/mods/services"
	"net"
	"sync"

	"github.com/xtaci/smux"
)

type Config struct {
	Net     string
	Addr    string
	Manager *services.ServicesManager

	EncFN  func(a *Auth, c net.Conn) (uint8, []byte, error)
	AuthFN func(c net.Conn) (*Auth, error)
}

type Auth struct {
	UserID  string
	PipeID  string
	Encrypt bool
}

type Server struct {
	conf     *Config
	listener net.Listener
	mu       sync.RWMutex
	users    map[string]*User
}

type User struct {
	id       string
	mu       sync.RWMutex
	pipes    map[string]*Pipe
	services map[uint16]*Service
}

type Pipe struct {
	userID string
	id     string
	active bool
	conn   net.Conn
	sess   *smux.Session
	ctrl   *smux.Stream
}

// Shared
type Service struct {
	ID    uint16   `json:"id"`
	Net   string   `json:"net"`
	Port  uint16   `json:"port"`
	Pipes []string `json:"pipes"`
}

type NetPort struct {
	Net  string
	Port uint16
}
