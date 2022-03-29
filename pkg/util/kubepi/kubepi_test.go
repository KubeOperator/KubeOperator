package kubepi

import (
	"fmt"
	"testing"

	"github.com/KubeOperator/KubeOperator/pkg/config"
)

func TestImport(t *testing.T) {
	config.Init()
	k := NewKubePi()
	opener, err := k.Open("songliu-test", "https://127.0.0.1:8443", "eyJhbGciOiJSUzI1NiIsImtpZCI6IjdBWkFJOTF0cnJRY1ljSGVvc3NxeXJHMXV3NGVMRWhqcmNqdU1vbU9HXzQifQ.eyJpc3MiOiJrdWJlcm5ldGVzL3NlcnZpY2VhY2NvdW50Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9uYW1lc3BhY2UiOiJrdWJlLXN5c3RlbSIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VjcmV0Lm5hbWUiOiJrby1hZG1pbi10b2tlbi12eG5qZiIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VydmljZS1hY2NvdW50Lm5hbWUiOiJrby1hZG1pbiIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VydmljZS1hY2NvdW50LnVpZCI6ImQyOWMwMzgzLTNkZGYtNDI1Yi1iYzEzLWIxNWQ2M2QxZDA3OSIsInN1YiI6InN5c3RlbTpzZXJ2aWNlYWNjb3VudDprdWJlLXN5c3RlbTprby1hZG1pbiJ9.GTVU_Mxwl7VyFjoV1UrT8VCQ88_9aPAgEEHMICuGGjBAzroOa-pp-SDMVoINxgPBFrMlARMCo6aEHW8gJ62YgKNA8bWRCikePSs6aCoRZvtja79ruKxtWwQcLu-wcBwJav9wCHJhdFz_Xaph2rGFPVHLblC5sun-zvXi3eMLwAWmKTep9cwwNSo9femSWP4CqaJXBBAmbBp0B7PhSCmRqfLokePjPbnaq_pyPYj1oU2jez9SxGrD2OisWT6IviZr2fTRlSbWex2lrhHcDDW73Xb2AuxyeqVPe6ycEenFCX88ZOdU8gXR59NE_fwgdoXNFPODEC5ZJoPbiwNSGXqISQ")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(opener.Redirect)
}

func TestDelete(t *testing.T) {
	config.Init()
	k := NewKubePi()
	err := k.Close("songliu-test", "https://127.0.0.1:8443")
	if err != nil {
		t.Error(err)
	}
}
