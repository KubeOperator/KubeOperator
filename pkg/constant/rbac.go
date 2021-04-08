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
			"/api/v1/license",
		},
		Method: []string{"POST", "DELETE", "PUT", "PATCH"},
		Permission: &grbac.Permission{
			AuthorizedRoles: []string{RoleAdmin},
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
			"/api/v1/projects",
			"/api/v1/projects/{**}",
			"/api/v1/projects/{**}/{**}",
			"/api/v1/projects/{**}/{**}/{**}",
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
			"/api/v1/hosts/{**}",
			"/api/v1/plans/{**}",
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
			"/api/v1/hosts",
			"/api/v1/hosts/{sync,upload}",
			"/api/v1/plans",
			"/api/v1/vmconfigs",
			"/api/v1/manifests/{**}",
		},
		Method: []string{"POST", "PATCH"},
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
}

var SystemRules = loader.AdvancedRules{
	{
		Host: []string{"*"},
		Path: []string{
			"/api/v1/hosts",
			"/api/v1/clusters",
			"/api/v1/clusters/{**}",
			"/api/v1/clusters/{**}/{**}",
			"/api/v1/clusters/{**}/{**}/{**}",
			"/api/v1/clusters/{**}/{**}/{**}/{**}",
			"/api/v1/message/{**}",
			"/api/v1/message/{**}/{**}",
			"/api/v1/events/npd/{**}/{**}",
			"/api/v1/users/change/password",
			"/api/v1/logs",
		},
		Method: []string{"GET", "POST", "DELETE", "PUT", "PATCH"},
		Permission: &grbac.Permission{
			AuthorizedRoles: []string{SystemRoleAdmin, SystemRoleUser},
			AllowAnyone:     false,
		},
	},
	{
		Host: []string{"*"},
		Path: []string{
			"/api/v1/users",
			"/api/v1/users/{**}",
			"/api/v1/users/batch",
			"/api/v1/regions",
			"/api/v1/regions/{**}",
			"/api/v1/regions/{**}/{**}",
			"/api/v1/zones",
			"/api/v1/zones/{**}",
			"/api/v1/zones/{**}/{**}",
			"/api/v1/vmconfigs/{**}",
			"/api/v1/hosts/{**}",
			"/api/v1/hosts/{**}/{**}",
			"/api/v1/multicluster/repositories",
			"/api/v1/multicluster/repositories/{**}",
			"/api/v1/multicluster/repositories/{**}/{**}",
			"/api/v1/multicluster/repositories/{**}/{**}/{**}",
			"/api/v1/multicluster/repositories/{**}/{**}/{**}/{**}",
			"/api/v1/ippools",
			"/api/v1/ippools/{**}",
			"/api/v1/ippools/{**}/{**}",
			"/api/v1/ippools/{**}/{**}/{**}",
		},
		Method: []string{"GET", "POST", "DELETE", "PUT", "PATCH"},
		Permission: &grbac.Permission{
			AuthorizedRoles: []string{SystemRoleAdmin},
			AllowAnyone:     false,
		},
	},
	{
		Host: []string{"*"},
		Path: []string{
			"/api/v1/license",
			"/api/v1/settings/{**}",
			"/api/v1/settings/{**}/{**}",
			"/api/v1/settings",
			"/api/v1/projects",
			"/api/v1/projects/{**}",
			"/api/v1/project/{**}",
			"/api/v1/manifests",
			"/api/v1/manifests/{**}",
			"/api/v1/plans",
			"/api/v1/plans/{**}",
			"/api/v1/plans/{**}/{**}",
			"/api/v1/backupaccounts",
			"/api/v1/vmconfigs",
			"/api/v1/credentials",
			"/api/v1/message/setting/{**}",
		},
		Method: []string{"GET"},
		Permission: &grbac.Permission{
			AuthorizedRoles: []string{SystemRoleAdmin, SystemRoleUser},
			AllowAnyone:     false,
		},
	},
	{
		Host: []string{"*"},
		Path: []string{"/api/v1/license",
			"/api/v1/settings/{**}",
			"/api/v1/settings",
			"/api/v1/settings/{**}/{**}",
			"/api/v1/projects",
			"/api/v1/projects/{**}",
			"/api/v1/plans",
			"/api/v1/plans/{**}",
			"/api/v1/plans/{**}/{**}",
			"/api/v1/backupaccounts",
			"/api/v1/backupaccounts/{**}",
			"/api/v1/vmconfigs",
			"/api/v1/manifests/{**}",
			"/api/v1/credentials/{**}",
			"/api/v1/credentials",
			"/api/v1/message/setting/{**}",
			"/api/v1/message/setting/{**}/{**}",
		},
		Method: []string{"POST", "PUT", "PATCH", "DELETE"},
		Permission: &grbac.Permission{
			AuthorizedRoles: []string{SystemRoleAdmin},
			AllowAnyone:     false,
		},
	},
	{
		Host: []string{"*"},
		Path: []string{
			"/api/v1/project/{**}",
			"/api/v1/project/{**}/{**}",
			"/api/v1/project/{**}/{**}/{**}",
			"/api/v1/project/{**}/cluster/{**}/{**}",
			"/api/v1/project/{**}/cluster/{**}/{**}/{**}",
		},
		Method: []string{"GET", "POST", "DELETE", "PUT", "PATCH"},
		Permission: &grbac.Permission{
			AllowAnyone: true,
		},
	},
}

var ProjectRules = loader.AdvancedRules{
	{
		Host: []string{"*"},
		Path: []string{
			"/api/v1/project/{**}/{**}",
			"/api/v1/project/{**}/{**}/{**}",
			"/api/v1/project/{**}/cluster/{**}/{**}",
			"/api/v1/project/{**}/cluster/{**}/{**}/{**}",
		},
		Method: []string{"POST", "DELETE", "PUT", "PATCH"},
		Permission: &grbac.Permission{
			AuthorizedRoles: []string{ProjectRoleProjectManager, SystemRoleAdmin},
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
		},
		Method: []string{"GET"},
		Permission: &grbac.Permission{
			AuthorizedRoles: []string{ProjectRoleProjectManager, SystemRoleAdmin},
			AllowAnyone:     false,
		},
	},
	{
		Host: []string{"*"},
		Path: []string{
			"/api/v1/project/{**}/{**}",
			"/api/v1/project/{**}/{**}/{**}",
			"/api/v1/project/{**}/cluster/{**}/{**}",
			"/api/v1/project/{**}/cluster/{**}/{**}/{**}",
		},
		Method: []string{"GET"},
		Permission: &grbac.Permission{
			AuthorizedRoles: []string{ProjectRoleProjectManager, SystemRoleAdmin, ProjectRoleClusterManager},
			AllowAnyone:     false,
		},
	},
}
