package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"slices"
	"strconv"
	"strings"

	"go.bug.st/serial"
	"go.bug.st/serial/enumerator"
)

var vendorIds = []string{"1209", "35ef"}

type batteryLevel struct {
	Left  int `json:"left"`
	Right int `json:"right"`
}

type WaybarOutput struct {
	Text       string `json:"text"`
	Tooltip    string `json:"tooltip"`
	Class      string `json:"class,omitempty"`
	Percentage int    `json:"percentage"`
}

func findKeyboardDev() (error, string) {
	ports, err := enumerator.GetDetailedPortsList()
	if err != nil {
		return err, ""
	}

	if len(ports) == 0 {
		return fmt.Errorf("dygma defy or raise 2 keyboard not found"), ""
	}

	for _, port := range ports {
		if slices.Contains(vendorIds, port.VID) {
			return nil, port.Name
		}
	}

	return fmt.Errorf("No device found"), ""
}

func readFromPort(ctx context.Context, port serial.Port, ch chan<- int, errCh chan<- error) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			buff := make([]byte, 4)
			n, err := port.Read(buff)
			if err != nil {
				if err.Error() != "EOF" {
					errCh <- fmt.Errorf("error reading from port: %w", err)
					return
				}
				continue
			}

			if n > 0 {
				response := strings.TrimSuffix(string(bytes.TrimSpace(buff[:n])), ".")
				if response == "" {
					continue
				}
				v, err := strconv.Atoi(response)
				if err != nil {
					errCh <- fmt.Errorf("failed to parse %q", response)
					return
				}
				ch <- v
			}
		}
	}
}

func getBatteryLevel(port serial.Port, side string, ch <-chan int, errCh <-chan error) (int, error) {
	command := "wireless.battery." + side + ".level\n"
	if _, err := port.Write([]byte(command)); err != nil {
		return 0, fmt.Errorf("failed to send command to keyboard: %v", err)
	}

	select {
	case level := <-ch:
		return level, nil
	case err := <-errCh:
		return 0, err
	}
}

func main() {
	err, dev := findKeyboardDev()
	if err != nil {
		log.Fatal("Could not find keyboard:", err)
	}
	mode := &serial.Mode{BaudRate: 9600}
	port, err := serial.Open(dev, mode)
	if err != nil {
		log.Fatal("failed to open port:", err)
	}
	defer port.Close()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ch := make(chan int)
	errCh := make(chan error, 1)
	go readFromPort(ctx, port, ch, errCh)

	battery := batteryLevel{}
	var level int
	var exit int
	level, err = getBatteryLevel(port, "left", ch, errCh)
	if err != nil {
		log.Print("failed to get left battery level:", err)
		exit = 1
	}
	battery.Left = level

	level, err = getBatteryLevel(port, "right", ch, errCh)
	if err != nil {
		log.Print("failed to get right battery level:", err)
		exit = 1
	}
	battery.Right = level

	lowestLevel := battery.Left
	if battery.Right < lowestLevel {
		lowestLevel = battery.Right
	}

	output := WaybarOutput{
		Text:       fmt.Sprintf("L:%d%% R:%d%%", battery.Left, battery.Right),
		Tooltip:    fmt.Sprintf("Left side: %d%%\rRight side: %d%%", battery.Left, battery.Right),
		Percentage: lowestLevel,
	}

	if battery.Left < 20 || battery.Right < 20 {
		output.Class = "critical"
	}

	jsonOutput, err := json.Marshal(output)
	if err != nil {
		log.Fatal("failed to marshal json:", err)
	}

	fmt.Println(string(jsonOutput))
	os.Exit(exit)
}
