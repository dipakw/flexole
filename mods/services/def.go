package services

import (
	"context"
	"net"
	"sync"
	"time"
)

const (
	MAX_UDP_PACKET_SIZE = 65535
)

type Info struct {
	Host string
	Port uint16
	Type string
}

type Service struct {
	Host    string
	Port    uint16
	Type    string
	Timeout time.Duration
	SrcFN   func(*Info) (net.Conn, error)

	// Internal.
	sock     string
	listener net.Listener
	udpConn  *net.UDPConn
	user     *User
	ctx      context.Context
	cancel   context.CancelFunc
}

type Config struct {
	Dir string
}

type ServicesManager struct {
	mu    sync.RWMutex
	conf  *Config
	tcp   map[uint16]bool
	udp   map[uint16]bool
	users map[string]*User
}

type User struct {
	id  string
	mu  sync.RWMutex
	mgr *ServicesManager
	dir string

	// UDP services
	udp map[uint16]*Service

	// TCP services
	tcp map[uint16]*Service

	// Unix services
	unix map[uint16]*Service
}
