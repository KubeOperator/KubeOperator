package cluster

import "math"

var MaxNodePodNumMap = map[int]int{
	24: 110,
	25: 64,
	26: 32,
	27: 16,
}

func GetNodeCIDRMaskSize(maxNodePodNum int) int {
	nodeCidrOccupy := math.Ceil(math.Log2(float64(maxNodePodNum)))
	nodeCIDRMaskSize := 32 - int(nodeCidrOccupy)
	return nodeCIDRMaskSize
}
