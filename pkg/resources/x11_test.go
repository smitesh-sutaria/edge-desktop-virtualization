package resources

import (
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"device-plugin/pkg/constants"
	pluginapi "k8s.io/kubelet/pkg/apis/deviceplugin/v1beta1"
)

func setupX11TestConstants(t *testing.T, baseDir string) {
	constants.X11ResourceName = "test.x11/resource"
	constants.X11SocketPath = filepath.Join(baseDir, "x11.sock")
	constants.X11DevicePath = filepath.Join(baseDir, "x11")
	constants.X11DeviceCount = 2
	constants.X11DevicePrefix = "x11"
	constants.Permissions = "rw"
}

func TestX11_GetResourceName(t *testing.T) {
	baseDir := t.TempDir()
	setupX11TestConstants(t, baseDir)
	r := &ResourceX11{}
	if got := r.GetResourceName(); got != constants.X11ResourceName {
		t.Errorf("GetResourceName() = %v, want %v", got, constants.X11ResourceName)
	}
}

func TestX11_GetSocketPath(t *testing.T) {
	baseDir := t.TempDir()
	setupX11TestConstants(t, baseDir)
	r := &ResourceX11{}
	if got := r.GetSocketPath(); got != constants.X11SocketPath {
		t.Errorf("GetSocketPath() = %v, want %v", got, constants.X11SocketPath)
	}
}

func TestX11_ListDevices_PathNotExist(t *testing.T) {
	baseDir := t.TempDir()
	setupX11TestConstants(t, baseDir)
	os.RemoveAll(constants.X11DevicePath)
	r := &ResourceX11{}
	devices := r.ListDevices()
	if len(devices) != 0 {
		t.Errorf("ListDevices() when path does not exist, want 0 devices, got %d", len(devices))
	}
}

func TestX11_ListDevices_ReturnsDevices(t *testing.T) {
	baseDir := t.TempDir()
	setupX11TestConstants(t, baseDir)
	os.MkdirAll(constants.X11DevicePath, 0755)
	r := &ResourceX11{}
	devices := r.ListDevices()
	if len(devices) != constants.X11DeviceCount {
		t.Errorf("ListDevices() = %d devices, want %d", len(devices), constants.X11DeviceCount)
	}
	for i, dev := range devices {
		wantID := constants.X11DevicePrefix + "-" + strconv.Itoa(i)
		if dev.ID != wantID {
			t.Errorf("Device ID = %v, want %v", dev.ID, wantID)
		}
		if dev.Health != pluginapi.Healthy {
			t.Errorf("Device Health = %v, want %v", dev.Health, pluginapi.Healthy)
		}
	}
}

func TestX11_Allocate(t *testing.T) {
	baseDir := t.TempDir()
	setupX11TestConstants(t, baseDir)
	r := &ResourceX11{}
	resp := r.Allocate([]string{"x11-0"})
	if len(resp) != 1 {
		t.Fatalf("Allocate() = %d responses, want 1", len(resp))
	}
	mounts := resp[0].Mounts
	if len(mounts) != 1 {
		t.Fatalf("Allocate() = %d mounts, want 1", len(mounts))
	}
	m := mounts[0]
	if m.HostPath != constants.X11DevicePath || m.ContainerPath != constants.X11DevicePath {
		t.Errorf("Mount paths = %v/%v, want %v", m.HostPath, m.ContainerPath, constants.X11DevicePath)
	}
	if m.ReadOnly {
		t.Errorf("Mount should not be readonly")
	}
}
