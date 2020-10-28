package dto

type ClusterManifest struct {
	Name        string        `json:"name"`
	Version     string        `json:"version"`
	IsActive    bool          `json:"isActive"`
	CoreVars    []NameVersion `json:"coreVars"`
	NetworkVars []NameVersion `json:"networkVars"`
	OtherVars   []NameVersion `json:"otherVars"`
}

type ClusterManifestUpdate struct {
	Name     string `json:"name"`
	IsActive bool   `json:"isActive"`
}

type NameVersion struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}
