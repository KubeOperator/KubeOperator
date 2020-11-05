package grafana

import (
	"encoding/json"
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	uuid "github.com/satori/go.uuid"
)

type Dashboard struct {
	Templating    map[string][]map[string]interface{} `json:"templating"`
	Annotations   map[string]interface{}              `json:"annotations"`
	Time          map[string]interface{}              `json:"time"`
	Timepicker    map[string]interface{}              `json:"timepicker"`
	Timezone      string                              `json:"timezone"`
	Title         string                              `json:"title"`
	Uid           string                              `json:"uid"`
	Editable      bool                                `json:"editable"`
	Links         []interface{}                       `json:"links"`
	Panels        []map[string]interface{}            `json:"panels"`
	Version       int                                 `json:"version"`
	SchemaVersion int                                 `json:"schema_version"`
	Tags          []string                            `json:"tags"`
}

func NewDashboard(dataSourceName string) *Dashboard {
	var dashboard Dashboard
	_ = json.Unmarshal([]byte(constant.DefaultDashboardTemplate), &dashboard)
	for i := range dashboard.Panels {
		dashboard.Panels[i]["datasource"] = dataSourceName
	}
	for _, v := range dashboard.Templating {
		for i := range v {
			v[i]["datasource"] = dataSourceName
		}
	}
	dashboard.Title = dataSourceName
	dashboard.Uid = uuid.NewV4().String()
	return &dashboard
}

type DataSource struct {
	Name      string `json:"name"`
	Type      string `json:"type"`
	Url       string `json:"url"`
	Access    string `json:"access"`
	BasicAuth bool   `json:"basic_auth"`
}

func NewDataSource(name string, url string) *DataSource {
	return &DataSource{
		Name:      name,
		Type:      "prometheus",
		Url:       url,
		Access:    "proxy",
		BasicAuth: false,
	}
}

type CreateDashboardRequest struct {
	Dashboard Dashboard
	Overwrite bool
}
