package prioritizer

import (
	"fmt"

	"github.com/uesyn/kubernetes-external-gpu-scheduler/util/logs"

	"k8s.io/api/core/v1"
	v1helper "k8s.io/kubernetes/pkg/apis/core/v1/helper"
	schedulerapi "k8s.io/kubernetes/pkg/scheduler/api"
	"k8s.io/kubernetes/pkg/scheduler/schedulercache"
)

type PrioritizeFunc func(pod v1.Pod, nodes []v1.Node) (*schedulerapi.HostPriorityList, error)

func NewPrioritizeFunc(targetResource string) {
	return func(pod v1.Pod, nodes []v1.Node) (*schedulerapi.HostPriorityList, error) {
	result := []schedulerapi.HostPriority{}
	for _, node := range nodes {
		r := schedulerapi.HostPriority{}
		var err error = nil
		r.Host = node.Name
		r.Score, err = calcNodeScore(pod, node, targetResource)
		if err != nil {
			return nil, err
		}
		result = append(result, r)
	}
	var results schedulerapi.HostPriorityList = result
	return &results, nil

}

func PrioritizeHandler(args schedulerapi.ExtenderArgs) (*schedulerapi.HostPriorityList, error) {
	return Prioritize(args.Pod, args.Nodes.Items)
}

func getResourceRequest(pod *v1.Pod) *schedulercache.Resource {
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

func getNodeInfo(node *v1.Node) (*schedulercache.NodeInfo, error) {
	nodeinfo := schedulercache.NewNodeInfo()
	err := nodeinfo.SetNode(node)
	if err != nil {
		return nil, err
	}
	return nodeinfo, nil
}

func calcNodeScore(pod *v1.Pod, node *v1.Node, targetResource string) (int, error) {
	if node == nil {
		return 0, fmt.Errorf("node not found")
	}
	nodeInfo, err := getNodeInfo(node)
	if err != nil {
		return 0, err
	}

	var podRequest *schedulercache.Resource
	podRequest = getResourceRequest(pod)

	allocatable := nodeInfo.AllocatableResource()

	ratio64 := int64(0)
	for rName, rQuant := range podRequest.ScalarResources {
		// Check whether resource name is extended resource
		if v1helper.IsExtendedResourceName(rName) {
			logs.Tracef("Extended Resource Name %s, Resouce Quantity %d.\n", rName.String(), rQuant)
		} else {
			continue
		}

		// Check whether extended resource is target.
		if rName.String() != targetResource {
			continue
		}

		ratio64 = rQuant + nodeInfo.RequestedResource().ScalarResources[rName]/allocatable.ScalarResources[rName]
		ratio64 = ratio64 * 10
		logs.Tracef("%s usage ratio is %d", rName.String(), ratio64)
		break
	}
	return int(ratio64), nil
}
