package k8sclient

import (
	"errors"

	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

var clientset *kubernetes.Clientset

func SetConfigInCluster() error {
	config, err := rest.InClusterConfig()
	if err != nil {
		return err
	}
	c, err := kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}
	clientset = c
	return nil
}

func SetConfigFromKubeconfig(kubeconfigpath string) error {
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfigpath)
	if err != nil {
		return err
	}
	c, err := kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}
	clientset = c
	return nil
}

func GetClientInCluster() (*kubernetes.Clientset, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}
	return kubernetes.NewForConfig(config)
}

func GetClient(kubeconfigpath string) (*kubernetes.Clientset, error) {
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfigpath)
	if err != nil {
		return nil, err
	}
	return kubernetes.NewForConfig(config)
}

func GetPodsOnNode(node *v1.Node) ([]v1.Pod, error) {
	if clientset == nil {
		return nil, errors.New("You must run SetConfigInCluster or SetConfigFromKubeconfig.")
	}
	opt := metav1.ListOptions{FieldSelector: "spec.nodeName=" + node.Name}
	pods, err := clientset.CoreV1().Pods("").List(opt)
	if err != nil {
		return nil, err
	}
	// maybe this is not needed...
	podlist := pods.DeepCopy()
	return podlist.Items, nil
}
