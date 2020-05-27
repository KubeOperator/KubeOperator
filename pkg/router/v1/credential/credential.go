package credential

import (
	"errors"
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/i18n"
	credentialModel "github.com/KubeOperator/KubeOperator/pkg/model/credential"
	"github.com/KubeOperator/KubeOperator/pkg/router/v1/credential/serializer"
	credentialService "github.com/KubeOperator/KubeOperator/pkg/service/credential"
	"github.com/gin-gonic/gin"
	"net/http"
)

var (
	invalidCredentialName = errors.New(i18n.Tr("invalid_credential_name", nil))
)

// ListCredential
// @Tags Credential
// @Summary Credential
// @Description List credentials
// @Accept  json
// @Produce json
// @Param pageNum query string false "page num"
// @Param pageSize query string false "page size"
// @Success 200 {object} serializer.ListCredentialResponse
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

	var resp = serializer.ListCredentialResponse{
		Items: []serializer.Credential{},
		Total: total,
	}
	for _, model := range models {
		resp.Items = append(resp.Items, serializer.FromModel(model))
	}

	ctx.JSON(http.StatusOK, resp)
}

// GetCredential
// @Tags Credential
// @Summary Credential
// @Description Get Credential
// @Accept  json
// @Produce json
// @Param credential_name path string true "credential name"
// @Success 200 {object} serializer.GetCredentialResponse
// @Router /credentials/{credential_name} [get]
func Get(ctx *gin.Context) {
	name := ctx.Param("name")
	if name == "test" {
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
	ctx.JSON(http.StatusOK, serializer.GetCredentialResponse{
		Item: serializer.FromModel(*model),
	})
}

// CreateCredential
// @Tags Credential
// @Summary Credential
// @Description Create a Credential
// @Accept  json
// @Produce json
// @Param request body serializer.CreateCredentialRequest true "credential"
// @Success 201 {object} serializer.Credential
// @Router /credentials/ [post]
func Create(ctx *gin.Context) {
	var req serializer.CreateCredentialRequest
	err := ctx.ShouldBind(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": err.Error(),
		})
		return
	}
	model := credentialModel.Credential{
		Name:       req.Name,
		Password:   req.Password,
		PrivateKey: req.PrivateKey,
		Type:       req.Type,
		Username:   req.Username,
	}
	err = credentialService.Save(&model)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusCreated, serializer.FromModel(model))
}

// UpdateCredential
// @Tags Credential
// @Summary Credential
// @Description Update a Credential
// @Accept  json
// @Produce json
// @Param request body serializer.UpdateCredentialRequest true "credential"
// @Param credential_name path string true "credential name"
// @Success 200 {object} serializer.Credential
// @Router /credentials/{credential_name} [patch]
func Update(ctx *gin.Context) {
	var req serializer.UpdateCredentialRequest
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
	ctx.JSON(http.StatusOK, serializer.FromModel(model))
}

// DeleteCredential
// @Tags Credential
// @Summary Credential
// @Description Delete a Credential
// @Accept  json
// @Produce json
// @Param credential_name path string true "credential name"
// @Success 200 {string} string
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
	ctx.JSON(http.StatusOK, name)
}

// BatchCredential
// @Tags Credential
// @Summary Credential
// @Description Batch Credentials
// @Accept  json
// @Produce json
// @Param request body serializer.BatchCredentialRequest true "Batch"
// @Success 200 {object} serializer.BatchCredentialResponse
// @Router /credentials/batch/ [post]
func Batch(ctx *gin.Context) {
	var req serializer.BatchCredentialRequest
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
	var resp serializer.BatchCredentialResponse
	for _, model := range models {
		resp.Items = append(resp.Items, serializer.FromModel(model))
	}
	ctx.JSON(http.StatusOK, resp)
}
