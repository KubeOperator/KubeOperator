package proxy

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/KubeOperator/KubeOperator/pkg/controller/log_save"

	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/dto"
	kubeUtil "github.com/KubeOperator/KubeOperator/pkg/util/kubernetes"
	"github.com/dgrijalva/jwt-go"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
)

func KubernetesClientProxy(ctx context.Context) {
	clusterName := ctx.Params().Get("cluster_name")
	proxyPath := ctx.Params().Get("p")
	endpoints, err := clusterService.GetApiServerEndpoints(clusterName)

	aliveHost, err := kubeUtil.SelectAliveHost(endpoints)
	if err != nil {
		_, _ = ctx.JSON(iris.StatusInternalServerError)
		return
	}
	u, err := url.Parse(fmt.Sprintf("https://%s", aliveHost))
	if err != nil {
		_, _ = ctx.JSON(iris.StatusInternalServerError)
		return
	}
	proxy := httputil.NewSingleHostReverseProxy(u)
	proxy.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	secret, err := clusterService.GetSecrets(clusterName)
	if err != nil {
		_, _ = ctx.JSON(iris.StatusInternalServerError)
		return
	}
	token := fmt.Sprintf("%s %s", keyPrefix, secret.KubernetesToken)
	ctx.Request().Header.Add(AuthorizationHeader, token)
	ctx.Request().URL.Path = proxyPath
	proxy.ModifyResponse = func(response *http.Response) error {
		if response.StatusCode == http.StatusUnauthorized {
			response.StatusCode = http.StatusInternalServerError
		}
		return nil
	}
	if ctx.Method() != "GET" {
		saveSystemLogs(ctx)
	}

	proxy.ServeHTTP(ctx.ResponseWriter(), ctx.Request())
}

func saveSystemLogs(ctx context.Context) {
	var (
		logStr     string
		bodyStruct interface{}
	)

	operator := getOperator(ctx)
	if err := ctx.ReadJSON(&bodyStruct); err != nil {
		fmt.Println(err)
	}

	buf, _ := json.Marshal(bodyStruct)
	ctx.Request().Body = ioutil.NopCloser(bytes.NewBuffer(buf))
	valueMap, _ := bodyStruct.(map[string]interface{})

	switch ctx.Params().Get("p") {
	case "apis/storage.k8s.io/v1/storageclasses":
		metadata, _ := valueMap["metadata"].(map[string]interface{})
		logStr = valueMap["provisioner"].(string) + "-" + metadata["name"].(string)
		go log_save.LogSave(operator, constant.CREATE_CLUSTER_STORAGE_CLASS, logStr)
	case "api/v1/namespaces":
		metadata, _ := valueMap["metadata"].(map[string]interface{})
		logStr = metadata["name"].(string)
		go log_save.LogSave(operator, constant.CREATE_CLUSTER_NAMESPACE, logStr)
	}
}

func getOperator(ctx context.Context) string {
	var u dto.SessionUser
	j := ctx.Values().Get("jwt")
	if j != nil {
		j := j.(*jwt.Token)
		foobar := j.Claims.(jwt.MapClaims)
		js, _ := json.Marshal(foobar)
		_ = json.Unmarshal(js, &u)
	} else {
		session := constant.Sess.Start(ctx)
		sessionUser := session.Get(constant.SessionUserKey)
		if sessionUser == nil {
			ctx.StatusCode(http.StatusUnauthorized)
			ctx.StopExecution()
			return ""
		}
		u = sessionUser.(*dto.Profile).User
	}
	return u.Name
}
