package services

func (s *Service) Info() *Info {
	return &Info{
		Host: s.Host,
		Port: s.Port,
		Type: s.Type,
	}
}
