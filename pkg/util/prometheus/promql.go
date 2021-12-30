package prometheus

import (
	"strings"

	"github.com/KubeOperator/KubeOperator/pkg/dto"
)

var promQLClusterTemplates = map[string]string{
	//cluster
	"cpu_used":    "avg by (value) (100 - (avg by (instance) (irate(node_cpu_seconds_total{mode='idle'}[5m])) * 100))",
	"memory_used": "(1 - sum(node_memory_MemFree_bytes+ node_memory_Cached_bytes + node_memory_Buffers_bytes + node_memory_SReclaimable_bytes)  /  sum(node_memory_MemTotal_bytes)) * 100",
	"load1":       `avg (node_load1)`,
	"load5":       `avg (node_load5)`,
	"load15":      `avg (node_load15)`,
	"disk_used":   `(sum(max(node_filesystem_size_bytes{device=~"/dev/.*", device!~"/dev/loop\\d+"} - node_filesystem_avail_bytes{device=~"/dev/.*", device!~"/dev/loop\\d+"}) by (device, instance))) / 1000000000`,

	"disk_inode_utilisation": `(1 - sum(max(node_filesystem_files_free{device=~"/dev/.*", device!~"/dev/loop\\d+"})) / sum(max(node_filesystem_files{device=~"/dev/.*", device!~"/dev/loop\\d+"}))) * 100`,
	"disk_read_throughput":   "sum(irate(node_disk_read_bytes_total[5m]) / 1000000)",
	"disk_write_throughput":  "sum(irate(node_disk_written_bytes_total[5m]) / 1000000)",
	"disk_read_iops":         "sum(irate(node_disk_reads_completed_total[5m]))",
	"disk_write_iops":        "sum(irate(node_disk_writes_completed_total[5m]))",
	"net_bytes_transmitted":  "sum(irate(node_network_transmit_bytes_total{device!~'^(cali.+|tunl.+|dummy.+|kube.+|flannel.+|cni.+|docker.+|veth.+|lo.*)'}[5m]) / 1000000)",
	"net_bytes_received":     "sum(irate(node_network_receive_bytes_total{device!~'^(cali.+|tunl.+|dummy.+|kube.+|flannel.+|cni.+|docker.+|veth.+|lo.*)'}[5m]) / 1000000)",
}

var promQLNodeTemplates = map[string]string{
	//node
	"cpu_used":    "avg by (value) (100 - (avg by (instance) (irate(node_cpu_seconds_total{node='$1', mode='idle'}[5m])) * 100))",
	"memory_used": "(1 - sum(node_memory_MemFree_bytes{node='$1'}+ node_memory_Cached_bytes{node='$1'} + node_memory_Buffers_bytes{node='$1'} + node_memory_SReclaimable_bytes{node='$1'})  /  sum(node_memory_MemTotal_bytes{node='$1'})) * 100",
	"load1":       `avg (node_load1{node='$1'})`,
	"load5":       `avg (node_load5{node='$1'})`,
	"load15":      `avg (node_load15{node='$1'})`,
	"disk_used":   `(sum(max(node_filesystem_size_bytes{node='$1', device=~"/dev/.*", device!~"/dev/loop\\d+"} - node_filesystem_avail_bytes{node='$1', device=~"/dev/.*", device!~"/dev/loop\\d+"}) by (device, instance))) / 1000000000`,

	"disk_inode_utilisation": `(1 - sum(max(node_filesystem_files_free{node='$1', device=~"/dev/.*", device!~"/dev/loop\\d+"})) / sum(max(node_filesystem_files{node='$1', device=~"/dev/.*", device!~"/dev/loop\\d+"}))) * 100`,
	"disk_read_throughput":   "sum(irate(node_disk_read_bytes_total{node='$1'}[5m]) / 1000000)",
	"disk_write_throughput":  "sum(irate(node_disk_written_bytes_total{node='$1'}[5m]) / 1000000)",
	"disk_read_iops":         "sum(irate(node_disk_reads_completed_total{node='$1'}[5m]))",
	"disk_write_iops":        "sum(irate(node_disk_writes_completed_total{node='$1'}[5m]))",
	"net_bytes_transmitted":  "sum(irate(node_network_transmit_bytes_total{node='$1', device!~'^(cali.+|tunl.+|dummy.+|kube.+|flannel.+|cni.+|docker.+|veth.+|lo.*)'}[5m]) / 1000000)",
	"net_bytes_received":     "sum(irate(node_network_receive_bytes_total{node='$1', device!~'^(cali.+|tunl.+|dummy.+|kube.+|flannel.+|cni.+|docker.+|veth.+|lo.*)'}[5m]) / 1000000)",
}

func makeExpr(metric string, opts *dto.QueryOptions) string {
	switch opts.Level {
	case "cluster":
		return promQLClusterTemplates[metric]
	case "node":
		return makeNodeMetricExpr(metric, opts)
	default:
		return promQLClusterTemplates[metric]
	}
}

func makeNodeMetricExpr(metric string, opts *dto.QueryOptions) string {
	nodeSelector := opts.NodeName
	return strings.Replace(promQLNodeTemplates[metric], "$1", nodeSelector, -1)
}
