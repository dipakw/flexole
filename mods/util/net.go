package util

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

func NetAddr(addr string, defaultPort string, minPort int, maxPort int) (string, error) {
	host, port, err := net.SplitHostPort(addr)

	if err != nil {
		if strings.Contains(err.Error(), "missing port") || port == "" {
			host = strings.TrimSuffix(addr, ":")
			port = defaultPort
		} else {
			return "", err
		}
	}

	if port == "" {
		port = defaultPort
	}

	portInt, err := strconv.Atoi(port)

	if err != nil {
		return "", fmt.Errorf("port \"%s\" is not numeric", port)
	}

	if minPort == 0 {
		minPort = 1
	}

	if maxPort == 0 {
		maxPort = 65535
	}

	if portInt < minPort || portInt > maxPort {
		return "", fmt.Errorf("port %d is not between %d and %d", portInt, minPort, maxPort)
	}

	return net.JoinHostPort(host, port), nil
}

func NetAddrDefault(network string, addr string, defaultHost string) (string, error) {
	if network != "tcp" && network != "udp" {
		return addr, nil
	}

	// Check if the addr a number.
	maybePort, err := strconv.Atoi(addr)

	if err == nil {
		if maybePort < 1 || maybePort > 65535 {
			return "", fmt.Errorf("invalid port: %d", maybePort)
		}

		return net.JoinHostPort(defaultHost, strconv.Itoa(maybePort)), nil
	}

	host, port, err := net.SplitHostPort(addr)

	if err != nil {
		return "", err
	}

	return net.JoinHostPort(host, port), nil
}

func StrOr(str string, or string) string {
	if str == "" {
		return or
	}

	return str
}
