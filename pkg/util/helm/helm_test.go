package helm

import (
	"fmt"
	"path"
	"testing"
)

func GetClient() (*Client, error) {
	return NewClient(Config{
		ApiServer:   "https://172.16.10.184:8443",
		BearerToken: "eyJhbGciOiJSUzI1NiIsImtpZCI6IlE5dVAxN2hTUjNzZ0pJcVdRU1ZtclBNb3JjNU5DeUt2UG5mVFVwNVpBRWsifQ.eyJpc3MiOiJrdWJlcm5ldGVzL3NlcnZpY2VhY2NvdW50Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9uYW1lc3BhY2UiOiJrdWJlLXN5c3RlbSIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VjcmV0Lm5hbWUiOiJrby1hZG1pbi10b2tlbi13aHQydCIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VydmljZS1hY2NvdW50Lm5hbWUiOiJrby1hZG1pbiIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VydmljZS1hY2NvdW50LnVpZCI6IjljNTFkMWQ4LWM2YzItNDNlMS1iYzk0LWYxMWQ5MDU3N2NkMSIsInN1YiI6InN5c3RlbTpzZXJ2aWNlYWNjb3VudDprdWJlLXN5c3RlbTprby1hZG1pbiJ9.b5ZwigSVZ4Yrs2YP0qnIviR0iRZNm_dPWJle6zFKxawZaZ2lIgYQe53RmeDcXdQMUMknfO2Ofgf5gPtN5gUfccZZkfXGe8v6ak1u7tH69MfUn3qohKqRzcHskCYzW1Q_CqsmH60VxeGdk_Iprmx7mJjSK4D7YqIIBfi5V9yeJWHX670OwwckBEXq0v7fiQdO4OQgtTyahULUqf4NM-9Wiv2sJpplRXSdq1xOpzHjptyZX5GpVkkbGGlf-R4KnHMi_RTm9OpZ5ZbKaf9dqgVLWu4paqVV8nThd5MvVG2mFfQDbY_an0DYucwGh16fkGE4TJBLHerzOoNkyQ761ZvbbA",
	})
}

func TestClient_List(t *testing.T) {
	h, err := GetClient()
	if err != nil {
		t.Error(err)
	}
	r, err := h.List()
	if err != nil {
		t.Error(err)
	}
	fmt.Println(r[0].Name)
}

func TestClient_Uninstall(t *testing.T) {
	h, err := GetClient()
	if err != nil {
		t.Error(err)
	}
	r, err := h.Uninstall("monitor")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(r.Info)
}

func TestClient_Install(t *testing.T) {
	h, err := GetClient()
	if err != nil {
		t.Error(err)
	}
	chart, err := LoadCharts(path.Join("resource/charts/prometheus-11.6.0.tgz"))
	if err != nil {
		t.Error(err)
	}
	values := map[string]interface{}{
		"alertmanager": map[string]interface{}{
			"enabled": false,
		},
		"server": map[string]interface{}{
			"persistentVolume": map[string]interface{}{
				"enabled": false,
			},
		},
	}
	r, err := h.Install("test", chart, values)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(r.Name)
}
