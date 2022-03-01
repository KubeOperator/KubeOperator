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
	"github.com/KubeOperator/KubeOperator/pkg/logger"
	"github.com/KubeOperator/KubeOperator/pkg/session"

	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/dto"
	kubeUtil "github.com/KubeOperator/KubeOperator/pkg/util/kubernetes"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
)

var log = logger.Default

func KubernetesClientProxy(ctx context.Context) {
	clusterName := ctx.Params().Get("cluster_name")
	proxyPath := ctx.Params().Get("p")
	endpoints, err := clusterService.GetApiServerEndpoints(clusterName)
	if err != nil {
		_, _ = ctx.JSON(iris.StatusInternalServerError)
		return
	}
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
		bodyStruct interface{}
		tmpPath    string
		askModule  string
		askParam   string
		valueMap   map[string]interface{}
	)

	operator := getOperator(ctx)
	_ = ctx.ReadJSON(&bodyStruct)
	if bodyStruct != nil {
		buf, err := json.Marshal(bodyStruct)
		if err != nil {
			log.Errorf("json marshal failed, %v", bodyStruct)
		}
		ctx.Request().Body = ioutil.NopCloser(bytes.NewBuffer(buf))
		valueMapItem, ok := bodyStruct.(map[string]interface{})
		if !ok {
			log.Errorf("type aassertion failed")
		}
		valueMap = valueMapItem
	}

	proxyPath := ctx.Params().Get("p")
	tmpPath = proxyPath[(strings.Index(proxyPath, "/v1/") + 4):]
	if strings.Contains(tmpPath, "/") {
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
		goSaveLogs(askParam, clusterName, operator, constant.DELETE_CLUSTER_STORAGE_CLASS, constant.CREATE_CLUSTER_STORAGE_CLASS, valueMap)
	case "namespaces":
		goSaveLogs(askParam, clusterName, operator, constant.DELETE_CLUSTER_NAMESPACE, constant.CREATE_CLUSTER_NAMESPACE, valueMap)
	case "persistentvolumes":
		goSaveLogs(askParam, clusterName, operator, constant.DELETE_CLUSTER_PVC, constant.CREATE_CLUSTER_PVC, valueMap)
	}
}

func goSaveLogs(askParam, clusterName, operator, deleteConstant, createConstant string, valueMap map[string]interface{}) {
	if len(askParam) != 0 {
		logStr := clusterName + "-" + askParam
		go kolog.Save(operator, deleteConstant, logStr)
	} else {
		if _, ok := valueMap["metadata"]; ok {
			metadata, isMap := valueMap["metadata"].(map[string]interface{})
			if isMap {
				if _, hasValue := metadata["name"]; hasValue {
					if _, isString := metadata["name"].(string); isString {
						logStr := clusterName + "-" + metadata["name"].(string)
						go kolog.Save(operator, createConstant, logStr)
					}
				}
			}
		}
	}
}

func getOperator(ctx context.Context) string {
	var sessionID = session.GloablSessionMgr.CheckCookieValid(ctx.ResponseWriter(), ctx.Request())
	if sessionID == "" {
		return ""
	}

	u, ok := session.GloablSessionMgr.GetSessionVal(sessionID, constant.SessionUserKey)
	if !ok {
		return ""
	}

	user, ok := u.(*dto.Profile)
	if !ok {
		return ""
	}
	return user.User.Name
}
