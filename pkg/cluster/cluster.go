package cluster

type Cluster struct {
	Name              string
	Nodes             []Node
	ApiServerAddress  string
	ApiServerBindPort string
}


