package middleware

import (
	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"github.com/storyicon/grbac"
	"net/http"
)

func RBACMiddleware() iris.Handler {
	var rules = grbac.Rules{
		{
			ID: 0,
			Resource: &grbac.Resource{
				Host:   "*",
				Path:   "**",
				Method: "*",
			},
			Permission: &grbac.Permission{
				AuthorizedRoles: []string{"*"},
				AllowAnyone:     false,
			},
		},

		//{
		//	ID: 1,
		//	Resource: &grbac.Resource{
		//		Host:   "*",
		//		Path:   "/api/v1/license",
		//		Method: "{POST}",
		//	},
		//	Permission: &grbac.Permission{
		//		AuthorizedRoles: []string{constant.SystemRoleAdmin},
		//		AllowAnyone:     false,
		//	},
		//},
		//{
		//	ID: 2,
		//	Resource: &grbac.Resource{
		//		Host:   "*",
		//		Path:   "/api/v1/license",
		//		Method: "{GET}",
		//	},
		//	Permission: &grbac.Permission{
		//		AuthorizedRoles: []string{constant.SystemRoleUser, constant.SystemRoleAdmin},
		//		AllowAnyone:     false,
		//	},
		//},
		//{
		//	ID: 3,
		//	Resource: &grbac.Resource{
		//		Host:   "*",
		//		Path:   "/api/v1/settings/{**}",
		//		Method: "{POST}",
		//	},
		//	Permission: &grbac.Permission{
		//		AuthorizedRoles: []string{constant.SystemRoleAdmin},
		//		AllowAnyone:     false,
		//	},
		//},
		//{
		//	ID: 4,
		//	Resource: &grbac.Resource{
		//		Host:   "*",
		//		Path:   "/api/v1/settings/{**}",
		//		Method: "{GET}",
		//	},
		//	Permission: &grbac.Permission{
		//		AuthorizedRoles: []string{constant.SystemRoleUser, constant.SystemRoleAdmin},
		//		AllowAnyone:     false,
		//	},
		//},
		//{
		//	ID: 0,
		//	Resource: &grbac.Resource{
		//		Host:   "*",
		//		Path:   "**",
		//		Method: "*",
		//	},
		//	Permission: &grbac.Permission{
		//		AuthorizedRoles: []string{constant.SystemRoleUser, constant.SystemRoleAdmin},
		//		AllowAnyone:     false,
		//	},
		//},
	}
	r, err := grbac.New(grbac.WithRules(rules))
	if err != nil {
		panic(err)
	}
	return func(c context.Context) {
		roles := querySystemRoles(c)
		state, err := r.IsRequestGranted(c.Request(), roles)
		if err != nil {
			c.StatusCode(http.StatusInternalServerError)
			c.StopExecution()
			return
		}
		if !state.IsGranted() {
			c.StatusCode(http.StatusForbidden)
			c.StopExecution()
			return
		}
		c.Next()
	}
}

func querySystemRoles(ctx context.Context) []string {
	u := ctx.Values().Get("user").(dto.SessionUser)
	var roles []string
	if u.IsAdmin {
		roles = append(roles, constant.SystemRoleAdmin)
	} else {
		roles = append(roles, constant.SystemRoleUser)
	}
	ctx.Values().Set("roles", roles)
	return roles
}
