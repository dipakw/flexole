package client

import (
	"fmt"
	"net"

	"github.com/xtaci/smux"
)

func (c *Client) relayUDP(service *Service, stream *smux.Stream) {
	defer stream.Close()

	addr, err := net.ResolveUDPAddr("udp", service.Local.Addr)

	if err != nil {
		fmt.Println("Error resolving address:", err)
		return
	}

	conn, err := net.DialUDP("udp", nil, addr)

	if err != nil {
		fmt.Println("Failed to dial UDP:", err)
		return
	}

	defer conn.Close()

	buf := make([]byte, MAX_UDP_PACKET_SIZE)

	n, err := stream.Read(buf)

	if err != nil {
		fmt.Println("Failed to read from stream:", err)
		return
	}

	_, err = conn.Write(buf[:n])

	if err != nil {
		fmt.Println("Failed to write to UDP service:", err)
		return
	}

	n, err = conn.Read(buf)

	if err != nil {
		fmt.Println("Failed to read from UDP:", err)
		return
	}

	_, err = stream.Write(buf[:n])

	if err != nil {
		fmt.Println("Failed to write to stream:", err)
		return
	}
}
