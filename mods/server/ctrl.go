package server

import (
	"encoding/json"
	"flexole/mods/cmd"
	"flexole/mods/services"
	"fmt"
	"net"
)

func (s *Server) listenCtrl(pipe *Pipe) {
	buf := make([]byte, 256)

	for {
		n, err := pipe.ctrl.Read(buf)

		if err != nil {
			fmt.Println("Failed to read control command:", err)
			break
		}

		command := (&cmd.Cmd{}).Unpack(buf[:n])

		if command == nil {
			fmt.Println("Failed to unpack command")
			continue
		}

		status := cmd.CMD_STATUS_UNKNOWN

		switch command.ID {
		case cmd.CMD_EXPOSE:
			status = s.cmdExpose(pipe.userID, command.Data)
		case cmd.CMD_DISPOSE:
			status = s.cmdDispose(pipe.userID, command.Data)
		default:
			status = cmd.CMD_INVALID_CMD
		}

		pipe.ctrl.Write(cmd.New(status, nil).Pack())
	}
}

func (s *Server) cmdExpose(userID string, data []byte) uint8 {
	var service Service

	if err := json.Unmarshal(data, &service); err != nil {
		return cmd.CMD_MALFORMED_DATA
	}

	user := s.conf.Manager.User(userID)

	if !user.Available(service.Net, service.Port) {
		return cmd.CMD_PORT_UNAVAIL
	}

	_, err := user.Start(true, &services.Service{
		ID:   service.ID,
		Host: "",
		Port: service.Port,
		Type: service.Net,

		SrcFN: func(info *services.Info) (net.Conn, error) {
			return s.srcfn(userID, info)
		},
	})

	if err != nil {
		return cmd.CMD_OP_FAILED
	}

	// Add service to server user services list.
	serverUser := s.User(userID)
	serverUser.mu.Lock()
	defer serverUser.mu.Unlock()
	serverUser.services[service.ID] = &service

	return cmd.CMD_STATUS_OK
}

func (s *Server) cmdDispose(userID string, data []byte) uint8 {
	var netPort NetPort

	if err := json.Unmarshal(data, &netPort); err != nil {
		return cmd.CMD_MALFORMED_DATA
	}

	user := s.conf.Manager.User(userID)

	if err := user.Stop(netPort.Net, netPort.Port); err != nil {
		return cmd.CMD_OP_FAILED
	}

	return cmd.CMD_STATUS_OK
}
