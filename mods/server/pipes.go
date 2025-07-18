package server

func (s *Server) AddPipe(pipe *Pipe) *Pipe {
	user := s.User(pipe.userID)

	user.mu.Lock()
	defer user.mu.Unlock()

	user.pipes[pipe.id] = pipe

	return pipe
}

func (s *Server) RemPipe(userID string, pipeID string) {
	user := s.User(userID)

	user.mu.Lock()
	defer user.mu.Unlock()

	delete(user.pipes, pipeID)
}
