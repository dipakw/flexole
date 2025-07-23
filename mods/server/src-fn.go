package server

import (
	"flexole/mods/cmd"
	"flexole/mods/services"
	"flexole/mods/util"
	"fmt"
	"net"
	"sync"
)

// UserID -> ServiceID -> Pipe Index.
var usedPipesMu = sync.RWMutex{}
var lastUsedPipe = map[string]map[uint16]int{}

func (s *Server) srcfn(userID string, info *services.Info) (net.Conn, error) {
	user := s.User(userID)

	user.mu.RLock()
	defer user.mu.RUnlock()

	service := user.services.getUnsafe(info.ID)

	if service == nil {
		s.conf.Log.Errf("Service not found => user: %s | id: %d", userID, info.ID)
		return nil, fmt.Errorf("service not found")
	}

	if len(service.Pipes) == 0 {
		s.conf.Log.Wrnf("Service has no pipes set => user: %s | id: %d", userID, info.ID)
		return nil, fmt.Errorf("no pipes set")
	}

	// Get the available pipes IDs.
	pipesIds := []string{}

	for _, pipeId := range service.Pipes {
		if user.pipes.hasUnsafe(pipeId) {
			pipesIds = append(pipesIds, pipeId)
		}
	}

	if len(pipesIds) == 0 {
		s.conf.Log.Wrnf("Service has no active pipes => user: %s | id: %d", userID, info.ID)
		return nil, fmt.Errorf("no active pipes")
	}

	// Get the last used pipe index.
	usedPipesMu.Lock()
	index, ok := lastUsedPipe[userID][info.ID]

	if ok {
		index++

		if index >= len(pipesIds) {
			index = 0
		}
	}

	if !ok {
		lastUsedPipe[userID] = make(map[uint16]int)
	}

	lastUsedPipe[userID][info.ID] = index
	usedPipesMu.Unlock()

	// Get the pipe.
	pipe := user.pipes.getUnsafe(pipesIds[index])

	// Open stream.
	stream, err := pipe.sess.OpenStream()

	if err != nil {
		s.conf.Log.Errf("Failed to open stream => user: %s | id: %d | index: %d | error: %s", userID, info.ID, index, err.Error())
		return nil, err
	}

	// Send connect command.
	if _, err := stream.Write(cmd.New(cmd.CMD_CONNECT, util.PackUint16(info.ID)).Pack()); err != nil {
		s.conf.Log.Errf("Failed to send connect command => user: %s | service: %d | pipe: %s | error: %s", userID, info.ID, pipe.id, err.Error())
		return nil, err
	}

	s.conf.Log.Inff("Requested service => user: %s | service: %d | pipe: %s", userID, info.ID, pipe.id)

	return stream, nil
}
