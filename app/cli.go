package app

import (
	"fmt"
	"os"
	"strings"
)

var cli_doc = strings.TrimSpace(`
Usage:
  flexole <command> [options]

Commands:
  server, s      Start the server
  client, c      Start the client
  generate, g    Generate sample config files
  help, h        Show this help message

Options (server):
  --config, -c    Path to server config file (default: server.yml)
  --quick, -q     Quick start with user key (overrides config)
  --log, -o       Log levels: i=info, w=warn, e=error (default: iwe)
  --host, -h      Server host (default: 0.0.0.0)
  --port, -p      Server port (default: 8887)

Options (client):
  --config, -c     Path to client config file (default: client.yml)
  --quick, -q      Quick start with user key (overrides config)
  --log, -o        Log levels: i=info, w=warn, e=error (default: iwe)
  --local, -l      Local address, e.g. tcp/localhost:8080
  --remote, -r     Remote, e.g. tcp/80@192.168.1.100
  --id             Service ID (0-65535)
  --encrypt, -e    Enable encryption (default: 1)

Options (generate):
  --client-config, -cc  Output client config file (default: client.yml)
  --server-config, -sc  Output server config file (default: server.yml)

Examples:
  flexole server --config=server.yml
  flexole client --quick=mykey --local=tcp/127.0.0.1:8080 --remote=tcp/80@192.168.1.100 --id=1
  flexole generate --client-config=client.yml --server-config=server.yml
  flexole help

Notes:
  - For log levels, combine letters (e.g., "iwe" for info, warn, error).
  - All options can use either --long or -short forms.
`)

var parseArgs = map[string]bool{
	"--quick":   true,
	"--config":  true,
	"--log":     true,
	"--host":    true,
	"--port":    true,
	"--local":   true,
	"--remote":  true,
	"--id":      true,
	"--encrypt": true,

	"--client-config": true,
	"--server-config": true,
}

var parseArgsShort = map[string]bool{
	"-q": true,
	"-c": true,
	"-o": true,
	"-h": true,
	"-p": true,
	"-l": true,
	"-r": true,
	"-i": true,
	"-e": true,

	"-cc": true,
	"-sc": true,
}

var mapShortToLong = map[string]string{
	"-q": "--quick",
	"-c": "--config",
	"-o": "--log",
	"-h": "--host",
	"-p": "--port",
	"-l": "--local",
	"-r": "--remote",
	"-i": "--id",
	"-e": "--encrypt",

	"-cc": "--client-config",
	"-sc": "--server-config",
}

func NewCli(defaultOpts map[string]string) *Cli {
	args := os.Args[1:]

	main := ""
	opts := map[string]*ValName{}

	if len(args) > 0 {
		main = args[0]

		for i := 1; i < len(args); i++ {
			parts := strings.SplitN(args[i], "=", 2)
			key := parts[0]
			val := ""

			if !parseArgs[key] && !parseArgsShort[key] {
				continue
			}

			if len(parts) == 2 {
				val = parts[1]
			}

			optKey := key

			if parseArgsShort[key] {
				optKey = mapShortToLong[key]
			}

			opts[optKey[2:]] = &ValName{
				Val:  val,
				Name: key,
			}
		}
	}

	return &Cli{
		main:        main,
		opts:        opts,
		defaultOpts: defaultOpts,
	}
}

func (c *Cli) Help(message string) {
	if message != "" {
		fmt.Println(message)
		fmt.Println()
	}

	fmt.Println(cli_doc)
}

func (c *Cli) Get(key string) *CliArg {
	val := ""
	name := "--" + key

	input, passed := c.opts[key]

	if passed {
		name = input.Name
		val = input.Val
	}

	return &CliArg{
		Passed:  passed,
		Input:   val,
		Default: c.defaultOpts[key],
		Name:    name,
	}
}

func (c *Cli) Gets(keys ...string) map[string]*CliArg {
	args := map[string]*CliArg{}

	for _, key := range keys {
		args[key] = c.Get(key)
	}

	return args
}

func (c *CliArg) Value() string {
	if c.Passed && c.Input != "" {
		return c.Input
	}

	return c.Default
}
