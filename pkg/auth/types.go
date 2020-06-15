package auth

type Credential struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type SessionUser struct {
	UserId   string `json:"userId"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Language string `json:"language"`
	IsActive bool   `json:"isActive"`
}

type JwtResponse struct {
	User  SessionUser `json:"user"`
	Token string      `json:"token"`
}
