package helm

import (
	"fmt"
	"path"
	"testing"
)

func GetClient() (*Client, error) {
	return NewClient(Config{
		ApiServer:   "https://172.16.10.184:8443",
		BearerToken: "eyJhbGciOiJSUzI1NiIsImtpZCI6IlRjVnRxaFFadGNrbVdSb3ZrZl9GYTlWUm9vVTBxeGVYZ09kSDBHbl84RU0ifQ.eyJpc3MiOiJrdWJlcm5ldGVzL3NlcnZpY2VhY2NvdW50Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9uYW1lc3BhY2UiOiJrdWJlLXN5c3RlbSIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VjcmV0Lm5hbWUiOiJrby1hZG1pbi10b2tlbi1ocjRsaCIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VydmljZS1hY2NvdW50Lm5hbWUiOiJrby1hZG1pbiIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VydmljZS1hY2NvdW50LnVpZCI6ImZlYTU2YjQ5LTg0ZjgtNDY0NC1hMDQxLTkzZWI1MGIxZWJlNSIsInN1YiI6InN5c3RlbTpzZXJ2aWNlYWNjb3VudDprdWJlLXN5c3RlbTprby1hZG1pbiJ9.fv8TAfrHh85-RjSMa1fbTalUWH_WCRzAqnBbxiCcNg7bP9qIfsGKMYLauE6O6eYWxwxfCCFoOJnaYZR1Qr-sHpEBEwKzuGonZekjYzzSzjkXc9eNchsMCkJ4klzra092_14-KdejgmvxRkF_vjcR6DrVB_P7E7s8UIbWM2TVn-EZ6tMIep8hq-3Qk5sh1WtILS4YF4BFKG9hNczkBIVBctExjLFrhzfQDDeTWMtLqb5v3QsJqWGA-ZiafsTfGlHpROhXPU4wKzInAD82BXO_i4RkpJqe8G0I_PhCuj_l2_tN4auqdpsQs31fyNLmXgOeNVmCBQSb5WYJnQq7tcHoZA",
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
	fmt.Println(len(r))
}

func TestClient_Uninstall(t *testing.T) {
	h, err := GetClient()
	if err != nil {
		t.Error(err)
	}
	r, err := h.Uninstall("test")
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
	chart, err := LoadCharts(path.Join("../../../resource/charts/prometheus-11.6.0.tgz"))
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
