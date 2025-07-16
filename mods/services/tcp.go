package services

import (
	"fmt"
	"net"
	"os"
	"path"
)

func (u *User) startTCPOrUnix(service *Service, background bool) error {
	addr := net.JoinHostPort(service.Host, fmt.Sprintf("%d", service.Port))

	if service.Type == "unix" {
		addr = path.Join(u.dir, fmt.Sprintf("%d.sock", service.Port))

		if _, err := os.Stat(addr); err == nil {
			if err := os.Remove(addr); err != nil {
				return err
			}
		}

		if err := os.MkdirAll(path.Dir(addr), 0755); err != nil {
			return err
		}

		service.sock = addr
	}

	var err error

	service.listener, err = net.Listen(service.Type, addr)

	if err != nil {
		return err
	}

	if service.Type == "tcp" {
		// Add service to the user's list.
		u.mu.Lock()
		u.tcp[service.Port] = service
		u.mu.Unlock()

		// Add service to the manager's list.
		u.mgr.mu.Lock()
		u.mgr.tcp[service.Port] = true
		u.mgr.mu.Unlock()
	}

	if service.Type == "unix" {
		u.mu.Lock()
		u.unix[service.Port] = service
		u.mu.Unlock()
	}

	info := service.Info()

	run := func() {
		for {
			conn, err := service.listener.Accept()

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
	}

	if background {
		go run()
	} else {
		run()
	}

	return nil
}
