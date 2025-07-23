package util

import (
	"net"
	"strings"
)

func NetAddr(addr string, defaultPort string) (string, error) {
	host, port, err := net.SplitHostPort(addr)

	if err != nil {
		if strings.Contains(err.Error(), "missing port") || port == "" {
			port = defaultPort
		} else {
			return "", err
		}
	}

	return net.JoinHostPort(host, port), nil
}
