package constants

import pluginapi "k8s.io/kubelet/pkg/apis/deviceplugin/v1beta1"

const (
	PluginType           = "PLUGIN_TYPE"
	X11PluginType        = "x11"
	UdmaPluginType       = "udma"
	IgpuPluginType       = "igpu"
	VfioPluginType       = "vfio"
	UsbPluginType        = "usb"
	X11ResourceName      = "intel.com/x11"
	X11DevicePrefix      = "x11"
	X11DevicePath        = "/tmp/.X11-unix"
	X11SockName          = "x11.sock"
	X11SocketPath        = pluginapi.DevicePluginPath + X11SockName
	X11DeviceCount       = 1000
	UdmaResourceName     = "intel.com/udma"
	UdmaDevicePrefix     = "udma"
	UdmaDevicePath       = "/dev/udmabuf"
	UdmaSockName         = "udma.sock"
	UdmaSocketPath       = pluginapi.DevicePluginPath + UdmaSockName
	UdmaDeviceCount      = 1000
	IGpuResourceName     = "intel.com/igpu"
	IGpuDevicePrefix     = "igpu"
	IGpuDevicePath       = "/dev/dri/card0"
	RenderD128DevicePath = "/dev/dri/renderD128"
	IGpuSockName         = "igpu.sock"
	IGpuSocketPath       = pluginapi.DevicePluginPath + IGpuSockName
	IGpuDeviceCount      = 1000
	VfioResourceName     = "intel.com/vfio"
	VfioDevicePrefix     = "vfio"
	VfioDevicePath       = "/dev/vfio"
	VfioSocketName       = "vfio.sock"
	VfioSocketPath       = pluginapi.DevicePluginPath + VfioSocketName
	VfioDeviceCount      = 1000
	UsbResourceName      = "intel.com/usb"
	UsbDevicePrefix      = "usb"
	UsbDevicePath        = "/dev/bus/usb"
	UsbSocketName        = "usb.sock"
	UsbSocketPath        = pluginapi.DevicePluginPath + UsbSocketName
	UsbDeviceCount       = 1000
	Permissions          = "rw"
	DisplayEnvKey        = "DISPLAY"
	DisplayEnvVal        = ":0"
)
