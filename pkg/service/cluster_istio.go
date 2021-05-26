package service

import (
	"context"
	"encoding/json"

	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/db"
	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/KubeOperator/KubeOperator/pkg/model"
	"github.com/KubeOperator/KubeOperator/pkg/service/cluster/istios"
	"github.com/KubeOperator/KubeOperator/pkg/util/helm"
	kubernetesUtil "github.com/KubeOperator/KubeOperator/pkg/util/kubernetes"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ClusterIstioService interface {
	List(clusterName string) ([]dto.ClusterIstio, error)
	Enable(clusterName string, istios []dto.ClusterIstio) ([]dto.ClusterIstio, error)
	Disable(clusterName string, istio []dto.ClusterIstio) ([]dto.ClusterIstio, error)
}

func NewClusterIstioService() ClusterIstioService {
	return &clusterIstioService{
		clusterService: NewClusterService(),
	}
}

type clusterIstioService struct {
	clusterService ClusterService
}

func (c clusterIstioService) List(clusterName string) ([]dto.ClusterIstio, error) {
	var (
		istioDtos []dto.ClusterIstio
		istios    []model.ClusterIstio
		cluster   model.Cluster
	)
	if err := db.DB.Where("name = ?", clusterName).Find(&cluster).Error; err != nil {
		return istioDtos, err
	}
	if err := db.DB.Where("cluster_id = ?", cluster.ID).Find(&istios).Error; err != nil {
		return istioDtos, err
	}
	for _, m := range istios {
		d := dto.ClusterIstio{ClusterIstio: m}
		d.Vars = map[string]interface{}{}
		_ = json.Unmarshal([]byte(m.Vars), &d.Vars)
		istioDtos = append(istioDtos, d)
	}
	return istioDtos, nil
}

func (c clusterIstioService) Enable(clusterName string, istioDtos []dto.ClusterIstio) ([]dto.ClusterIstio, error) {
	cluster, endpoints, secret, err := c.getBaseParams(clusterName)
	if err != nil {
		return istioDtos, err
	}
	if err := getNs(endpoints, secret, constant.IstioNamespace); err != nil {
		return istioDtos, err
	}

	helminfo, err := NewIstioHelmInfo(cluster.Cluster, endpoints, secret.ClusterSecret, constant.IstioNamespace)
	if err != nil {
		return istioDtos, err
	}

	// base chart必须最先安装
	for i := 0; i < len(istioDtos); i++ {
		if istioDtos[i].ClusterIstio.Name == "base" {
			buf, _ := json.Marshal(&istioDtos[i].Vars)
			istioDtos[i].ClusterIstio.Vars = string(buf)
			istioDtos[i].ClusterIstio.ClusterID = cluster.ID
			istioDtos[i].ClusterIstio.Status = constant.ClusterInitializing
			base := istios.NewBaseInterface(&istioDtos[i].ClusterIstio, helminfo)
			if err = saveIstio(&istioDtos[i].ClusterIstio); err != nil {
				return istioDtos, err
			}
			c.doInstall(base, &istioDtos[i].ClusterIstio)
			break
		}
	}

	var ct istios.IstioInterface
	for i := 0; i < len(istioDtos); i++ {
		buf, _ := json.Marshal(&istioDtos[i].Vars)
		istioDtos[i].ClusterIstio.Vars = string(buf)
		istioDtos[i].ClusterIstio.ClusterID = cluster.ID
		switch istioDtos[i].ClusterIstio.Name {
		case "base":
			continue
		case "pilot":
			ct = istios.NewPilotInterface(&istioDtos[i].ClusterIstio, helminfo)
		case "ingress":
			ct = istios.NewIngressInterface(&istioDtos[i].ClusterIstio, helminfo)
		case "egress":
			ct = istios.NewEgressInterface(&istioDtos[i].ClusterIstio, helminfo)
		}
		if err != nil {
			return istioDtos, err
		}
		if istioDtos[i].Operation == "enable" {
			istioDtos[i].ClusterIstio.Status = constant.ClusterInitializing
			go c.doInstall(ct, &istioDtos[i].ClusterIstio)
			if err = saveIstio(&istioDtos[i].ClusterIstio); err != nil {
				return istioDtos, err
			}
		} else if istioDtos[i].Operation == "disable" {
			istioDtos[i].ClusterIstio.Status = constant.ClusterTerminated
			go c.doUninstall(ct, &istioDtos[i].ClusterIstio)
			if err = saveIstio(&istioDtos[i].ClusterIstio); err != nil {
				return istioDtos, err
			}
		}
	}
	return istioDtos, err
}

func (c clusterIstioService) Disable(clusterName string, istioDtos []dto.ClusterIstio) ([]dto.ClusterIstio, error) {
	cluster, endpoints, secret, err := c.getBaseParams(clusterName)
	if err != nil {
		return istioDtos, err
	}

	helminfo, err := NewIstioHelmInfo(cluster.Cluster, endpoints, secret.ClusterSecret, constant.IstioNamespace)
	if err != nil {
		return istioDtos, err
	}

	var ct istios.IstioInterface
	for i := 0; i < len(istioDtos); i++ {
		buf, _ := json.Marshal(&istioDtos[i].Vars)
		istioDtos[i].ClusterIstio.Vars = string(buf)
		istioDtos[i].ClusterIstio.ClusterID = cluster.ID
		switch istioDtos[i].ClusterIstio.Name {
		case "base":
			ct = istios.NewBaseInterface(&istioDtos[i].ClusterIstio, helminfo)
		case "pilot":
			ct = istios.NewPilotInterface(&istioDtos[i].ClusterIstio, helminfo)
		case "ingress":
			ct = istios.NewIngressInterface(&istioDtos[i].ClusterIstio, helminfo)
		case "egress":
			ct = istios.NewEgressInterface(&istioDtos[i].ClusterIstio, helminfo)
		}
		if err != nil {
			return istioDtos, err
		}
		istioDtos[i].ClusterIstio.Status = constant.ClusterTerminated
		go c.doUninstall(ct, &istioDtos[i].ClusterIstio)
		if err = saveIstio(&istioDtos[i].ClusterIstio); err != nil {
			return istioDtos, err
		}
	}
	return istioDtos, nil
}

func (c clusterIstioService) getBaseParams(clusterName string) (dto.Cluster, []kubernetesUtil.Host, dto.ClusterSecret, error) {
	var (
		cluster   dto.Cluster
		endpoints []kubernetesUtil.Host
		secret    dto.ClusterSecret
		err       error
	)
	if err := db.DB.Where("name = ?", clusterName).Preload("Spec").Find(&cluster).Error; err != nil {
		return cluster, endpoints, secret, err
	}

	endpoints, err = c.clusterService.GetApiServerEndpoints(clusterName)
	if err != nil {
		return cluster, endpoints, secret, err
	}
	secret, err = c.clusterService.GetSecrets(clusterName)
	if err != nil {
		return cluster, endpoints, secret, err
	}

	return cluster, endpoints, secret, nil
}

func (c clusterIstioService) doInstall(p istios.IstioInterface, istio *model.ClusterIstio) {
	err := p.Install()
	if err != nil {
		istio.Status = constant.ClusterFailed
		istio.Message = err.Error()
	} else {
		istio.Status = constant.ClusterRunning
	}
	_ = saveIstio(istio)
}

func (c clusterIstioService) doUninstall(p istios.IstioInterface, istio *model.ClusterIstio) {
	_ = p.Uninstall()
	istio.Status = constant.ClusterWaiting
	_ = saveIstio(istio)
}

func saveIstio(istio *model.ClusterIstio) error {
	var item model.ClusterIstio
	notFound := db.DB.Where("cluster_id = ? AND name = ?", istio.ClusterID, istio.Name).First(&item).RecordNotFound()
	if notFound {
		if err := db.DB.Create(istio).Error; err != nil {
			return err
		}
	} else {
		istio.ID = item.ID
		if err := db.DB.Save(istio).Error; err != nil {
			return err
		}
	}
	return nil
}

func getNs(endpoints []kubernetesUtil.Host, secret dto.ClusterSecret, namespace string) error {
	kubeClient, err := kubernetesUtil.NewKubernetesClient(&kubernetesUtil.Config{
		Hosts: endpoints,
		Token: secret.KubernetesToken,
	})
	if err != nil {
		return err
	}

	ns, _ := kubeClient.CoreV1().Namespaces().Get(context.TODO(), namespace, metav1.GetOptions{})
	if ns.ObjectMeta.Name == "" {
		n := &v1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name: namespace,
			},
		}
		_, err = kubeClient.CoreV1().Namespaces().Create(context.TODO(), n, metav1.CreateOptions{})
		if err != nil {
			return err
		}
	}
	return nil
}

func NewIstioHelmInfo(cluster model.Cluster, endpoints []kubernetesUtil.Host, secret model.ClusterSecret, namespace string) (istios.IstioHelmInfo, error) {
	var p istios.IstioHelmInfo
	p.LocalhostName = constant.LocalRepositoryDomainName
	helmClient, err := helm.NewClient(&helm.Config{
		Hosts:         endpoints,
		BearerToken:   secret.KubernetesToken,
		OldNamespace:  namespace,
		Namespace:     namespace,
		Architectures: cluster.Spec.Architectures,
	})
	if err != nil {
		return p, err
	}
	p.HelmClient = helmClient
	kubeClient, _ := kubernetesUtil.NewKubernetesClient(&kubernetesUtil.Config{
		Hosts: endpoints,
		Token: secret.KubernetesToken,
	})
	p.KubeClient = kubeClient
	return p, nil
}
