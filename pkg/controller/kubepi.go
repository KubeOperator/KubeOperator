package controller

import (
	"net/http"

	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/KubeOperator/KubeOperator/pkg/service"
	"github.com/KubeOperator/KubeOperator/pkg/util/encrypt"
	"github.com/KubeOperator/KubeOperator/pkg/util/kubepi"
	"github.com/kataras/iris/v12/context"
)

type KubePiController struct {
	Ctx            context.Context
	KubePiService  service.KubepiService
	ClusterService service.ClusterService
}

func NewKubePiController() *KubePiController {
	return &KubePiController{
		KubePiService:  service.NewKubepiService(),
		ClusterService: service.NewClusterService(),
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

func (p KubePiController) GetJumpBy(name string) (*dto.Dashboard, error) {
	ss, err := p.KubePiService.LoadInfo(name)
	if err != nil {
		return nil, err
	}
	secrets, err := p.ClusterService.GetSecrets(name)
	if err != nil {
		return nil, err
	}
	apiServer, err := p.ClusterService.GetApiServerEndpoint(name)
	if err != nil {
		return nil, err
	}
	kubepiClient := kubepi.GetClient()
	username := ss.Name
	password, err := encrypt.StringDecrypt(ss.Password)
	if err != nil {
		return nil, err
	}
	if username != "" && password != "" {
		kubepiClient = kubepi.GetClient(kubepi.WithUsernameAndPassword(username, password))
	}
	opener, err := kubepiClient.Open(name, string(apiServer), secrets.KubernetesToken)
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
