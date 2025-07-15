package services

import (
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

func (s *Services) Start(service *Service) error {
	service.key = fmt.Sprintf("%d", service.Port)

	if service.Type == "tcp" || service.Type == "unix" {
		return s.startTCPOrUnix(service, s.conf.Dir)
	}

	if service.Type == "udp" {
		return s.startUDP(service)
	}

	return fmt.Errorf("invalid service type: %s", service.Type)
}
