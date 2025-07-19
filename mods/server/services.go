package server

import (
	"flexole/mods/cmd"
	"flexole/mods/services"
	"net"
)

func (ss *Services) add(service *Service) (*Service, uint8) {
	ss.user.mu.Lock()
	defer ss.user.mu.Unlock()

	user := ss.server.conf.Manager.User(ss.user.id)

	if !user.Available(service.Net, service.Port) {
		return nil, cmd.CMD_PORT_UNAVAIL
	}

	ss.user.servicesList[service.ID] = service

	_, err := user.Start(&services.Service{
		ID:   service.ID,
		Host: "",
		Port: service.Port,
		Type: service.Net,

		SrcFN: func(info *services.Info) (net.Conn, error) {
			return ss.server.srcfn(ss.user.id, info)
		},
	})

	if err != nil {
		ss.server.conf.Log.Errf("Failed to start service: %s", err.Error())
		return nil, cmd.CMD_OP_FAILED
	}

	return service, cmd.CMD_STATUS_OK
}

func (ss *Services) rem(id uint16) (*Service, uint8) {
	ss.user.mu.Lock()
	defer ss.user.mu.Unlock()
	service, ok := ss.user.servicesList[id]

	if !ok || service == nil {
		return nil, cmd.CMD_OP_FAILED
	}

	// Stop service.
	if _, err := ss.server.conf.Manager.User(ss.user.id).Stop(service.Net, service.Port); err != nil {
		return nil, cmd.CMD_OP_FAILED
	}

	// Remove service from user services list.
	delete(ss.user.servicesList, id)

	return service, cmd.CMD_STATUS_OK
}

func (ss *Services) purge() error {
	ss.user.mu.RLock()
	defer ss.user.mu.RUnlock()

	// Request manager to stop all services.
	ss.server.conf.Manager.User(ss.user.id).Reset()

	// Remove all services from user services list.
	ss.user.servicesList = map[uint16]*Service{}

	return nil
}
