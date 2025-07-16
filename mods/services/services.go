package services

import (
	"path"
	"sync"
)

func Manager(c *Config) *ServicesManager {
	instance := &ServicesManager{
		mu:    sync.RWMutex{},
		conf:  c,
		users: map[string]*User{}, // List of users.
		tcp:   map[uint16]bool{},  // Used TCP ports.
		udp:   map[uint16]bool{},  // Used UDP ports.
	}

	return instance
}

func (s *ServicesManager) User(id string) *User {
	if id == "" {
		return nil
	}

	s.mu.RLock()
	user, ok := s.users[id]
	s.mu.RUnlock()

	if !ok || user == nil {
		s.mu.Lock()
		defer s.mu.Unlock()

		user = &User{
			id:   id,
			mgr:  s,
			dir:  path.Join(s.conf.Dir, id),
			udp:  map[uint16]*Service{},
			tcp:  map[uint16]*Service{},
			unix: map[uint16]*Service{},
		}

		s.users[id] = user
	}

	return user
}

func (s *ServicesManager) HasUser(id string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	user, ok := s.users[id]
	return ok && user != nil
}

func (s *ServicesManager) Reset() []error {
	errs := []error{}

	for _, user := range s.users {
		if user == nil {
			continue
		}

		errs = append(errs, user.Reset()...)
	}

	return errs
}
