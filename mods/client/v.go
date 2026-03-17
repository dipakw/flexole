package client

import (
	"flexole/mods/serve"
	"net"

	"github.com/dipakw/byrate/handle"
	"github.com/xtaci/smux"
)

var byrateConfig = &handle.Config{
	Version: "0.2.1",
}

func (c *Client) vSpeed(stream *smux.Stream) {
	conn1, conn2 := net.Pipe()

	defer conn1.Close()
	defer conn2.Close()

	go handle.Handle(conn1, byrateConfig)

	relay(c.ctx, conn2, stream)
}

func (c *Client) vServe(serviceID uint16, dir string, stream *smux.Stream) error {
	var err error

	c.mu.RLock()
	_, exists := c.serves[serviceID]
	c.mu.RUnlock()

	if !exists {
		c.mu.Lock()
		c.serves[serviceID], err = serve.New(&serve.Config{
			Dir: dir,
		})

		if err != nil {
			delete(c.serves, serviceID)
			c.mu.Unlock()
			return err
		}

		c.mu.Unlock()
	}

	conn1, conn2 := net.Pipe()

	defer conn1.Close()
	defer conn2.Close()

	go c.serves[serviceID].Handle(conn1)

	relay(c.ctx, conn2, stream)

	return nil
}
