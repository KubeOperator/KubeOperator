package permission

type Permission struct {
	ResourceType  string          `json:"resourceType"`
	OperationAuth []OperationAuth `json:"operationAuth"`
}

type OperationAuth struct {
	Operation string   `json:"operation"`
	Roles     []string `json:"roles"`
}

type UserPermission struct {
	ProjectId           string               `json:"projectId"`
	ProjectName         string               `json:"projectName"`
	UserPermissionRoles []UserPermissionRole `json:"userPermissionRoles"`
	ProjectRole         string               `json:"projectRole"`
}

type UserPermissionRole struct {
	ResourceType string   `json:"operation"`
	Roles        []string `json:"roles"`
}

var PermissionRoles = `
[
    {
      "resourceType": "PROJECT",
      "operationAuth": [
        {
          "operation": "LIST",
          "roles": [
            "PROJECT_MANAGER",
            "CLUSTER_MANAGER"
          ]
        },
        {
          "operation": "CREATE",
          "roles": []
        },
        {
          "operation": "UPDATE",
          "roles": [
            "PROJECT_MANAGER"
          ]
        },
        {
          "operation": "DELETE",
          "roles": []
        }
      ]
    },
    {
      "resourceType": "PROJECT-MEMBER",
      "operationAuth": [
        {
          "operation": "LIST",
          "roles": [
            "PROJECT_MANAGER",
            "CLUSTER_MANAGER"
          ]
        },
        {
          "operation": "CREATE",
          "roles": [
            "PROJECT_MANAGER"
          ]
        },
        {
          "operation": "UPDATE",
          "roles": [
            "PROJECT_MANAGER"
          ]
        },
        {
          "operation": "DELETE",
          "roles": [
            "PROJECT_MANAGER"
          ]
        }
      ]
    },
    {
      "resourceType": "PROJECT-RESOURCE",
      "operationAuth": [
        {
          "operation": "LIST",
          "roles": [
            "PROJECT_MANAGER"
          ]
        },
        {
          "operation": "CREATE",
          "roles": [
            "PROJECT_MANAGER"
          ]
        },
        {
          "operation": "UPDATE",
          "roles": [
            "PROJECT_MANAGER"
          ]
        },
        {
          "operation": "DELETE",
          "roles": [
            "PROJECT_MANAGER"
          ]
        }
      ]
    },
    {
      "resourceType": "CLUSTER",
      "operationAuth": [
        {
          "operation": "LIST",
          "roles": [
            "PROJECT_MANAGER",
            "CLUSTER_MANAGER"
          ]
        },
        {
          "operation": "CREATE",
          "roles": [
            "PROJECT_MANAGER",
            "CLUSTER_MANAGER"
          ]
        },
        {
          "operation": "UPDATE",
          "roles": [
            "PROJECT_MANAGER",
            "CLUSTER_MANAGER"
          ]
        },
        {
          "operation": "DELETE",
          "roles": [
            "PROJECT_MANAGER",
            "CLUSTER_MANAGER"
          ]
        }
      ]
    }
  ]
`
