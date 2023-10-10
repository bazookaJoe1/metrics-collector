package net_params_validator

import (
	"fmt"
	"net"
	"strconv"
)

func ValidatePort(port string) error {
	intPort, err := strconv.ParseInt(port, 10, 64)
	if err != nil {
		return err
	}

	if intPort <= 0 || intPort > 65535 {
		return fmt.Errorf("invalid port number: %s", port)
	}

	return nil
}

func ValidateIP(address string) error {
	ip := net.ParseIP(address)
	if ip == nil && address != "localhost" {
		return fmt.Errorf("invalid ip: %s", address)
	}
	return nil
}
