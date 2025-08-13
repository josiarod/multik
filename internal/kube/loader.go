package kube

import (
	"fmt"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type ClusterClient struct {
	Name   string
	Client *kubernetes.Clientset
}

func LoadAllContexts() ([]ClusterClient, error) {
	rules := clientcmd.NewDefaultClientConfigLoadingRules() // respects KUBECONFIG/defaults
	cfg, err := rules.Load()
	if err != nil {
		return nil, fmt.Errorf("load kubeconfig: %w", err)
	}

	var out []ClusterClient
	for name := range cfg.Contexts {
		cc := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
			rules,
			&clientcmd.ConfigOverrides{CurrentContext: name},
		)
		restCfg, err := cc.ClientConfig()
		if err != nil {
			continue // skip broken contexts
		}

		// etiquette (safe defaults)

		restCfg.QPS = 5
		restCfg.Burst = 10

		cs, err := kubernetes.NewForConfig(restCfg)
		if err != nil {
			continue
		}
		out = append(out, ClusterClient{Name: name, Client: cs})
	}

	return out, nil
}
