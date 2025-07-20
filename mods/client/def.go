package client

import (
	"context"
	"net"
	"sync"

	"github.com/xtaci/smux"
)

type Config struct {
	ID     []byte
	Key    []byte
	Server *Addr
}

type Client struct {
	conf         *Config
	mu           sync.RWMutex
	wg           sync.WaitGroup
	ctx          context.Context
	cancel       context.CancelFunc
	pipesList    map[string]*connPipe
	servicesList map[uint16]*Service

	Pipes    *Pipes
	Services *Services
}

type Addr struct {
	Net  string
	Addr string
}

type Pipe struct {
	ID      string
	Encrypt bool
}

type Pipes struct {
	c *Client
}

type Services struct {
	c *Client
}

type Local struct {
	Net  string
	Addr string
}

type Remote struct {
	ID    uint16   `json:"id"`
	Net   string   `json:"net"`
	Port  uint16   `json:"port"`
	Pipes []string `json:"pipes"`
}

type Service struct {
	Local  *Local
	Remote *Remote
}

type connPipe struct {
	id     string
	active bool
	conn   net.Conn
	sess   *smux.Session
	ctrl   *smux.Stream
}
