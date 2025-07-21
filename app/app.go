package app

import (
	"os"
)

func Run() {
	if len(os.Args) < 2 {
		showHelp("No command provided")
		return
	}

	cmd, opts := cliArgs()

	switch cmd {
	case "server":
		config, err := getServerConfig(opts)

		if err != nil {
			showHelp(err.Error())
			return
		}

		startServer(config)

	case "client":
		config, err := getClientConfig(opts)

		if err != nil {
			showHelp(err.Error())
			return
		}

		startClient(config)

	case "help":
		showHelp("")

	default:
		showHelp("Unknown command: " + cmd)
	}
}
