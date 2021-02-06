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
	"strings"

	"github.com/KubeOperator/KubeOperator/pkg/controller/kolog"

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
		saveSystemLogs(ctx, clusterName)
	}

	proxy.ServeHTTP(ctx.ResponseWriter(), ctx.Request())
}

func saveSystemLogs(ctx context.Context, clusterName string) {
	var (
		logStr     string
		bodyStruct interface{}
		tmpPath    string
		askModule  string
		askParam   string
	)

	operator := getOperator(ctx)
	if err := ctx.ReadJSON(&bodyStruct); err != nil {
		fmt.Println(err)
	}

	buf, _ := json.Marshal(bodyStruct)
	ctx.Request().Body = ioutil.NopCloser(bytes.NewBuffer(buf))
	valueMap, _ := bodyStruct.(map[string]interface{})

	proxyPath := ctx.Params().Get("p")
	tmpPath = proxyPath[(strings.Index(proxyPath, "/v1/") + 4):]
	if strings.Index(tmpPath, "/") != -1 {
		itemvalue := strings.Split(tmpPath, "/")
		askModule = itemvalue[0]
		if len(itemvalue) == 2 {
			askParam = itemvalue[1]
		}
	} else {
		askModule = tmpPath
	}

	switch askModule {
	case "storageclasses":
		if len(askParam) != 0 {
			logStr = clusterName + "-" + askParam
			go kolog.Save(operator, constant.DELETE_CLUSTER_STORAGE_CLASS, logStr)
		} else {
			metadata, isMap := valueMap["metadata"].(map[string]interface{})
			if isMap {
				_, hasValue := metadata["name"]
				if hasValue {
					_, isString := metadata["name"].(string)
					if isString {
						logStr = clusterName + "-" + metadata["name"].(string)
						go kolog.Save(operator, constant.CREATE_CLUSTER_STORAGE_CLASS, logStr)
					}
				}
			}
		}
	case "namespaces":
		if len(askParam) != 0 {
			logStr = clusterName + "-" + askParam
			go kolog.Save(operator, constant.DELETE_CLUSTER_NAMESPACE, logStr)
		} else {
			metadata, isMap := valueMap["metadata"].(map[string]interface{})
			if isMap {
				_, hasValue := metadata["name"]
				if hasValue {
					_, isString := metadata["name"].(string)
					if isString {
						logStr = clusterName + "-" + metadata["name"].(string)
						go kolog.Save(operator, constant.CREATE_CLUSTER_NAMESPACE, logStr)
					}
				}
			}
		}
	case "persistentvolumes":
		if len(askParam) != 0 {
			logStr = clusterName + "-" + askParam
			go kolog.Save(operator, constant.DELETE_CLUSTER_PVC, logStr)
		} else {
			metadata, isMap := valueMap["metadata"].(map[string]interface{})
			if isMap {
				_, hasValue := metadata["name"]
				if hasValue {
					_, isString := metadata["name"].(string)
					if isString {
						logStr = clusterName + "-" + metadata["name"].(string)
						go kolog.Save(operator, constant.CREATE_CLUSTER_PVC, logStr)
					}
				}
			}
		}
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
