package services

import (
	"os"
)

func (s *Service) Info() *Info {
	return &Info{
		Host: s.Host,
		Port: s.Port,
		Type: s.Type,
	}
}

func (s *Service) Stop() error {
	s.user.mu.Lock()
	defer s.user.mu.Unlock()
	return s.stop()
}

/**
 * PRIVATE METHODS BELOW.
 */
func (s *Service) stop() error {
	s.cancel()

	if s.Type == "unix" {
		if s.sock != "" {
			return os.Remove(s.sock)
		}

		delete(s.user.unix, s.Port)
	}

	if s.Type == "tcp" {
		// Remove from user's list.
		delete(s.user.tcp, s.Port)

		// Remove from manager's list.
		delete(s.user.mgr.tcp, s.Port)
	}

	if s.Type == "udp" {
		// Remove from user's list.
		delete(s.user.udp, s.Port)

		// Remove from manager's list.
		delete(s.user.mgr.udp, s.Port)
	}

	return nil
}
