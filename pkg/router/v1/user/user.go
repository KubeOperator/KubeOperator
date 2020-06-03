package user

import (
	"errors"
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/i18n"
	userModel "github.com/KubeOperator/KubeOperator/pkg/model/user"
	"github.com/KubeOperator/KubeOperator/pkg/router/v1/user/serializer"
	userService "github.com/KubeOperator/KubeOperator/pkg/service/user"
	"github.com/gin-gonic/gin"
	"net/http"
)

var (
	invalidUserName  = errors.New(i18n.Tr("invalid_user_name", nil))
	deleteUserFailed = errors.New(i18n.Tr("delete_user_failed", nil))
)

// ListUser
// @Tags User
// @Summary User
// @Description List users
// @Accept  json
// @Produce json
// @Param pageNum query string false "page num"
// @Param pageSize query string false "page size"
// @Success 200 {object} serializer.ListUserResponse
// @Router /users/ [get]
func List(ctx *gin.Context) {
	page := ctx.GetBool("page")
	var models []userModel.User
	total := 0
	if page {
		pageNum := ctx.GetInt(constant.PageNumQueryKey)
		pageSize := ctx.GetInt(constant.PageSizeQueryKey)
		m, t, err := userService.Page(pageNum, pageSize)
		models = m
		total = t
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"msg": err.Error(),
			})
			return
		}
	} else {
		ms, err := userService.List()
		models = ms
		total = len(ms)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"msg": err.Error(),
			})
			return
		}
	}
	var resp = serializer.ListUserResponse{
		Items: []serializer.User{},
		Total: total,
	}
	for _, model := range models {
		resp.Items = append(resp.Items, serializer.FromModel(model))
	}

	ctx.JSON(http.StatusOK, resp)
}

// GetUser
// @Tags User
// @Summary User
// @Description Get User
// @Accept  json
// @Produce json
// @Param credential_name path string true "user name"
// @Success 200 {object} serializer.GetUserResponse
// @Router /users/{user_name} [get]
func Get(ctx *gin.Context) {
	name := ctx.Param("name")
	if name == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": invalidUserName.Error(),
		})
	}
	model, err := userService.Get(name)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"mgs": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, serializer.GetUserResponse{
		Item: serializer.FromModel(model),
	})
}

// CreateUser
// @Tags User
// @Summary User
// @Description Create a User
// @Accept  json
// @Produce json
// @Param request body serializer.CreateUserRequest true "user"
// @Success 201 {object} serializer.User
// @Router /users/ [post]
func Create(ctx *gin.Context) {
	var req serializer.CreateUserRequest
	err := ctx.ShouldBind(req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": err.Error(),
		})
		return
	}
	model := userModel.User{
		Name:     req.Name,
		Password: req.Password,
		Email:    req.Email,
		IsActive: true,
		Language: userModel.ZH,
	}

	err = userService.Save(&model)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusCreated, serializer.FromModel(model))
}

// UpdateUser
// @Tags User
// @Summary User
// @Description Update a User
// @Accept  json
// @Produce json
// @Param request body serializer.UpdateUserRequest true "user"
// @Param user_name path string true "user name"
// @Success 200 {object} serializer.User
// @Router /users/{user_name} [patch]
func Update(ctx *gin.Context) {
	var req serializer.UpdateUserRequest
	err := ctx.ShouldBind(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": err.Error(),
		})
		return
	}
	user := serializer.User{
		Name:     req.Name,
		Password: req.Password,
		Email:    req.Email,
		IsActive: req.IsActive,
		Language: req.Language,
	}
	model := serializer.ToModel(user)
	err = userService.Save(&model)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, serializer.FromModel(model))
}

// DeleteUser
// @Tags User
// @Summary User
// @Description Delete a User
// @Accept  json
// @Produce json
// @Param user_name path string true "user name"
// @Success 200 {string} string
// @Router /users/{user_name} [delete]
func Delete(ctx *gin.Context) {
	name := ctx.Param("name")
	if name == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": invalidUserName.Error(),
		})
		return
	}
	err := userService.Delete(name)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, name)
}

// BatchUser
// @Tags User
// @Summary User
// @Description Batch users
// @Accept  json
// @Produce json
// @Param request body serializer.BatchUserRequest true "Batch"
// @Success 200 {object} serializer.BatchUserResponse
// @Router /users/batch/ [post]
func Batch(ctx *gin.Context) {
	var req serializer.BatchUserRequest
	err := ctx.ShouldBind(req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": err.Error(),
		})
		return
	}
	models := make([]userModel.User, 0)
	for _, item := range req.Items {
		models = append(models, serializer.ToModel(item))
	}
	models, err = userService.Batch(req.Operation, models)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		return
	}
	if len(models) == 0 {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": deleteUserFailed.Error(),
		})
		return
	}
	var resp serializer.BatchUserResponse
	for _, model := range models {
		resp.Items = append(resp.Items, serializer.FromModel(model))
	}
	ctx.JSON(http.StatusOK, resp)
}
