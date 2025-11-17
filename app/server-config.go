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
		"dir":    os.TempDir(),
	})

	args := cli.Gets(
		"quick",
		"log",
		"config",
		"host",
		"port",
		"user",
		"dir",
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

	if err := yaml.Unmarshal(fileBytes, &config); err != nil {
		return nil, err
	}

	if err := normalizeAndValidateServerConfig(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

func prepareQuickServerConfig(args map[string]*CliArg) (*ServerConfig, error) {
	config := &ServerConfig{
		Version: "1.0.0",

		Bind: &Addr{
			Net:  "tcp",
			Addr: net.JoinHostPort(args["host"].Value(), args["port"].Value()),
		},

		Dir: args["dir"].Value(),

		Logs: &Logs{
			Allow: util.LogShortToKinds(args["log"].Value()),

			Outs: []*LogOut{
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

func normalizeAndValidateServerConfig(config *ServerConfig) error {
	if config.Bind == nil {
		config.Bind = &Addr{
			Net:  "tcp",
			Addr: net.JoinHostPort(DEFAULT_HOST, DEFAULT_PORT),
		}
	}

	if config.Logs == nil {
		config.Logs = &Logs{
			Allow: []string{"info", "warn", "error"},
			Outs: []*LogOut{
				{
					To:    "stdout",
					Color: true,
				},
			},
		}
	}

	if config.Bind.Net == "" {
		config.Bind.Net = "tcp"
	}

	if config.Bind.Addr == "" {
		config.Bind.Addr = net.JoinHostPort(DEFAULT_HOST, DEFAULT_PORT)
	}

	if config.Bind.Net != "tcp" && config.Bind.Net != "unix" {
		return fmt.Errorf("server only supports tcp and unix networks, got: %s", config.Bind.Net)
	}

	if config.Bind.Net == "tcp" {
		var err error = nil
		config.Bind.Addr, err = util.NetAddr(config.Bind.Addr, DEFAULT_PORT, 0, 0)

		if err != nil {
			return fmt.Errorf("invalid server address: %w", err)
		}
	}

	return nil
}
