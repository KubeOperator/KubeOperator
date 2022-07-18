package controller

import (
	"fmt"
	"net/http"

	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/KubeOperator/KubeOperator/pkg/repository"
	"github.com/KubeOperator/KubeOperator/pkg/service"
	"github.com/KubeOperator/KubeOperator/pkg/util/encrypt"
	"github.com/KubeOperator/KubeOperator/pkg/util/kubepi"
	"github.com/kataras/iris/v12/context"
)

type KubePiController struct {
	Ctx            context.Context
	KubePiService  service.KubepiService
	ClusterService service.ClusterService
	clusterRepo    repository.ClusterRepository
}

func NewKubePiController() *KubePiController {
	return &KubePiController{
		KubePiService:  service.NewKubepiService(),
		ClusterService: service.NewClusterService(),
		clusterRepo:    repository.NewClusterRepository(),
	}
}

func (u *KubePiController) GetUser() (interface{}, error) {
	users, err := u.KubePiService.GetKubePiUser()
	return users, err
}

func (p KubePiController) PostBind() error {
	var req dto.BindKubePI
	err := p.Ctx.ReadJSON(&req)
	if err != nil {
		return err
	}

	if err := p.KubePiService.BindKubePi(req); err != nil {
		return err
	}
	return nil
}

func (p KubePiController) PostSearch() (*dto.BindResponse, error) {
	var req dto.SearchBind
	err := p.Ctx.ReadJSON(&req)
	if err != nil {
		return nil, err
	}

	bind, err := p.KubePiService.GetKubePiBind(req)
	if err != nil {
		return nil, err
	}
	return bind, nil
}

func (p KubePiController) PostCheckConn() error {
	var req dto.CheckConn
	err := p.Ctx.ReadJSON(&req)
	if err != nil {
		return err
	}

	return p.KubePiService.CheckConn(req)
}

func (p KubePiController) GetJumpBy(project string, clusterName string) (*dto.Dashboard, error) {
	user := p.Ctx.Values().Get("user")
	roleStr, _ := user.(dto.SessionUser)
	userInfo, err := p.KubePiService.LoadInfo(project, clusterName, roleStr.IsAdmin)
	if err != nil {
		return nil, err
	}
	cluster, err := p.clusterRepo.GetWithPreload(clusterName, []string{"SpecConf", "Secret", "Nodes", "Nodes.Host", "Nodes.Host.Credential"})
	if err != nil {
		return nil, fmt.Errorf("load cluster info failed, err: %v", err.Error())
	}
	kubepiClient := kubepi.GetClient()
	username := userInfo.Name
	password, err := encrypt.StringDecrypt(userInfo.Password)
	if err != nil {
		return nil, err
	}
	if username != "" && password != "" {
		kubepiClient = kubepi.GetClient(kubepi.WithUsernameAndPassword(username, password))
	}

	conn := kubepi.ConnInfo{
		Name:               cluster.Name,
		ApiServer:          fmt.Sprintf("https://%s:%d", cluster.SpecConf.LbKubeApiserverIp, cluster.SpecConf.KubeApiServerPort),
		AuthenticationMode: cluster.SpecConf.AuthenticationMode,
		KubernetesToken:    cluster.Secret.KubernetesToken,
		KeyDataStr:         cluster.Secret.KeyDataStr,
		CertDataStr:        cluster.Secret.CertDataStr,
		ConfigContent:      cluster.Secret.ConfigContent,
	}
	opener, err := kubepiClient.Open(conn)
	if err != nil {
		return nil, err
	}
	p.Ctx.SetCookie(&http.Cookie{
		Name:     opener.SessionCookie.Name,
		Value:    opener.SessionCookie.Value,
		Path:     opener.SessionCookie.Path,
		Expires:  opener.SessionCookie.Expires,
		HttpOnly: opener.SessionCookie.HttpOnly,
		SameSite: opener.SessionCookie.SameSite,
		MaxAge:   opener.SessionCookie.MaxAge,
	})
	return &dto.Dashboard{Url: opener.Redirect}, nil
}
