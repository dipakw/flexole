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
	service.manager = s

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
	return s.stop(key)
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

	for key := range s.list {
		if err := s.stop(key); err != nil {
			errors = append(errors, err)
		}
	}

	return errors
}

/**
 * PRIVATE METHODS BELOW.
 */
func (s *Services) stop(key string) error {
	service, ok := s.list[key]

	if !ok || service == nil {
		return errors.New("service not found")
	}

	if err := service.stop(); err != nil {
		return err
	}

	delete(s.list, key)

	return nil
}
