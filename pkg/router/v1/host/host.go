package host

import (
	"errors"
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	hostModel "github.com/KubeOperator/KubeOperator/pkg/model/host"
	"github.com/KubeOperator/KubeOperator/pkg/router/v1/host/serializer"
	credentialService "github.com/KubeOperator/KubeOperator/pkg/service/credential"
	hostService "github.com/KubeOperator/KubeOperator/pkg/service/host"
	"github.com/gin-gonic/gin"
	"net/http"
)

var (
	invalidHostName = errors.New("invalid host name")
)

// ListHost
// @Tags Host
// @Summary Host
// @Description List hosts
// @Accept  json
// @Produce json
// @Param pageNum query string false "page num"
// @Param pageSize query string false "page size"
// @Success 200 {object} serializer.ListHostResponse
// @Router /hosts/ [get]
func List(ctx *gin.Context) {
	page := ctx.GetBool("page")
	var models []hostModel.Host
	total := 0
	if page {
		pageNum := ctx.GetInt(constant.PageNumQueryKey)
		pageSize := ctx.GetInt(constant.PageSizeQueryKey)
		m, t, err := hostService.Page(pageNum, pageSize)
		models = m
		total = t
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"msg": err.Error(),
			})
			return
		}
	} else {
		ms, err := hostService.List()
		models = ms
		total = len(models)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"msg": err.Error(),
			})
			return
		}
	}

	var resp = serializer.ListHostResponse{
		Items: []serializer.Host{},
		Total: total,
	}
	for _, model := range models {
		resp.Items = append(resp.Items, serializer.FromModel(model))
	}

	ctx.JSON(http.StatusOK, resp)
}

// GetHost
// @Tags Host
// @Summary Host
// @Description Get Host
// @Accept  json
// @Produce json
// @Param host_name path string true "host name"
// @Success 200 {object} serializer.GetHostResponse
// @Router /hosts/{host_name} [get]
func Get(ctx *gin.Context) {
	name := ctx.Param("name")
	if name == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": invalidHostName.Error(),
		})
		return
	}
	model, err := hostService.Get(name)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, serializer.GetHostResponse{
		Item: serializer.FromModel(*model),
	})
}

// CreateHost
// @Tags Host
// @Summary Host
// @Description Create a Host
// @Accept  json
// @Produce json
// @Param request body serializer.CreateHostRequest true "host"
// @Success 201 {object} serializer.Host
// @Router /hosts/ [post]
func Create(ctx *gin.Context) {
	var req serializer.CreateHostRequest
	err := ctx.ShouldBind(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": err.Error(),
		})
		return
	}
	credential, err := credentialService.GetById(req.CredentialID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": err.Error(),
		})
		return
	}
	model := hostModel.Host{
		Name:         req.Name,
		Ip:           req.Ip,
		Port:         req.Port,
		CredentialID: req.CredentialID,
		Credential:   credential,
	}
	err = hostService.GetHostGpu(&model)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		return
	}
	go hostService.GetHostConfig(&model)

	err = hostService.Save(&model)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusCreated, serializer.FromModel(model))
}

// UpdateHost
// @Tags Host
// @Summary Host
// @Description Update a Host
// @Accept  json
// @Produce json
// @Param request body serializer.UpdateHostRequest true "host"
// @Param host_name path string true "host_name"
// @Success 200 {object} serializer.Host
// @Router /hosts/{host_name} [patch]
func Update(ctx *gin.Context) {
	var req serializer.UpdateHostRequest
	err := ctx.ShouldBind(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": err.Error(),
		})
		return
	}
	model := serializer.ToModel(req.Item)
	err = hostService.Save(&model)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, serializer.FromModel(model))
}

// DeleteHost
// @Tags Host
// @Summary Host
// @Description Delete a Host
// @Accept  json
// @Produce json
// @Param host_name path string true "host name"
// @Success 200 {string} string
// @Router /hosts/{host_name} [delete]
func Delete(ctx *gin.Context) {
	name := ctx.Param("name")
	if name == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": invalidHostName.Error(),
		})
		return
	}
	err := hostService.Delete(name)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, name)
}

// BatchHost
// @Tags Host
// @Summary Host
// @Description Batch Host
// @Accept  json
// @Produce json
// @Param request body serializer.BatchHostRequest true "Batch"
// @Success 200 {object} serializer.BatchHostResponse
// @Router /hosts/batch/ [post]
func Batch(ctx *gin.Context) {
	var req serializer.BatchHostRequest
	err := ctx.ShouldBind(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": err.Error(),
		})
		return
	}
	models := make([]hostModel.Host, 0)
	for _, item := range req.Items {
		models = append(models, serializer.ToModel(item))
	}
	models, err = hostService.Batch(req.Operation, models)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		return
	}
	var resp serializer.BatchHostResponse
	for _, model := range models {
		resp.Items = append(resp.Items, serializer.FromModel(model))
	}
	ctx.JSON(http.StatusOK, resp)
}
