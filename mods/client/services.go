package client

import (
	"encoding/json"
	"flexole/mods/cmd"
	"flexole/mods/util"
	"fmt"
)

func (ss *Services) Add(s *Service) (uint16, error) {
	jsonBytes, err := json.Marshal(s.Remote)

	if err != nil {
		return 0, err
	}

	if err := ss.c.sendCtrlCommand(true, cmd.CMD_ADD_SERVICE, jsonBytes); err != nil {
		return 0, err
	}

	ss.c.mu.Lock()
	defer ss.c.mu.Unlock()
	ss.c.servicesList[s.Remote.ID] = s

	return s.Remote.ID, nil
}

func (ss *Services) Rem(id uint16) error {
	ss.c.mu.Lock()
	defer ss.c.mu.Unlock()

	if service, ok := ss.c.servicesList[id]; !ok || service == nil {
		return fmt.Errorf("service not found: %d", id)
	}

	if err := ss.c.sendCtrlCommand(false, cmd.CMD_REM_SERVICE, util.PackUint16(id)); err != nil {
		return err
	}

	delete(ss.c.servicesList, id)

	return nil
}
