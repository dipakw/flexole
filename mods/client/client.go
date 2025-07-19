package client

import (
	"context"
	"encoding/json"
	"flexole/mods/cmd"
	"sync"
)

func New(c *Config) (*Client, error) {
	instance := &Client{
		conf:     c,
		mu:       sync.RWMutex{},
		pipes:    map[string]*aPipe{},
		services: map[uint16]*Service{},
	}

	instance.Pipes = &Pipes{
		c: instance,
	}

	instance.ctx, instance.cancel = context.WithCancel(context.Background())

	return instance, nil
}

func (c *Client) Wait() {
	c.wg.Wait()
}

func (c *Client) Expose(s *Service) (uint16, error) {
	jsonBytes, err := json.Marshal(s.Remote)

	if err != nil {
		return 0, err
	}

	if err := c.sendCtrlCommand(true, cmd.CMD_EXPOSE, jsonBytes); err != nil {
		return 0, err
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	c.services[s.Remote.ID] = s

	return s.Remote.ID, nil
}

func (c *Client) Dispose(id uint16) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	service, ok := c.services[id]

	if !ok {
		return nil
	}

	jsonBytes, err := json.Marshal(&NetPort{
		Net:  service.Remote.Net,
		Port: service.Remote.Port,
	})

	if err != nil {
		return err
	}

	if err := c.sendCtrlCommand(false, cmd.CMD_DISPOSE, jsonBytes); err != nil {
		return err
	}

	delete(c.services, id)

	return nil
}

func (c *Client) Shutdown() error {
	if err := c.sendCtrlCommand(true, cmd.CMD_SHUTDOWN, nil); err != nil {
		return err
	}

	c.cancel()
	c.mu.Lock()
	defer c.mu.Unlock()

	c.services = map[uint16]*Service{}

	return nil
}
