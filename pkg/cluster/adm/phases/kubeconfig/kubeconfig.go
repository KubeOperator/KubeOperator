package kubeconfig

import (
	"bytes"
	"k8s.io/apimachinery/pkg/runtime"
	"ko3-gin/pkg/util/ssh"
	clientcmdlatest "k8s.io/client-go/tools/clientcmd/api/latest"

)

type Option struct {
	MasterEndpoint string
	ClusterName    string
	CACert         []byte
	Token          string
}

func Install(s ssh.Interface, option *Option) error {
	config := CreateWithToken(option.MasterEndpoint, option.ClusterName, "kubernetes-admin", option.CACert, option.Token)
	data, err := runtime.Encode(clientcmdlatest.Codec, config)
	if err != nil {
		return err
	}
	err = s.WriteFile(bytes.NewReader(data), "/root/.kube/config")
	if err != nil {
		return err
	}

	return nil
}
