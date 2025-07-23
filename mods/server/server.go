package server

import (
	"context"
	"fmt"
	"net"
	"sync"
)

func New(c *Config) (*Server, error) {
	instance := &Server{
		conf:  c,
		mu:    sync.RWMutex{},
		users: map[string]*User{},
		wg:    sync.WaitGroup{},
	}

	instance.ctx, instance.cancel = context.WithCancel(context.Background())

	return instance, nil
}

func (s *Server) User(id string) *User {
	s.mu.RLock()
	user, ok := s.users[id]
	s.mu.RUnlock()

	if !ok || user == nil {
		user = &User{
			id:           id,
			mu:           sync.RWMutex{},
			pipesList:    map[string]*Pipe{},
			servicesList: map[uint16]*Service{},
		}

		user.pipes = &Pipes{
			user:   user,
			server: s,
		}

		user.services = &Services{
			user:   user,
			server: s,
		}

		s.mu.Lock()
		s.users[id] = user
		s.mu.Unlock()
	}

	return user
}

func (s *Server) Start() error {
	var err error

	s.listener, err = net.Listen(s.conf.Net, s.conf.Addr)

	if err != nil {
		return err
	}

	s.wg.Add(1)

	go func() {
		defer s.wg.Done()
		defer s.listener.Close()

		for {
			select {
			case <-s.ctx.Done():
				return
			default:
				conn, err := s.listener.Accept()

				if err != nil {
					fmt.Println("ERR:", err)
					continue
				}

				go s.handle(s.ctx, conn)
			}
		}
	}()

	return nil
}

func (s *Server) Addr() string {
	return s.conf.Addr
}

func (s *Server) Net() string {
	return s.conf.Net
}

func (s *Server) Wait() {
	s.wg.Wait()
}

func (s *Server) Stop() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.cancel()
	return nil
}
