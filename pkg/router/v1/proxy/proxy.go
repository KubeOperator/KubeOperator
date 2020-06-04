package proxy

import (
	"crypto/tls"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func KubernetesClientProxy(ctx *gin.Context) {
	u, err := url.Parse("https://172.16.10.179:6443")
	if err != nil {
		return
	}
	proxy := httputil.NewSingleHostReverseProxy(u)
	proxy.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	t := "eyJhbGciOiJSUzI1NiIsImtpZCI6IlhsUWRlZE5TZDF0Yl9TbGlrbElNZVdqVW5RSnRRZzNCaU8xNVg3ZXhtUVUifQ.eyJpc3MiOiJrdWJlcm5ldGVzL3NlcnZpY2VhY2NvdW50Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9uYW1lc3BhY2UiOiJrdWJlLXN5c3RlbSIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VjcmV0Lm5hbWUiOiJ0aWxsZXItdG9rZW4tcjIyZDkiLCJrdWJlcm5ldGVzLmlvL3NlcnZpY2VhY2NvdW50L3NlcnZpY2UtYWNjb3VudC5uYW1lIjoidGlsbGVyIiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9zZXJ2aWNlLWFjY291bnQudWlkIjoiNDQ3NTQ2NDEtYjZkOS00ODEwLTlhYjUtYWY2ZWViYzVmZWJlIiwic3ViIjoic3lzdGVtOnNlcnZpY2VhY2NvdW50Omt1YmUtc3lzdGVtOnRpbGxlciJ9.ltP24ykjKGsZBtiSWMA5sxVBVAGytrUrzrXTWTHnUzVPQk73pfP1dq6vjaInYGiTgVGdNVbTQI_CH4Q6nzO_dDyqsnEVyxWm74Pr-cUx0q-257J0rGVU4iL21yDI-fTezPt3PLJTt2vB7BpHiKzPhqOyIqpw2K7OK56hrgojzeUh5BbNGeLNjC0aaRhjgyx-r0z49E5tMnR2tJb6URxyOsJ0z9X87-te43sdTDlr0FUWWf6Q0yJy5vZO9FigYvFDXBj5IFOMiyYvjWjf5nwFJFsrNn__q8Xx9WkCQxI1la4CLy_2WfXTnsYsiQbLW0fmwtpImHQQs7tXfDqYzn3WWA"
	token := fmt.Sprintf("Bearer %s", t)
	ctx.Request.Header.Add("Authorization", token)
	path := ctx.Param("path")
	ctx.Request.URL = &url.URL{
		Path: path,
	}
	proxy.ServeHTTP(ctx.Writer, ctx.Request)

}
