package constant

import "github.com/kataras/iris/v12/sessions"

const (
	SessionUserKey         = "user"
	cookieNameForSessionID = "ksession"

	AuthMethodSession = "session"
	AuthMethodJWT     = "jwt"
)

var Sess = sessions.New(sessions.Config{Cookie: cookieNameForSessionID})
