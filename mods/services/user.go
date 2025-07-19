package services

import (
	"context"
	"fmt"
)

func (u *User) Start(service *Service) (*Service, error) {
	service.user = u

	if !u.Available(service.Type, service.Port) {
		return nil, fmt.Errorf("port %d is not available", service.Port)
	}

	service.ctx, service.cancel = context.WithCancel(context.Background())

	if service.Type == "tcp" || service.Type == "unix" {
		return service, u.startTCPOrUnix(service)
	}

	if service.Type == "udp" {
		return service, u.startUDP(service)
	}

	return nil, fmt.Errorf("invalid service type: %s", service.Type)
}

func (u *User) Available(network string, port uint16) bool {
	if network == "unix" {
		u.mu.RLock()
		defer u.mu.RUnlock()
		service, ok := u.unix[port]
		return !ok || service == nil
	}

	if network == "tcp" {
		u.mgr.mu.RLock()
		defer u.mgr.mu.RUnlock()
		status, ok := u.mgr.tcp[port]
		return !ok || !status
	}

	if network == "udp" {
		u.mgr.mu.RLock()
		defer u.mgr.mu.RUnlock()
		status, ok := u.mgr.udp[port]
		return !ok || !status
	}

	return false
}

func (u *User) Stop(network string, port uint16) error {
	u.mu.Lock()
	u.mgr.mu.Lock()
	defer u.mu.Unlock()
	defer u.mgr.mu.Unlock()

	var service *Service
	var ok bool

	switch network {
	case "unix":
		service, ok = u.unix[port]
	case "tcp":
		service, ok = u.tcp[port]
	case "udp":
		service, ok = u.udp[port]
	}

	if !ok || service == nil {
		return nil
	}

	return service.stop()
}

func (u *User) Reset() []error {
	u.mu.Lock()
	u.mgr.mu.Lock()
	defer u.mu.Unlock()
	defer u.mgr.mu.Unlock()

	errs := []error{}

	types := []map[uint16]*Service{
		u.unix,
		u.tcp,
		u.udp,
	}

	for _, services := range types {
		for _, service := range services {
			if service == nil {
				continue
			}

			if err := service.stop(); err != nil {
				errs = append(errs, err)
			}
		}
	}

	return errs
}
