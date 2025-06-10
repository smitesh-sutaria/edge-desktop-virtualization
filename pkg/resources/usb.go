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

// scanUSBBus finds all USB devices and returns them
func scanUSBBus() ([]*pluginapi.ContainerAllocateResponse, error) {
	containerResponse := new(pluginapi.ContainerAllocateResponse)

	busDirs, err := ioutil.ReadDir(constants.UsbDevicePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read USB bus directory: %v", err)
	}

	var devices []*pluginapi.DeviceSpec
	for _, busDir := range busDirs {
		if busDir.IsDir() {
			busPath := filepath.Join(constants.UsbDevicePath, busDir.Name())

			deviceFiles, err := ioutil.ReadDir(busPath)
			if err != nil {
				continue
			}

			for _, deviceFile := range deviceFiles {
				if !deviceFile.IsDir() {
					dev := new(pluginapi.DeviceSpec)
					devicePath := filepath.Join(busPath, deviceFile.Name())
					dev.HostPath = devicePath
					dev.ContainerPath = devicePath
					dev.Permissions = constants.Permissions
					devices = append(devices, dev)
				}
			}
		}
	}

	// Add devices to the container response
	containerResponse.Devices = devices

	// Add container response to the main response
	return []*pluginapi.ContainerAllocateResponse{{Devices: devices}}, nil
}

// Allocate provides the required mounts for USB
func (r *ResourceUSB) Allocate(deviceIDs []string) []*pluginapi.ContainerAllocateResponse {
	log.Printf("usb deviceIds - %v\n", deviceIDs)
	response, err := scanUSBBus()
	if err != nil {
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
	return response
}