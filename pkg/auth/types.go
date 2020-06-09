package auth

type Credential struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type SessionUser struct {
	UserId   string
	Name     string
	Email    string
	Language string
	IsActive bool
}
