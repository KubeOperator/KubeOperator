package permission

type Permission struct {
	ResourceType  string          `json:"resourceType"`
	OperationAuth []OperationAuth `json:"operationAuth"`
}

type OperationAuth struct {
	Operation string   `json:"operation"`
	Roles     []string `json:"roles"`
}

var Permissions = `
{
  "resourceType": "cluster",
  "operationAuth": [
	{
	  "operation": "READ",
	  "roles": [
		"SYSTEMADMIN",
		"PROJECTMANAGER",
		"CLUSTERMANAGER"
	  ]
	}
  ]
}
`
