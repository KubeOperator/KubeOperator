package dto

type LoginCredential struct {
	Username   string `json:"username"`
	Password   string `json:"password"`
	Language   string `json:"language"`
	CaptchaId  string `json:"captchaId"`
	Code       string `json:"code"`
	AuthMethod string `json:"authMethod"`
}

type SessionUser struct {
	UserId         string   `json:"userId"`
	Name           string   `json:"name"`
	Email          string   `json:"email"`
	Language       string   `json:"language"`
	IsActive       bool     `json:"isActive"`
	IsAdmin        bool     `json:"isAdmin"`
	Roles          []string `json:"roles"`
	CurrentProject string   `json:"currentProject"`
}

type Profile struct {
	User  SessionUser `json:"user"`
	Token string      `json:"token,omitempty"`
}

type Captcha struct {
	Image     string `json:"image"`
	CaptchaId string `json:"captchaId"`
}

type SessionStatus struct {
	IsLogin bool `json:"isLogin"`
}
