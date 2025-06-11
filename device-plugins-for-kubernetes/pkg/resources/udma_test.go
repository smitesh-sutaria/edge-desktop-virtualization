/*
 *  Copyright (C) 2025 Intel Corporation
 *  SPDX-License-Identifier: Apache-2.0
 */
package resources

import (
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"device-plugin/pkg/constants"

	pluginapi "k8s.io/kubelet/pkg/apis/deviceplugin/v1beta1"
)

func setupUDMATestConstants(t *testing.T, baseDir string) {
	constants.UdmaResourceName = "test.udma/resource"
	constants.UdmaSocketPath = filepath.Join(baseDir, "udma.sock")
	constants.UdmaDevicePath = filepath.Join(baseDir, "udma")
	constants.UdmaDeviceCount = 2
	constants.UdmaDevicePrefix = "udma"
	constants.VfioDevicePath = filepath.Join(baseDir, "vfio")
	constants.Permissions = "rw"
}

func TestUDMA_GetResourceName(t *testing.T) {
	baseDir := t.TempDir()
	setupUDMATestConstants(t, baseDir)
	r := &ResourceUDMA{}
	if got := r.GetResourceName(); got != constants.UdmaResourceName {
		t.Errorf("GetResourceName() = %v, want %v", got, constants.UdmaResourceName)
	}
}

func TestUDMA_GetSocketPath(t *testing.T) {
	baseDir := t.TempDir()
	setupUDMATestConstants(t, baseDir)
	r := &ResourceUDMA{}
	if got := r.GetSocketPath(); got != constants.UdmaSocketPath {
		t.Errorf("GetSocketPath() = %v, want %v", got, constants.UdmaSocketPath)
	}
}

func TestUDMA_ListDevices_PathNotExist(t *testing.T) {
	baseDir := t.TempDir()
	setupUDMATestConstants(t, baseDir)
	os.RemoveAll(constants.UdmaDevicePath)
	os.RemoveAll(constants.VfioDevicePath)
	r := &ResourceUDMA{}
	devices := r.ListDevices()
	if len(devices) != 0 {
		t.Errorf("ListDevices() when path does not exist, want 0 devices, got %d", len(devices))
	}
}

func TestUDMA_ListDevices_ReturnsDevices(t *testing.T) {
	baseDir := t.TempDir()
	setupUDMATestConstants(t, baseDir)
	os.MkdirAll(constants.UdmaDevicePath, 0755)
	os.MkdirAll(constants.VfioDevicePath, 0755)
	r := &ResourceUDMA{}
	devices := r.ListDevices()
	if len(devices) != constants.UdmaDeviceCount {
		t.Errorf("ListDevices() = %d devices, want %d", len(devices), constants.UdmaDeviceCount)
	}
	for i, dev := range devices {
		wantID := constants.UdmaDevicePrefix + "-" + strconv.Itoa(i)
		if dev.ID != wantID {
			t.Errorf("Device ID = %v, want %v", dev.ID, wantID)
		}
		if dev.Health != pluginapi.Healthy {
			t.Errorf("Device Health = %v, want %v", dev.Health, pluginapi.Healthy)
		}
	}
}

func TestUDMA_Allocate(t *testing.T) {
	baseDir := t.TempDir()
	setupUDMATestConstants(t, baseDir)
	r := &ResourceUDMA{}
	resp := r.Allocate([]string{"udma-0"})
	if len(resp) != 1 {
		t.Fatalf("Allocate() = %d responses, want 1", len(resp))
	}
	devSpecs := resp[0].Devices
	if len(devSpecs) != 1 {
		t.Fatalf("Allocate() = %d device specs, want 1", len(devSpecs))
	}
	dev := devSpecs[0]
	if dev.HostPath != constants.UdmaDevicePath || dev.ContainerPath != constants.UdmaDevicePath {
		t.Errorf("Device paths = %v/%v, want %v", dev.HostPath, dev.ContainerPath, constants.UdmaDevicePath)
	}
	if dev.Permissions != constants.Permissions {
		t.Errorf("Device permissions = %v, want %v", dev.Permissions, constants.Permissions)
	}
}
