package util

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

func NetAddr(addr string, defaultPort string) (string, error) {
	host, port, err := net.SplitHostPort(addr)

	if err != nil {
		if strings.Contains(err.Error(), "missing port") || port == "" {
			host = addr
			port = defaultPort
		} else {
			return "", err
		}
	}

	return net.JoinHostPort(host, port), nil
}

func NetAddrDefault(network string, addr string, defaultHost string) (string, error) {
	if network == "tcp" || network == "udp" {
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
