package constant

import (
	"github.com/storyicon/grbac"
	"github.com/storyicon/grbac/pkg/loader"
)

const (
	RoleAdmin          = "ADMIN"
	RoleProjectManager = "PROJECT_MANAGER"
	RoleClusterManager = "CLUSTER_MANAGER"
)

const (
	SystemRoleAdmin = "admin"
	SystemRoleUser  = "user"
)

const (
	ProjectRoleProjectManager = "PROJECT_MANAGER"
	ProjectRoleClusterManager = "CLUSTER_MANAGER"
)

var Roles = loader.AdvancedRules{
	{
		Host: []string{"*"},
		Path: []string{
			"/api/v1/license",
			"/api/v1/projects",
		},
		Method: []string{"GET"},
		Permission: &grbac.Permission{
			AuthorizedRoles: []string{RoleAdmin, RoleProjectManager, RoleClusterManager},
			AllowAnyone:     false,
		},
	},
	{
		Host: []string{"*"},
		Path: []string{
			"/api/v1/users/{**}",
		},
		Method: []string{"PATCH"},
		Permission: &grbac.Permission{
			AuthorizedRoles: []string{RoleAdmin, RoleProjectManager, RoleClusterManager},
			AllowAnyone:     false,
		},
	},
	{
		Host: []string{"*"},
		Path: []string{
			"/api/v1/users/change/password",
		},
		Method: []string{"POST"},
		Permission: &grbac.Permission{
			AuthorizedRoles: []string{RoleAdmin, RoleProjectManager, RoleClusterManager},
			AllowAnyone:     false,
		},
	},
	{
		Host: []string{"*"},
		Path: []string{
			"/api/v1/clusters",
			"/api/v1/clusters/{**}",
			"/api/v1/clusters/{**}/{**}",
			"/api/v1/clusters/{**}/{**}/{**}",
			"/api/v1/clusters/{**}/{**}/{**}/{**}",
			"/api/v1/users/change/password",
			"/api/v1/logs",
			"/api/v1/message/{**}",
			"/api/v1/message/{**}/check/{**}",
			"/api/v1/message/{**}/{**}",
		},
		Method: []string{"GET", "POST", "DELETE", "PUT", "PATCH"},
		Permission: &grbac.Permission{
			AuthorizedRoles: []string{RoleAdmin, RoleProjectManager, RoleClusterManager},
			AllowAnyone:     false,
		},
	},
	{
		Host: []string{"*"},
		Path: []string{
			"/api/v1/projects/{**}",
			"/api/v1/projects/{**}/{resources,members}/{**}",
			"/api/v1/projects/{**}/clusters/{**}/{**}",
			"/api/v1/projects/{**}/clusters/{**}/{**}/{**}",
			"/api/v1/multicluster/repositories",
			"/api/v1/multicluster/repositories/{**}",
			"/api/v1/multicluster/repositories/{**}/{**}",
			"/api/v1/multicluster/repositories/{**}/{**}/{**}",
			"/api/v1/multicluster/repositories/{**}/{**}/{**}/{**}",
		},
		Method: []string{"GET", "POST", "DELETE", "PUT", "PATCH"},
		Permission: &grbac.Permission{
			AuthorizedRoles: []string{RoleAdmin, RoleProjectManager},
			AllowAnyone:     false,
		},
	},
	{
		Host: []string{"*"},
		Path: []string{
			"/api/v1/plans",
			"/api/v1/plans/{**}",
			"/api/v1/plans/{**}/{**}",
			"/api/v1/vmconfigs",
			"/api/v1/vmconfigs/{**}",
			"/api/v1/hosts",
			"/api/v1/hosts/{**}",
			"/api/v1/manifests",
			"/api/v1/manifests/{active,group}",
			"/api/v1/backupaccounts",
			"/api/v1/backupaccounts/{**}",
			"/api/v1/projects/{**}/{resources,members}",
		},
		Method: []string{"GET"},
		Permission: &grbac.Permission{
			AuthorizedRoles: []string{RoleAdmin, RoleProjectManager},
			AllowAnyone:     false,
		},
	},
	{
		Host: []string{"*"},
		Path: []string{
			"/api/v1/plans/search",
			"/api/v1/vmconfigs/search",
			"/api/v1/hosts/search",
			"/api/v1/backupaccounts/search",
		},
		Method: []string{"POST"},
		Permission: &grbac.Permission{
			AuthorizedRoles: []string{RoleAdmin, RoleProjectManager},
			AllowAnyone:     false,
		},
	},
	{
		Host: []string{"*"},
		Path: []string{
			"/api/v1/regions",
			"/api/v1/regions/{**}",
			"/api/v1/regions/{**}/{**}",
			"/api/v1/zones",
			"/api/v1/zones/{**}",
			"/api/v1/zones/{**}/{**}",
			"/api/v1/ippools",
			"/api/v1/ippools/{**}",
			"/api/v1/ippools/{**}/{**}",
			"/api/v1/ippools/{**}/{**}/{**}",
			"/api/v1/credentials",
			"/api/v1/credentials/{**}",
			"/api/v1/settings",
			"/api/v1/settings/{**}",
			"/api/v1/settings/{**}/{**}",
		},
		Method: []string{"GET", "POST", "DELETE", "PUT", "PATCH"},
		Permission: &grbac.Permission{
			AuthorizedRoles: []string{RoleAdmin},
			AllowAnyone:     false,
		},
	},
	{
		Host: []string{"*"},
		Path: []string{
			"/api/v1/users",
			"/api/v1/users/{**}",
		},
		Method: []string{"GET"},
		Permission: &grbac.Permission{
			AuthorizedRoles: []string{RoleAdmin},
			AllowAnyone:     false,
		},
	},
	{
		Host: []string{"*"},
		Path: []string{
			"/api/v1/projects",
			"/api/v1/users",
			"/api/v1/users/{**}",
			"/api/v1/license",
			"/api/v1/hosts",
			"/api/v1/hosts/{sync,upload}",
			"/api/v1/plans",
			"/api/v1/vmconfigs",
			"/api/v1/backupaccounts",
			"/api/v1/backupaccounts/buckets",
			"/api/v1/projects/{**}/{resources,members}",
		},
		Method: []string{"POST"},
		Permission: &grbac.Permission{
			AuthorizedRoles: []string{RoleAdmin},
			AllowAnyone:     false,
		},
	},
	{
		Host: []string{"*"},
		Path: []string{
			"/api/v1/users/{**}",
			"/api/v1/hosts/{**}",
			"/api/v1/plans/{**}",
			"/api/v1/vmconfigs/{**}",
			"/api/v1/backupaccounts/{**}",
			"/api/v1/projects/{**}/{resources,members}/{**}",
		},
		Method: []string{"DELETE"},
		Permission: &grbac.Permission{
			AuthorizedRoles: []string{RoleAdmin},
			AllowAnyone:     false,
		},
	},
	{
		Host: []string{"*"},
		Path: []string{
			"/api/v1/vmconfigs/{**}",
			"/api/v1/manifests/{**}",
			"/api/v1/backupaccounts/{**}",
			"/api/v1/plans/{**}",
		},
		Method: []string{"PATCH"},
		Permission: &grbac.Permission{
			AuthorizedRoles: []string{RoleAdmin},
			AllowAnyone:     false,
		},
	},
}
