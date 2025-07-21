package app

import (
	"flexole/mods/server"
	"flexole/mods/services"
	"fmt"
	"os"

	"github.com/dipakw/logs"
)

func startServer(conf *ServerConfig) {
	manager := services.Manager(&services.Config{
		Dir: "./tmp", // TODO: make it configurable
	})

	defer manager.Reset()

	logger := logs.New(&logs.Config{
		Allow: logs.ALL,

		Outs: []*logs.Out{
			{
				Target: os.Stdout,
				Color:  true,
			},
		},
	})

	users := map[string]*User{}

	for _, user := range conf.Users {
		users[user.ID] = &user
	}

	flexole, err := server.New(&server.Config{
		Net:     conf.Config.Net,
		Addr:    conf.Config.Addr,
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
		logger.Err(err)
		return
	}

	if err := flexole.Start(); err != nil {
		logger.Err(err)
		return
	}

	fmt.Printf("Server started: %s://%s\n", conf.Config.Net, conf.Config.Addr)

	flexole.Wait()
}
