package constant

const (
	SINGLE = "SINGLE"
)

type VmConfig struct {
	ID     int `json:"id"`
	Cpu    int `json:"cpu"`
	Memory int `json:"memory"`
	Disk   int `json:"disk"`
}
