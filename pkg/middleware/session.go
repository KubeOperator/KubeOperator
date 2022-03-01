package middleware

import (
	"errors"
	"net/http"

	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/db"
	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/KubeOperator/KubeOperator/pkg/model"
	"github.com/KubeOperator/KubeOperator/pkg/session"
	"github.com/kataras/iris/v12/context"
)

func SessionMiddleware(ctx context.Context) {
	var sessionID = session.GloablSessionMgr.CheckCookieValid(ctx.ResponseWriter(), ctx.Request())
	if sessionID == "" {
		errorHandler(ctx, errors.New("session invalid !"))
		return
	}

	u, ok := session.GloablSessionMgr.GetSessionVal(sessionID, constant.SessionUserKey)
	if !ok {
		errorHandler(ctx, errors.New("session invalid !"))
		return
	}

	user, ok := u.(*dto.Profile)
	if ok {
		roles, err := getUserRole(&user.User)
		if err != nil {
			log.Errorf("get user %s role failed failed, %v", user.User.Name, err)
		}
		user.User.Roles = roles
		ctx.Values().Set("user", user.User)
		ctx.Values().Set("operator", user.User.Name)
		ctx.Next()
		return
	}
	ctx.Next()
}

func errorHandler(ctx context.Context, err error) {
	if err == nil {
		return
	}
	ctx.StopExecution()
	response := &dto.Response{
		Msg: err.Error(),
	}
	ctx.StatusCode(http.StatusUnauthorized)
	_, _ = ctx.JSON(response)
}

func getUserRole(user *dto.SessionUser) ([]string, error) {
	roles := []string{}
	if user.IsAdmin {
		roles = append(roles, constant.SystemRoleAdmin)
		return roles, nil
	}
	var projectMember []model.ProjectMember
	if err := db.DB.Model(&model.ProjectMember{}).Where("user_id = ?", user.UserId).Find(&projectMember).Error; err != nil {
		return roles, nil
	}
	isProjectManage := false
	isClusterManage := false
	for _, memeber := range projectMember {
		if memeber.Role == constant.ProjectRoleProjectManager && !isProjectManage {
			isProjectManage = true
			continue
		}
		if memeber.Role == constant.ProjectRoleClusterManager && !isClusterManage {
			isClusterManage = true
			continue
		}
	}
	if isProjectManage {
		roles = append(roles, constant.ProjectRoleProjectManager)
	}
	if isClusterManage {
		roles = append(roles, constant.ProjectRoleClusterManager)
	}
	if !isClusterManage && !isProjectManage {
		roles = append(roles, constant.SystemRoleUser)
	}
	return roles, nil
}
