package kotf

import (
	"github.com/KubeOperator/kotf/api"
	kotfClient "github.com/KubeOperator/kotf/pkg/client"
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

func (k *Kotf) Init(cloudType string, provider string, cloudRegion string, hosts string) (*api.Result, error) {
	return k.Client.Init(k.Cluster, cloudType, provider, cloudRegion, hosts)

}
func (k *Kotf) Apply() (*api.Result, error) {
	return k.Client.Apply(k.Cluster)
}

func (k *Kotf) Destroy() (*api.Result, error) {
	return k.Client.Destroy(k.Cluster)
}
