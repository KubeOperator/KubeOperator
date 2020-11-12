package constant

import (
	"github.com/kataras/iris/v12/sessions"
	"time"
)

const (
	SessionUserKey         = "user"
	cookieNameForSessionID = "ksession"

	AuthMethodSession = "session"
	AuthMethodJWT     = "jwt"
)

var Sess = sessions.New(sessions.Config{Cookie: cookieNameForSessionID, Expires: 12 * time.Hour})
