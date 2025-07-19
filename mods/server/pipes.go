package server

import "net"

func (pp *Pipes) add(pipe *Pipe) *Pipe {
	pp.user.mu.Lock()
	defer pp.user.mu.Unlock()
	pp.user.pipesList[pipe.id] = pipe

	return pipe
}

func (pp *Pipes) rem(id string) *Pipe {
	pp.user.mu.Lock()
	defer pp.user.mu.Unlock()
	pipe, ok := pp.user.pipesList[id]

	if !ok || pipe == nil {
		return nil
	}

	delete(pp.user.pipesList, id)

	return pipe
}

func (pp *Pipes) purge() error {
	pp.user.mu.RLock()

	conns := make(map[string]net.Conn)

	for i, pipe := range pp.user.pipesList {
		conns[i] = pipe.conn
	}

	pp.user.mu.RUnlock()

	for _, conn := range conns {
		conn.Close()
	}

	return nil
}
