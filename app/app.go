package app

import (
	"fmt"
	"os"

	"github.com/dipakw/logs"
)

func Run(config *Config) {
	if len(os.Args) < 2 {
		fmt.Println(cli_doc)
		return
	}

	cmd := os.Args[1]

	switch cmd {
	case "version", "v":
		fmt.Println("Version: " + config.Version)

	case "server", "s":
		logger := logs.New(&logs.Config{
			Allow: logs.ALL,

			Outs: []*logs.Out{
				{
					Target: os.Stdout,
					Color:  true,
				},
			},
		})

		config, err := getServerConfig()

		if err != nil {
			logger.Err(err)
			os.Exit(1)
		}

		server, err := startServer(config)

		if err != nil {
			logger.Err(err)
			os.Exit(1)
		}

		logger.Inff("Server started: %s://%s\n", server.Net(), server.Addr())

		server.Wait()

	case "client", "c":
		config, err := getClientConfig()

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		startClient(config)

	case "generate", "g":
		if err := generateConfig(config.Samples); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

	case "help", "h":
		fmt.Println(cli_doc)

	default:
		fmt.Println("Unknown command: " + cmd)
	}
}
