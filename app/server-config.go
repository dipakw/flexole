package app

import (
	"flexole/mods/util"
	"fmt"
	"net"
	"os"

	"gopkg.in/yaml.v3"
)

func getServerConfig() (*ServerConfig, error) {
	cli := NewCli(map[string]string{
		"quick":  "quick",
		"log":    "iwe",
		"config": "server.yml",
		"host":   "0.0.0.0",
		"port":   DEFAULT_PORT,
		"user":   "quick",
	})

	args := cli.Gets(
		"quick",
		"log",
		"config",
		"host",
		"port",
		"user",
	)

	if args["quick"].Passed && !args["config"].Passed {
		return prepareQuickServerConfig(args)
	}

	return loadServerConfigFile(args["config"].Value())
}

func loadServerConfigFile(file string) (*ServerConfig, error) {
	if file == "" {
		return nil, fmt.Errorf("No config file provided")
	}

	var config ServerConfig

	fileBytes, err := os.ReadFile(file)

	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(fileBytes, &config)

	if err != nil {
		return nil, err
	}

	return &config, nil
}

func prepareQuickServerConfig(args map[string]*CliArg) (*ServerConfig, error) {
	config := &ServerConfig{
		Version: "1.0.0",

		Config: &Addr{
			Net:  "tcp",
			Addr: net.JoinHostPort(args["host"].Value(), args["port"].Value()),
		},

		Logs: &Logs{
			Allow: util.LogShortToKinds(args["log"].Value()),

			Outs: []LogOut{
				{
					To:    "stdout",
					Color: true,
				},
			},
		},

		Users: []User{
			{
				ID:       args["user"].Value(),
				Enabled:  true,
				Key:      args["quick"].Value(),
				MaxPipes: 15,
				MaxServices: MaxServices{
					Unix: 5,
					TCP:  5,
					UDP:  5,
				},
			},
		},
	}

	return config, nil
}
