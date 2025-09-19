//go:build !darwin

package main

import (
	"fmt"
	"go.bug.st/serial/enumerator"
	"slices"
)

func findKeyboardDev() (string, error) {
	ports, err := enumerator.GetDetailedPortsList()
	if err != nil {
		return "", err
	}

	if len(ports) == 0 {
		return "", fmt.Errorf("dygma defy or raise 2 keyboard not found")
	}

	for _, port := range ports {
		if slices.Contains(vendorIds, port.VID) {
			return port.Name, err
		}
	}

	return "", fmt.Errorf("no device found")
}
