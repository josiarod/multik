package types

import "time"

type ListOpts struct {
	Namespace     string
	AllNamespaces bool
	LabelSelector string
}

type PodRow struct {
	Cluster   string        `json:"cluster"`
	Namespace string        `json:"namespace"`
	Name      string        `json:"name"`
	Ready     string        `json:"ready"`
	Status    string        `json:"status"`
	Restarts  int32         `json:"restarts"`
	Age       time.Duration `json:"age"`
	Node      string        `json:"node"`
	Image     string        `json:"image"`
}
