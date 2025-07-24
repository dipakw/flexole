package server

import (
	"flexole/mods/cmd"
	"flexole/mods/services"
	"net"
	"strings"
	"time"
)

func (ss *Services) add(service *Service) (*Service, uint8) {
	ss.server.conf.Log.Inff("Adding service => user: %s | net: %s | port: %d | id: %d", ss.user.id, service.Net, service.Port, service.ID)

	ss.user.mu.Lock()
	defer ss.user.mu.Unlock()

	user := ss.server.conf.Manager.User(ss.user.id)

	if !user.Available(service.Net, service.Port) {
		return nil, cmd.CMD_PORT_UNAVAIL
	}

	_, err := user.Start(&services.Service{
		ID:      service.ID,
		Host:    "",
		Port:    service.Port,
		Type:    service.Net,
		Timeout: 10 * time.Second, // TODO: Make this configurable.

		SrcFN: func(info *services.Info) (net.Conn, error) {
			return ss.server.srcfn(ss.user.id, info)
		},
	})

	if err != nil {
		if strings.Contains(err.Error(), "address already in use") {
			return nil, cmd.CMD_PORT_UNAVAIL
		}

		ss.server.conf.Log.Errf("Failed to start service => user: %s | net: %s | port: %d | id: %d | error: %s", ss.user.id, service.Net, service.Port, service.ID, err.Error())
		return nil, cmd.CMD_OP_FAILED
	}

	ss.user.servicesList[service.ID] = service

	ss.server.conf.Log.Inff("Service added => user: %s | net: %s | port: %d | id: %d", ss.user.id, service.Net, service.Port, service.ID)

	return service, cmd.CMD_STATUS_OK
}

func (ss *Services) rem(id uint16) (*Service, uint8) {
	ss.user.mu.RLock()
	defer ss.user.mu.RUnlock()

	return ss.remUnsafe(id)
}

func (ss *Services) getUnsafe(id uint16) *Service {
	return ss.user.servicesList[id]
}

func (ss *Services) remUnsafe(id uint16) (*Service, uint8) {
	ss.server.conf.Log.Inff("Removing service => user: %s | id: %d", ss.user.id, id)

	service, ok := ss.user.servicesList[id]

	if !ok || service == nil {
		return nil, cmd.CMD_OP_FAILED
	}

	// Stop service.
	if _, err := ss.server.conf.Manager.User(ss.user.id).Stop(service.Net, service.Port); err != nil {
		ss.server.conf.Log.Errf("Failed to stop service => user: %s | net: %s | port: %d | id: %d | error: %s", ss.user.id, service.Net, service.Port, id, err.Error())
		return nil, cmd.CMD_OP_FAILED
	}

	// Remove service from user services list.
	delete(ss.user.servicesList, id)

	ss.server.conf.Log.Inff("Service removed => user: %s | net: %s | port: %d | id: %d", ss.user.id, service.Net, service.Port, id)

	return service, cmd.CMD_STATUS_OK
}

func (ss *Services) purge() error {
	ss.server.conf.Log.Inff("Purging services => user: %s", ss.user.id)

	ss.user.mu.RLock()

	// Get the ids of services.
	ids := map[uint16]bool{}

	for id := range ss.user.servicesList {
		ids[id] = true
	}

	ss.user.mu.RUnlock()

	for id := range ids {
		ss.rem(id)
	}

	return nil
}

func (ss *Services) count(kind string) int {
	ss.user.mu.RLock()
	defer ss.user.mu.RUnlock()

	if kind == "*" {
		return len(ss.user.servicesList)
	}

	count := 0

	for _, s := range ss.user.servicesList {
		if s.Net == kind {
			count++
		}
	}

	return count
}

// Get the ids of services that have no pipes.
func (ss *Services) unpipedUnsafe() []uint16 {
	ids := []uint16{}

	for _, s := range ss.user.servicesList {
		isUnpiped := true

		for _, p := range s.Pipes {
			if pipe, ok := ss.user.pipesList[p]; ok && pipe != nil {
				isUnpiped = false
				break
			}
		}

		if isUnpiped {
			ids = append(ids, s.ID)
		}
	}

	return ids
}
