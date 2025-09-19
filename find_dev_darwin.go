//go:build darwin

package main

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"slices"
	"strings"
)

//WARNING:  THIS CODE IS NOT MINE!!! THIS SECTION IS ENTIRRELY VIDE-CODED ...
// YES AI SPIT THIS THING OUT AND I AHVE NO CLUE ABOUT IT AS I HAVE NO MAC AND HAVE NO CLUE HOW MAC WORKS!!! JUST TRYING MY BEST TO CROSS COMPILE IT FOR ALL SYSTEMS!

// usbDevice struct now correctly models the potentially nested JSON data.
// The SubItems field will capture the array in the "_items" key.
type usbDevice struct {
	VendorID  string      `json:"vendor_id"`
	SerialNum string      `json:"serial_num"`
	Product   string      `json:"_name"`
	SubItems  []usbDevice `json:"_items,omitempty"`
}

type systemProfilerOutput struct {
	SPUSBDataType []usbDevice `json:"SPUSBDataType"`
}

// findKeyboardDev finds the keyboard on macOS by shelling out to system_profiler.
func findKeyboardDev() (string, error) {
	cmd := exec.Command("system_profiler", "SPUSBDataType", "-json")
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to run system_profiler: %w", err)
	}

	var data systemProfilerOutput
	if err := json.Unmarshal(out, &data); err != nil {
		return "", fmt.Errorf("failed to parse system_profiler JSON: %w", err)
	}

	// Create a single, flat list of all devices and sub-devices.
	var allDevices []usbDevice
	for _, device := range data.SPUSBDataType {
		allDevices = append(allDevices, device)
		if len(device.SubItems) > 0 {
			// Use '...' to append all elements of the sub-slice.
			allDevices = append(allDevices, device.SubItems...)
		}
	}

	// Now search through the flattened list.
	for _, d := range allDevices {
		// Get the vendor ID like "0x1209" and remove the "0x" prefix.
		rawVid := strings.Split(d.VendorID, " ")[0]
		cleanVid := strings.TrimPrefix(rawVid, "0x")

		// Check if the cleaned ID is in our shared 'vendorIds' slice.
		if slices.Contains(vendorIds, cleanVid) {
			// The serial number field for USB modems on macOS often contains the device path.
			if strings.HasPrefix(d.SerialNum, "/dev/tty.usbmodem") {
				return d.SerialNum, nil
			}
		}
	}

	return "", fmt.Errorf("dygma defy or raise 2 keyboard not found")
}
