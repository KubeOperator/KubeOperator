package constant

import (
	"github.com/kataras/iris/v12/sessions"
	"time"
)

const (
	SessionUserKey         = "user"
	CookieNameForSessionID = "ksession"

	AuthMethodSession = "session"
	AuthMethodJWT     = "jwt"
)

var Sess = sessions.New(sessions.Config{Cookie: CookieNameForSessionID, Expires: 12 * time.Hour})
