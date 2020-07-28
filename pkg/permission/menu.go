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
      "menu": "PROJECT",
      "roles": [
        "PROJECT_MANAGER",
        "CLUSTER_MANAGER"
      ]
    },
    {
      "menu": "PROJECT-MEMBER",
      "roles": [
        "PROJECT_MANAGER"
      ]
    },
    {
      "menu": "PROJECT-RESOURCE",
      "roles": [
        "PROJECT_MANAGER"
      ]
    },
    {
      "menu": "CLUSTER",
      "roles": [
        "PROJECT_MANAGER",
        "CLUSTER_MANAGER"
      ]
    },
    {
      "menu": "CLUSTER-OVERVIEW",
      "roles": [
        "CLUSTER_MANAGER",
        "PROJECT_MANAGER"
      ]
    },
    {
      "menu": "CLUSTER-NODE",
      "roles": [
        "CLUSTER_MANAGER",
        "PROJECT_MANAGER"
      ]
    },
    {
      "menu": "CLUSTER-NAMESPACE",
      "roles": [
        "CLUSTER_MANAGER",
        "PROJECT_MANAGER"
      ]
    },
    {
      "menu": "CLUSTER-STORAGE",
      "roles": [
        "CLUSTER_MANAGER",
        "PROJECT_MANAGER"
      ]
    },
    {
      "menu": "CLUSTER-LOGGING",
      "roles": [
        "CLUSTER_MANAGER",
        "PROJECT_MANAGER"
      ]
    },
    {
      "menu": "CLUSTER-MONITOR",
      "roles": [
        "CLUSTER_MANAGER",
        "PROJECT_MANAGER"
      ]
    },
    {
      "menu": "CLUSTER-CATALOG",
      "roles": [
        "CLUSTER_MANAGER",
        "PROJECT_MANAGER"
      ]
    },
    {
      "menu": "CLUSTER-REPOSITORY",
      "roles": [
        "CLUSTER_MANAGER",
        "PROJECT_MANAGER"
      ]
    },
    {
      "menu": "CLUSTER-REPOSITORY-CHARTMUSEUM",
      "roles": [
        "CLUSTER_MANAGER",
        "PROJECT_MANAGER"
      ]
    },
    {
      "menu": "CLUSTER-REPOSITORY-REGISTRY",
      "roles": [
        "CLUSTER_MANAGER",
        "PROJECT_MANAGER"
      ]
    },
    {
      "menu": "CLUSTER-TOOL",
      "roles": [
        "CLUSTER_MANAGER",
        "PROJECT_MANAGER"
      ]
    },
    {
      "menu": "CLUSTER-DASHBOARD",
      "roles": [
        "CLUSTER_MANAGER",
        "PROJECT_MANAGER"
      ]
    },
    {
      "menu": "HOST",
      "roles": [
      ]
    },
    {
      "menu": "SETTING",
      "roles": [
      ]
    },
    {
      "menu": "SETTING-SYSTEM",
      "roles": [
      ]
    },
    {
      "menu": "SETTING-CREDENTIAL",
      "roles": [
      ]
    },
    {
      "menu": "DEPLOY",
      "roles": [
      ]
    },
    {
      "menu": "DEPLOY-REGION",
      "roles": [
      ]
    },
    {
      "menu": "DEPLOY-ZONE",
      "roles": [
      ]
    },
    {
      "menu": "DEPLOY-PLAN",
      "roles": [
      ]
    },
    {
      "menu": "USER",
      "roles": [
      ]
    }
  ]
`
