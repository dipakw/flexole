package services

import (
	"fmt"
	"net"
	"time"
)

func (u *User) startUDP(service *Service, background bool) error {
	var err error

	service.udpConn, err = net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.ParseIP(service.Host),
		Port: int(service.Port),
	})

	if err != nil {
		return err
	}

	// Add service to the user's list.
	u.mu.Lock()
	u.udp[service.Port] = service
	u.mu.Unlock()

	// Add service to the manager's list.
	u.mgr.mu.Lock()
	u.mgr.udp[service.Port] = true
	u.mgr.mu.Unlock()

	info := service.Info()

	run := func() {
		for {
			buffer := make([]byte, MAX_UDP_PACKET_SIZE)
			n, addr, err := service.udpConn.ReadFromUDP(buffer)

			if err != nil {
				continue
			}

			src, err := service.SrcFN(info)

			if err != nil {
				continue
			}

			go func() {
				defer src.Close()

				src.SetReadDeadline(time.Now().Add(service.Timeout))

				_, err := src.Write(buffer[:n])

				if err != nil {
					return
				}

				buf := make([]byte, MAX_UDP_PACKET_SIZE)
				n, err = src.Read(buf)

				if err != nil {
					return
				}

				if _, err = service.udpConn.WriteToUDP(buf[:n], addr); err != nil {
					fmt.Println("Failed to write to", "::", addr, "::", err)
				}
			}()
		}
	}

	if background {
		go run()
	} else {
		run()
	}

	return nil
}
