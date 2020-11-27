package tools

import (
	"encoding/json"
	"fmt"

	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/model"
)

const (
	FluentedElasticsearchImageName = "fluentd_elasticsearch/fluentd"
	FluentedElasticsearchTag       = "v2.8.0"
	ElasticSearchImageName         = "elasticsearch/elasticsearch"
	ElasticSearchTag               = "7.6.2"
)

type EFK struct {
	Cluster       *Cluster
	Tool          *model.ClusterTool
	LocalHostName string
}

func (c EFK) setDefaultValue() {
	values := map[string]interface{}{}
	_ = json.Unmarshal([]byte(c.Tool.Vars), &values)
	values["fluentd-elasticsearch.image.repository"] = fmt.Sprintf("%s:%d/%s", c.LocalHostName, constant.LocalDockerRepositoryPort, FluentedElasticsearchImageName)
	values["fluentd-elasticsearch.imageTag"] = FluentedElasticsearchTag
	values["elasticsearch.image"] = fmt.Sprintf("%s:%d/%s", c.LocalHostName, constant.LocalDockerRepositoryPort, ElasticSearchImageName)
	values["elasticsearch.imageTag"] = ElasticSearchTag

	if _, ok := values["elasticsearch.esJavaOpts.item"]; !ok {
		values["elasticsearch.esJavaOpts.item"] = 1
	}
	values["elasticsearch.esJavaOpts"] = fmt.Sprintf("-Xmx%vg -Xms%vg", values["elasticsearch.esJavaOpts.item"], values["elasticsearch.esJavaOpts.item"])
	delete(values, "elasticsearch.esJavaOpts.item")

	if _, ok := values["elasticsearch.volumeClaimTemplate.resources.requests.storage"]; ok {
		values["elasticsearch.volumeClaimTemplate.resources.requests.storage"] = fmt.Sprintf("%vGi", values["elasticsearch.volumeClaimTemplate.resources.requests.storage"])
	}
	str, _ := json.Marshal(&values)
	c.Tool.Vars = string(str)
}

func NewEFK(cluster *Cluster, localhostName string, tool *model.ClusterTool) (*EFK, error) {
	p := &EFK{
		Tool:          tool,
		Cluster:       cluster,
		LocalHostName: localhostName,
	}
	return p, nil
}

func (c EFK) Install() error {
	c.setDefaultValue()
	if err := installChart(c.Cluster.HelmClient, c.Tool, constant.LoggingChartName); err != nil {
		return err
	}
	if err := createRoute(c.Cluster.Namespace, constant.DefaultLoggingIngressName, constant.DefaultLoggingIngress, constant.DefaultLoggingServiceName, 9200, c.Cluster.KubeClient); err != nil {
		return err
	}
	if err := waitForStatefulSetsRunning(c.Cluster.Namespace, constant.DefaultLoggingStateSetsfulName, 1, c.Cluster.KubeClient); err != nil {
		return err
	}
	return nil
}

func (c EFK) Uninstall() error {
	return uninstall(c.Cluster.Namespace, c.Tool, constant.DefaultLoggingIngressName, c.Cluster.HelmClient, c.Cluster.KubeClient)
}
