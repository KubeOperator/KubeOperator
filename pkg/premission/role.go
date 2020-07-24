package premission

const (
	SYSTEMADMIN    = "System Admin"
	PROJECTMANAGER = "Project Manager"
	CLUSTERMANAGER = "Cluster Manager"
)

type Role struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Resources string `json:"resources"`
}

type Menu struct {
	ID  int    `json:"id"`
	Key string `json:"key"`
}
