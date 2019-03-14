// Copyright (c) 2017, NVIDIA CORPORATION. All rights reserved.

package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"path"
	"strings"
	"time"

	"github.com/NVIDIA/gpu-monitoring-tools/bindings/go/nvml"
	clientset "github.com/caicloud/clientset/kubernetes"
	"github.com/caicloud/clientset/pkg/apis/resource/v1beta1"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	pluginapi "k8s.io/kubernetes/pkg/kubelet/apis/deviceplugin/v1beta1"
)

const (
	resourceName           = "nvidia.com/gpu"
	serverSock             = pluginapi.DevicePluginPath + "nvidia.sock"
	envDisableHealthChecks = "DP_DISABLE_HEALTHCHECKS"
	allHealthChecks        = "xids"
	formatResourceName     = "nvidia.com"
)

// NvidiaDevicePlugin implements the Kubernetes device plugin API
type NvidiaDevicePlugin struct {
	devs   []*pluginapi.Device
	socket string

	stop   chan interface{}
	health chan *pluginapi.Device

	server *grpc.Server

	resourceClient *clientset.Clientset
}

// NewNvidiaDevicePlugin returns an initialized NvidiaDevicePlugin
func NewNvidiaDevicePlugin(resourceClient *clientset.Clientset) *NvidiaDevicePlugin {
	nodeName := os.Getenv("NODE_NAME")
	devs, nvmlDevs := getDevices()

	nvidiaDevicePlugin := &NvidiaDevicePlugin{
		devs:   devs,
		socket: serverSock,

		stop:   make(chan interface{}),
		health: make(chan *pluginapi.Device),

		resourceClient: resourceClient,
	}

	nvidiaDevicePlugin.generateExtendedResources(nvmlDevs, nodeName)

	return nvidiaDevicePlugin
}

func (m *NvidiaDevicePlugin) GetDevicePluginOptions(context.Context, *pluginapi.Empty) (*pluginapi.DevicePluginOptions, error) {
	return &pluginapi.DevicePluginOptions{}, nil
}

// dial establishes the gRPC communication with the registered device plugin.
func dial(unixSocketPath string, timeout time.Duration) (*grpc.ClientConn, error) {
	c, err := grpc.Dial(unixSocketPath, grpc.WithInsecure(), grpc.WithBlock(),
		grpc.WithTimeout(timeout),
		grpc.WithDialer(func(addr string, timeout time.Duration) (net.Conn, error) {
			return net.DialTimeout("unix", addr, timeout)
		}),
	)

	if err != nil {
		return nil, err
	}

	return c, nil
}

// Start starts the gRPC server of the device plugin
func (m *NvidiaDevicePlugin) Start() error {
	err := m.cleanup()
	if err != nil {
		return err
	}

	sock, err := net.Listen("unix", m.socket)
	if err != nil {
		return err
	}

	m.server = grpc.NewServer([]grpc.ServerOption{}...)
	pluginapi.RegisterDevicePluginServer(m.server, m)

	go m.server.Serve(sock)

	// Wait for server to start by launching a blocking connexion
	conn, err := dial(m.socket, 5*time.Second)
	if err != nil {
		return err
	}
	conn.Close()

	go m.healthcheck()

	return nil
}

// Stop stops the gRPC server
func (m *NvidiaDevicePlugin) Stop() error {
	if m.server == nil {
		return nil
	}

	m.server.Stop()
	m.server = nil
	close(m.stop)

	return m.cleanup()
}

// Register registers the device plugin for the given resourceName with Kubelet.
func (m *NvidiaDevicePlugin) Register(kubeletEndpoint, resourceName string) error {
	conn, err := dial(kubeletEndpoint, 5*time.Second)
	if err != nil {
		return err
	}
	defer conn.Close()

	client := pluginapi.NewRegistrationClient(conn)
	reqt := &pluginapi.RegisterRequest{
		Version:      pluginapi.Version,
		Endpoint:     path.Base(m.socket),
		ResourceName: resourceName,
	}

	_, err = client.Register(context.Background(), reqt)
	if err != nil {
		return err
	}
	return nil
}

// ListAndWatch lists devices and update that list according to the health status
func (m *NvidiaDevicePlugin) ListAndWatch(e *pluginapi.Empty, s pluginapi.DevicePlugin_ListAndWatchServer) error {
	s.Send(&pluginapi.ListAndWatchResponse{Devices: m.devs})

	for {
		select {
		case <-m.stop:
			return nil
		case d := <-m.health:
			// FIXME: there is no way to recover from the Unhealthy state.
			d.Health = pluginapi.Unhealthy
			m.updateExtendedResource(d.ID, v1beta1.ExtendedResourcePending)
			s.Send(&pluginapi.ListAndWatchResponse{Devices: m.devs})
		}
	}
}

func (m *NvidiaDevicePlugin) unhealthy(dev *pluginapi.Device) {
	m.health <- dev
}

// Allocate which return list of devices.
func (m *NvidiaDevicePlugin) Allocate(ctx context.Context, reqs *pluginapi.AllocateRequest) (*pluginapi.AllocateResponse, error) {
	devs := m.devs
	responses := pluginapi.AllocateResponse{}
	for _, req := range reqs.ContainerRequests {
		response := pluginapi.ContainerAllocateResponse{
			Envs: map[string]string{
				"NVIDIA_VISIBLE_DEVICES": strings.Join(req.DevicesIDs, ","),
			},
		}

		for _, id := range req.DevicesIDs {
			if !deviceExists(devs, id) {
				return nil, fmt.Errorf("invalid allocation request: unknown device: %s", id)
			}
		}

		responses.ContainerResponses = append(responses.ContainerResponses, &response)
	}
	return &responses, nil
}

func (m *NvidiaDevicePlugin) PreStartContainer(context.Context, *pluginapi.PreStartContainerRequest) (*pluginapi.PreStartContainerResponse, error) {
	return &pluginapi.PreStartContainerResponse{}, nil
}

func (m *NvidiaDevicePlugin) cleanup() error {
	if err := os.Remove(m.socket); err != nil && !os.IsNotExist(err) {
		return err
	}

	return nil
}

func (m *NvidiaDevicePlugin) healthcheck() {
	disableHealthChecks := strings.ToLower(os.Getenv(envDisableHealthChecks))
	if disableHealthChecks == "all" {
		disableHealthChecks = allHealthChecks
	}

	ctx, cancel := context.WithCancel(context.Background())

	var xids chan *pluginapi.Device
	if !strings.Contains(disableHealthChecks, "xids") {
		xids = make(chan *pluginapi.Device)
		go watchXIDs(ctx, m.devs, xids)
	}

	for {
		select {
		case <-m.stop:
			cancel()
			return
		case dev := <-xids:
			m.unhealthy(dev)
		}
	}
}

// Serve starts the gRPC server and register the device plugin to Kubelet
func (m *NvidiaDevicePlugin) Serve() error {
	err := m.Start()
	if err != nil {
		log.Printf("Could not start device plugin: %s", err)
		return err
	}
	log.Println("Starting to serve on", m.socket)

	err = m.Register(pluginapi.KubeletSocket, resourceName)
	if err != nil {
		log.Printf("Could not register device plugin: %s", err)
		m.Stop()
		return err
	}
	log.Println("Registered device plugin with Kubelet")

	return nil
}

func (m *NvidiaDevicePlugin) generateExtendedResources(devs []*nvml.Device, nodeName string) {
	if nodeName == "" {
		log.Fatalf("nodeName cannot be empty.")
	}

	for _, dev := range devs {
		extendedResourceName := formatExtendedName(dev.UUID)
		extendedResource := &v1beta1.ExtendedResource{
			ObjectMeta: metav1.ObjectMeta{
				Name: extendedResourceName,
				Labels: map[string]string{
					"hostname":  nodeName,
					"model":     strings.Replace(convertString(dev.Model), " ", "-", -1),
					"cores":     convertCoreUint(dev.Clocks.Cores),
					"memory":    convertUint64(dev.Memory),
					"bandwidth": convertBandwidthUint(dev.PCI.Bandwidth),
				},
			},
			Spec: v1beta1.ExtendedResourceSpec{
				RawResourceName: resourceName,
				DeviceID:        dev.UUID,
				NodeName:        nodeName,
				Properties: map[string]string{
					"Model":     convertString(dev.Model),
					"Cores":     convertCoreUint(dev.Clocks.Cores),
					"Memory":    convertUint64(dev.Memory),
					"Bandwidth": convertBandwidthUint(dev.PCI.Bandwidth),
				},
			},
			Status: v1beta1.ExtendedResourceStatus{
				Phase: v1beta1.ExtendedResourceAvailable,
			},
		}

		extendedResourceOrigin, err := m.resourceClient.ResourceV1beta1().ExtendedResources().Get(extendedResourceName, metav1.GetOptions{})
		if err == nil && extendedResourceOrigin != nil {
			extendedResource.ResourceVersion = extendedResourceOrigin.ResourceVersion
			_, err := m.resourceClient.ResourceV1beta1().ExtendedResources().Update(extendedResource)
			if err != nil {
				log.Printf("Update ExtendedResource: %+v", err)
				continue
			}
		}

		if errors.IsNotFound(err) {
			_, err := m.resourceClient.ResourceV1beta1().ExtendedResources().Create(extendedResource)
			if err != nil {
				log.Printf("Create ExtendedResource: %+v", err)
				continue
			}
		}
		if err != nil {
			log.Printf("Get ExtendedResource: %+v", err)
		}
	}
}

func (m *NvidiaDevicePlugin) updateExtendedResource(deviceID string, phase v1beta1.ExtendedResourcePhase) {
	computeResourceName := formatExtendedName(deviceID)
	computeResourceOrigin, err := m.resourceClient.ResourceV1beta1().ExtendedResources().Get(computeResourceName, metav1.GetOptions{})
	if errors.IsNotFound(err) || computeResourceOrigin == nil {
		return
	}
	if err != nil {
		log.Printf("Cannot get ExtendedResource: %+v", err)
		return
	}

	computeResource := computeResourceOrigin.DeepCopy()

	computeResource.Status.Phase = phase
	_, err = m.resourceClient.ResourceV1beta1().ExtendedResources().Update(computeResource)
	if err != nil {
		log.Printf("Cannot update ExtendedResource: %+v", err)
		return
	}
}

func formatExtendedName(deviceID string) string {
	return fmt.Sprintf("%s-%s", formatResourceName, strings.ToLower(deviceID))
}

func convertString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func convertUint64(i *uint64) string {
	if i == nil {
		return ""
	}
	return fmt.Sprintf("%dMiB", *i)
}

func convertCoreUint(i *uint) string {
	if i == nil {
		return ""
	}
	return fmt.Sprintf("%dMHz", *i)
}

func convertBandwidthUint(i *uint) string {
	if i == nil {
		return ""
	}
	return fmt.Sprintf("%dMB", *i)
}
