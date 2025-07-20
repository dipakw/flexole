package server

import (
	"net"
)

func (pp *Pipes) add(pipe *Pipe) *Pipe {
	pp.user.mu.Lock()
	defer pp.user.mu.Unlock()
	pp.user.pipesList[pipe.id] = pipe

	pp.server.conf.Log.Inff("Pipe added => user: %s | addr: %s | id: %s", pp.user.id, pipe.conn.RemoteAddr(), pipe.id)

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

	pp.server.conf.Log.Inff("Pipe removed => user: %s | addr: %s | id: %s", pp.user.id, pipe.conn.RemoteAddr(), pipe.id)

	unpipedServicesIds := pp.user.services.unpipedUnsafe()

	if len(unpipedServicesIds) > 0 {
		pp.server.conf.Log.Inff("Removing unpiped services => user: %s | ids: %v", pp.user.id, unpipedServicesIds)

		for _, id := range unpipedServicesIds {
			pp.user.services.remUnsafe(id)
		}
	}

	return pipe
}

func (pp *Pipes) purge() error {
	pp.server.conf.Log.Inff("Purging pipes => user: %s", pp.user.id)

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

func (pp *Pipes) hasUnsafe(id string) bool {
	if pipe, ok := pp.user.pipesList[id]; ok && pipe != nil {
		return true
	}

	return false
}

func (pp *Pipes) getUnsafe(id string) *Pipe {
	return pp.user.pipesList[id]
}
