package cluster

import "ko3-gin/pkg/host"

type Node struct {
	Host   host.Host
	Labels map[string]string
}

func (n *Node) LabelValue(name string) string {
	val, ok := n.Labels[name]
	if !ok {
		return ""
	}
	return val
}

type Label struct {
	Name  string
	Value string
}


