package client

import (
	"encoding/json"
	"flexole/mods/cmd"
	"flexole/mods/util"
	"fmt"
)

func (ss *Services) Add(s *Service) (uint16, error) {
	if ss.Has(s.Remote.ID) {
		return 0, fmt.Errorf("id %d already exists", s.Remote.ID)
	}

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

	ss.c.wg.Add(1)

	return s.Remote.ID, nil
}

func (ss *Services) Rem(id uint16) error {
	ss.c.mu.Lock()
	defer ss.c.mu.Unlock()

	if !ss.hasUnsafe(id) {
		return nil
	}

	if err := ss.c.sendCtrlCommand(false, cmd.CMD_REM_SERVICE, util.PackUint16(id)); err != nil {
		return err
	}

	delete(ss.c.servicesList, id)

	return nil
}

func (ss *Services) Has(id uint16) bool {
	ss.c.mu.RLock()
	defer ss.c.mu.RUnlock()
	return ss.hasUnsafe(id)
}

func (ss *Services) hasUnsafe(id uint16) bool {
	s, ok := ss.c.servicesList[id]
	return ok && s != nil
}
