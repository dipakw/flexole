package app

import (
	"flexole/mods/util"
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
		"user":    "quick",
	})

	args := cli.Gets(
		"config",
		"quick",
		"log",
		"local",
		"remote",
		"id",
		"encrypt",
		"user",
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
	quick, err := parseClientQuickArg(args["quick"])

	if err != nil {
		return nil, err
	}

	id, err := parseClientIdArg(args["id"])

	if err != nil {
		return nil, err
	}

	local, err := parseClientLocalArg(args["local"])

	if err != nil {
		return nil, err
	}

	server, remote, err := parseClientRemoteArg(args["remote"])

	if err != nil {
		return nil, err
	}

	config := &ClientConfig{
		Version: "1.0.0",

		Auth: &Auth{
			ID:  args["user"].Value(),
			Key: quick,
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

func parseClientLocalArg(arg *CliArg) (*Addr, error) {
	val := arg.Value()

	if val == "" {
		return nil, fmt.Errorf("Required argument %s is missing", arg.Name)
	}

	invalidFormat := fmt.Errorf(`Invalid argument %s, format: [protocol]/[address], example: tcp/localhost:8080`, arg.Name)
	parts := strings.SplitN(val, "/", 2)

	if len(parts) != 2 {
		return nil, invalidFormat
	}

	addr := &Addr{
		Net:  parts[0],
		Addr: parts[1],
	}

	return addr, nil
}

func parseClientRemoteArg(arg *CliArg) (*Addr, *NetPort, error) {
	val := arg.Value()

	if val == "" {
		return nil, nil, fmt.Errorf("Required argument %s is missing", arg.Name)
	}

	invalidFormat := fmt.Errorf(`Invalid argument %s, format: [protocol]/[port]@[server], example: tcp/80@192.168.1.100`, arg.Name)
	parts := strings.SplitN(val, "@", 2)

	if len(parts) != 2 {
		return nil, nil, invalidFormat
	}

	netPort, serverAddr := parts[0], parts[1]

	netPortParts := strings.SplitN(netPort, "/", 2)

	if len(netPortParts) != 2 {
		return nil, nil, invalidFormat
	}

	serverAddr, err := util.NetAddr(serverAddr, DEFAULT_PORT)

	if err != nil {
		return nil, nil, fmt.Errorf("Invalid server address: %s", err.Error())
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

func parseClientIdArg(arg *CliArg) (uint16, error) {
	val := arg.Value()

	if val == "" {
		return 0, fmt.Errorf("Required argument %s is missing", arg.Name)
	}

	invalidFormat := fmt.Errorf(`Invalid argument %s, it should be a number between 0 and 65535`, arg.Name)

	id, err := strconv.Atoi(val)

	if err != nil || id < 0 || id > 65535 {
		return 0, invalidFormat
	}

	return uint16(id), nil
}

func parseClientQuickArg(arg *CliArg) (string, error) {
	val := arg.Value()

	if val == "" {
		return "", fmt.Errorf("Required argument %s is missing", arg.Name)
	}

	return val, nil
}
