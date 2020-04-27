package model

type HostInfo struct {
	Hostname string
}

type Cri struct {
	Runtime string
	Version string
}

type Memory struct {
	Size int64
}

type CPU struct {
	Name       string
	MHz        float32
	PhysicalId int
	CoreNum    int
}
