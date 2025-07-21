package app

import (
	"fmt"
	"net"
	"os"
	"strconv"

	"gopkg.in/yaml.v3"
)

func getServerConfig(opts map[string]string) (*ServerConfig, error) {
	configFile, configPassed := opts["config"]

	if !configPassed {
		configFile = "server.yml"
	}

	quickCode, quickPassed := opts["quick"]

	if !quickPassed {
		configPassed = true
	}

	isQuick := quickPassed && !configPassed

	if !isQuick && configFile == "" {
		return nil, fmt.Errorf("no config file provided")
	}

	if isQuick && quickCode == "" {
		quickCode = "quick"
		opts["quick"] = quickCode
	}

	if isQuick {
		port, ok := opts["port"]

		if !ok || port == "" {
			port = "8887"
			opts["port"] = port
		}

		portInt, err := strconv.Atoi(port)

		if err != nil || portInt < 1 || portInt > 65535 {
			return nil, fmt.Errorf("invalid port")
		}
	}

	if isQuick {
		return getServerQuickConfig(opts)
	}

	return readServerConfig(configFile)
}

func readServerConfig(file string) (*ServerConfig, error) {
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

func getServerQuickConfig(opts map[string]string) (*ServerConfig, error) {
	config := &ServerConfig{
		Version: "1.0.0",

		Config: Addr{
			Net:  "tcp",
			Addr: net.JoinHostPort("0.0.0.0", opts["port"]),
		},

		Logs: Logs{
			Allow: []string{"info", "warn", "error"},

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
				Key:      opts["quick"],
				MaxPipes: 1,
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
