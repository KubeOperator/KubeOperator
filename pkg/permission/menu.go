package permission

type MenuRole struct {
	Menu  string   `json:"menu"`
	Roles []string `json:"roles"`
}

type UserMenu struct {
	ProjectName string   `json:"projectName"`
	ProjectId   string   `json:"projectId"`
	Menus       []string `json:"menus"`
}

var MenuRoles = `
[
    {
      "menu": "project",
      "roles": [
        "PROJECT_MANAGER",
        "CLUSTER_MANAGER"
      ]
    },
    {
      "menu": "project-member",
      "roles": [
        "PROJECT_MANAGER"
      ]
    },
    {
      "menu": "project-resource",
      "roles": [
        "PROJECT_MANAGER"
      ]
    },
    {
      "menu": "cluster",
      "roles": [
        "PROJECT_MANAGER",
        "CLUSTER_MANAGER"
      ]
    },
    {
      "menu": "cluster",
      "roles": [
        "CLUSTER_MANAGER",
        "PROJECT_MANAGER"
      ]
    },
    {
      "menu": "cluster-overview",
      "roles": [
        "CLUSTER_MANAGER",
        "PROJECT_MANAGER"
      ]
    },
    {
      "menu": "cluster-node",
      "roles": [
        "CLUSTER_MANAGER",
        "PROJECT_MANAGER"
      ]
    },
    {
      "menu": "cluster-namespace",
      "roles": [
        "CLUSTER_MANAGER",
        "PROJECT_MANAGER"
      ]
    },
    {
      "menu": "cluster-storage",
      "roles": [
        "CLUSTER_MANAGER",
        "PROJECT_MANAGER"
      ]
    },
    {
      "menu": "cluster-logging",
      "roles": [
        "CLUSTER_MANAGER",
        "PROJECT_MANAGER"
      ]
    },
    {
      "menu": "cluster-monitor",
      "roles": [
        "CLUSTER_MANAGER",
        "PROJECT_MANAGER"
      ]
    },
    {
      "menu": "cluster-catalog",
      "roles": [
        "CLUSTER_MANAGER",
        "PROJECT_MANAGER"
      ]
    },
    {
      "menu": "cluster-repository",
      "roles": [
        "CLUSTER_MANAGER",
        "PROJECT_MANAGER"
      ]
    },
    {
      "menu": "cluster-repository-chartmuseum",
      "roles": [
        "CLUSTER_MANAGER",
        "PROJECT_MANAGER"
      ]
    },
    {
      "menu": "cluster-repository-registry",
      "roles": [
        "CLUSTER_MANAGER",
        "PROJECT_MANAGER"
      ]
    },
    {
      "menu": "cluster-tool",
      "roles": [
        "CLUSTER_MANAGER",
        "PROJECT_MANAGER"
      ]
    },
    {
      "menu": "cluster-dashboard",
      "roles": [
        "CLUSTER_MANAGER",
        "PROJECT_MANAGER"
      ]
    },
    {
      "menu": "host",
      "roles": [
        "CLUSTER_MANAGER",
        "PROJECT_MANAGER"
      ]
    },
    {
      "menu": "setting",
      "roles": [
        "CLUSTER_MANAGER",
        "PROJECT_MANAGER"
      ]
    },
    {
      "menu": "setting-system",
      "roles": [
        "CLUSTER_MANAGER",
        "PROJECT_MANAGER"
      ]
    },
    {
      "menu": "setting-credential",
      "roles": [
        "CLUSTER_MANAGER",
        "PROJECT_MANAGER"
      ]
    },
    {
      "menu": "setting-deploy",
      "roles": [
        "CLUSTER_MANAGER",
        "PROJECT_MANAGER"
      ]
    },
    {
      "menu": "setting-region",
      "roles": [
        "CLUSTER_MANAGER",
        "PROJECT_MANAGER"
      ]
    },
    {
      "menu": "setting-zone",
      "roles": [
        "CLUSTER_MANAGER",
        "PROJECT_MANAGER"
      ]
    },
    {
      "menu": "setting-plan",
      "roles": [
        "CLUSTER_MANAGER",
        "PROJECT_MANAGER"
      ]
    },
    {
      "menu": "user",
      "roles": [
        "CLUSTER_MANAGER",
        "PROJECT_MANAGER"
      ]
    }
  ]
`
