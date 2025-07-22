package app

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

func getClientConfig() (*ClientConfig, error) {
	cli := NewCli(map[string]string{
		"quick":   "quick",
		"log":     "iwe",
		"config":  "client.yml",
		"encrypt": "1",
	})

	args := cli.Gets(
		"config",
		"quick",
		"log",
		"local",
		"remote",
		"id",
		"encrypt",
	)

	if args["quick"].Passed && !args["config"].Passed {
		return prepareQuickClientConfig(args)
	}

	return loadClientConfigFile(args["config"].Value())
}

func loadClientConfigFile(file string) (*ClientConfig, error) {
	if file == "" {
		return nil, fmt.Errorf("No config file provided")
	}

	var config ClientConfig

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

func prepareQuickClientConfig(args map[string]*CliArg) (*ClientConfig, error) {
	quick, err := parseClientQuickArg(args["quick"].Value())

	if err != nil {
		return nil, err
	}

	id, err := parseClientIdArg(args["id"].Value())

	if err != nil {
		return nil, err
	}

	local, err := parseClientLocalArg(args["local"].Value())

	if err != nil {
		return nil, err
	}

	server, remote, err := parseClientRemoteArg(args["remote"].Value())

	if err != nil {
		return nil, err
	}

	config := &ClientConfig{
		Version: "1.0.0",

		Auth: &Auth{
			ID:  "quick",
			Key: quick,
		},

		Server: server,

		Pipes: []*Pipe{
			{
				ID:      fmt.Sprintf("pipe:%d", id),
				Enabled: true,
				Encrypt: args["encrypt"].Value() == "1",
			},
		},

		Services: []*Service{
			{
				ID:      id,
				Enabled: true,
				Local:   local,
				Remote:  remote,

				Pipes: []string{
					fmt.Sprintf("pipe:%d", id),
				},
			},
		},
	}

	return config, nil
}

func parseClientLocalArg(arg string) (*Addr, error) {
	if arg == "" {
		return nil, fmt.Errorf("Required argument --local is missing")
	}

	invalidFormat := fmt.Errorf(`Invalid argument --local, format: [protocol]/[address], example: tcp/localhost:8080`)
	parts := strings.SplitN(arg, "/", 2)

	if len(parts) != 2 {
		return nil, invalidFormat
	}

	addr := &Addr{
		Net:  parts[0],
		Addr: parts[1],
	}

	return addr, nil
}

func parseClientRemoteArg(arg string) (*Addr, *NetPort, error) {
	if arg == "" {
		return nil, nil, fmt.Errorf("Required argument --remote is missing")
	}

	invalidFormat := fmt.Errorf(`Invalid argument --remote, format: [protocol]/[port]@[server], example: tcp/80@192.168.1.100`)
	parts := strings.SplitN(arg, "@", 2)

	if len(parts) != 2 {
		return nil, nil, invalidFormat
	}

	netPort, serverAddr := parts[0], parts[1]

	netPortParts := strings.SplitN(netPort, "/", 2)

	if len(netPortParts) != 2 {
		return nil, nil, invalidFormat
	}

	addr := &Addr{
		Net:  "tcp",
		Addr: serverAddr,
	}

	port, err := strconv.Atoi(netPortParts[1])

	if err != nil || port < 0 || port > 65535 {
		return nil, nil, fmt.Errorf(`Invalid remote port '%s', it should be a number between 0 and 65535`, netPortParts[1])
	}

	netPortConf := &NetPort{
		Net:  netPortParts[0],
		Port: uint16(port),
	}

	return addr, netPortConf, nil
}

func parseClientIdArg(arg string) (uint16, error) {
	if arg == "" {
		return 0, fmt.Errorf("Required argument --id is missing")
	}

	invalidFormat := fmt.Errorf(`Invalid argument --id, it should be a number between 0 and 65535`)

	id, err := strconv.Atoi(arg)

	if err != nil || id < 0 || id > 65535 {
		return 0, invalidFormat
	}

	return uint16(id), nil
}

func parseClientQuickArg(arg string) (string, error) {
	if arg == "" {
		return "", fmt.Errorf("Required argument --quick is missing")
	}

	return arg, nil
}
