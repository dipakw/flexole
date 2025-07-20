package server

import (
	"context"
	"flexole/mods/services"
	"net"
	"sync"

	"github.com/dipakw/logs"
	"github.com/xtaci/smux"
)

type Config struct {
	Net     string
	Addr    string
	Manager *services.ServicesManager
	Log     logs.Log
	KeyFN   func(id []byte) ([]byte, error)
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
	ctx      context.Context
	cancel   context.CancelFunc
	wg       sync.WaitGroup
}

type User struct {
	id           string
	mu           sync.RWMutex
	pipesList    map[string]*Pipe
	servicesList map[uint16]*Service
	pipes        *Pipes
	services     *Services
}

type Services struct {
	user   *User
	server *Server
}

type Pipes struct {
	user   *User
	server *Server
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
