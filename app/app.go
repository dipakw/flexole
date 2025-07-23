package app

import (
	"embed"
	"fmt"
	"os"
)

func Run(samples *embed.FS) {
	if len(os.Args) < 2 {
		fmt.Println(cli_doc)
		return
	}

	cmd := os.Args[1]

	switch cmd {
	case "server", "s":
		config, err := getServerConfig()

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		server, err := startServer(config)

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Printf("Server started: %s://%s\n", server.Net(), server.Addr())

		server.Wait()

	case "client", "c":
		config, err := getClientConfig()

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		startClient(config)

	case "generate", "g":
		if err := generateConfig(samples); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

	case "help", "h":
		fmt.Println(cli_doc)

	default:
		fmt.Println("Unknown command: " + cmd)
	}
}
