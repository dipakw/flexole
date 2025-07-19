package services

import (
	"fmt"
	"net"
	"os"
	"path"
	"strings"
	"time"
)

func (u *User) startTCPOrUnix(service *Service) error {
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

	service.wg.Add(1)

	if service.Type == "tcp" {
		u.mu.Lock()
		u.mgr.mu.Lock()

		u.tcp[service.Port] = service
		u.mgr.tcp[service.Port] = true

		u.mu.Unlock()
		u.mgr.mu.Unlock()
	}

	if service.Type == "unix" {
		u.mu.Lock()
		u.unix[service.Port] = service
		u.mu.Unlock()
	}

	info := service.Info()

	go func() {
		defer service.wg.Done()
		defer time.Sleep(10 * time.Millisecond)
		defer service.listener.Close()

		for {
			select {
			case <-service.ctx.Done():
				return
			default:
				conn, err := service.listener.Accept()

				if err != nil {
					if strings.Contains(err.Error(), "closed") {
						return
					}

					continue
				}

				src, err := service.SrcFN(info)

				if err != nil {
					conn.Close()
					continue
				}

				go relay(service.ctx, conn, src)
			}
		}
	}()

	return nil
}
