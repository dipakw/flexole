package server

import (
	"flexole/mods/cmd"
	"flexole/mods/services"
	"fmt"
	"net"
)

func (s *Server) srcfn(userID string, info *services.Info) (net.Conn, error) {
	user := s.User(userID)

	var pipe *Pipe

	for _, p := range user.pipes {
		if p.active {
			pipe = p
			break
		}
	}

	if pipe == nil {
		fmt.Println("No active pipe")
		return nil, fmt.Errorf("no active pipe")
	}

	// Open stream.
	stream, err := pipe.sess.OpenStream()

	if err != nil {
		fmt.Println("Failed to open stream:", err)
		return nil, err
	}

	// Send connect command.
	_, err = stream.Write(cmd.New(cmd.CMD_CONNECT, cmd.PackUint16(info.ID)).Pack())

	if err != nil {
		fmt.Println("Failed to send connect command:", err)
		return nil, err
	}

	return stream, nil
}
