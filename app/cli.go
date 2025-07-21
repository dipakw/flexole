package app

import (
	"fmt"
	"os"
	"strings"
)

var cli_doc = strings.TrimSpace(`
Usage:
  flexole <command> [options]
`)

func cliArgs() (string, map[string]string) {
	args := os.Args[1:]
	if len(args) == 0 {
		return "", map[string]string{}
	}

	main := args[0]
	opts := map[string]string{}

	for i := 1; i < len(args); i++ {
		arg := args[i]
		if strings.HasPrefix(arg, "--") {
			arg = arg[2:]
			parts := strings.SplitN(arg, "=", 2)
			key := parts[0]
			value := ""
			if len(parts) == 2 {
				value = parts[1]
			}
			opts[key] = value
		}
	}

	return main, opts
}

func showHelp(message string) {
	if message != "" {
		fmt.Println(message)
		fmt.Println()
	}

	fmt.Println(cli_doc)
}
