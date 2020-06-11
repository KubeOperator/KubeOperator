package v1

import (
	"github.com/KubeOperator/KubeOperator/pkg/router/v1/cluster"
	"github.com/KubeOperator/KubeOperator/pkg/router/v1/cluster/tools"
	"github.com/KubeOperator/KubeOperator/pkg/router/v1/credential"
	"github.com/KubeOperator/KubeOperator/pkg/router/v1/host"
	"github.com/KubeOperator/KubeOperator/pkg/router/v1/user"
	"github.com/gin-gonic/gin"
)

func V1(root *gin.RouterGroup) *gin.RouterGroup {
	v1Api := root.Group("v1")
	{
		v1HostApi := v1Api.Group("/hosts")
		{
			v1HostApi.GET("/", host.List)
			v1HostApi.POST("/", host.Create)
			v1HostApi.GET("/:name/", host.Get)
			v1HostApi.PATCH("/:name/", host.Update)
			v1HostApi.DELETE("/:name/", host.Delete)
			v1HostApi.POST("/batch/", host.Batch)
			v1HostApi.POST("/sync/:name/", host.Sync)
		}
		v1ClusterApi := v1Api.Group("/clusters")
		{
			v1ClusterApi.GET("/", cluster.List)
			v1ClusterApi.POST("/", cluster.Create)
			v1ClusterApi.GET("/:name/", cluster.Get)
			v1ClusterApi.DELETE("/:name/", cluster.Delete)
			v1ClusterApi.POST("/batch/", cluster.Batch)
			v1ClusterApi.GET("/:name/status/", cluster.Status)
			v1ClusterApi.POST("/init/:name/", cluster.Init)
		}
		v1ToolsApi := v1Api.Group("/tools/:cluster_name/")
		{
			//query ingress tools status
			v1ToolsApi.GET("/ingress", tools.GetIngress)
			// install ingress tools
			v1ToolsApi.POST("/ingress", tools.InstallIngress)
		}
		v1CredentialApi := v1Api.Group("/credentials")
		{
			v1CredentialApi.GET("/", credential.List)
			v1CredentialApi.POST("/", credential.Create)
			v1CredentialApi.GET("/:name/", credential.Get)
			v1CredentialApi.PATCH("/:name/", credential.Update)
			v1CredentialApi.DELETE("/:name/", credential.Delete)
			v1CredentialApi.POST("/batch/", credential.Batch)
		}
		v1UserApi := v1Api.Group("/users")
		{
			v1UserApi.GET("/", user.List)
			v1UserApi.POST("/", user.Create)
			v1UserApi.GET("/:name/", user.Get)
			v1UserApi.PATCH("/:name/", user.Update)
			v1UserApi.DELETE("/:name/", user.Delete)
			v1UserApi.POST("/batch/", user.Batch)
			v1UserApi.POST("/changePassword/:name/", user.ChangeUserPassword)
		}
	}
	return v1Api
}
