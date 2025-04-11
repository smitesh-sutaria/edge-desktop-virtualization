package resources

import (
	"device-plugin/pkg/constants"
	"fmt"
	pluginapi "k8s.io/kubelet/pkg/apis/deviceplugin/v1beta1"
	"log"
	"os"
)

type ResourceUSB struct{}

// GetResourceName returns the resource name for registration
func (r *ResourceUSB) GetResourceName() string {
	return constants.UsbResourceName
}

// GetSocketPath returns the socket path for USB
func (r *ResourceUSB) GetSocketPath() string { return constants.UsbSocketPath }

// ListDevices returns available USB devices
func (r *ResourceUSB) ListDevices() []*pluginapi.Device {
	// Check if mount path exists and then allocate.
	if _, err := os.Stat(constants.UsbDevicePath); os.IsNotExist(err) {
		return []*pluginapi.Device{} // Return an empty list if the path doesn't exist
	}
	devices := make([]*pluginapi.Device, constants.UsbDeviceCount)
	for i := 0; i < constants.UsbDeviceCount; i++ {
		devices[i] = &pluginapi.Device{
			ID:     fmt.Sprintf("%s-%d", constants.UsbDevicePrefix, i),
			Health: pluginapi.Healthy,
		}
	}
	return devices
}

// Allocate provides the required mounts for USB
func (r *ResourceUSB) Allocate(deviceIDs []string) []*pluginapi.ContainerAllocateResponse {
	log.Printf("usb deviceIds - %v\n", deviceIDs)
	return []*pluginapi.ContainerAllocateResponse{
		{
			Mounts: []*pluginapi.Mount{
				{
					ContainerPath: constants.UsbDevicePath,
					HostPath:      constants.UsbDevicePath,
					ReadOnly:      false,
				},
			},
		},
	}
}
