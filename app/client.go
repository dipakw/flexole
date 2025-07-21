package app

import (
	"flexole/mods/client"
	"os"

	"github.com/dipakw/logs"
)

func startClient(config *ClientConfig) {
	logger := logs.New(&logs.Config{
		Allow: logs.ALL,

		Outs: []*logs.Out{
			{
				Target: os.Stdout,
				Color:  true,
			},
		},
	})

	tunnel, err := client.New(&client.Config{
		ID:  []byte(config.Auth.ID),
		Key: []byte(config.Auth.Key),

		Server: &client.Addr{
			Net:  config.Auth.Server.Net,
			Addr: config.Auth.Server.Addr,
		},
	})

	if err != nil {
		logger.Err(err)
		return
	}

	// Add pipes.
	for _, pipe := range config.Pipes {
		if err := tunnel.Pipes.Add(pipe.ID, pipe.Encrypt); err != nil {
			logger.Err(err)
		}
	}

	// Add services.
	for _, service := range config.Services {
		_, err := tunnel.Services.Add(&client.Service{
			Local: &client.Local{
				Net:  service.Local.Net,
				Addr: service.Local.Addr,
			},

			Remote: &client.Remote{
				ID:    service.ID,
				Net:   service.Remote.Net,
				Port:  uint16(service.Remote.Port),
				Pipes: service.Remote.Pipes,
			},
		})

		if err != nil {
			logger.Err(err)
		}
	}

	tunnel.Wait()
}
