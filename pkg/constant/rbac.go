package constant

import (
	"github.com/storyicon/grbac"
	"github.com/storyicon/grbac/pkg/loader"
)

const (
	SystemRoleAdmin = "admin"
	SystemRoleUser  = "user"
)

var Roles = loader.AdvancedRules{
	{
		Host: []string{"*"},
		Path: []string{
			"/api/v1/hosts",
			"/api/v1/clusters",
			"/api/v1/clusters/{**}",
			"/api/v1/clusters/{**}/{**}",
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
			"/api/v1/users/batch",
			"/api/v1/regions",
			"/api/v1/zones",
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
			"/api/v1/projects",
			"/api/v1/projects/{**}",
			"/api/v1/manifests",
			"/api/v1/plans",
			"/api/v1/backupaccounts",
			"/api/v1/vm/configs",
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
			"/api/v1/projects",
			"/api/v1/plans",
			"/api/v1/backupaccounts",
			"/api/v1/vm/configs",
		},
		Method: []string{"POST", "PUT", "DELETE"},
		Permission: &grbac.Permission{
			AuthorizedRoles: []string{SystemRoleAdmin},
			AllowAnyone:     false,
		},
	},
}
