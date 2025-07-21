package app

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

func getClientConfig(opts map[string]string) (*ClientConfig, error) {
	configFile, configPassed := opts["config"]

	if !configPassed {
		configFile = "client.yml"
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
		return getClientQuickConfig(opts)
	}

	return readClientConfig(configFile)
}

func readClientConfig(file string) (*ClientConfig, error) {
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

func getClientQuickConfig(opts map[string]string) (*ClientConfig, error) {
	local, ok := opts["local"]

	if !ok {
		return nil, fmt.Errorf("no local provided")
	}

	// Split at the first /
	localParts := strings.SplitN(local, "#", 2)

	if len(localParts) != 2 {
		return nil, fmt.Errorf("invalid local format")
	}

	localNetAddr := localParts[0]
	localId := localParts[1]

	if localNetAddr == "" {
		return nil, fmt.Errorf("no local net/addr provided")
	}

	if localId == "" {
		return nil, fmt.Errorf("no local id provided")
	}

	localParts = strings.SplitN(localNetAddr, "/", 2)

	if len(localParts) != 2 {
		return nil, fmt.Errorf("invalid local format")
	}

	localNet := localParts[0]
	localAddr := localParts[1]

	if localNet != "tcp" && localNet != "udp" && localNet != "unix" {
		return nil, fmt.Errorf("invalid local net: %s", localNet)
	}

	id, err := strconv.Atoi(localId)

	if err != nil || id < 1 || id > 65535 {
		return nil, fmt.Errorf("invalid local id: %s", localId)
	}

	remote, ok := opts["remote"]

	if !ok {
		return nil, fmt.Errorf("no remote provided")
	}

	remoteParts := strings.SplitN(remote, "@", 2)

	if len(remoteParts) != 2 {
		return nil, fmt.Errorf("invalid remote format")
	}

	netPort := remoteParts[0]
	remoteAddr := remoteParts[1]

	netPortParts := strings.Split(netPort, "/")

	if len(netPortParts) != 2 {
		return nil, fmt.Errorf("invalid remote format")
	}

	remoteNet := netPortParts[0]
	remotePort, err := strconv.Atoi(netPortParts[1])

	if remoteNet != "tcp" && remoteNet != "udp" && remoteNet != "unix" {
		return nil, fmt.Errorf("invalid remote net: %s", remoteNet)
	}

	if err != nil {
		return nil, fmt.Errorf("invalid remote port")
	}

	if localNet == "udp" && remoteNet != "udp" {
		return nil, fmt.Errorf("local udp can only be used with udp remote")
	}

	// Check if the remote addr has a port
	serverHost, serverPort, err := net.SplitHostPort(remoteAddr)

	if strings.HasSuffix(remoteAddr, ":") || serverPort == "0" {
		return nil, fmt.Errorf("invalid remote addr: %s", remoteAddr)
	}

	if err != nil {
		serverPort = DEFAULT_PORT
		err = nil
	}

	if err != nil {
		return nil, fmt.Errorf("invalid remote addr: %s", remoteAddr)
	}

	config := &ClientConfig{
		Version: "1.0.0",

		Auth: Auth{
			ID:  "quick",
			Key: opts["quick"],

			Server: ServerAddr{
				Net:  "tcp",
				Addr: net.JoinHostPort(serverHost, serverPort),
			},
		},

		Pipes: []Pipe{
			{
				ID:      fmt.Sprintf("pipe:%d", id),
				Enabled: true,
				Encrypt: true,
			},
		},

		Services: []Service{
			{
				ID:      uint16(id),
				Enabled: true,

				Local: Addr{
					Net:  localNet,
					Addr: localAddr,
				},

				Remote: RemoteConf{
					Net:  remoteNet,
					Port: remotePort,
					Pipes: []string{
						fmt.Sprintf("pipe:%d", id),
					},
				},
			},
		},
	}

	return config, nil
}
