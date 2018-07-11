package main

import (
	"fmt"

	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/sets"
	v1helper "k8s.io/kubernetes/pkg/apis/core/v1/helper"
	schedulerapi "k8s.io/kubernetes/pkg/scheduler/api"
	"k8s.io/kubernetes/pkg/scheduler/schedulercache"
)

type Prioritize struct {
	Name string
	Func func(pod v1.Pod, nodes []v1.Node) (*schedulerapi.HostPriorityList, error)
}

func (p Prioritize) Handler(args schedulerapi.ExtenderArgs) (*schedulerapi.HostPriorityList, error) {
	return p.Func(args.Pod, args.Nodes.Items)
}

func GetResourceRequest(pod *v1.Pod) *schedulercache.Resource {
	result := &schedulercache.Resource{}
	for _, container := range pod.Spec.Containers {
		result.Add(container.Resources.Requests)
	}

	// take max_resource(sum_pod, any_init_container)
	for _, container := range pod.Spec.InitContainers {
		for rName, rQuantity := range container.Resources.Requests {
			switch rName {
			case v1.ResourceMemory:
				if mem := rQuantity.Value(); mem > result.Memory {
					result.Memory = mem
				}
			case v1.ResourceEphemeralStorage:
				if ephemeralStorage := rQuantity.Value(); ephemeralStorage > result.EphemeralStorage {
					result.EphemeralStorage = ephemeralStorage
				}
			case v1.ResourceCPU:
				if cpu := rQuantity.MilliValue(); cpu > result.MilliCPU {
					result.MilliCPU = cpu
				}
				//			case v1.ResourceNvidiaGPU:
				//				if gpu := rQuantity.Value(); gpu > result.NvidiaGPU {
				//					result.NvidiaGPU = gpu
				//				}
			default:
				if v1helper.IsScalarResourceName(rName) {
					value := rQuantity.Value()
					if value > result.ScalarResources[rName] {
						result.SetScalar(rName, value)
					}
				}
			}
		}
	}

	return result
}

func GetNodeInfo(node *v1.Node) (*schedulercache.NodeInfo, error) {
	nodeinfo := schedulercache.NewNodeInfo()
	err := nodeinfo.SetNode(node)
	if err != nil {
		return nil, err
	}
	return nodeinfo, nil
}

func ExternalResourcePrioritizer(pod *v1.Pod, node *v1.Node) (bool, error) {
	if node == nil {
		return false, fmt.Errorf("node not found")
	}
	nodeInfo, err := GetNodeInfo(node)
	if err != nil {
		return false, err
	}

	// No extended resources should be ignored by default.
	ignoredExtendedResources := sets.NewString()

	var podRequest *schedulercache.Resource
	podRequest = GetResourceRequest(pod)

	allocatable := nodeInfo.AllocatableResource()

	for rName, rQuant := range podRequest.ScalarResources {
		if v1helper.IsExtendedResourceName(rName) {
			// If this resource is one of the extended resources that should be
			// ignored, we will skip checking it.
			if ignoredExtendedResources.Has(string(rName)) {
				continue
			}
		}
		if allocatable.ScalarResources[rName] < rQuant+nodeInfo.RequestedResource().ScalarResources[rName] {
			//			predicateFails = append(predicateFails, NewInsufficientResourceError(rName, podRequest.ScalarResources[rName], nodeInfo.RequestedResource().ScalarResources[rName], allocatable.ScalarResources[rName]))
		}
	}
	return true, nil
}
