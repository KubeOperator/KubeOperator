package constant

const (
	SMALL     = "small"
	MEDIUM    = "medium"
	LARGE     = "large"
	XLARGE    = "xlarge"
	XXLARGE   = "2xlarge"
	XXXXLARGE = "4xlarge"
	SINGLE    = "SINGLE"
	MULTIPLE  = "MULTIPLE"
)

type VmConfig struct {
	ID     int `json:"id"`
	Cpu    int `json:"cpu"`
	Memory int `json:"memory"`
	Disk   int `json:"disk"`
}

var VmConfigList = map[string]VmConfig{
	SMALL: {
		Cpu:    2,
		Memory: 8,
		Disk:   50,
	},
	MEDIUM: {
		Cpu:    4,
		Memory: 16,
		Disk:   50,
	},
	LARGE: {
		Cpu:    8,
		Memory: 32,
		Disk:   50,
	},
	XLARGE: {
		Cpu:    16,
		Memory: 64,
		Disk:   50,
	},
	XXLARGE: {
		Cpu:    32,
		Memory: 128,
		Disk:   50,
	},
	XXXXLARGE: {
		Cpu:    64,
		Memory: 256,
		Disk:   50,
	},
}
