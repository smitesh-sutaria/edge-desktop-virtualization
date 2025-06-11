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

type ResourceX11 struct{}

// GetResourceName returns the resource name for registration
func (r *ResourceX11) GetResourceName() string {
	return constants.X11ResourceName
}

// GetSocketPath returns the socket path for x11
func (r *ResourceX11) GetSocketPath() string { return constants.X11SocketPath }

// ListDevices returns available X11 devices
func (r *ResourceX11) ListDevices() []*pluginapi.Device {
	// Check if mount path exists and then allocate.
	if _, err := os.Stat(constants.X11DevicePath); os.IsNotExist(err) {
		return []*pluginapi.Device{} // Return an empty list if the path doesn't exist
	}
	devices := make([]*pluginapi.Device, constants.X11DeviceCount)
	for i := 0; i < constants.X11DeviceCount; i++ {
		devices[i] = &pluginapi.Device{
			ID:     fmt.Sprintf("%s-%d", constants.X11DevicePrefix, i),
			Health: pluginapi.Healthy,
		}
	}
	return devices
}

// Allocate provides the required mounts for X11
func (r *ResourceX11) Allocate(deviceIDs []string) []*pluginapi.ContainerAllocateResponse {
	log.Printf("x11 deviceIds - %v\n", deviceIDs)
	return []*pluginapi.ContainerAllocateResponse{
		{
			Mounts: []*pluginapi.Mount{
				{
					ContainerPath: constants.X11DevicePath,
					HostPath:      constants.X11DevicePath,
					ReadOnly:      false,
				},
			},
		},
	}
}
