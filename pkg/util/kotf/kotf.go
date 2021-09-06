package kotf

import (
	"fmt"

	"github.com/KubeOperator/kotf/api"
	kotfClient "github.com/KubeOperator/kotf/pkg/client"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

type Config struct {
	Cluster string
}

type Kotf struct {
	Cluster string
	Client  *kotfClient.KotfClient
}

func NewTerraform(c *Config) *Kotf {
	host := viper.GetString("kotf.host")
	port := viper.GetInt("kotf.port")
	return &Kotf{
		Cluster: c.Cluster,
		Client:  kotfClient.NewKotfClient(host, port),
	}
}

func (k *Kotf) Init(cloudType string, provider string, cloudRegion string, hosts string) (*api.KotfResult, error) {
	result, err := k.Client.Init(k.Cluster, cloudType, provider, cloudRegion, hosts)
	if err != nil {
		return result, errors.Wrap(err, fmt.Sprintf("terraform init failed: %v", err))
	}
	return result, nil

}
func (k *Kotf) Apply(regionVars string) (*api.KotfResult, error) {
	result, err := k.Client.Apply(k.Cluster, regionVars)
	if err != nil {
		return result, errors.Wrap(err, fmt.Sprintf("terraform apply failed: %v", err))
	}
	return result, nil
}

func (k *Kotf) Destroy(regionVars string) (*api.KotfResult, error) {
	result, err := k.Client.Destroy(k.Cluster, regionVars)
	if err != nil {
		return result, errors.Wrap(err, fmt.Sprintf("terraform destory failed: %v", err))
	}
	return result, nil
}
