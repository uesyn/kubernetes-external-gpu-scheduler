package prioritizer

import (
	"fmt"
	"log"

	"github.com/uesyn/kubernetes-external-gpu-scheduler/k8sclient"
	"github.com/uesyn/kubernetes-external-gpu-scheduler/util/logs"

	"k8s.io/api/core/v1"
	v1helper "k8s.io/kubernetes/pkg/apis/core/v1/helper"
	schedulerapi "k8s.io/kubernetes/pkg/scheduler/api"
	"k8s.io/kubernetes/pkg/scheduler/schedulercache"
)

type Prioritizer interface {
	Prioritize(*v1.Pod, []v1.Node) (*schedulerapi.HostPriorityList, error)
}

type PrioritizeFunc func(pod *v1.Pod, nodes []v1.Node) (*schedulerapi.HostPriorityList, error)

type ExtendedResourcePrioritizer struct {
	TargetResource string
	Func           PrioritizeFunc
}

func NewExtendedResourcePrioritizer(targetResource string) *ExtendedResourcePrioritizer {
	fn := newExtendedResourcePrioritizeFuncFactory(targetResource)
	erp := ExtendedResourcePrioritizer{
		TargetResource: targetResource,
		Func:           fn,
	}
	return &erp
}

func (erp *ExtendedResourcePrioritizer) Prioritize(pod *v1.Pod, nodes []v1.Node) (*schedulerapi.HostPriorityList, error) {
	return erp.Func(pod, nodes)
}

func newExtendedResourcePrioritizeFuncFactory(targetResource string) PrioritizeFunc {
	return func(pod *v1.Pod, nodes []v1.Node) (*schedulerapi.HostPriorityList, error) {
		results := []schedulerapi.HostPriority{}
		logs.Debugln("Extended Resource Prioritize Started...")
		for _, node := range nodes {
			logs.Debugln("Target Node:", node.Name)
			r := schedulerapi.HostPriority{}
			var err error = nil
			r.Host = node.Name
			r.Score, err = calcNodeScore(pod, &node, targetResource)
			if err != nil {
				return nil, err
			}
			logs.Debugln(r.Host, "score is", r.Score)
			results = append(results, r)
		}
		var resultlist schedulerapi.HostPriorityList = results
		return &resultlist, nil
	}
}

func getResourceRequest(pod *v1.Pod) *schedulercache.Resource {
	result := &schedulercache.Resource{}
	for _, container := range pod.Spec.Containers {
		result.Add(container.Resources.Requests)
	}

	for _, container := range pod.Spec.InitContainers {
		for rName, rQuantity := range container.Resources.Requests {
			if v1helper.IsScalarResourceName(rName) {
				value := rQuantity.Value()
				if value > result.ScalarResources[rName] {
					result.SetScalar(rName, value)
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

	pods, err := k8sclient.GetPodsOnNode(node)
	if err != nil {
		return nil, err
	}
	for _, pod := range pods {
		nodeinfo.AddPod(&pod)
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

	ratio := int(0)
	for rName, rQuant := range podRequest.ScalarResources {
		// Check whether resource name is extended resource
		if v1helper.IsExtendedResourceName(rName) {
			log.Printf("info: Extended Resource Name %s, Resouce Quantity %d.\n", rName.String(), rQuant)
		} else {
			continue
		}

		// Check whether extended resource is target.
		logs.Debugf("rName is %s\n", rName.String())
		if rName.String() != targetResource {
			logs.Debugf("%s is not target resource.\n", rName.String())
			continue
		}

		logs.Debugf("Node Requested Resource is %d\n", nodeInfo.RequestedResource().ScalarResources[rName])
		logs.Debugf("Node Allocatable Resource is %d\n", allocatable.ScalarResources[rName])
		ratio = int(float64(rQuant) + float64(nodeInfo.RequestedResource().ScalarResources[rName])/float64(allocatable.ScalarResources[rName])*float64(10))
		logs.Infof("%s usage ratio is %d\n", rName.String(), int(ratio))
		break
	}
	return ratio, nil
}
