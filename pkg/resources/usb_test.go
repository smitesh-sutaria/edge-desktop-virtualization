package resources

import (
	"device-plugin/pkg/constants"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"testing"
)

const (
	testUsbResourceName = "test.usb/resource"
	testUsbSocketPath   = "/tmp/test-usb.sock"
	testUsbDevicePath   = "/tmp/test-usb-devices"
	testUsbDeviceCount  = 2
	testUsbDevicePrefix = "usb"
	testPermissions     = "rw"
)

func setupUsbTestConstants() {
	constants.UsbResourceName = testUsbResourceName
	constants.UsbSocketPath = testUsbSocketPath
	constants.UsbDevicePath = testUsbDevicePath
	constants.UsbDeviceCount = testUsbDeviceCount
	constants.UsbDevicePrefix = testUsbDevicePrefix
	constants.Permissions = testPermissions
}

func TestUsbGetResourceName(t *testing.T) {
	setupUsbTestConstants()
	r := &ResourceUSB{}
	if got := r.GetResourceName(); got != testUsbResourceName {
		t.Errorf("GetResourceName() = %v, want %v", got, testUsbResourceName)
	}
}

func TestUsbGetSocketPath(t *testing.T) {
	setupUsbTestConstants()
	r := &ResourceUSB{}
	if got := r.GetSocketPath(); got != testUsbSocketPath {
		t.Errorf("GetSocketPath() = %v, want %v", got, testUsbSocketPath)
	}
}

func TestUsbListDevices_PathNotExist(t *testing.T) {
	setupUsbTestConstants()
	os.RemoveAll(testUsbDevicePath) // Ensure path does not exist
	r := &ResourceUSB{}
	devices := r.ListDevices()
	if len(devices) != 0 {
		t.Errorf("ListDevices() when path does not exist, want 0 devices, got %d", len(devices))
	}
}

func TestUsbListDevices_ReturnsDevices(t *testing.T) {
	setupUsbTestConstants()
	os.MkdirAll(testUsbDevicePath, 0755)
	defer os.RemoveAll(testUsbDevicePath)

	r := &ResourceUSB{}
	devices := r.ListDevices()
	if len(devices) != testUsbDeviceCount {
		t.Errorf("ListDevices() = %d devices, want %d", len(devices), testUsbDeviceCount)
	}
	for i, dev := range devices {
		wantID := testUsbDevicePrefix + "-" + strconv.Itoa(i)
		if dev.ID != wantID {
			t.Errorf("Device ID = %v, want %v", dev.ID, wantID)
		}
	}
}

func TestUsbAllocate_Success(t *testing.T) {
	setupUsbTestConstants()
	os.MkdirAll(filepath.Join(testUsbDevicePath, "001"), 0755)
	defer os.RemoveAll(testUsbDevicePath)

	deviceFile := filepath.Join(testUsbDevicePath, "001", "dev1")
	ioutil.WriteFile(deviceFile, []byte{}, 0644)

	r := &ResourceUSB{}
	resp := r.Allocate([]string{"usb-0"})
	if len(resp) == 0 || len(resp[0].Devices) == 0 {
		t.Errorf("Allocate() = %+v, want at least one device", resp)
	}
}

func TestUsbAllocate_ErrorPath(t *testing.T) {
	setupUsbTestConstants()
	os.RemoveAll(testUsbDevicePath)

	r := &ResourceUSB{}
	resp := r.Allocate([]string{"usb-0"})
	if len(resp) == 0 || len(resp[0].Mounts) == 0 {
		t.Errorf("Allocate() error path, want mount fallback, got %+v", resp)
	}
	if resp[0].Mounts[0].HostPath != testUsbDevicePath {
		t.Errorf("Allocate() mount HostPath = %v, want %v", resp[0].Mounts[0].HostPath, testUsbDevicePath)
	}
}
