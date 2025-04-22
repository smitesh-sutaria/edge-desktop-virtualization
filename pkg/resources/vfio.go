package resources

import (
	"device-plugin/pkg/constants"
	"fmt"
	"io/ioutil"
	pluginapi "k8s.io/kubelet/pkg/apis/deviceplugin/v1beta1"
	"log"
	"os"
	"path/filepath"
)

type ResourceVFIO struct{}

// GetResourceName returns the resource name for registration
func (r *ResourceVFIO) GetResourceName() string {
	return constants.VfioResourceName
}

// GetSocketPath returns the socket path for VFIO
func (r *ResourceVFIO) GetSocketPath() string { return constants.VfioSocketPath }

// ListDevices returns available VFIO devices
func (r *ResourceVFIO) ListDevices() []*pluginapi.Device {
	// Check if mount path exists and then allocate.
	if _, err := os.Stat(constants.VfioDevicePath); os.IsNotExist(err) {
		return []*pluginapi.Device{} // Return an empty list if the path doesn't exist
	}
	devices := make([]*pluginapi.Device, constants.VfioDeviceCount)
	for i := 0; i < constants.VfioDeviceCount; i++ {
		devices[i] = &pluginapi.Device{
			ID:     fmt.Sprintf("%s-%d", constants.VfioDevicePrefix, i),
			Health: pluginapi.Healthy,
		}
	}
	return devices
}

// AllocateDevicesFromVfioFolder creates an AllocateResponse containing all files
// from vfio folder as device specifications
func AllocateDevicesFromVfioFolder() ([]*pluginapi.ContainerAllocateResponse, error) {
	containerResponse := new(pluginapi.ContainerAllocateResponse)

	files, err := ioutil.ReadDir(constants.VfioDevicePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory %s: %v", constants.VfioDevicePath, err)
	}

	var devices []*pluginapi.DeviceSpec
	for _, file := range files {
		if !file.IsDir() {
			dev := new(pluginapi.DeviceSpec)
			devicePath := filepath.Join(constants.VfioDevicePath, file.Name())
			dev.HostPath = devicePath
			dev.ContainerPath = devicePath
			dev.Permissions = constants.Permissions
			devices = append(devices, dev)
		}
	}

	// Add devices to the container response
	containerResponse.Devices = devices

	// Add container response to the main response
	return []*pluginapi.ContainerAllocateResponse{{Devices: devices}}, nil
}

// Allocate provides the required devices for VFIO
func (r *ResourceVFIO) Allocate(deviceIDs []string) []*pluginapi.ContainerAllocateResponse {
	log.Printf("vfio deviceIds - %v\n", deviceIDs)
	response, err := AllocateDevicesFromVfioFolder()
	if err != nil {
		return []*pluginapi.ContainerAllocateResponse{
			{
				Mounts: []*pluginapi.Mount{
					{
						ContainerPath: constants.VfioDevicePath,
						HostPath:      constants.VfioDevicePath,
						ReadOnly:      false,
					},
				},
			},
		}
	}
	return response
}
