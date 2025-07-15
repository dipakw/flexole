package services

import (
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
	Group   string
	Host    string
	Port    uint16
	Type    string
	Timeout time.Duration
	SrcFN   func(*Info) (net.Conn, error)

	// Internal.
	key      string
	sock     string
	listener net.Listener
	udpConn  *net.UDPConn
	manager  *Services
}

type Config struct {
	Dir string
}

type Services struct {
	mutex sync.RWMutex
	conf  *Config
	list  map[string]*Service
}
