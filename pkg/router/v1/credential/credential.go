package credential

import (
	"errors"
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	commonModel "github.com/KubeOperator/KubeOperator/pkg/model/common"
	credentialModel "github.com/KubeOperator/KubeOperator/pkg/model/credential"
	serializer "github.com/KubeOperator/KubeOperator/pkg/router/v1/credential/serializer"
	credentialService "github.com/KubeOperator/KubeOperator/pkg/service/credential"
	"github.com/gin-gonic/gin"
	"net/http"
)

var (
	invalidCredentialName = errors.New("invalid credential name")
)

// ListCredential
// @Summary Credential
// @Description List credentials
// @Accept  json
// @Produce json
// @Param pageNum query string false "page num"
// @Param pageSize query string false "page size"
// @Success 200 {object} serializer.ListResponse
// @Router /credentials/ [get]
func List(ctx *gin.Context) {
	page := ctx.GetBool("page")
	var models []credentialModel.Credential
	total := 0
	if page {
		pageNum := ctx.GetInt(constant.PageNumQueryKey)
		pageSize := ctx.GetInt(constant.PageSizeQueryKey)
		m, t, err := credentialService.Page(pageNum, pageSize)
		models = m
		total = t
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"msg": err.Error(),
			})
			return
		}
	} else {
		ms, err := credentialService.List()
		models = ms
		total = len(ms)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"msg": err.Error(),
			})
			return
		}

	}

	var resp = serializer.ListResponse{
		Items: []serializer.Credential{},
		Total: total,
	}
	for _, model := range models {
		resp.Items = append(resp.Items, serializer.FromModel(model))
	}

	ctx.JSON(http.StatusOK, resp)
}

// GetCredential
// @Summary Credential
// @Description Get Credential
// @Accept  json
// @Produce json
// @Param credential_name path string true "credential name"
// @Success 200 {object} serializer.GetResponse
// @Router /credentials/{credential_name} [get]
func Get(ctx *gin.Context) {
	name := ctx.Param("name")
	if name == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": invalidCredentialName.Error(),
		})
		return
	}
	model, err := credentialService.Get(name)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, serializer.GetResponse{
		Item: serializer.FromModel(*model),
	})
}

// CreateCredential
// @Summary Credential
// @Description Create a Credential
// @Accept  json
// @Produce json
// @Param request body serializer.CreateRequest true "credential"
// @Success 201 {object} serializer.CreateResponse
// @Router /credentials/ [post]
func Create(ctx *gin.Context) {
	var req serializer.CreateRequest
	err := ctx.ShouldBind(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": err.Error(),
		})
		return
	}
	model := credentialModel.Credential{
		BaseModel: commonModel.BaseModel{
			Name: req.Name,
		},
	}
	err = credentialService.Save(&model)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusCreated, serializer.CreateResponse{Item: serializer.FromModel(model)})
}

// UpdateCredential
// @Summary Credential
// @Description Update a Credential
// @Accept  json
// @Produce json
// @Param request body serializer.UpdateRequest true "credential"
// @Param credential_name path string true "credential name"
// @Success 200 {object} serializer.UpdateResponse
// @Router /credentials/{credential_name} [patch]
func Update(ctx *gin.Context) {
	var req serializer.UpdateRequest
	err := ctx.ShouldBind(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": err.Error(),
		})
		return
	}
	model := serializer.ToModel(req.Item)
	err = credentialService.Save(&model)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, serializer.UpdateResponse{
		Item: serializer.FromModel(model),
	})
}

// DeleteCredential
// @Summary Cluster
// @Description Delete a Credential
// @Accept  json
// @Produce json
// @Param credential_name path string true "credential name"
// @Success 200 {object} serializer.DeleteResponse
// @Router /credentials/{credential_name} [delete]
func Delete(ctx *gin.Context) {
	name := ctx.Param("name")
	if name == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": invalidCredentialName.Error(),
		})
		return
	}
	err := credentialService.Delete(name)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, serializer.DeleteResponse{})
}

// BatchCredential
// @Summary Cluster
// @Description Batch Credentials
// @Accept  json
// @Produce json
// @Param request body serializer.BatchRequest true "Batch"
// @Success 200 {object} serializer.BatchResponse
// @Router /credentials/batch/ [post]
func Batch(ctx *gin.Context) {
	var req serializer.BatchRequest
	err := ctx.ShouldBind(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": err.Error(),
		})
		return
	}
	models := make([]credentialModel.Credential, 0)
	for _, item := range req.Items {
		models = append(models, serializer.ToModel(item))
	}
	models, err = credentialService.Batch(req.Operation, models)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		return
	}
	var resp serializer.BatchResponse
	for _, model := range models {
		resp.Items = append(resp.Items, serializer.FromModel(model))
	}
	ctx.JSON(http.StatusOK, resp)
}
