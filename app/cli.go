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

func NewCli(defaultOpts map[string]string) *Cli {
	args := os.Args[1:]

	main := ""
	opts := map[string]string{}

	if len(args) > 0 {
		main = args[0]

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
	input, passed := c.opts[key]

	return &CliArg{
		Passed:  passed,
		Input:   input,
		Default: c.defaultOpts[key],
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
