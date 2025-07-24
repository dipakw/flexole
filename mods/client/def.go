package client

import (
	"context"
	"flexole/mods/cmd"
	"net"
	"sync"

	"github.com/dipakw/logs"
	"github.com/xtaci/smux"
)

const (
	MAX_UDP_PACKET_SIZE = 65535
)

var MESSAGES = map[uint8]string{
	cmd.CMD_STATUS_UNKNOWN: "unknown status",
	cmd.CMD_INVALID_CMD:    "invalid command",
	cmd.CMD_MALFORMED_DATA: "malformed data",
	cmd.CMD_PORT_UNAVAIL:   "port not available",
	cmd.CMD_OP_FAILED:      "operation failed",
	cmd.CMD_SERVICES_LIMIT: "services limit reached",
	cmd.CMD_PIPES_LIMIT:    "pipes limit reached",
	cmd.CMD_NOT_AVAILABLE:  "not available",
}

type Config struct {
	ID     []byte
	Key    []byte
	Server *Addr
	Log    logs.Log
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
