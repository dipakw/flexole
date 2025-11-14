package app

import (
	"flexole/mods/client"
	"flexole/mods/util"
	"os"
	"strings"

	"github.com/dipakw/logs"
)

func startClient(conf *ClientConfig) {
	logger := logs.New(&logs.Config{
		Allow: util.LogKindsToFlag(conf.Logs.Allow),

		Outs: []*logs.Out{
			{
				Target: os.Stdout,
				Color:  true,
			},
		},
	})

	tunnel, err := client.New(&client.Config{
		ID:  []byte(conf.Auth.ID),
		Key: []byte(conf.Auth.Key),
		Log: logger,

		Server: &client.Addr{
			Net:  conf.Server.Net,
			Addr: conf.Server.Addr,
		},
	})

	if err != nil {
		logger.Err(err)
		return
	}

	stop := false

	// Add pipes.
	for _, pipe := range conf.Pipes {
		if !pipe.Enabled {
			continue
		}

		if err := tunnel.Pipes.Add(pipe.ID, pipe.Encrypt); err != nil {
			stop = true

			if strings.Contains(err.Error(), "connection refused") {
				logger.Errf(`Failed to connect to server: %s`, tunnel.ServerAddr())
				break
			}

			logger.Err(err)
		}
	}

	if stop {
		return
	}

	// Add services.
	for _, service := range conf.Services {
		if !service.Enabled {
			continue
		}

		_, err := tunnel.Services.Add(&client.Service{
			Local: &client.Local{
				Net:  service.Local.Net,
				Addr: service.Local.Addr,
			},

			Remote: &client.Remote{
				ID:    service.ID,
				Net:   service.Remote.Net,
				Port:  uint16(service.Remote.Port),
				Pipes: service.Pipes,
			},
		})

		if err != nil {
			logger.Errf("Failed to add service => id: %d | error: %s", service.ID, err.Error())
		}
	}

	tunnel.Wait()
}
