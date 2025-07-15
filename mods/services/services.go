package services

import (
	"errors"
	"fmt"
	"sync"
)

func Manager(c *Config) *Services {
	instance := &Services{
		mutex: sync.RWMutex{},
		conf:  c,
		list:  map[string]*Service{},
	}

	return instance
}

func (s *Services) Start(service *Service) (*Service, error) {
	service.key = fmt.Sprintf("%d", service.Port)

	if service.Type == "tcp" || service.Type == "unix" {
		return service, s.startTCPOrUnix(service, s.conf.Dir)
	}

	if service.Type == "udp" {
		return service, s.startUDP(service)
	}

	return nil, fmt.Errorf("invalid service type: %s", service.Type)
}

func (s *Services) Stop(key string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	service, ok := s.list[key]

	if !ok || service == nil {
		return errors.New("service not found")
	}

	return service.Stop()
}

func (s *Services) Has(key string) bool {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	if service, ok := s.list[key]; !ok || service == nil {
		return false
	}

	return true
}

func (s *Services) Reset() []error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	var errors = []error{}

	for _, service := range s.list {
		if err := service.Stop(); err != nil {
			errors = append(errors, err)
		}
	}

	return errors
}
