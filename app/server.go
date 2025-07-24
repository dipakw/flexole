package app

import (
	"flexole/mods/server"
	"flexole/mods/services"
	"flexole/mods/util"
	"fmt"
	"os"

	"github.com/dipakw/logs"
)

func startServer(conf *ServerConfig) (*server.Server, error) {
	dir := util.StrOr(conf.Dir, os.TempDir())

	manager := services.Manager(&services.Config{
		Dir: dir,
	})

	defer manager.Reset()

	logger := logs.New(&logs.Config{
		Allow: util.LogKindsToFlag(conf.Logs.Allow),

		Outs: []*logs.Out{
			{
				Target: os.Stdout,
				Color:  true,
			},
		},
	})

	addr, err := util.NetAddr(conf.Bind.Addr, DEFAULT_PORT)

	if err != nil {
		return nil, err
	}

	users := map[string]*User{}

	for _, user := range conf.Users {
		users[user.ID] = &user
	}

	flexole, err := server.New(&server.Config{
		Net:     conf.Bind.Net,
		Addr:    addr,
		Manager: manager,
		Log:     logger,

		KeyFN: func(id string) ([]byte, error) {
			user, ok := users[id]

			if !ok || user == nil {
				return nil, fmt.Errorf(`user "%s" not found`, id)
			}

			if !user.Enabled {
				return nil, fmt.Errorf(`user "%s" not enabled`, id)
			}

			if user.Key == "" {
				return nil, fmt.Errorf(`user "%s" has no key`, id)
			}

			return []byte(user.Key), nil
		},

		LimitFN: func(id string, kind string) int {
			if user, ok := users[id]; ok && user != nil {
				switch kind {
				case "pipes":
					return user.MaxPipes
				case "service:tcp":
					return user.MaxServices.TCP
				case "service:udp":
					return user.MaxServices.UDP
				case "service:unix":
					return user.MaxServices.Unix
				default:
					return 0
				}
			}

			return 0
		},
	})

	if err != nil {
		return nil, err
	}

	if err := flexole.Start(); err != nil {
		return nil, err
	}

	return flexole, nil
}
