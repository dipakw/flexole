package server

import (
	"context"
	"encoding/json"
	"flexole/mods/cmd"
	"io"
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
				s.conf.Log.Wrn("Failed to read control command:", err)
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
	case cmd.CMD_EXPOSE:
		status = s.cmdExpose(pipe.userID, command.Data)
	case cmd.CMD_DISPOSE:
		status = s.cmdDispose(pipe.userID, command.Data)
	case cmd.CMD_SHUTDOWN:
		status = s.cmdShutdown(pipe.userID)
	default:
		status = cmd.CMD_INVALID_CMD
	}

	pipe.ctrl.Write(cmd.New(status, nil).Pack())
}

func (s *Server) cmdExpose(userID string, data []byte) uint8 {
	var service Service

	if err := json.Unmarshal(data, &service); err != nil {
		return cmd.CMD_MALFORMED_DATA
	}

	// Add service to server user services list.
	_, status := s.User(userID).services.add(&service)

	return status
}

func (s *Server) cmdDispose(userID string, data []byte) uint8 {
	id := cmd.UnpackUint16(data)
	_, status := s.User(userID).services.rem(id)

	return status
}

func (s *Server) cmdShutdown(userID string) uint8 {
	// Get user.
	user := s.User(userID)

	// Stop all services.
	user.services.purge()

	// Close all pipes.
	user.pipes.purge()

	return cmd.CMD_STATUS_OK
}
