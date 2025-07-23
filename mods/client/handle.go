package client

import (
	"errors"
	"flexole/mods/cmd"
	"fmt"
	"io"
	"net"

	"github.com/xtaci/smux"
)

func (c *Client) setupCtrlChan(sess *smux.Session) (*smux.Stream, error) {
	stream, err := sess.AcceptStream()

	if err != nil {
		return nil, err
	}

	buf := make([]byte, 8)

	n, err := stream.Read(buf)

	if err != nil {
		return nil, err
	}

	if (&cmd.Cmd{}).Unpack(buf[:n]).ID != cmd.CMD_OPEN_CTRL_CHAN {
		return nil, errors.New("invalid ctrl chan command")
	}

	return stream, nil
}

func (c *Client) listen(pipeId string) {
	pipe := c.pipesList[pipeId]

	defer c.wg.Done()
	defer c.Pipes.Rem(pipeId)

	for {
		select {
		case <-c.ctx.Done():
			return
		default:
			stream, err := pipe.sess.AcceptStream()

			if err != nil {
				fmt.Println("ERR:", err, err == io.EOF)
				return
			}

			go c.handle(stream)
		}
	}
}

func (c *Client) handle(stream *smux.Stream) {
	defer stream.Close()

	buf := make([]byte, 8)
	n, err := stream.Read(buf)

	if err != nil {
		fmt.Println("Failed to read from stream:", err)
		return
	}

	command := (&cmd.Cmd{}).Unpack(buf[:n])

	if command == nil {
		fmt.Println("Failed to unpack command")
		return
	}

	if command.ID != cmd.CMD_CONNECT {
		fmt.Println("Invalid command:", command.ID)
		return
	}

	serviceID := cmd.UnpackUint16(command.Data)

	service, ok := c.servicesList[serviceID]

	if !ok {
		fmt.Println("Service not found:", serviceID)
		return
	}

	if service.Local.Net == "tcp" || service.Local.Net == "unix" {
		conn, err := net.Dial(service.Local.Net, service.Local.Addr)

		if err != nil {
			fmt.Println("Failed to dial service:", err)
			return
		}

		relay(c.ctx, conn, stream)

		return
	}

	if service.Local.Net == "udp" {
		c.relayUDP(service, stream)
		return
	}

	fmt.Println("Service not supported:", service.Local.Net)

	return
}
