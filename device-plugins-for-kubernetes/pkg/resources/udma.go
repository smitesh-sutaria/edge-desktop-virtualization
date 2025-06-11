/*
 *  Copyright (C) 2025 Intel Corporation
 *  SPDX-License-Identifier: Apache-2.0
 */
package resources

import (
	"device-plugin/pkg/constants"
	"fmt"
	"log"
	"os"

	pluginapi "k8s.io/kubelet/pkg/apis/deviceplugin/v1beta1"
)

type ResourceUDMA struct{}

// GetResourceName returns the resource name for registration
func (r *ResourceUDMA) GetResourceName() string {
	return constants.UdmaResourceName
}

// GetSocketPath returns the socket path for udma
func (r *ResourceUDMA) GetSocketPath() string { return constants.UdmaSocketPath }

// ListDevices returns available UDMA devices
func (r *ResourceUDMA) ListDevices() []*pluginapi.Device {
	// Check if mount path exists and then allocate.
	if _, err := os.Stat(constants.UdmaDevicePath); os.IsNotExist(err) {
		return []*pluginapi.Device{} // Return an empty list if the path doesn't exist
	}
	if _, err := os.Stat(constants.VfioDevicePath); os.IsNotExist(err) {
		return []*pluginapi.Device{} // Return an empty list if the path doesn't exist
	}
	devices := make([]*pluginapi.Device, constants.UdmaDeviceCount)
	for i := 0; i < constants.UdmaDeviceCount; i++ {
		devices[i] = &pluginapi.Device{
			ID:     fmt.Sprintf("%s-%d", constants.UdmaDevicePrefix, i),
			Health: pluginapi.Healthy,
		}
	}
	return devices
}

// Allocate provides the required devices for UDMA
func (r *ResourceUDMA) Allocate(deviceIDs []string) []*pluginapi.ContainerAllocateResponse {
	log.Printf("udma deviceIds - %v\n", deviceIDs)
	return []*pluginapi.ContainerAllocateResponse{
		{
			Devices: []*pluginapi.DeviceSpec{
				{
					HostPath:      constants.UdmaDevicePath,
					ContainerPath: constants.UdmaDevicePath,
					Permissions:   constants.Permissions,
				},
			},
		},
	}
}
