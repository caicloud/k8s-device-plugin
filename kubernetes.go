package main

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/NVIDIA/gpu-monitoring-tools/bindings/go/nvml"
	resourcev1alpha1 "github.com/caicloud/mantle/pkg/apis/resource/v1alpha1"
	erclientset "github.com/caicloud/mantle/pkg/client/clientset/versioned"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog/v2"
)

const (
	formatExtendedResourceName = "nvidia.com-%s"

	// GPUMemory, in bytes. (1Gi = 1GiB = 1 * 1024 * 1024 * 1024)
	GPUMemory         = "gpu.memory"
	GPUThread         = "gpu.thread"

	HostnameLabel = "kubernetes.io/hostname"
	NameLabel = "nvidia.com/name"
	MemoryLabel = "nvidia.com/memory"

)

type KubeClient struct {
	node string
	client *erclientset.Clientset
}

// NewKubeClient return a new KubeClient
func NewKubeClient(kubeConfig, nodeName string) (*KubeClient, error){
	if len(nodeName) == 0 {
		nodeName = os.Getenv("NODE_NAME")
		if len(nodeName) == 0 {
			return nil, fmt.Errorf("NODE_NAME env or --node-name must be set and be pod.spec.nodeName")
		}
	}

	rest, err := clientcmd.BuildConfigFromFlags("", kubeConfig)
	if err != nil {
		return nil, err
	}
	client, err := erclientset.NewForConfig(rest)
	if err != nil {
		return nil, err
	}

	return &KubeClient{
		node: nodeName,
		client: client,
	}, nil
}

// SyncExtendedResources sync gpu to ExtendedResources
func (kc *KubeClient) SyncExtendedResources(ctx context.Context, devices []*Device, resourceName string) error {
	var err error
	var erList *resourcev1alpha1.ExtendedResourceList
	deviceSet := sets.NewString()
	for _, d := range devices {
		deviceSet.Insert(fmt.Sprintf(formatExtendedResourceName, strings.ToLower(d.ID)))
		if err = kc.createOrUpdateExtendedResource(ctx, d, resourceName); err != nil {
			return err
		}
	}

	// clean up non-existent resources
	erList, err = kc.client.ResourceV1alpha1().ExtendedResources().List(
		ctx,
		metav1.ListOptions{
		LabelSelector: labels.FormatLabels(map[string]string{
			HostnameLabel: kc.node,
		}),
	})
	if err != nil {
		return err
	}

	// delete non-existent device
	for _, er := range erList.Items {
		if deviceSet.Has(er.Name) {
			continue
		}
		// TODO: mark resource is unavailable, clean up by a controller.
		err = kc.client.ResourceV1alpha1().ExtendedResources().Delete(ctx, er.Name, metav1.DeleteOptions{})
		if err != nil {
			// Note: do not return an error when deleting er fails
			klog.Warning("Delete ExtendedResources %s fails: %s", er.Name, err)
		}
	}
	return nil
}

// CleanupExtendedResources cleanup ers
func (kc *KubeClient) CleanupExtendedResources(ctx context.Context, devices []*Device) {
	for _, d := range devices {
		if err := kc.DeleteExtendedResourcesFromUUID(ctx, d.ID); err != nil {
			// Note: do not return an error when delete er fails
			klog.Warning("Delete ExtendedResources for %s fails: %s", d.ID , err)
		}
	}
}

// MarkDeviceUnhealthy mark device unhealthy
func (kc *KubeClient) MarkDeviceUnhealthy(ctx context.Context, uuid string) {
	erName := fmt.Sprintf(formatExtendedResourceName, strings.ToLower(uuid))
	er, err := kc.client.ResourceV1alpha1().ExtendedResources().Get(ctx, erName, metav1.GetOptions{})
	if err != nil {
		klog.Warning("Get ExtendedResource failed: %s", err)
	}
	er.Status.Phase = resourcev1alpha1.ExtendedResourcePending
	er.Status.Conditions = append(er.Status.Conditions, metav1.Condition{
		Type: "HealthCheck",
		Status: "Unhealthy",
		Message: "Device unhealthy",
	})

	_, err = kc.client.ResourceV1alpha1().ExtendedResources().UpdateStatus(ctx, er, metav1.UpdateOptions{})
	if err != nil {
		klog.Warning("update ExtendedResource failed: %s", err)
	}
}

// DeleteExtendedResourcesFromUUID delete extended resource from device uuid
func (kc *KubeClient) DeleteExtendedResourcesFromUUID(ctx context.Context, uuid string) error {
	erName := fmt.Sprintf(formatExtendedResourceName, strings.ToLower(uuid))
	return kc.DeleteExtendedResources(ctx, erName)
}

// DeleteExtendedResources delete extended resource
func (kc *KubeClient) DeleteExtendedResources(ctx context.Context, erName string) error {
	er, err := kc.client.ResourceV1alpha1().ExtendedResources().Get(ctx, erName, metav1.GetOptions{})
	if err != nil {
		if apierrors.IsNotFound(err) {
			return nil
		}
		return err
	}
	if er.Status.Phase == resourcev1alpha1.ExtendedResourceBound {
		return fmt.Errorf("could not delete bound device")
	}
	return kc.client.ResourceV1alpha1().ExtendedResources().Delete(ctx, erName, metav1.DeleteOptions{})
}

func (kc *KubeClient) createOrUpdateExtendedResource(ctx context.Context, device *Device, resouceName string) error {
	newER := deviceToExtendedResource(device.nvDevice, kc.node, resouceName)
	oldER, err := kc.client.ResourceV1alpha1().ExtendedResources().Get(ctx, newER.Name, metav1.GetOptions{})
	if apierrors.IsNotFound(err) {
		_, err = kc.client.ResourceV1alpha1().ExtendedResources().Create(ctx, newER, metav1.CreateOptions{})
	}
	if err != nil {
		return err
	}

	if cmpExtendedResources(oldER, newER) {
		return nil
	}

	// update device information
	if oldER.Labels == nil {
		oldER.Labels = make(map[string]string)
	}
	for key, val := range newER.Labels {
		oldER.Labels[key] = val
	}
	oldER.Spec = newER.Spec
	oldER.Status.Capacity = newER.Status.Capacity
	// if the device is unhealthy, update to healthy
	if oldER.Status.Phase == resourcev1alpha1.ExtendedResourcePending {
		oldER.Status.Phase = resourcev1alpha1.ExtendedResourceAvailable
	}
	_, err = kc.client.ResourceV1alpha1().ExtendedResources().Update(ctx, oldER, metav1.UpdateOptions{})
	return err
}


func deviceToExtendedResource (device *nvml.Device, nodeName, resoureName string) *resourcev1alpha1.ExtendedResource {
	erName := fmt.Sprintf(formatExtendedResourceName, strings.ToLower(device.UUID))
	name := strings.Replace(stringPtrToString(device.Model), " ", "-", -1)
	gpuMemory := resource.NewQuantity(int64(uint64PtrToUint64(device.Memory) * 1024 * 1024), resource.BinarySI)
	bar1Memory := resource.NewQuantity(int64(uint64PtrToUint64(device.PCI.BAR1) * 1024 * 1024), resource.BinarySI)
	gpuThread := resource.NewQuantity(int64(100), resource.DecimalSI)


	return &resourcev1alpha1.ExtendedResource{
		ObjectMeta: metav1.ObjectMeta{
			Name: erName,
			Labels: map[string]string{
				HostnameLabel: nodeName,
				NameLabel: name,
				MemoryLabel: gpuMemory.String(),
			},
		},
		Spec: resourcev1alpha1.ExtendedResourceSpec{
			RawResourceName: resoureName,
			DeviceID: device.UUID,
			NodeName: nodeName,
			Properties: map[string]string{
				"Name":           name,
				"UUID":           device.UUID,
				"Memory":         gpuMemory.String(),
				"BAR1Memory":     bar1Memory.String(),
				"Path":           device.Path,
				"Power":          strconv.FormatUint(uintPtrToUint64(device.Power), 10) + "W",
				"CPUAffinity":    strconv.FormatUint(uintPtrToUint64(device.CPUAffinity), 10),
				"Bandwidth":      strconv.FormatUint(uintPtrToUint64(device.PCI.Bandwidth), 10) + "MB/s",
				"PCIBusID":       device.PCI.BusID,
				"MaxSMClock":     strconv.FormatUint(uintPtrToUint64(device.Clocks.Cores), 10) + "MHz",
				"MaxMemoryClock": strconv.FormatUint(uintPtrToUint64(device.Clocks.Memory), 10) + "MHz",
			},
		},
		Status: resourcev1alpha1.ExtendedResourceStatus{
			Phase: resourcev1alpha1.ExtendedResourceAvailable,
			Capacity: corev1.ResourceList{
				GPUMemory: *gpuMemory,
				GPUThread: *gpuThread,
			},
			Allocatable: corev1.ResourceList{
				GPUMemory: *gpuMemory,
				GPUThread: *gpuThread,
			},
		},
	}
}

// camER return false if two er different.
func cmpExtendedResources (new, old *resourcev1alpha1.ExtendedResource) bool {
	// compare labels
	if new.Labels[HostnameLabel] != old.Labels[HostnameLabel] ||
		new.Labels[NameLabel] != old.Labels[NameLabel] ||
		new.Labels[MemoryLabel] != old.Labels[MemoryLabel] {
		return false
	}

	// compare spec
	if new.Spec.RawResourceName != old.Spec.RawResourceName ||
		new.Spec.NodeName != old.Spec.NodeName ||
		len(new.Spec.Properties) != len(old.Spec.Properties) ||
		new.Spec.DeviceID != old.Spec.DeviceID {
		return false
	}
	
	for key, val := range new.Spec.Properties {
		if old.Spec.Properties[key] != val {
			return false
		}
	}
	
	// compare status
	if len(new.Status.Capacity) != len(old.Status.Capacity) {
		return false
	}
	for key, val := range new.Status.Capacity {
		if old.Status.Capacity[key] != val {
			return false
		}
	}
	
	return true
}

func uintPtrToUint64 (i *uint) uint64 {
	if i == nil {
		return 0
	}
	return uint64(*i)
}

func uint64PtrToUint64 (i *uint64) uint64 {
	if i == nil {
		return 0
	}
	return *i
}

func stringPtrToString ( s *string) string {
	if s == nil {
		return ""
	}
	return *s
}