package middleware

import (
	"crypto/md5"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

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
			log.Errorf("get user %s role failed, %v", user.User.Name, err)
		}
		user.User.Roles = roles
		ip := GetClientPublicIP(ctx.Request())
		if len(ip) == 0 {
			ip = GetClientIP(ctx.Request())
		}
		ctx.Values().Set("user", user.User)
		ctx.Values().Set("ipfrom", ip)
		if ctx.Request().Method == "POST" {
			ctx.Values().Set("hasCsrf", getClientCsrf(ctx.Request()))
		}
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

func getClientCsrf(r *http.Request) bool {
	cs := r.Header.Get("X-CSRF-TOKEN")
	timeNow := time.Now().UTC().Format("01-02 15:04")
	timeAMinuteAgo := time.Now().UTC().Add(-1 * time.Minute).Format("01-02 15:04")

	return (md5Str("kubeoperator"+timeNow) == cs || md5Str("kubeoperator"+timeAMinuteAgo) == cs)
}

func md5Str(str string) string {
	data := []byte(str)
	has := md5.Sum(data)
	md5str := fmt.Sprintf("%x", has)
	return md5str
}

func GetClientIP(r *http.Request) string {
	ip := strings.TrimSpace(strings.Split(r.Header.Get("X-Forwarded-For"), ",")[0])
	if ip != "" {
		return ip
	}
	ip = strings.TrimSpace(r.Header.Get("X-Real-Ip"))
	if ip != "" {
		return ip
	}
	if ip, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr)); err == nil {
		return ip
	}
	return ""
}

func GetClientPublicIP(r *http.Request) string {
	var ip string
	for _, ip = range strings.Split(r.Header.Get("X-Forwarded-For"), ",") {
		if ip = strings.TrimSpace(ip); ip != "" && !hasLocalIPAddr(ip) {
			return ip
		}
	}
	if ip = strings.TrimSpace(r.Header.Get("X-Real-Ip")); ip != "" && !hasLocalIPAddr(ip) {
		return ip
	}
	if ip = remoteIP(r); !hasLocalIPAddr(ip) {
		return ip
	}
	return ""
}

func remoteIP(r *http.Request) string {
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		fmt.Println(err)
	}
	return ip
}
func hasLocalIPAddr(ip string) bool {
	return hasLocalIP(net.ParseIP(ip))
}

func hasLocalIP(ip net.IP) bool {
	if ip.IsLoopback() {
		return true
	}

	ip4 := ip.To4()
	if ip4 == nil {
		return false
	}

	return ip4[0] == 10 || // 10.0.0.0/8
		(ip4[0] == 172 && ip4[1] >= 16 && ip4[1] <= 31) || // 172.16.0.0/12
		(ip4[0] == 169 && ip4[1] == 254) || // 169.254.0.0/16
		(ip4[0] == 192 && ip4[1] == 168) // 192.168.0.0/16
}
