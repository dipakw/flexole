package client

import (
	"errors"
	"flexole/mods/cmd"
	"fmt"
)

func (c *Client) sendCtrlCommand(safe bool, id uint8, payload []byte) error {
	if safe {
		c.mu.Lock()
		defer c.mu.Unlock()
	}

	if len(c.pipesList) == 0 {
		return errors.New("no pipes available")
	}

	command := cmd.New(id, payload).Pack()

	var n int
	var err error

	for _, pipe := range c.pipesList {
		if pipe == nil {
			err = errors.New("pipe is nil")
			continue
		}

		n, err = pipe.ctrl.Write(command)

		// Command sent successfully.
		if err == nil && n == len(command) {
			buf := make([]byte, 8)

			// Wait for response.
			n, err = pipe.ctrl.Read(buf)

			if err == nil {
				response := cmd.New(0, nil).Unpack(buf)

				fmt.Println("Response:", response)

				if response.ID != cmd.CMD_STATUS_OK {
					return fmt.Errorf("status code: %d", response.ID)
				}
			}

			break
		}

		if n != len(command) {
			err = errors.New("partial command write")
		}
	}

	return err
}
