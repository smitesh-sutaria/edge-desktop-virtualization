package grpcserver

import (
	"context"
	"device-plugin/pkg/resources"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net"
	"os"
	"path/filepath"
	"time"

	pluginapi "k8s.io/kubelet/pkg/apis/deviceplugin/v1beta1"
)

type DevicePluginServer struct {
	pluginapi.UnimplementedDevicePluginServer
	resource resources.Resource
}

// GetDevicePluginOptions provides plugin options (empty response is fine)
func (s *DevicePluginServer) GetDevicePluginOptions(ctx context.Context, req *pluginapi.Empty) (*pluginapi.DevicePluginOptions, error) {
	return &pluginapi.DevicePluginOptions{}, nil
}

// GetPreferredAllocation is optional (returns empty response)
func (s *DevicePluginServer) GetPreferredAllocation(ctx context.Context, req *pluginapi.PreferredAllocationRequest) (*pluginapi.PreferredAllocationResponse, error) {
	return &pluginapi.PreferredAllocationResponse{}, nil
}

// ListAndWatch streams available devices to Kubelet
func (s *DevicePluginServer) ListAndWatch(req *pluginapi.Empty, stream pluginapi.DevicePlugin_ListAndWatchServer) error {
	if err := stream.Send(&pluginapi.ListAndWatchResponse{Devices: s.resource.ListDevices()}); err != nil {
		return fmt.Errorf("failed to send initial device list: %v", err)
	}

	ticker := time.NewTicker(20 * time.Second)
	defer ticker.Stop()

	done := make(chan bool)
	defer close(done)

	go func() {
		<-stream.Context().Done()
		log.Printf("Stream context done, kubelet may have restarted")
		done <- true
	}()

	for {
		select {
		case <-ticker.C:
			devices := s.resource.ListDevices()

			log.Printf("Reporting %d devices as available", len(devices))

			if err := stream.Send(&pluginapi.ListAndWatchResponse{Devices: devices}); err != nil {
				return fmt.Errorf("failed to send device update: %v", err)
			}

		case <-done:
			log.Printf("ListAndWatch stream context canceled")
			return nil
		}
	}
}

// Allocate handles resource allocation requests from Kubernetes
func (s *DevicePluginServer) Allocate(ctx context.Context, req *pluginapi.AllocateRequest) (*pluginapi.AllocateResponse, error) {
	var responses []*pluginapi.ContainerAllocateResponse
	res := s.resource.Allocate(req.ContainerRequests[0].DevicesIDs)
	responses = append(responses, res...)

	return &pluginapi.AllocateResponse{ContainerResponses: responses}, nil
}

// RegisterWithKubelet registers the resource with kubelet
func RegisterWithKubelet(socket string, resourceName string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	conn, err := grpc.DialContext(ctx, pluginapi.KubeletSocket, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithContextDialer(func(ctx context.Context, addr string) (net.Conn, error) {
		d := &net.Dialer{}
		return d.DialContext(ctx, "unix", addr)
	}))
	if err != nil {
		return err
	}
	defer conn.Close()
	log.Printf("connection created for %v", resourceName)
	client := pluginapi.NewRegistrationClient(conn)
	log.Printf("client created for %v", resourceName)
	_, err = client.Register(ctx, &pluginapi.RegisterRequest{
		Version:      pluginapi.Version,
		Endpoint:     filepath.Base(socket),
		ResourceName: resourceName,
	})
	if err != nil {
		return err
	}
	log.Printf("client registered for %v", resourceName)
	return nil
}

// CleanupSocket removes existing socket file if it exists
func CleanupSocket(socket string) error {
	if _, err := os.Stat(socket); err == nil {
		if err := os.Remove(socket); err != nil {
			return err
		}
	}
	return nil
}

// StartDevicePlugin starts server with restart capability
func StartDevicePlugin(socket string, res resources.Resource) {
	for {
		if err := Serve(socket, res); err != nil {
			log.Printf("Device plugin server failed: %v", err)
		}

		time.Sleep(10 * time.Second)

		log.Printf("Restarting device plugin for resource %s", res.GetResourceName())
	}
}

// Serve starts the gRPC server for each resource
func Serve(socket string, res resources.Resource) error {
	// Cleanup any existing socket file before binding
	if err := CleanupSocket(socket); err != nil {
		log.Fatalf("Failed to remove existing socket file: %v", err)
	}

	lis, err := net.Listen("unix", socket)
	if err != nil {
		return err
	}

	s := grpc.NewServer()
	pluginapi.RegisterDevicePluginServer(s, &DevicePluginServer{resource: res})

	log.Printf("Device plugin server starting for socket - %v\n", socket)

	errChan := make(chan error, 1)

	go func() {
		if err := s.Serve(lis); err != nil {
			errChan <- err
		}
	}()

	// Wait for server to start
	time.Sleep(2 * time.Second)

	// Register resource
	log.Printf("registering with kubelet for %v\n", res.GetResourceName())
	if err = RegisterWithKubelet(socket, res.GetResourceName()); err != nil {
		s.Stop()
		return fmt.Errorf("failed to register resource %s with Kubelet: %v", res.GetResourceName(), err)
	}

	// Detect kubelet restarts
	socketWatcher := make(chan bool)
	go monitorSocket(socket, socketWatcher)

	// Wait for server error, socket deletion, or signal
	select {
	case err := <-errChan:
		return fmt.Errorf("gRPC server error: %v", err)
	case <-socketWatcher:
		return fmt.Errorf("socket was deleted, kubelet likely restarted")
	}

	return nil
}

// monitorSocket monitors socket file for deletion (indicates kubelet restart)
func monitorSocket(socket string, notify chan<- bool) {
	for {
		if _, err := os.Stat(socket); os.IsNotExist(err) {
			log.Printf("Socket file %s was deleted", socket)
			notify <- true
			return
		}
		time.Sleep(1 * time.Second)
	}
}
