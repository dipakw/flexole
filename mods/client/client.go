package client

import (
	"context"
	"flexole/mods/cmd"
	"sync"
)

func New(c *Config) (*Client, error) {
	instance := &Client{
		conf:         c,
		mu:           sync.RWMutex{},
		pipesList:    map[string]*connPipe{},
		servicesList: map[uint16]*Service{},
	}

	instance.Pipes = &Pipes{
		c: instance,
	}

	instance.Services = &Services{
		c: instance,
	}

	instance.ctx, instance.cancel = context.WithCancel(context.Background())

	return instance, nil
}

func (c *Client) Wait() {
	c.wg.Wait()
}

func (c *Client) Shutdown() error {
	if _, err := c.sendCtrlCommand(true, cmd.CMD_SHUTDOWN, nil); err != nil {
		return err
	}

	c.cancel()
	c.mu.Lock()
	defer c.mu.Unlock()

	c.servicesList = map[uint16]*Service{}

	return nil
}

func (c *Client) ServerAddr() string {
	return c.conf.Server.Net + "/" + c.conf.Server.Addr
}
