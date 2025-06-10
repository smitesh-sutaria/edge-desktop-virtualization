package resources

import (
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"device-plugin/pkg/constants"
	pluginapi "k8s.io/kubelet/pkg/apis/deviceplugin/v1beta1"
)

// Helper to set up test constants
func setupVfioTestConstants(t *testing.T, baseDir string) {
	constants.VfioResourceName = "test.vfio/resource"
	constants.VfioSocketPath = filepath.Join(baseDir, "vfio.sock")
	constants.VfioDevicePath = filepath.Join(baseDir, "vfio")
	constants.VfioDeviceCount = 2
	constants.VfioDevicePrefix = "vfio"
	constants.Permissions = "rw"
}

func TestVfioGetResourceName(t *testing.T) {
	baseDir := t.TempDir()
	setupVfioTestConstants(t, baseDir)
	r := &ResourceVFIO{}
	if got := r.GetResourceName(); got != constants.VfioResourceName {
		t.Errorf("GetResourceName() = %v, want %v", got, constants.VfioResourceName)
	}
}

func TestVfioGetSocketPath(t *testing.T) {
	baseDir := t.TempDir()
	setupVfioTestConstants(t, baseDir)
	r := &ResourceVFIO{}
	if got := r.GetSocketPath(); got != constants.VfioSocketPath {
		t.Errorf("GetSocketPath() = %v, want %v", got, constants.VfioSocketPath)
	}
}

func TestVfioListDevices_PathNotExist(t *testing.T) {
	baseDir := t.TempDir()
	setupVfioTestConstants(t, baseDir)
	os.RemoveAll(constants.VfioDevicePath) // Ensure path does not exist
	r := &ResourceVFIO{}
	devices := r.ListDevices()
	if len(devices) != 0 {
		t.Errorf("ListDevices() when path does not exist, want 0 devices, got %d", len(devices))
	}
}

func TestVfioListDevices_ReturnsDevices(t *testing.T) {
	baseDir := t.TempDir()
	setupVfioTestConstants(t, baseDir)
	os.MkdirAll(constants.VfioDevicePath, 0755)
	r := &ResourceVFIO{}
	devices := r.ListDevices()
	if len(devices) != constants.VfioDeviceCount {
		t.Errorf("ListDevices() = %d devices, want %d", len(devices), constants.VfioDeviceCount)
	}
	for i, dev := range devices {
		wantID := constants.VfioDevicePrefix + "-" + strconv.Itoa(i)
		if dev.ID != wantID {
			t.Errorf("Device ID = %v, want %v", dev.ID, wantID)
		}
		if dev.Health != pluginapi.Healthy {
			t.Errorf("Device Health = %v, want %v", dev.Health, pluginapi.Healthy)
		}
	}
}

func TestVfioAllocateDevicesFromVfioFolder_Success(t *testing.T) {
	baseDir := t.TempDir()
	setupVfioTestConstants(t, baseDir)
	os.MkdirAll(constants.VfioDevicePath, 0755)
	// Create fake device files
	for i := 0; i < 2; i++ {
		f, err := os.Create(filepath.Join(constants.VfioDevicePath, "dev"+strconv.Itoa(i)))
		if err != nil {
			t.Fatalf("Failed to create device file: %v", err)
		}
		f.Close()
	}
	responses, err := AllocateDevicesFromVfioFolder()
	if err != nil {
		t.Fatalf("AllocateDevicesFromVfioFolder() error = %v, want nil", err)
	}
	if len(responses) != 1 {
		t.Errorf("AllocateDevicesFromVfioFolder() = %d responses, want 1", len(responses))
	}
	if len(responses[0].Devices) != 2 {
		t.Errorf("AllocateDevicesFromVfioFolder() = %d devices, want 2", len(responses[0].Devices))
	}
	for _, dev := range responses[0].Devices {
		if dev.Permissions != constants.Permissions {
			t.Errorf("Device.Permissions = %v, want %v", dev.Permissions, constants.Permissions)
		}
	}
}

func TestVfioAllocateDevicesFromVfioFolder_DirNotExist(t *testing.T) {
	baseDir := t.TempDir()
	setupVfioTestConstants(t, baseDir)
	os.RemoveAll(constants.VfioDevicePath) // Ensure path does not exist
	_, err := AllocateDevicesFromVfioFolder()
	if err == nil {
		t.Errorf("AllocateDevicesFromVfioFolder() error = nil, want error")
	}
}

func TestVfioAllocate_Success(t *testing.T) {
	baseDir := t.TempDir()
	setupVfioTestConstants(t, baseDir)
	os.MkdirAll(constants.VfioDevicePath, 0755)
	// Create fake device files
	for i := 0; i < 2; i++ {
		f, err := os.Create(filepath.Join(constants.VfioDevicePath, "dev"+strconv.Itoa(i)))
		if err != nil {
			t.Fatalf("Failed to create device file: %v", err)
		}
		f.Close()
	}
	r := &ResourceVFIO{}
	resp := r.Allocate([]string{"vfio-0"})
	if len(resp) == 0 || len(resp[0].Devices) != 2 {
		t.Errorf("Allocate() = %+v, want 2 devices", resp)
	}
}

func TestVfioAllocate_ErrorPath(t *testing.T) {
	baseDir := t.TempDir()
	setupVfioTestConstants(t, baseDir)
	os.RemoveAll(constants.VfioDevicePath) // Ensure path does not exist
	r := &ResourceVFIO{}
	resp := r.Allocate([]string{"vfio-0"})
	if len(resp) == 0 || len(resp[0].Mounts) == 0 {
		t.Errorf("Allocate() error path, want mount fallback, got %+v", resp)
	}
	if resp[0].Mounts[0].HostPath != constants.VfioDevicePath {
		t.Errorf("Allocate() mount HostPath = %v, want %v", resp[0].Mounts[0].HostPath, constants.VfioDevicePath)
	}
}
