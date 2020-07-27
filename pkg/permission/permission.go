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
	ProjectId   string `json:"projectId"`
	ProjectName string `json:"projectName"`
	Roles       string `json:"roles"`
}

var Permissions = `
[
    {
      "resourceType": "project",
      "operationAuth": [
        {
          "operation": "LIST",
          "roles": [
            "PROJECT_MANAGER",
            "CLUSTER_MANAGER"
          ]
        }
      ]
    },
    {
      "resourceType": "project",
      "operationAuth": [
        {
          "operation": "ADD",
          "roles": []
        }
      ]
    },
    {
      "resourceType": "project",
      "operationAuth": [
        {
          "operation": "UPDATE",
          "roles": [
            "PROJECT_MANAGER"
          ]
        }
      ]
    },
    {
      "resourceType": "project",
      "operationAuth": [
        {
          "operation": "DELETE",
          "roles": []
        }
      ]
    },
    {
      "resourceType": "project-member",
      "operationAuth": [
        {
          "operation": "LIST",
          "roles": [
            "PROJECT_MANAGER",
            "CLUSTER_MANAGER"
          ]
        }
      ]
    },
    {
      "resourceType": "project-member",
      "operationAuth": [
        {
          "operation": "ADD",
          "roles": [
            "PROJECT_MANAGER"
          ]
        }
      ]
    },
    {
      "resourceType": "project-member",
      "operationAuth": [
        {
          "operation": "UPDATE",
          "roles": [
            "PROJECT_MANAGER"
          ]
        }
      ]
    },
    {
      "resourceType": "project-member",
      "operationAuth": [
        {
          "operation": "DELETE",
          "roles": [
            "PROJECT_MANAGER"
          ]
        }
      ]
    },
    {
      "resourceType": "project-resource",
      "operationAuth": [
        {
          "operation": "LIST",
          "roles": [
            "PROJECT_MANAGER"
          ]
        }
      ]
    },
    {
      "resourceType": "project-resource",
      "operationAuth": [
        {
          "operation": "ADD",
          "roles": [
            "PROJECT_MANAGER"
          ]
        }
      ]
    },
    {
      "resourceType": "project-resource",
      "operationAuth": [
        {
          "operation": "UPDATE",
          "roles": [
            "PROJECT_MANAGER"
          ]
        }
      ]
    },
    {
      "resourceType": "project-resource",
      "operationAuth": [
        {
          "operation": "DELETE",
          "roles": [
            "PROJECT_MANAGER"
          ]
        }
      ]
    },
    {
      "resourceType": "cluster",
      "operationAuth": [
        {
          "operation": "LIST",
          "roles": [
            "PROJECT_MANAGER",
            "CLUSTER_MANAGER"
          ]
        }
      ]
    },
    {
      "resourceType": "cluster",
      "operationAuth": [
        {
          "operation": "ADD",
          "roles": [
            "PROJECT_MANAGER",
            "CLUSTER_MANAGER"
          ]
        }
      ]
    },
    {
      "resourceType": "cluster",
      "operationAuth": [
        {
          "operation": "UPDATE",
          "roles": [
            "PROJECT_MANAGER",
            "CLUSTER_MANAGER"
          ]
        }
      ]
    },
    {
      "resourceType": "cluster",
      "operationAuth": [
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
