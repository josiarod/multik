package kube

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ServerVersion returns the gitVersion string for a given clientset.
func ServerVersion(ctx context.Context, cc ClusterClient) (string, error) {
	// Two options: Discovery().ServerVersion() or /version via REST.
	// Discovery is sipler here
	info, err := cc.Client.Discovery().ServerVersion()
	if err != nil {
		return "", err
	}
	return info.GitVersion, nil
}

// Quick sanity call to ensure API is reachable
func APIServerLive(ctx context.Context, cc ClusterClient) error {
	_, err := cc.Client.CoreV1().Namespaces().List(ctx, metav1.ListOptions{Limit: 1})
	return err
}
