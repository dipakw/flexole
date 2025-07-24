package server

import (
	"context"
	"encoding/json"
	"flexole/mods/cmd"
	"flexole/mods/util"
	"fmt"
	"io"
	"strings"
)

func (s *Server) listenCtrl(ctx context.Context, pipe *Pipe) {
	defer pipe.ctrl.Close()

	buf := make([]byte, 256)

	for {
		select {
		case <-ctx.Done():
			return
		default:
			n, err := pipe.ctrl.Read(buf)

			if err == io.EOF {
				return
			}

			if err != nil {
				if !strings.Contains(err.Error(), "use of closed network connection") {
					s.conf.Log.Err("Failed to read control command:", err)
				}

				return
			}

			command := (&cmd.Cmd{}).Unpack(buf[:n])

			if command == nil {
				s.conf.Log.Err("Failed to unpack command:", buf[:20])
				continue
			}

			s.handleCommand(command, pipe)
		}
	}
}

func (s *Server) handleCommand(command *cmd.Cmd, pipe *Pipe) {
	status := cmd.CMD_STATUS_UNKNOWN
	response := []byte{}

	switch command.ID {
	case cmd.CMD_ADD_SERVICE:
		var service *Service
		status, service = s.cmdAddService(pipe.userID, command.Data)

		if status == cmd.CMD_STATUS_OK && service != nil {
			response = s.conf.EvtAddService(s, pipe.userID, service)
		}

	case cmd.CMD_REM_SERVICE:
		status = s.cmdRemService(pipe.userID, command.Data)

	case cmd.CMD_SHUTDOWN:
		status = s.cmdShutdown(pipe.userID)

	default:
		s.conf.Log.Errf("Invalid command => user: %s | id: %d", pipe.userID, command.ID)
		status = cmd.CMD_INVALID_CMD
	}

	pipe.ctrl.Write(cmd.New(status, response).Pack())
}

func (s *Server) cmdAddService(userID string, data []byte) (uint8, *Service) {
	var service Service

	if err := json.Unmarshal(data, &service); err != nil {
		s.conf.Log.Errf("Malformed command [ADD_SERVICE] => user: %s | error: %s", userID, err.Error())
		return cmd.CMD_MALFORMED_DATA, nil
	}

	user := s.User(userID)
	limit := s.conf.LimitFN(userID, fmt.Sprintf("service:%s", service.Net))

	if limit == 0 {
		s.conf.Log.Inff("Services limit is 0 => user: %s | net: %s", userID, service.Net)
		return cmd.CMD_NOT_AVAILABLE, nil
	}

	if user.services.count(service.Net) >= limit {
		s.conf.Log.Inff("Services limit reached => user: %s | net: %s | limit: %d", userID, service.Net, s.conf.LimitFN(userID, service.Net))
		return cmd.CMD_SERVICES_LIMIT, nil
	}

	s.conf.Log.Inff("Command [ADD_SERVICE] => user: %s | net: %s | port: %d | id: %d", userID, service.Net, service.Port, service.ID)

	// Add service to server user services list.
	_, status := user.services.add(&service)

	return status, &service
}

func (s *Server) cmdRemService(userID string, data []byte) uint8 {
	id := util.UnpackUint16(data)

	s.conf.Log.Inff("Command [REM_SERVICE] => user: %s | id: %d", userID, id)

	_, status := s.User(userID).services.rem(id)

	return status
}

func (s *Server) cmdShutdown(userID string) uint8 {
	s.conf.Log.Inff("Command [SHUTDOWN] => user: %s", userID)

	// Get user.
	user := s.User(userID)

	// Stop all services.
	user.services.purge()

	// Close all pipes.
	user.pipes.purge()

	return cmd.CMD_STATUS_OK
}
