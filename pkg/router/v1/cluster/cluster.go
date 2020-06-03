package cluster

import (
	"errors"
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	clusterModel "github.com/KubeOperator/KubeOperator/pkg/model/cluster"
	hostModel "github.com/KubeOperator/KubeOperator/pkg/model/host"
	"github.com/KubeOperator/KubeOperator/pkg/router/v1/cluster/serializer"
	clusterService "github.com/KubeOperator/KubeOperator/pkg/service/cluster"
	"github.com/gin-gonic/gin"
	"net/http"
)

var (
	invalidClusterNameError = errors.New("invalid cluster name")
)

// ListCluster
// @Tags Cluster
// @Summary Cluster
// @Description List clusters
// @Accept  json
// @Produce json
// @Param pageNum query string false "page num"
// @Param pageSize query string false "page size"
// @Success 200 {object} serializer.ListClusterResponse
// @Router /clusters/ [get]
func List(ctx *gin.Context) {
	page := ctx.GetBool("page")
	var models []clusterModel.Cluster
	total := 0
	if page {
		pageNum := ctx.GetInt(constant.PageNumQueryKey)
		pageSize := ctx.GetInt(constant.PageSizeQueryKey)
		m, t, err := clusterService.Page(pageNum, pageSize)
		models = m
		total = t
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"msg": err.Error(),
			})
			return
		}
	} else {
		ms, err := clusterService.List()
		models = ms
		total = len(ms)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"msg": err.Error(),
			})
			return
		}
	}
	var resp = serializer.ListClusterResponse{
		Items: []serializer.Cluster{},
		Total: total,
	}
	for _, model := range models {
		resp.Items = append(resp.Items, serializer.FromModel(model))
	}
	ctx.JSON(http.StatusOK, resp)
}

// GetCluster
// @Tags Cluster
// @Summary Cluster
// @Description Get Cluster
// @Accept  json
// @Produce json
// @Param cluster_name path string true "cluster name"
// @Success 200 {object} serializer.GetClusterResponse
// @Router /clusters/{cluster_name} [get]
func Get(ctx *gin.Context) {
	name := ctx.Param("name")
	if name == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": invalidClusterNameError,
		})
		return
	}
	model, err := clusterService.Get(name)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, serializer.GetClusterResponse{Item: serializer.FromModel(model)})
}

// CreateCluster
// @Tags Cluster
// @Summary Cluster
// @Description Create a Cluster
// @Accept  json
// @Produce json
// @Param request body serializer.CreateClusterRequest true "cluster"
// @Success 201 {object} serializer.Cluster
// @Router /clusters/ [post]
func Create(ctx *gin.Context) {
	var req serializer.CreateClusterRequest
	err := ctx.ShouldBind(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": err.Error(),
		})
		return
	}
	model := clusterModel.Cluster{
		Name: req.Name,
		Spec: clusterModel.Spec{
			Version:     req.Version,
			RuntimeType: req.RuntimeType,
			NetworkType: req.NetworkType,
		},
		Status: clusterModel.Status{},
		Nodes:  []clusterModel.Node{},
	}
	for _, n := range req.Nodes {
		model.Nodes = append(model.Nodes, clusterModel.Node{
			Role: n.Role,
			Host: hostModel.Host{Name: n.HostName},
		})
	}
	err = clusterService.Save(&model)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusCreated, serializer.FromModel(model))
}

// DeleteCluster
// @Tags Cluster
// @Summary Cluster
// @Description Delete a Cluster
// @Accept  json
// @Produce json
// @Param cluster_name path string true "cluster name"
// @Success 200 {string} string
// @Router /clusters/{cluster_name} [delete]
func Delete(ctx *gin.Context) {
	name := ctx.Param("name")
	if name == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": invalidClusterNameError.Error(),
		})
		return
	}
	err := clusterService.Delete(name)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, name)
}

// BatchCluster
// @Tags Cluster
// @Summary Cluster
// @Description Batch Clusters
// @Accept  json
// @Produce json
// @Param request body serializer.BatchClusterRequest true "Batch"
// @Success 200 {object} serializer.BatchClusterResponse
// @Router /clusters/batch/ [post]
func Batch(ctx *gin.Context) {
	var req serializer.BatchClusterRequest
	err := ctx.ShouldBind(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": err.Error(),
		})
		return
	}
	models := make([]clusterModel.Cluster, 0)
	for _, item := range req.Items {
		models = append(models, serializer.ToModel(item))
	}
	models, err = clusterService.Batch(req.Operation, models)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		return
	}
	var resp serializer.BatchClusterResponse
	for _, model := range models {
		resp.Items = append(resp.Items, serializer.FromModel(model))
	}
	ctx.JSON(http.StatusOK, resp)
}

// GetClusterStatus
// @Tags Cluster
// @Summary Cluster
// @Description Get Clusters Status
// @Accept  json
// @Produce json
// @Param cluster_name path string true "cluster name"
// @Success 200 {object} serializer.ClusterStatusResponse
// @Router /clusters/{cluster_name}/status/ [get]
func Status(ctx *gin.Context) {
	clusterName := ctx.Param("name")
	if clusterName == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": invalidClusterNameError.Error(),
		})
		return
	}
	model, err := clusterService.GetClusterStatus(clusterName)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		return
	}
	var resp serializer.ClusterStatusResponse
	resp.Status = serializer.Status{
		Phase:      model.Phase,
		Conditions: make([]serializer.Condition, 0),
	}
	for _, condition := range model.Conditions {
		resp.Status.Conditions = append(resp.Status.Conditions, serializer.Condition{
			Name:    condition.Name,
			Status:  condition.Status,
			Message: condition.Message,
		})
	}
	ctx.JSON(http.StatusOK, resp)
}
