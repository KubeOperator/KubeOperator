package middleware

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/db"
	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/KubeOperator/KubeOperator/pkg/model"
	"github.com/kataras/iris/v12/context"
	"github.com/storyicon/grbac"
)

func ProjectMiddleware(ctx context.Context) {
	r, err := grbac.New(grbac.WithAdvancedRules(constant.ProjectRules))
	if err != nil {
		panic(err)
	}
	roles, err := queryProjectRoles(ctx)
	if err != nil {
		ctx.StatusCode(http.StatusBadRequest)
		ctx.StopExecution()
		return
	}
	state, err := r.IsRequestGranted(ctx.Request(), roles)
	if err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.StopExecution()
		return
	}
	if !state.IsGranted() {
		ctx.StatusCode(http.StatusForbidden)
		ctx.StopExecution()
		return
	}
	ctx.Next()
}

func queryProjectRoles(ctx context.Context) ([]string, error) {
	var roles []string
	u := ctx.Values().Get("user").(dto.SessionUser)
	projectNameUnDecode := ctx.Request().Header.Get("project")
	projectName, err := url.QueryUnescape(projectNameUnDecode)
	if err != nil {
		return nil, fmt.Errorf("decode error: %s", projectName)
	}
	var project model.Project
	notFound := db.DB.Where("name = ?", projectName).First(&project).RecordNotFound()
	if notFound {
		return nil, fmt.Errorf("project: %s not found", projectName)
	}
	ctx.Values().Set("project", projectName)
	// admin 拥有一切权限
	if u.IsAdmin {
		return []string{constant.SystemRoleAdmin}, nil
	}
	var member model.ProjectMember
	db.DB.Where("project_id = ? AND user_id = ?", project.ID, u.UserId).First(&member)
	if member.ID != "" {
		roles = append(roles, member.Role)
	}
	return roles, nil
}
