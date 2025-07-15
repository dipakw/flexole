package services

import (
	"fmt"
	"net"
	"time"
)

func (s *Services) startUDP(service *Service) error {
	conn, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.ParseIP(service.Host),
		Port: int(service.Port),
	})

	if err != nil {
		return err
	}

	// Add service to the list.
	s.mutex.Lock()
	s.list[service.key] = service
	s.mutex.Unlock()

	go func() {
		for {
			buffer := make([]byte, MAX_UDP_PACKET_SIZE)
			n, addr, err := conn.ReadFromUDP(buffer)

			if err != nil {
				continue
			}

			src, err := service.SrcFN(service.Info())

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

				if _, err = conn.WriteToUDP(buf[:n], addr); err != nil {
					fmt.Println("Failed to write to", "::", addr, "::", err)
				}
			}()
		}
	}()

	return nil
}
