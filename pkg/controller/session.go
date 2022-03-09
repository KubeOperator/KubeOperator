package controller

import (
	"errors"

	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/controller/kolog"
	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/KubeOperator/KubeOperator/pkg/middleware"
	"github.com/KubeOperator/KubeOperator/pkg/model"
	"github.com/KubeOperator/KubeOperator/pkg/service"
	"github.com/KubeOperator/KubeOperator/pkg/session"
	"github.com/KubeOperator/KubeOperator/pkg/util/captcha"
	"github.com/go-playground/validator/v10"
	"github.com/kataras/iris/v12/context"
	"github.com/spf13/viper"
)

type SessionController struct {
	Ctx         context.Context
	UserService service.UserService
}

func NewSessionController() *SessionController {
	return &SessionController{
		UserService: service.NewUserService(),
	}
}

func (s *SessionController) Get() (*dto.Profile, error) {
	var sessionID = session.GloablSessionMgr.CheckCookieValid(s.Ctx.ResponseWriter(), s.Ctx.Request())
	if len(sessionID) == 0 {
		return nil, errors.New("session invalid !")
	}

	u, ok := session.GloablSessionMgr.GetSessionVal(sessionID, constant.SessionUserKey)
	if !ok {
		session.GloablSessionMgr.EndSessionBy(sessionID)
		return nil, errors.New("session invalid !")
	}

	user, ok := u.(*dto.Profile)
	if !ok {
		return nil, errors.New("type aassertion failed")
	}
	return user, nil
}

// Login
// @Tags auth
// @Summary Login
// @Description Login
// @Param request body dto.LoginCredential true "request"
// @Accept  json
// @Produce  json
// @Success 200 {object} dto.Profile
// @Router /auth/session/ [post]
func (s *SessionController) Post() (*dto.Profile, error) {
	aul := dto.LoginCredential{}
	if err := s.Ctx.ReadJSON(&aul); err != nil {
		return nil, err
	}
	validate := validator.New()
	if err := validate.Struct(aul); err != nil {
		return nil, err
	}

	enable := viper.GetBool("validate.enable")
	if enable {
		if err := captcha.VerifyCode(aul.CaptchaId, aul.Code); err != nil {
			return nil, err
		}
	}

	p, err := s.handleLogin(aul.Username, []byte(aul.Password), false)
	if err != nil {
		return nil, err
	}

	ip := middleware.GetClientPublicIP(s.Ctx.Request())
	if len(ip) == 0 {
		ip = middleware.GetClientIP(s.Ctx.Request())
	}
	s.Ctx.Values().Set("operator", aul.Username)
	s.Ctx.Values().Set("ipfrom", ip)
	go kolog.Save(s.Ctx, constant.LOGIN, "-")
	return p, nil
}

// Login by system
// @Tags auth
// @Summary Login by system user
// @Description Login by system user
// @Param request body dto.LoginCredentialSystem true "request"
// @Accept  json
// @Produce  json
// @Success 200 {object} dto.Profile
// @Router /auth/session/system [post]
func (s *SessionController) PostSystem() (*dto.Profile, error) {
	aul := dto.LoginCredentialSystem{}
	if err := s.Ctx.ReadJSON(&aul); err != nil {
		return nil, err
	}

	validate := validator.New()
	if err := validate.Struct(aul); err != nil {
		return nil, err
	}

	p, err := s.handleLogin(aul.Username, []byte(aul.Password), true)
	if err != nil {
		return nil, err
	}

	ip := middleware.GetClientPublicIP(s.Ctx.Request())
	if len(ip) == 0 {
		ip = middleware.GetClientIP(s.Ctx.Request())
	}
	s.Ctx.Values().Set("operator", aul.Username)
	s.Ctx.Values().Set("ipfrom", ip)
	go kolog.Save(s.Ctx, constant.LOGIN, "-")
	return p, nil
}

// Logout
// @Tags auth
// @Summary Logout
// @Description Logout
// @Accept  json
// @Produce  json
// @Router /auth/session/ [delete]
func (s *SessionController) Delete() error {
	var sessionID = session.GloablSessionMgr.CheckCookieValid(s.Ctx.ResponseWriter(), s.Ctx.Request())
	if len(sessionID) == 0 {
		session.GloablSessionMgr.EndSession(s.Ctx.ResponseWriter(), s.Ctx.Request())
		return nil
	}

	u, ok := session.GloablSessionMgr.GetSessionVal(sessionID, constant.SessionUserKey)
	if !ok {
		session.GloablSessionMgr.EndSessionBy(sessionID)
		return nil
	}

	user, ok := u.(*dto.Profile)
	if !ok {
		session.GloablSessionMgr.EndSessionBy(sessionID)
		return nil
	}
	session.GloablSessionMgr.EndSessionBy(sessionID)

	ip := middleware.GetClientPublicIP(s.Ctx.Request())
	if len(ip) == 0 {
		ip = middleware.GetClientIP(s.Ctx.Request())
	}
	s.Ctx.Values().Set("operator", user.User.Name)
	s.Ctx.Values().Set("ipfrom", ip)
	go kolog.Save(s.Ctx, constant.LOGOUT, "-")
	return nil
}

func toSessionUser(u model.User) dto.SessionUser {
	return dto.SessionUser{
		UserId:   u.ID,
		Name:     u.Name,
		Language: u.Language,
		IsActive: u.IsActive,
		IsAdmin:  u.IsAdmin,
	}
}

func (s *SessionController) handleLogin(username string, password []byte, isSystem bool) (*dto.Profile, error) {
	p := &dto.Profile{}
	u, err := s.UserService.UserAuth(username, password, isSystem)
	if err != nil && err != service.WithoutChangePwd {
		return nil, err
	}
	p.User = toSessionUser(*u)

	sId := s.Ctx.GetCookie(constant.CookieNameForSessionID)
	if sId != "" {
		s.Ctx.RemoveCookie(constant.CookieNameForSessionID)
	}

	var sessionID = session.GloablSessionMgr.StartSession(s.Ctx.ResponseWriter(), s.Ctx.Request())

	var onlineSessionIDList = session.GloablSessionMgr.GetSessionIDList()
	for _, onlineSessionID := range onlineSessionIDList {
		if userInfo, ok := session.GloablSessionMgr.GetSessionVal(onlineSessionID, constant.SessionUserKey); ok {
			if value, ok := userInfo.(*dto.Profile); ok {
				if value.User.UserId == p.User.UserId {
					session.GloablSessionMgr.EndSessionBy(onlineSessionID)
				}
			}
		}
	}
	session.GloablSessionMgr.SetSessionVal(sessionID, constant.SessionUserKey, p)

	if err == service.WithoutChangePwd {
		return nil, err
	}
	return p, nil
}
