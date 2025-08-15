package kube

import (
	"context"
	"fmt"
	"time"

	"github.com/josiarod/multik/internal/types"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func ListPods(ctx context.Context, cc ClusterClient, opts types.ListOpts) (*corev1.PodList, error) {
	ns := opts.Namespace
	if opts.AllNamespaces {
		ns = "" // "" == all namespaces
	}
	listOpts := metav1.ListOptions{LabelSelector: opts.LabelSelector}
	pods, err := cc.Client.CoreV1().Pods(ns).List(ctx, listOpts)
	if err != nil {
		return nil, fmt.Errorf("%s: list pods: %w", cc.Name, err)
	}
	return pods, nil
}

func ToPodRows(cluster string, pods *corev1.PodList) []types.PodRow {
	out := make([]types.PodRow, 0, len(pods.Items))
	now := time.Now()

	for _, p := range pods.Items {
		ready, total := 0, 0
		for _, cs := range p.Status.ContainerStatuses {
			total++
			if cs.Ready {
				ready++
			}
		}
		image := ""
		if len(p.Spec.Containers) > 0 {
			image = p.Spec.Containers[0].Image
		}
		out = append(out, types.PodRow{
			Cluster:   cluster,
			Namespace: p.Namespace,
			Name:      p.Name,
			Ready:     fmt.Sprintf("%d/%d", ready, total),
			Status:    string(p.Status.Phase),
			Restarts:  sumRestarts(p.Status.ContainerStatuses),
			Age:       now.Sub(p.CreationTimestamp.Time),
			Node:      p.Spec.NodeName,
			Image:     image,
		})
	}
	return out
}

func sumRestarts(st []corev1.ContainerStatus) int32 {
	var n int32
	for _, c := range st {
		n += c.RestartCount
	}
	return n
}
