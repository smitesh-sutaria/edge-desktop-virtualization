package resources

import (
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"device-plugin/pkg/constants"
	pluginapi "k8s.io/kubelet/pkg/apis/deviceplugin/v1beta1"
)

func setupIGPUTestConstants(t *testing.T, baseDir string) {
	constants.IGpuResourceName = "test.igpu/resource"
	constants.IGpuSocketPath = filepath.Join(baseDir, "igpu.sock")
	constants.IGpuDevicePath = filepath.Join(baseDir, "igpu")
	constants.IGpuDeviceCount = 2
	constants.IGpuDevicePrefix = "igpu"
	constants.RenderD128DevicePath = filepath.Join(baseDir, "renderD128")
	constants.Permissions = "rw"
}

func TestIGPU_GetResourceName(t *testing.T) {
	baseDir := t.TempDir()
	setupIGPUTestConstants(t, baseDir)
	r := &ResourceIGPU{}
	if got := r.GetResourceName(); got != constants.IGpuResourceName {
		t.Errorf("GetResourceName() = %v, want %v", got, constants.IGpuResourceName)
	}
}

func TestIGPU_GetSocketPath(t *testing.T) {
	baseDir := t.TempDir()
	setupIGPUTestConstants(t, baseDir)
	r := &ResourceIGPU{}
	if got := r.GetSocketPath(); got != constants.IGpuSocketPath {
		t.Errorf("GetSocketPath() = %v, want %v", got, constants.IGpuSocketPath)
	}
}

func TestIGPU_ListDevices_PathNotExist(t *testing.T) {
	baseDir := t.TempDir()
	setupIGPUTestConstants(t, baseDir)
	os.RemoveAll(constants.IGpuDevicePath) // Ensure path does not exist
	r := &ResourceIGPU{}
	devices := r.ListDevices()
	if len(devices) != 0 {
		t.Errorf("ListDevices() when path does not exist, want 0 devices, got %d", len(devices))
	}
}

func TestIGPU_ListDevices_ReturnsDevices(t *testing.T) {
	baseDir := t.TempDir()
	setupIGPUTestConstants(t, baseDir)
	os.MkdirAll(constants.IGpuDevicePath, 0755)
	r := &ResourceIGPU{}
	devices := r.ListDevices()
	if len(devices) != constants.IGpuDeviceCount {
		t.Errorf("ListDevices() = %d devices, want %d", len(devices), constants.IGpuDeviceCount)
	}
	for i, dev := range devices {
		wantID := constants.IGpuDevicePrefix + "-" + strconv.Itoa(i)
		if dev.ID != wantID {
			t.Errorf("Device ID = %v, want %v", dev.ID, wantID)
		}
		if dev.Health != pluginapi.Healthy {
			t.Errorf("Device Health = %v, want %v", dev.Health, pluginapi.Healthy)
		}
	}
}

func TestIGPU_Allocate(t *testing.T) {
	baseDir := t.TempDir()
	setupIGPUTestConstants(t, baseDir)
	r := &ResourceIGPU{}
	resp := r.Allocate([]string{"igpu-0"})
	if len(resp) != 1 {
		t.Fatalf("Allocate() = %d responses, want 1", len(resp))
	}
	devSpecs := resp[0].Devices
	if len(devSpecs) != 2 {
		t.Fatalf("Allocate() = %d device specs, want 2", len(devSpecs))
	}
	foundIGPU, foundRender := false, false
	for _, dev := range devSpecs {
		switch dev.HostPath {
		case constants.IGpuDevicePath:
			foundIGPU = true
			if dev.ContainerPath != constants.IGpuDevicePath {
				t.Errorf("IGPU ContainerPath = %v, want %v", dev.ContainerPath, constants.IGpuDevicePath)
			}
			if dev.Permissions != constants.Permissions {
				t.Errorf("IGPU Permissions = %v, want %v", dev.Permissions, constants.Permissions)
			}
		case constants.RenderD128DevicePath:
			foundRender = true
			if dev.ContainerPath != constants.RenderD128DevicePath {
				t.Errorf("RenderD128 ContainerPath = %v, want %v", dev.ContainerPath, constants.RenderD128DevicePath)
			}
			if dev.Permissions != constants.Permissions {
				t.Errorf("RenderD128 Permissions = %v, want %v", dev.Permissions, constants.Permissions)
			}
		default:
			t.Errorf("Unexpected HostPath: %v", dev.HostPath)
		}
	}
	if !foundIGPU {
		t.Error("Did not find device spec for IGpuDevicePath")
	}
	if !foundRender {
		t.Error("Did not find device spec for RenderD128DevicePath")
	}
}
