package server

import (
	"context"
	"encoding/json"
	"flexole/mods/cmd"
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

	switch command.ID {
	case cmd.CMD_ADD_SERVICE:
		status = s.cmdAddService(pipe.userID, command.Data)
	case cmd.CMD_REM_SERVICE:
		status = s.cmdRemService(pipe.userID, command.Data)
	case cmd.CMD_SHUTDOWN:
		status = s.cmdShutdown(pipe.userID)
	default:
		s.conf.Log.Errf("Invalid command => user: %s | id: %d", pipe.userID, command.ID)
		status = cmd.CMD_INVALID_CMD
	}

	pipe.ctrl.Write(cmd.New(status, nil).Pack())
}

func (s *Server) cmdAddService(userID string, data []byte) uint8 {
	var service Service

	if err := json.Unmarshal(data, &service); err != nil {
		s.conf.Log.Errf("Malformed command [ADD_SERVICE] => user: %s | error: %s", userID, err.Error())
		return cmd.CMD_MALFORMED_DATA
	}

	s.conf.Log.Inff("Command [ADD_SERVICE] => user: %s | net: %s | port: %d | id: %d", userID, service.Net, service.Port, service.ID)

	// Add service to server user services list.
	_, status := s.User(userID).services.add(&service)

	return status
}

func (s *Server) cmdRemService(userID string, data []byte) uint8 {
	id := cmd.UnpackUint16(data)

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
