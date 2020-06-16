package controller

import (
	"github.com/KubeOperator/KubeOperator/pkg/service"
	"github.com/KubeOperator/KubeOperator/pkg/service/dto"
	"github.com/kataras/iris/v12"
)

type userController struct {
	ctx         iris.Context
	userService service.UserService
}

func (u userController) Get() ([]dto.User, error) {
	return u.userService.List()
}

func (u userController) GetBy(name string) (dto.User, error) {
	return u.userService.Get(name)
}

func (u userController) Post() error {
	var req dto.UserCreate
	err := u.ctx.ReadJSON(&req)
	if err != nil {
		return err
	}
	return u.userService.Create(req)
}

func (u userController) Delete(name string) error {
	return u.userService.Delete(name)
}

//func (u userController) Batch(operation string, items []dto.User) error {
//	_, err := u.userService.Batch(operation, items)
//	return err
//}
