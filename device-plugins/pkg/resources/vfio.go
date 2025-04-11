package resources

import (
	"device-plugin/pkg/constants"
	"fmt"
	pluginapi "k8s.io/kubelet/pkg/apis/deviceplugin/v1beta1"
	"log"
	"os"
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

// Allocate provides the required mounts for VFIO
func (r *ResourceVFIO) Allocate(deviceIDs []string) []*pluginapi.ContainerAllocateResponse {
	log.Printf("vfio deviceIds - %v\n", deviceIDs)
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
