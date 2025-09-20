//go:build darwin

package main

import (
	"fmt"
	"os/exec"
	"regexp"
)

// Darwin (MacOS) specific code to detect the keyboard device path.
func findKeyboardDev() (string, error) {
	cmd := exec.Command("ioreg", "-n", "DEFY", "-r", "-c", "IOUSBHostDevice", "-l")
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to run ioreg: %w", err)
	}

	if out == nil {
		return "", fmt.Errorf("dygma defy or raise 2 keyboard not found")
	}

	strout := string(out)
	re := regexp.MustCompile(`(usbmodem[^"]+)`)
	matches := re.FindStringSubmatch(strout)

	if len(matches) < 1 {
		return "", fmt.Errorf("dygma defy or raise 2 keyboard not found")
	}

	return fmt.Sprintf("/dev/cu.%s", matches[0]), nil
}
