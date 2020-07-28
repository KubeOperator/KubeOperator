package middleware

import (
	"encoding/json"
	"github.com/KubeOperator/KubeOperator/pkg/auth"
	"github.com/KubeOperator/KubeOperator/pkg/db"
	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/KubeOperator/KubeOperator/pkg/model"
	"github.com/KubeOperator/KubeOperator/pkg/permission"
	"github.com/KubeOperator/KubeOperator/pkg/service"
	"github.com/dgrijalva/jwt-go"
	jwtmiddleware "github.com/iris-contrib/middleware/jwt"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"github.com/spf13/viper"
	"time"
)

var (
	secretKey       []byte
	exp             int
	UserIsNotActive = "USER_IS_NOT_ACTIVE"
)

func JWTMiddleware() *jwtmiddleware.Middleware {
	secretKey = []byte(viper.GetString("jwt.secret"))
	exp = viper.GetInt("jwt.exp")
	return jwtmiddleware.New(jwtmiddleware.Config{
		Extractor: jwtmiddleware.FromAuthHeader,
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			//自己加密的秘钥或者说盐值
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

func LoginHandler(ctx context.Context) {
	aul := new(auth.Credential)
	if err := ctx.ReadJSON(&aul); err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		_, _ = ctx.JSON(dto.Response{Msg: err.Error()})
		return
	}

	data, err := CheckLogin(aul.Username, aul.Password)
	if err != nil {
		ctx.StatusCode(iris.StatusUnauthorized)
		_, _ = ctx.JSON(dto.Response{Msg: ctx.Tr(err.Error())})
		return
	}
	ctx.StatusCode(iris.StatusOK)
	_, _ = ctx.JSON(data)
	return
}

func CheckLogin(username string, password string) (*auth.JwtResponse, error) {
	user, err := service.UserAuth(username, password)
	if err != nil {
		return nil, err
	}
	token, err := CreateToken(user)
	if err != nil {
		return nil, err
	}
	resp := new(auth.JwtResponse)
	resp.Token = token
	resp.User = *user
	return resp, err
}

func CreateToken(user *auth.SessionUser) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"name":     user.Name,
		"email":    user.Email,
		"userId":   user.UserId,
		"isActive": user.IsActive,
		"language": user.Language,
		"isAdmin":  user.IsAdmin,
		"iat":      time.Now().Unix(),
		"exp":      time.Now().Add(time.Minute * time.Duration(exp)).Unix(),
	})
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}
	return tokenString, err
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
			userPermission.UserPermissionRoles = userPermissionRoles
			userPermissions = append(userPermissions, userPermission)
		}
		resp.RoleMenus = userMenus
		resp.Permissions = userPermissions
	}
	ctx.JSON(resp)
}
