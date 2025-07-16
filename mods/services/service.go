package services

import "os"

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
	if err := s.stopListener(); err != nil {
		return err
	}

	if s.Type == "unix" {
		s.user.mu.Lock()
		delete(s.user.unix, s.Port)
		s.user.mu.Unlock()
	}

	if s.Type == "tcp" {
		// Remove from user's list.
		s.user.mu.Lock()
		delete(s.user.tcp, s.Port)
		s.user.mu.Unlock()

		// Remove from manager's list.
		s.user.mgr.mu.Lock()
		delete(s.user.mgr.tcp, s.Port)
		s.user.mgr.mu.Unlock()
	}

	if s.Type == "udp" {
		// Remove from user's list.
		s.user.mu.Lock()
		delete(s.user.udp, s.Port)
		s.user.mu.Unlock()

		// Remove from manager's list.
		s.user.mgr.mu.Lock()
		delete(s.user.mgr.udp, s.Port)
		s.user.mgr.mu.Unlock()
	}

	return nil
}

func (s *Service) stopListener() error {
	if s.Type == "udp" {
		return s.udpConn.Close()
	}

	if s.Type == "unix" || s.Type == "tcp" {
		if err := s.listener.Close(); err != nil {
			return err
		}

		if s.sock != "" {
			return os.Remove(s.sock)
		}
	}

	return nil
}
