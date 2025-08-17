package types

import "time"

type ContainerInfo struct {
	Name   string   `json:"name"`
	Image  string   `json:"image"`
	Ports  []string `json:"ports,omitempty"`
	EnvLen int      `json:"env_len,omitempty"`
}

type PodDescribe struct {
	Cluster    string          `json:"cluster"`
	Namespace  string          `json:"namespace"`
	Name       string          `json:"name"`
	Node       string          `json:"node"`
	Status     string          `json:"status"`
	Ready      string          `json:"ready"`
	Restarts   int32           `json:"restarts"`
	Age        time.Duration   `json:"age"`
	IP         string          `json:"ip,omitempty"`
	QoSClass   string          `json:"qosclass,omitempty"`
	Containers []ContainerInfo `json:"containers"`
}
