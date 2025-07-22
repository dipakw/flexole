package app

import (
	"fmt"
	"os"
)

func Run() {
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

		startServer(config)

	case "client", "c":
		config, err := getClientConfig()

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		startClient(config)

	case "help", "h":
		fmt.Println(cli_doc)

	default:
		fmt.Println("Unknown command: " + cmd)
	}
}
