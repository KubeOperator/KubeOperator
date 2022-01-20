package dto

import "time"

type LoginCredential struct {
	Username   string `json:"username" validate:"required"`
	Password   string `json:"password" validate:"required"`
	Language   string `json:"language" validate:"required"`
	CaptchaId  string `json:"captchaId" validate:"required"`
	Code       string `json:"code" validate:"required"`
	AuthMethod string `json:"authMethod" validate:"-"`
}

type SessionUser struct {
	UserId   string `json:"userId"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Language string `json:"language"`
	IsActive bool   `json:"isActive"`
	IsAdmin  bool   `json:"isAdmin"`
}

type Profile struct {
	User    SessionUser `json:"user"`
	Timeout time.Time   `json:"timeout"`
	Token   string      `json:"token,omitempty"`
}

type Captcha struct {
	Image     string `json:"image"`
	CaptchaId string `json:"captchaId"`
}
