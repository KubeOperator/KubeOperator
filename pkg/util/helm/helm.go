package helm

import (
	"helm.sh/helm/v3/pkg/chart"
	"k8s.io/client-go/kubernetes"
)

type Interface interface {
	Install(chart *chart.Chart, values map[string]interface{})
}

type Config struct {
	KubeClientSet kubernetes.Clientset
}

type Client struct {
	Config Config
}

func (c Client) Install(chart *chart.Chart, values map[string]interface{}) {
	//cf := genericclioptions.NewConfigFlags(true)
	//apiServer := "https://172.16.10.184:8443"
	//token := "eyJhbGciOiJSUzI1NiIsImtpZCI6Im5QWVVaVDhONmFMVXBTTHZtRjdERk1aY1lEcTUtQURBZ19UODRLOHhlNncifQ.eyJpc3MiOiJrdWJlcm5ldGVzL3NlcnZpY2VhY2NvdW50Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9uYW1lc3BhY2UiOiJrdWJlLXN5c3RlbSIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VjcmV0Lm5hbWUiOiJrby1hZG1pbi10b2tlbi1iODQ0ZCIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VydmljZS1hY2NvdW50Lm5hbWUiOiJrby1hZG1pbiIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VydmljZS1hY2NvdW50LnVpZCI6ImJjZjY5MjA2LTcwOWEtNDBmYS04NTZlLWNjZGMwNzJlNGEzMiIsInN1YiI6InN5c3RlbTpzZXJ2aWNlYWNjb3VudDprdWJlLXN5c3RlbTprby1hZG1pbiJ9.ugvH3kXa87OvxnNQaenSyyCQYYkE7VZ3wbuV6RAL3QfuUZvj-2CI4SJY6thLddEjD9dQ7r-tyq4FQIYOl0uj1-oqm_pNCC6Ya2Hby2O296d6StPHKzVsiG-sYDKW9nengc_GJptMZF9S51Jlb5MvNpkx6pFw1Gty8n9jpdBN_5l7qyeBwGGoJSa0sgcJPnOgSy5j8Y905fv_eT6tcJSBY0q-cptNEMsLTngZ_ikZqye5UoM6P8EvT7GtWYMPHqv8DYXVb_BEu97Xv9vC9ZF8sT9GVkbQIJLN1E_Tt9CvqlVKEPUEAEhdiWeds8-FLcutDP_x56AtMG2Lk7ltRJHszg"
	//inscure := true
	//cf.APIServer = &apiServer
	//cf.BearerToken = &token
	//cf.Insecure = &inscure
	//cfg, err := cf.ToRESTConfig()
	//if err != nil {
	//	return
	//}
	//api, err := kubernetes.NewForConfig(cfg)
	//if err!=nil{
	//	return
	//}
	//ns, err := api.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	//fmt.Println()
	//
	//installClient := action.NewInstall(&action.Configuration{
	//	KubeClient:     nil,
	//	RegistryClient: nil,
	//	Capabilities:   nil,
	//	Log:            nil,
	//})
	//_, _ = installClient.Run(chart, values)
}
