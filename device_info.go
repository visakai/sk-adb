package adb

import (
	"bufio"
	"strings"

	"github.com/matt-e/go-adb/internal/errors"
)

type DeviceInfo struct {
	// Always set.
	Serial string
	State  DeviceState

	// Product, device, and model are not set in the short form.
	Product    string
	Model      string
	DeviceInfo string

	// Only set for devices connected via USB.
	Usb string

	// Set to 'device' usually, but have observed 'unauthorized' or 'offline'
	Flag string
}

// IsUsb returns true if the device is connected via USB.
func (d *DeviceInfo) IsUsb() bool {
	return d.Usb != ""
}

func newDevice(serial string, state DeviceState, attrs map[string]string) (*DeviceInfo, error) {
	if serial == "" {
		return nil, errors.AssertionErrorf("device serial cannot be blank")
	}

	return &DeviceInfo{
		Serial:     serial,
		State:      state,
		Product:    attrs["product"],
		Model:      attrs["model"],
		DeviceInfo: attrs["device"],
		Usb:        attrs["usb"],
	}, nil
}

func parseDeviceList(list string, lineParseFunc func(string) (*DeviceInfo, error)) ([]*DeviceInfo, error) {
	devices := []*DeviceInfo{}
	scanner := bufio.NewScanner(strings.NewReader(list))

	for scanner.Scan() {
		device, err := lineParseFunc(scanner.Text())
		if err != nil {
			return nil, err
		}
		devices = append(devices, device)
	}

	return devices, nil
}

func parseDeviceShort(line string) (*DeviceInfo, error) {
	fields := strings.Fields(line)
	if len(fields) != 2 {
		return nil, errors.Errorf(errors.ParseError,
			"malformed device line, expected 2 fields but found %d", len(fields))
	}

	state, err := parseDeviceState(fields[1])
	if err != nil {
		return nil, err
	}

	return newDevice(fields[0], state, map[string]string{})
}

func parseDeviceLong(line string) (*DeviceInfo, error) {
	fields := strings.Fields(line)
	if len(fields) < 2 {
		return nil, errors.Errorf(errors.ParseError,
			"malformed device line, expected at least 5 fields but found %d", len(fields))
	}

	state, err := parseDeviceState(fields[1])
	if err != nil {
		return nil, err
	}

	attrs := parseDeviceAttributes(fields[2:])
	return newDevice(fields[0], state, attrs)
}

// samples of parameter fields passed in
// [product:PRODUCT model:MODEL device:DEVICE] for emulators
// [usb:1234 product:PRODUCT model:MODEL device:DEVICE] for physical device, adb usb
// [2-1.1 product:PRODUCT model:MODEL device:DEVICE] for physical device, libusb
func parseDeviceAttributes(fields []string) map[string]string {
	attrs := map[string]string{}
	if len(fields) == 0 {
		return attrs
	}
	if !strings.HasPrefix(fields[0], "product:") {
		usb := fields[0] // libusb's USB backend implementation
		if strings.HasPrefix(usb, "usb:") {
			usb = usb[4:] // ADB's USB backend implementation
		}
		attrs["usb"] = usb
		fields = fields[1:]
	}
	for _, field := range fields {
		key, val := parseKeyVal(field)
		attrs[key] = val
	}
	return attrs
}

// Parses a key:val pair and returns key, val.
func parseKeyVal(pair string) (string, string) {
	split := strings.Split(pair, ":")
	return split[0], split[1]
}
