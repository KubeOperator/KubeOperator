package kubepi

import (
	"fmt"
	"github.com/KubeOperator/KubeOperator/pkg/config"
	"testing"
)

func TestImport(t *testing.T) {
	config.Init()
	k := NewKubePi()
	opener, err := k.Open("aaaa", "https://127.0.0.1:60104", "eyJhbGciOiJSUzI1NiIsImtpZCI6ImVST3R4djA1TXhsVWdyeW1LemNTSFJ6Z3BaVEt0Q3pETm9aWndZSlNsQTQifQ.eyJpc3MiOiJrdWJlcm5ldGVzL3NlcnZpY2VhY2NvdW50Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9uYW1lc3BhY2UiOiJrdWJlLXN5c3RlbSIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VjcmV0Lm5hbWUiOiJhZG1pbi10b2tlbi1kcDU5NyIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VydmljZS1hY2NvdW50Lm5hbWUiOiJhZG1pbiIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VydmljZS1hY2NvdW50LnVpZCI6ImY0NmZhNjMwLTYzYmMtNDZkNC1hYzFmLTZiYWM1YWRhNTMyNSIsInN1YiI6InN5c3RlbTpzZXJ2aWNlYWNjb3VudDprdWJlLXN5c3RlbTphZG1pbiJ9.h3X"+
		"GpW4lxF7YykNdzr743M0GtwAqraYyWfWd00_zaWc4myYlCrIUlrfVH9qkkUCODLLlK1EA1V29fg5TB6DxuPSoBf-qdTSHUlQPgdFStgQRJGRZtozOGAhBs4KyO1rSVujcuaAOJGP7h9xi_gF1Jozu_mQvsn62oNQi"+
		"gde4CYeWiqf40avZwYQZhrXD4BspWbUZXsSK7dIN95Aduw-bB0bXFVIz9anDVrcMiXF_hmvuBsGZAU997rMPdejnOQJf7BmxI7npx3gJzqbLnDNSgGRX_2hMvQhQqzf3ffvjQDCujDtGxh2QrsXk9RXzrAdMM1B9l-YCr8m"+
		"bWIpdzIdcBw")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(opener.Redirect)
}
