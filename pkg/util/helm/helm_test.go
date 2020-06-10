package helm

import (
	"testing"
)

func debug(format string, v ...interface{}) {}

func TestClient_Install(t *testing.T) {
	//cf := genericclioptions.NewConfigFlags(true)
	//inscure := true
	//cf.APIServer = &apiServer
	//cf.BearerToken = &token
	//cf.Insecure = &inscure
	//actionConfig := new(action.Configuration)
	//err := actionConfig.Init(cf, "default", "memory", debug)
	//if err != nil {
	//	return
	//}
	//chartRequested, err := loader.Load("harbor-1.4.0.tgz")
	//client := action.NewInstall(actionConfig)
	//client.ReleaseName = "harbor"
	//r, err := client.Run(chartRequested, map[string]interface{}{})
	//if err != nil {
	//	fmt.Println(err.Error())
	//	return
	//}
	//fmt.Println(r.Info)

}
