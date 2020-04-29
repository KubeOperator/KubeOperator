package util

import "ko3-gin/pkg/cluster"

type nodeSelector struct {
	nodes []cluster.Node
}

func NewNodeSelector(nodes []cluster.Node) *nodeSelector {
	return &nodeSelector{nodes: nodes}
}

func (ns *nodeSelector) SelectNodes(label string, value string) []cluster.Node {
	var result []cluster.Node
	for _, node := range ns.nodes {
		if node.LabelValue(label) == value {
			result = append(result, node)
		}
	}
	return result
}

func (ns *nodeSelector) SelectFistNode(label string, value string) cluster.Node {
	var node cluster.Node
	for _, n := range ns.nodes {
		if n.LabelValue(label) == value {
			node = n
			break
		}
	}
	return node
}
