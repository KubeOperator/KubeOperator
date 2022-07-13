package helm

import (
	"fmt"
	"testing"

	"github.com/KubeOperator/KubeOperator/pkg/config"
	"github.com/KubeOperator/KubeOperator/pkg/db"
	"github.com/KubeOperator/KubeOperator/pkg/logger"
	"github.com/spf13/viper"
)

func GetClient() (*Client, error) {
	return NewClient(&Config{
		Host:        "172.16.10.49:8443",
		BearerToken: "eyJhbGciOiJSUzI1NiIsImtpZCI6IlhOd1dFbkJDVEk2WjJ0Q0pDcGt2Y250M2x6dXBTY29zUkZ3Z0Ezcjd3U00ifQ.eyJpc3MiOiJrdWJlcm5ldGVzL3NlcnZpY2VhY2NvdW50Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9uYW1lc3BhY2UiOiJrdWJlLXN5c3RlbSIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VjcmV0Lm5hbWUiOiJrby1hZG1pbi10b2tlbi0ydG1qcCIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VydmljZS1hY2NvdW50Lm5hbWUiOiJrby1hZG1pbiIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VydmljZS1hY2NvdW50LnVpZCI6ImIxYjVlZDAzLTBkZmQtNDkxNi04NTJjLWMyODUzMDZhNmEwMSIsInN1YiI6InN5c3RlbTpzZXJ2aWNlYWNjb3VudDprdWJlLXN5c3RlbTprby1hZG1pbiJ9.eaRiGuZXqCwC9Eb4VQZvahcYCEfwrEBq6_nR7_F3lX2yQcWsLV-gEYVb04euZTfsvFaDYCRAquaAgEPS25g4lZF26JODGXpoIFFYpdwVMbjMDlDy3QB7LA3tXMXPU6w40jAE_IApmrq5pw-VUFrNdQU6KlO0mZsK3nDYnNTHK6KC1dyc69qb6awxYMb6xXPfINiksdHyUzq5mYo6PGAoLaJM4Vs_dz3iUIdQYBpdd3vZ0XlBkW0cz3ye7vDLbDeVx89E1tDZj7Et0a8pxMaE_YOwm-qCtJPEqw2Wjv-z33CD42AZaZ17td20oWfq3Lgl2Hr4769Xec21nfYkATijsg",
	})
}

func DbInit() {
	config.Init()
	dbi := db.InitDBPhase{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetInt("db.port"),
		Name:     viper.GetString("db.name"),
		User:     viper.GetString("db.user"),
		Password: viper.GetString("db.password"),
	}
	err := dbi.Init()
	if err != nil {
		logger.Log.Fatal(err)
	}
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
	r, err := h.Uninstall("efk")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(r.Info)
}

// func TestClient_Install(t *testing.T) {
// 	DbInit()
// 	h, err := GetClient()
// 	if err != nil {
// 		t.Error(err)
// 	}

// 	valueMap := map[string]interface{}{}
// 	var valueStrings []string
// 	for k, v := range valueMap {
// 		str := fmt.Sprintf("%s=%v", k, v)
// 		valueStrings = append(valueStrings, str)
// 	}
// 	valueMap = map[string]interface{}{}
// 	for _, str := range valueStrings {
// 		_ = strvals.ParseInto(str, valueMap)
// 	}
// 	r, err := h.Install("dashboard", "nexus/kubernetes-dashboard", valueMap)
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	fmt.Println(r.Name)
// }

// func TestClient_AddRepo(t *testing.T) {
// 	DbInit()
// 	err := updateRepo()
// 	if err != nil {
// 		logger.Log.Fatal(err)
// 	}
// }
