package client

import (
	"errors"
	"flexole/mods/cmd"
	"flexole/mods/util"
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
		return nil, errors.New("received non-open-ctrl-chan command")
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
				if err != io.EOF {
					c.conf.Log.Errf("Failed to accept stream => pipe: %s | error: %s", pipeId, err.Error())
				}

				return
			}

			go c.handle(pipeId, stream)
		}
	}
}

func (c *Client) handle(pipeId string, stream *smux.Stream) {
	defer stream.Close()

	buf := make([]byte, 8)
	n, err := stream.Read(buf)

	if err != nil {
		if err != io.EOF {
			c.conf.Log.Errf("Failed to read from stream => pipe: %s | error: %s", pipeId, err.Error())
		}

		return
	}

	command := (&cmd.Cmd{}).Unpack(buf[:n])

	if command == nil {
		c.conf.Log.Errf("Failed to unpack command => pipe: %s | data: %v", pipeId, buf[:n])
		return
	}

	if command.ID != cmd.CMD_CONNECT {
		c.conf.Log.Errf("Received non-connect command => pipe: %s | command: %d", pipeId, command.ID)
		return
	}

	serviceID := util.UnpackUint16(command.Data)

	service, ok := c.servicesList[serviceID]

	if service != nil {
		c.conf.Log.Inff("Requested service => id: %d | net: %s | addr: %s | pipe: %s", serviceID, service.Local.Net, service.Local.Addr, pipeId)
	} else {
		c.conf.Log.Inff("Requested service => id: %d | pipe: %s", serviceID, pipeId)
	}

	if !ok {
		c.conf.Log.Wrnf("Service not found => pipe: %s | id: %d", pipeId, serviceID)
		return
	}

	if service.Local.Net == "tcp" || service.Local.Net == "unix" {
		conn, err := net.Dial(service.Local.Net, service.Local.Addr)

		if err != nil {
			c.conf.Log.Errf("Failed to connect => pipe: %s | error: %s", pipeId, err.Error())
			return
		}

		relay(c.ctx, conn, stream)

		return
	}

	if service.Local.Net == "udp" {
		c.relayUDP(service, stream)
		return
	}

	if service.Local.Net == "v" {
		if service.Local.Addr == "speed" {
			c.vSpeed(stream)
			return
		}

		c.conf.Log.Errf("Service not found => pipe: %s | id: %d | net: %s | addr: %s", pipeId, serviceID, service.Local.Net, service.Local.Addr)

		return
	}

	c.conf.Log.Errf("Unsupported service => pipe: %s | id: %d | net: %s | addr: %s", pipeId, serviceID, service.Local.Net, service.Local.Addr)

	return
}
