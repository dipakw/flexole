package services

import (
	"fmt"
	"net"
	"os"
	"path"
)

func (s *Services) startTCPOrUnix(service *Service, dir string) error {
	addr := fmt.Sprintf("%s:%d", service.Host, service.Port)

	if service.Type == "unix" {
		addr = path.Join(dir, fmt.Sprintf("%d.sock", service.Port))

		if _, err := os.Stat(addr); err == nil {
			if err := os.Remove(addr); err != nil {
				return err
			}
		}
	}

	listener, err := net.Listen("unix", addr)

	if err != nil {
		return err
	}

	// Add service to the list.
	s.mutex.Lock()
	s.list[service.key] = service
	s.mutex.Unlock()

	info := service.Info()

	go func() {
		for {
			conn, err := listener.Accept()

			if err != nil {
				continue
			}

			src, err := service.SrcFN(info)

			if err != nil {
				conn.Close()
				continue
			}

			go relay(conn, src)
		}
	}()

	return nil
}
