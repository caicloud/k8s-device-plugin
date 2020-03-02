/*
 * Copyright (c) 2019, NVIDIA CORPORATION.  All rights reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"path"
	"strings"
	"time"

	"github.com/NVIDIA/gpu-monitoring-tools/bindings/go/nvml"
	resourcev1beta1 "github.com/caicloud/clientset/customclient/typed/resource/v1beta1"
	"github.com/caicloud/clientset/pkg/apis/resource/v1beta1"
	"google.golang.org/grpc"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	pluginapi "k8s.io/kubelet/pkg/apis/deviceplugin/v1beta1"
)

const (
	formatResourceName = "nvidia.com"

	// GPUMemory, in bytes. (1Gi = 1GiB = 1 * 1024 * 1024 * 1024)
	GPUMemory         = "gpu.memory"
	GPUThread         = "gpu.thread"
	GPUThreadCapacity = int64(100)

	mpsPath = "/tmp/nvidia-mps"
)

// NvidiaDevicePlugin implements the Kubernetes device plugin API
type NvidiaDevicePlugin struct {
	ResourceManager
	clientset      resourcev1beta1.ExtendedResourceInterface
	resourceName   string
	allocateEnvvar string
	socket         string

	server        *grpc.Server
	cachedDevices []*pluginapi.Device
	nvmlDevs      []*nvml.Device
	health        chan *pluginapi.Device
	stop          chan interface{}
	nodeName      string
}

// NewNvidiaDevicePlugin returns an initialized NvidiaDevicePlugin
func NewNvidiaDevicePlugin(clientset resourcev1beta1.ExtendedResourceInterface, resourceName string, resourceManager ResourceManager, allocateEnvvar string, socket string) *NvidiaDevicePlugin {
	nodeName := os.Getenv("NODE_NAME")
	if len(nodeName) == 0 {
		log.Fatalf("`NODE_NAME` env must be set and be `pod.spec.nodeName`.")
	}

	return &NvidiaDevicePlugin{
		clientset:       clientset,
		ResourceManager: resourceManager,
		resourceName:    resourceName,
		allocateEnvvar:  allocateEnvvar,
		socket:          socket,

		// These will be reinitialized every
		// time the plugin server is restarted.
		cachedDevices: nil,
		nvmlDevs:      nil,
		server:        nil,
		health:        nil,
		stop:          nil,
		nodeName:      nodeName,
	}
}

func (m *NvidiaDevicePlugin) initialize() {
	m.cachedDevices, m.nvmlDevs = m.Devices()
	m.server = grpc.NewServer([]grpc.ServerOption{}...)
	m.health = make(chan *pluginapi.Device)
	m.stop = make(chan interface{})

	m.createOrUpdateERs()
}

func (m *NvidiaDevicePlugin) cleanup() {
	close(m.stop)
	m.cachedDevices = nil
	m.nvmlDevs = nil
	m.server = nil
	m.health = nil
	m.stop = nil

	// delete all ers on the node when cleanup
	m.deleteAllERs()
}

// Start starts the gRPC server, registers the device plugin with the Kubelet,
// and starts the device healthchecks.
func (m *NvidiaDevicePlugin) Start() error {
	m.initialize()

	err := m.Serve()
	if err != nil {
		log.Printf("Could not start device plugin for '%s': %s", m.resourceName, err)
		m.cleanup()
		return err
	}
	log.Printf("Starting to serve '%s' on %s", m.resourceName, m.socket)

	err = m.Register()
	if err != nil {
		log.Printf("Could not register device plugin: %s", err)
		m.Stop()
		return err
	}
	log.Printf("Registered device plugin for '%s' with Kubelet", m.resourceName)

	go m.CheckHealth(m.stop, m.cachedDevices, m.health)

	return nil
}

// Stop stops the gRPC server.
func (m *NvidiaDevicePlugin) Stop() error {
	if m == nil || m.server == nil {
		return nil
	}
	log.Printf("Stopping to serve '%s' on %s", m.resourceName, m.socket)
	m.server.Stop()
	if err := os.Remove(m.socket); err != nil && !os.IsNotExist(err) {
		return err
	}
	m.cleanup()
	return nil
}

// Serve starts the gRPC server of the device plugin.
func (m *NvidiaDevicePlugin) Serve() error {
	sock, err := net.Listen("unix", m.socket)
	if err != nil {
		return err
	}

	pluginapi.RegisterDevicePluginServer(m.server, m)

	go func() {
		lastCrashTime := time.Now()
		restartCount := 0
		for {
			log.Printf("Starting GRPC server for '%s'", m.resourceName)
			err := m.server.Serve(sock)
			if err == nil {
				break
			}

			log.Printf("GRPC server for '%s' crashed with error: %v", m.resourceName, err)

			// restart if it has not been too often
			// i.e. if server has crashed more than 5 times and it didn't last more than one hour each time
			if restartCount > 5 {
				// quit
				log.Fatal("GRPC server for '%s' has repeatedly crashed recently. Quitting", m.resourceName)
			}
			timeSinceLastCrash := time.Since(lastCrashTime).Seconds()
			lastCrashTime = time.Now()
			if timeSinceLastCrash > 3600 {
				// it has been one hour since the last crash.. reset the count
				// to reflect on the frequency
				restartCount = 1
			} else {
				restartCount += 1
			}
		}
	}()

	// Wait for server to start by launching a blocking connexion
	conn, err := m.dial(m.socket, 5*time.Second)
	if err != nil {
		return err
	}
	conn.Close()

	return nil
}

// Register registers the device plugin for the given resourceName with Kubelet.
func (m *NvidiaDevicePlugin) Register() error {
	conn, err := m.dial(pluginapi.KubeletSocket, 5*time.Second)
	if err != nil {
		return err
	}
	defer conn.Close()

	client := pluginapi.NewRegistrationClient(conn)
	reqt := &pluginapi.RegisterRequest{
		Version:      pluginapi.Version,
		Endpoint:     path.Base(m.socket),
		ResourceName: m.resourceName,
	}

	_, err = client.Register(context.Background(), reqt)
	if err != nil {
		return err
	}
	return nil
}

func (m *NvidiaDevicePlugin) GetDevicePluginOptions(context.Context, *pluginapi.Empty) (*pluginapi.DevicePluginOptions, error) {
	return &pluginapi.DevicePluginOptions{}, nil
}

// ListAndWatch lists devices and update that list according to the health status
func (m *NvidiaDevicePlugin) ListAndWatch(e *pluginapi.Empty, s pluginapi.DevicePlugin_ListAndWatchServer) error {
	s.Send(&pluginapi.ListAndWatchResponse{Devices: m.cachedDevices})

	for {
		select {
		case <-m.stop:
			return nil
		case d := <-m.health:
			// FIXME: there is no way to recover from the Unhealthy state.
			d.Health = pluginapi.Unhealthy
			log.Printf("'%s' device marked unhealthy: %s", m.resourceName, d.ID)
			s.Send(&pluginapi.ListAndWatchResponse{Devices: m.cachedDevices})
			m.updateER(d, v1beta1.ExtendedResourcePending)
		}
	}
}

// Allocate which return list of devices.
func (m *NvidiaDevicePlugin) Allocate(ctx context.Context, reqs *pluginapi.AllocateRequest) (*pluginapi.AllocateResponse, error) {
	responses := pluginapi.AllocateResponse{}
	for _, req := range reqs.ContainerRequests {
		response := pluginapi.ContainerAllocateResponse{
			Envs: map[string]string{
				m.allocateEnvvar: strings.Join(req.DevicesIDs, ","),
			},
			Mounts: []*pluginapi.Mount{
				{
					ContainerPath: mpsPath,
					HostPath:      mpsPath,
				},
			},
		}

		for _, id := range req.DevicesIDs {
			if !m.deviceExists(m.cachedDevices, id) {
				return nil, fmt.Errorf("invalid allocation request for '%s': unknown device: %s", m.resourceName, id)
			}
		}

		responses.ContainerResponses = append(responses.ContainerResponses, &response)
	}

	return &responses, nil
}

func (m *NvidiaDevicePlugin) PreStartContainer(context.Context, *pluginapi.PreStartContainerRequest) (*pluginapi.PreStartContainerResponse, error) {
	return &pluginapi.PreStartContainerResponse{}, nil
}

// dial establishes the gRPC communication with the registered device plugin.
func (m *NvidiaDevicePlugin) dial(unixSocketPath string, timeout time.Duration) (*grpc.ClientConn, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if c, err := grpc.DialContext(ctx, unixSocketPath, grpc.WithInsecure(), grpc.WithBlock(),
		grpc.WithContextDialer(func(_ context.Context, addr string) (net.Conn, error) {
			return net.DialTimeout("unix", addr, timeout)
		}),
	); err == nil {
		return c, nil
	} else {
		return nil, err
	}
}

func (m *NvidiaDevicePlugin) deviceExists(devs []*pluginapi.Device, id string) bool {
	for _, d := range devs {
		if d.ID == id {
			return true
		}
	}
	return false
}

// createOrUpdateERs is create or update ExtendedResources
func (m *NvidiaDevicePlugin) createOrUpdateERs() {
	for _, dev := range m.nvmlDevs {
		erName := formatExtendedName(dev.UUID)
		cores := appendString(convertUint(dev.Clocks.Cores), "MHz")
		bandwidth := appendString(convertUint(dev.PCI.Bandwidth), "MB")
		memory := appendString(convertUint64(dev.Memory), "MiB")
		gpuMemory := convertInt64(dev.Memory) * 1024 * 1024

		extendedResource := &v1beta1.ExtendedResource{
			ObjectMeta: metav1.ObjectMeta{
				Name: erName,
				Labels: map[string]string{
					"hostname":  m.nodeName,
					"model":     strings.Replace(convertString(dev.Model), " ", "-", -1),
					"cores":     cores,
					"memory":    memory,
					"bandwidth": bandwidth,
				},
			},
			Spec: v1beta1.ExtendedResourceSpec{
				RawResourceName: m.resourceName,
				DeviceID:        dev.UUID,
				NodeName:        m.nodeName,
				Properties: map[string]string{
					"Model":                 convertString(dev.Model),
					"Cores":                 cores,
					"Memory":                memory,
					"Bandwidth":             bandwidth,
					"MaxPcieLinkWidth":      convertUint(dev.PCI.MaxPcieLinkWidth),
					"MaxPcieLinkGeneration": convertUint(dev.PCI.MaxPcieLinkGeneration),
					"Brand":                 dev.Brand,
					"Name":                  dev.Name,
					"MaxGraphicsClock":      appendString(convertUint(dev.Clocks.Graphics), "MHz"),
					"MaxMemClock":           appendString(convertUint(dev.Clocks.Memory), "MHz"),
					"MaxVideoClock":         appendString(convertUint(dev.Clocks.Video), "MHz"),
				},
			},
			Status: v1beta1.ExtendedResourceStatus{
				Phase: v1beta1.ExtendedResourceAvailable,
				Capacity: corev1.ResourceList{
					GPUMemory: *resource.NewQuantity(gpuMemory, resource.BinarySI),
					GPUThread: *resource.NewQuantity(GPUThreadCapacity, resource.DecimalSI),
				},
				Allocatable: corev1.ResourceList{
					GPUMemory: *resource.NewQuantity(gpuMemory, resource.BinarySI),
					GPUThread: *resource.NewQuantity(GPUThreadCapacity, resource.DecimalSI),
				},
			},
		}

		if er, err := m.clientset.Get(erName, metav1.GetOptions{}); err == nil {
			er.Labels = extendedResource.Labels
			er.Spec = extendedResource.Spec
			er.Status = extendedResource.Status

			if _, err := m.clientset.Update(er); err != nil {
				log.Printf("Update extended resource failed: %+v", err)
			} else {
				log.Printf("Update extended resource %v succeed", erName)
			}
			continue
		} else if !errors.IsNotFound(err) {
			log.Printf("Unexpected error: %+v", err)
			continue
		}

		// er is not found, so create it.
		if _, err := m.clientset.Create(extendedResource); err != nil {
			log.Printf("Create extended resource failed: %+v", err)
		} else {
			log.Printf("Create extended resource %v succeed", erName)
		}
	}
}

func (m *NvidiaDevicePlugin) updateER(d *pluginapi.Device, phase v1beta1.ExtendedResourcePhase) {
	er, err := m.clientset.Get(formatExtendedName(d.ID), metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return
		}
		log.Printf("Get extended resource failed: %+v", err)
		return
	}

	er.Status.Phase = phase
	if _, err = m.clientset.Update(er); err != nil {
		log.Printf("Update extended resource failed: %+v", err)
	}
}

// deleteAllERs delete all ExtendedResources on the node
func (m *NvidiaDevicePlugin) deleteAllERs() {
	ers, err := m.clientset.List(metav1.ListOptions{
		LabelSelector: fmt.Sprintf("hostname=%s", m.nodeName),
	})
	if err != nil {
		log.Printf("List extended resource failed: %+v", err)
		return
	}

	for _, er := range ers.Items {
		if err := m.clientset.Delete(er.Name, &metav1.DeleteOptions{}); err != nil {
			log.Printf("Delete extended resource: %v", err)
		}
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

func appendString(s string, suffix string) string {
	return fmt.Sprintf("%s%s", s, suffix)
}

func convertUint(i *uint) string {
	if i == nil {
		return "0"
	}
	return fmt.Sprintf("%d", *i)
}

func convertUint64(i *uint64) string {
	if i == nil {
		return "0"
	}
	return fmt.Sprintf("%d", *i)
}

func convertInt64(i *uint64) int64 {
	if i == nil {
		return 0
	}
	return int64(*i)
}
