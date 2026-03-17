package serve

import (
	"errors"
	"net"
	"os"
)

type Config struct {
	Dir string
}

func New(cfg *Config) (*Serve, error) {
	if cfg.Dir == "" {
		return nil, errors.New("dir is required")
	}

	if info, err := os.Stat(cfg.Dir); os.IsNotExist(err) {
		return nil, errors.New("dir does not exist")
	} else if !info.IsDir() {
		return nil, errors.New("dir is not a directory")
	}

	instance := &Serve{
		cfg: cfg,
		fs:  os.DirFS(cfg.Dir),
	}

	instance.vs = instance.newVirtualServer()

	return instance, nil
}

func (s *Serve) Handle(conn net.Conn) {
	s.vs.serve(conn)
}
