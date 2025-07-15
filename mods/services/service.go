package services

import "os"

func (s *Service) Info() *Info {
	return &Info{
		Host: s.Host,
		Port: s.Port,
		Type: s.Type,
	}
}

func (s *Service) Key() string {
	return s.key
}

func (s *Service) stop() error {
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

func (s *Service) Stop() error {
	return s.manager.Stop(s.key)
}
