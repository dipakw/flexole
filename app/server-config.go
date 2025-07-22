package app

import (
	"fmt"
	"net"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

func getServerConfig() (*ServerConfig, error) {
	cli := NewCli(map[string]string{
		"quick":  "quick",
		"log":    "iwe",
		"config": "server.yml",
		"host":   "0.0.0.0",
		"port":   DEFAULT_PORT,
	})

	args := cli.Gets(
		"quick",
		"log",
		"config",
		"host",
		"port",
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
	allowLogs := []string{}
	logConf := args["log"].Value()

	if logConf != "off" {
		if strings.Contains(logConf, "i") {
			allowLogs = append(allowLogs, "info")
		}

		if strings.Contains(logConf, "w") {
			allowLogs = append(allowLogs, "warn")
		}

		if strings.Contains(logConf, "e") {
			allowLogs = append(allowLogs, "error")
		}
	}

	config := &ServerConfig{
		Version: "1.0.0",

		Config: &Addr{
			Net:  "tcp",
			Addr: net.JoinHostPort(args["host"].Value(), args["port"].Value()),
		},

		Logs: &Logs{
			Allow: allowLogs,

			Outs: []LogOut{
				{
					To:    "stdout",
					Color: true,
				},
			},
		},

		Users: []User{
			{
				ID:       "quick",
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
