package dto

type ClusterManifest struct {
	Name        string        `json:"name"`
	Version     string        `json:"version"`
	IsActive    bool          `json:"isActive"`
	CoreVars    []NameVersion `json:"coreVars"`
	NetworkVars []NameVersion `json:"networkVars"`
	ToolVars    []NameVersion `json:"toolVars"`
	OtherVars   []NameVersion `json:"otherVars"`
}

type ClusterManifestUpdate struct {
	Name     string `json:"name"`
	Version  string `json:"version"`
	IsActive bool   `json:"isActive"`
}

type NameVersion struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

func (c ClusterManifest) GetVars() map[string]string {
	vars := make(map[string]string)
	for _, v := range c.CoreVars {
		if v.Name == "version" {
			vars["kube_version"] = v.Version
		}
		if v.Name == "docker" {
			vars["docker_version"] = v.Version
		}
		if v.Name == "etcd" {
			vars["etcd_version"] = v.Version
		}
		if v.Name == "containerd" {
			vars["containerd_version"] = v.Version
		}
	}
	for _, v := range c.NetworkVars {
		if v.Name == "calico" {
			vars["calico_version"] = v.Version
		}
		if v.Name == "flanneld" {
			vars["flannel_version"] = v.Version
		}
	}
	for _, v := range c.OtherVars {
		if v.Name == "coredns" {
			vars["coredns_version"] = v.Version
		}
		if v.Name == "helm-v2" {
			vars["helm_v2_version"] = v.Version
		}
		if v.Name == "helm-v3" {
			vars["helm_v3_version"] = v.Version
		}
		if v.Name == "ingress-nginx" {
			vars["nginx_ingress_version"] = v.Version
		}
		if v.Name == "traefik" {
			vars["traefik_ingress_version"] = v.Version
		}
		if v.Name == "metrics-server" {
			vars["metrics_server_version"] = v.Version
		}
	}
	return vars
}
