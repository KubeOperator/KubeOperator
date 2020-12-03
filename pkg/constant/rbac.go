package constant

import (
	"github.com/storyicon/grbac"
	"github.com/storyicon/grbac/pkg/loader"
)

const (
	SystemRoleAdmin = "admin"
	SystemRoleUser  = "user"
)

const (
	ProjectRoleProjectManager = "PROJECT_MANAGER"
	ProjectRoleClusterManager = "CLUSTER_MANAGER"
)

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
			"/api/v1/vm/configs/{**}",
			"/api/v1/hosts/{**}",
			"/api/v1/hosts/{**}/{**}",
			"/api/v1/multicluster/repositories",
			"/api/v1/multicluster/repositories/{**}",
			"/api/v1/multicluster/repositories/{**}/{**}",
			"/api/v1/multicluster/repositories/{**}/{**}/{**}",
			"/api/v1/multicluster/repositories/{**}/{**}/{**}/{**}",
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
			"/api/v1/vm/configs",
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
			"/api/v1/vm/configs",
			"/api/v1/manifests/{**}",
			"/api/v1/credentials/{**}",
			"/api/v1/credentials",
			"/api/v1/message/setting/{**}",
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
			"/api/v1/project/resources",
			"/api/v1/project/members",
			"/api/v1/project/resources/{**}",
			"/api/v1/project/members/{**}",
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
			"/api/v1/project/resources",
			"/api/v1/project/members",
			"/api/v1/project/resources/{**}",
			"/api/v1/project/members/{**}",
		},
		Method: []string{"GET"},
		Permission: &grbac.Permission{
			AuthorizedRoles: []string{ProjectRoleProjectManager, SystemRoleAdmin, ProjectRoleClusterManager},
			AllowAnyone:     false,
		},
	},
}
