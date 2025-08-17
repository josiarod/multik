package kube

import (
	"context"
	"fmt"
	"time"

	"github.com/josiarod/multik/internal/types"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GetPod fetches one Pod by namespace/name using a type client
func GetPod(ctx context.Context, cc ClusterClient, namespace, name string) (*corev1.Pod, error) {
	pod, err := cc.Client.CoreV1().Pods(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("%s: get pod %s/%s: %w", cc.Name, namespace, name, err)
	}
	return pod, nil
}

func ToPodDescribe(cluster string, p *corev1.Pod) types.PodDescribe {
	now := time.Now()

	ready, total := 0, 0
	for _, cs := range p.Status.ContainerStatuses {
		total++
		if cs.Ready {
			ready++
		}
	}

	containers := make([]types.ContainerInfo, 0, len(p.Spec.Containers))
	for _, c := range p.Spec.Containers {
		ports := make([]string, 0, len(c.Ports))
		for _, prt := range c.Ports {
			proto := string(prt.Protocol)
			ports = append(ports, fmt.Sprintf("%d%s", prt.ContainerPort, proto))
		}
		containers = append(containers, types.ContainerInfo{
			Name:   c.Name,
			Image:  c.Image,
			Ports:  ports,
			EnvLen: len(c.Env),
		})
	}

	return types.PodDescribe{
		Cluster:    cluster,
		Namespace:  p.Namespace,
		Name:       p.Name,
		Node:       p.Spec.NodeName,
		Status:     string(p.Status.Phase),
		Ready:      fmt.Sprintf("%d/%d", ready, total),
		Restarts:   sumRestarts(p.Status.ContainerStatuses),
		Age:        now.Sub(p.CreationTimestamp.Time),
		IP:         p.Status.PodIP,
		QoSClass:   string(p.Status.QOSClass),
		Containers: containers,
	}
}
