package server

import (
	"fmt"
	"net"
	"sync"
)

func New(c *Config) (*Server, error) {
	if c.AuthFN == nil {
		return nil, fmt.Errorf("option 'AuthFN' is required")
	}

	if c.EncFN == nil {
		return nil, fmt.Errorf("option 'EncFN' is required")
	}

	instance := &Server{
		conf:  c,
		mu:    sync.RWMutex{},
		users: map[string]*User{},
	}

	return instance, nil
}

func (s *Server) User(id string) *User {
	s.mu.RLock()
	user, ok := s.users[id]
	s.mu.RUnlock()

	if !ok || user == nil {
		user = &User{
			id:       id,
			mu:       sync.RWMutex{},
			pipes:    map[string]*Pipe{},
			services: map[uint16]*Service{},
		}

		s.mu.Lock()
		s.users[id] = user
		s.mu.Unlock()
	}

	return user
}

func (s *Server) Start(background bool) error {
	var err error

	s.listener, err = net.Listen(s.conf.Net, s.conf.Addr)

	if err != nil {
		return nil
	}

	var run = func() error {
		for {
			conn, err := s.listener.Accept()

			if err != nil {
				fmt.Println("ERR:", err)
				continue
			}

			go s.handle(conn)
		}
	}

	if background {
		go run()
		return nil
	}

	return run()
}

func (s *Server) Stop() error {
	if s.listener != nil {
		return s.listener.Close()
	}

	return nil
}
