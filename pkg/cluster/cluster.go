package cluster

import "time"

type Cluster struct {
	Name   string
	Spec   Spec
	Status Status
}

type Spec struct {
	Version                    string
	NetworkType                NetworkType
	Machines                   []Machine
	ClusterCIDR                string
	ServiceCIDR                string
	DNSDomain                  string
	DockerExtraArgs            map[string]string
	KubeletExtraArgs           map[string]string
	APIServerExtraArgs         map[string]string
	ControllerManagerExtraArgs map[string]string
	SchedulerExtraArgs         map[string]string
}

type Machine struct {
	IP         string
	Port       int
	Username   string
	Password   []byte
	PrivateKey []byte
	PassPhrase []byte
	Labels     map[string]string
}

type Status struct {
	Version     string
	Phase       Phase
	Message     string
	Reason      string
	Addresses   []Address
	ServiceCIDR string
	Conditions  []Condition
}

type NetworkType string

const (
	// NetworkFlannel indicates the communication network using the flannel to establish the pod between nodes.
	NetworkFlannel NetworkType = "Flannel"
	// NetworkCalico indicates the communication network using the calico to establish the pod between nodes.
	NetworkCalico NetworkType = "Calico"
	// NetworkIPIP indicates the communication network using the IPIP to establish the pod between nodes.
	NetworkIPIP NetworkType = "IPIP"
)

type Phase string

const (
	// ClusterRunning is the normal running phase.
	ClusterRunning Phase = "Running"
	// ClusterInitializing is the initialize phase.
	ClusterInitializing Phase = "Initializing"
	// ClusterFailed is the failed phase.
	ClusterFailed Phase = "Failed"
	// ClusterTerminating means the cluster is undergoing graceful termination.
	ClusterTerminating Phase = "Terminating"
)

type Address struct {
	Host string
	Port int
}
type Condition struct {
	Type               string
	Status             ConditionStatus
	Reason             string
	Message            string
	LastProbeTime      time.Time
	LastTransitionTime time.Time
}

type ConditionStatus string

const (
	ConditionTrue    ConditionStatus = "True"
	ConditionFalse   ConditionStatus = "False"
	ConditionUnknown ConditionStatus = "Unknown"
)
