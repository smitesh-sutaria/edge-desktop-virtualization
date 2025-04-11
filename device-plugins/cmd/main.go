package main

import (
	"device-plugin/pkg/constants"
	"device-plugin/pkg/grpcserver"
	"device-plugin/pkg/resources"
	"log"
	"os"
	"os/signal"
	"syscall"
)

// Main function
func main() {
	var res resources.Resource
	var socketPath string
	pluginType := os.Getenv(constants.PluginType)

	switch pluginType {
	case constants.X11PluginType:
		res = &resources.ResourceX11{}
	case constants.UdmaPluginType:
		res = &resources.ResourceUDMA{}
	case constants.IgpuPluginType:
		res = &resources.ResourceIGPU{}
	case constants.VfioPluginType:
		res = &resources.ResourceVFIO{}
	case constants.UsbPluginType:
		res = &resources.ResourceUSB{}
	default:
		log.Fatalf("unknown PLUGIN_TYPE: %s", pluginType)
	}

	socketPath = res.GetSocketPath()

	log.Println("Starting device plugin")

	grpcserver.StartDevicePlugin(socketPath, res)

	log.Printf("Started device plugin in socket %v\n", socketPath)

	// Block the main function to keep the program running
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs

	log.Println("Shutting down device plugin")

	err := grpcserver.CleanupSocket(socketPath)
	if err != nil {
		log.Printf("Could not clean up socket...Exiting")
	}
}
