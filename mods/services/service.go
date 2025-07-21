package services

import (
	"os"
)

func (s *Service) Info() *Info {
	return &Info{
		ID:   s.ID,
		Host: s.Host,
		Port: s.Port,
		Type: s.Type,
	}
}

func (s *Service) Stop() error {
	s.user.mu.Lock()
	s.user.mgr.mu.Lock()
	defer s.user.mu.Unlock()
	defer s.user.mgr.mu.Unlock()

	return s.stop()
}

/**
 * PRIVATE METHODS BELOW.
 */
func (s *Service) stop() error {
	s.cancel()

	if s.Type == "unix" {
		s.listener.Close()

		if s.sock != "" {
			os.Remove(s.sock)
		}

		delete(s.user.unix, s.Port)
	}

	if s.Type == "tcp" {
		s.listener.Close()

		// Remove from user's list.
		delete(s.user.tcp, s.Port)

		// Remove from manager's list.
		delete(s.user.mgr.tcp, s.Port)
	}

	if s.Type == "udp" {
		s.udpConn.Close()

		// Remove from user's list.
		delete(s.user.udp, s.Port)

		// Remove from manager's list.
		delete(s.user.mgr.udp, s.Port)
	}

	return nil
}
