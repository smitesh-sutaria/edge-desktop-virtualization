package main

import (
	"device-plugin/pkg/grpcserver"
	"device-plugin/pkg/resources"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func fetchResources() []resources.Resource {
	var res []resources.Resource

	res = append(res, &resources.ResourceX11{})
	res = append(res, &resources.ResourceUDMA{})
	res = append(res, &resources.ResourceIGPU{})
	res = append(res, &resources.ResourceVFIO{})
	res = append(res, &resources.ResourceUSB{})
	return res
}

// Main function
func main() {
	log.Println("Starting device plugin")

	grpcserver.StartDevicePlugin(fetchResources(),
		grpcserver.Serve,
		grpcserver.RegisterWithKubelet)

	log.Println("Started device plugin")

	// Block the main function to keep the program running
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs

	log.Println("Shutting down device plugin")

	err := grpcserver.CleanupSockets(fetchResources())
	if err != nil {
		log.Printf("Could not clean up socket...Exiting")
	}
}
