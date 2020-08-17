package middleware

import (
	"encoding/json"
	"errors"
	"github.com/KubeOperator/KubeOperator/pkg/auth"
	"github.com/KubeOperator/KubeOperator/pkg/db"
	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/KubeOperator/KubeOperator/pkg/model"
	"github.com/KubeOperator/KubeOperator/pkg/permission"
	"github.com/dgrijalva/jwt-go"
	jwtmiddleware "github.com/iris-contrib/middleware/jwt"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"github.com/spf13/viper"
)

var (
	UserIsNotRelatedProject = "USER_IS_NOT_RELATED_PROJECT"
)

func JWTMiddleware() *jwtmiddleware.Middleware {
	secretKey := []byte(viper.GetString("jwt.secret"))
	return jwtmiddleware.New(jwtmiddleware.Config{
		Extractor: jwtmiddleware.FromAuthHeader,
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return secretKey, nil
		},
		SigningMethod: jwt.SigningMethodHS256,
		ErrorHandler:  ErrorHandler,
	})
}

func ErrorHandler(ctx context.Context, err error) {
	if err == nil {
		return
	}
	ctx.StopExecution()
	response := &dto.Response{
		Msg: err.Error(),
	}
	ctx.StatusCode(iris.StatusInternalServerError)
	ctx.JSON(response)
}

func GetAuthUser(ctx context.Context) {
	user := ctx.Values().Get("jwt").(*jwt.Token)
	foobar := user.Claims.(jwt.MapClaims)
	sessionUserJson, _ := json.Marshal(foobar)
	sessionUserJsonStr := string(sessionUserJson)
	var sessionUser auth.SessionUser
	json.Unmarshal([]byte(sessionUserJsonStr), &sessionUser)
	resp := new(auth.JwtResponse)
	resp.User = sessionUser
	resp.Token = user.Raw

	if !sessionUser.IsAdmin {
		var user model.User
		err := db.DB.Model(model.User{}).Where(model.User{Name: sessionUser.Name, Email: sessionUser.Email}).First(&user).Error
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			_, _ = ctx.JSON(dto.Response{Msg: err.Error()})
			return
		}
		var projectMembers []dto.ProjectMember
		err = db.DB.Model(model.ProjectMember{}).Where(model.ProjectMember{UserID: sessionUser.UserId}).Find(&projectMembers).Error
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			_, _ = ctx.JSON(dto.Response{Msg: err.Error()})
			return
		}

		if len(projectMembers) == 0 {
			ctx.StatusCode(iris.StatusInternalServerError)
			_, _ = ctx.JSON(dto.Response{Msg: ctx.Tr(errors.New(UserIsNotRelatedProject).Error())})
			return
		}

		var userMenus []permission.UserMenu
		var menuRoles []permission.MenuRole
		_ = json.Unmarshal([]byte(permission.MenuRoles), &menuRoles)

		var userPermissions []permission.UserPermission
		var permissionRoles []permission.Permission
		_ = json.Unmarshal([]byte(permission.PermissionRoles), &permissionRoles)

		for _, pm := range projectMembers {
			var userMenu permission.UserMenu
			var menus []string
			for _, menuRole := range menuRoles {
				for _, role := range menuRole.Roles {
					if role == pm.Role {
						menus = append(menus, menuRole.Menu)
						break
					}
				}
			}
			userMenu.ProjectId = pm.ProjectID
			userMenu.Menus = menus
			userMenus = append(userMenus, userMenu)

			var userPermission permission.UserPermission
			var userPermissionRoles []permission.UserPermissionRole
			for _, up := range permissionRoles {
				var userPermissionRole permission.UserPermissionRole
				var roles []string
				for _, opAuths := range up.OperationAuth {
					for _, role := range opAuths.Roles {
						if role == pm.Role {
							roles = append(roles, opAuths.Operation)
							break
						}
					}
				}
				if len(roles) > 0 {
					userPermissionRole.ResourceType = up.ResourceType
					userPermissionRole.Roles = roles
					userPermissionRoles = append(userPermissionRoles, userPermissionRole)
				}
			}
			userPermission.ProjectId = pm.ProjectID

			var project model.Project
			err := db.DB.Model(model.Project{}).Where(model.Project{ID: pm.ProjectID}).First(&project).Error
			if err != nil {
				ctx.StatusCode(iris.StatusInternalServerError)
				_, _ = ctx.JSON(dto.Response{Msg: err.Error()})
				return
			}
			userPermission.ProjectName = project.Name
			userPermission.ProjectRole = pm.Role
			userPermission.UserPermissionRoles = userPermissionRoles
			userPermissions = append(userPermissions, userPermission)
		}
		resp.RoleMenus = userMenus
		resp.Permissions = userPermissions
	}
	ctx.JSON(resp)
}
