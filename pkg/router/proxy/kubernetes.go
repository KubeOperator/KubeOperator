package proxy

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/KubeOperator/KubeOperator/pkg/controller/kolog"
	"github.com/KubeOperator/KubeOperator/pkg/repository"
	clusterUtil "github.com/KubeOperator/KubeOperator/pkg/util/cluster"
	"k8s.io/client-go/rest"

	"github.com/KubeOperator/KubeOperator/pkg/constant"
	"github.com/KubeOperator/KubeOperator/pkg/dto"
	"github.com/dgrijalva/jwt-go"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
)

func KubernetesClientProxy(ctx context.Context) {
	clusterName := ctx.Params().Get("cluster_name")
	proxyPath := ctx.Params().Get("p")

	cluteterRepo := repository.NewClusterRepository()
	cluster, err := cluteterRepo.GetWithPreload(clusterName, []string{"SpecConf", "Secret", "Nodes", "Nodes.Host", "Nodes.Host.Credential"})
	if err != nil {
		_, _ = ctx.JSON(iris.StatusInternalServerError)
		return
	}

	availableHost, err := clusterUtil.LoadAvailableHost(&cluster)
	if err != nil {
		_, _ = ctx.JSON(iris.StatusInternalServerError)
		return
	}
	u, err := url.Parse(fmt.Sprintf("https://%s", availableHost))
	if err != nil {
		_, _ = ctx.JSON(iris.StatusInternalServerError)
		return
	}
	conf, err := clusterUtil.LoadConnConf(&cluster, availableHost)
	if err != nil {
		_, _ = ctx.JSON(iris.StatusInternalServerError)
		return
	}
	tls2, err := rest.TransportFor(conf)
	if err != nil {
		_, _ = ctx.JSON(iris.StatusInternalServerError)
		return
	}
	httpClient := http.Client{Transport: tls2}
	proxy := httputil.NewSingleHostReverseProxy(u)
	ctx.Request().URL.Path = proxyPath
	proxy.Transport = httpClient.Transport
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
		buf, _ := json.Marshal(bodyStruct)
		ctx.Request().Body = ioutil.NopCloser(bytes.NewBuffer(buf))
		valueMap, _ = bodyStruct.(map[string]interface{})
	}

	proxyPath := ctx.Params().Get("p")
	tmpPath = proxyPath[(strings.Index(proxyPath, "/v1/") + 4):]
	if strings.Contains(tmpPath, "/") {
		itemvalue := strings.Split(tmpPath, "/")
		if len(itemvalue) < 4 {
			askModule = itemvalue[0]
			if len(itemvalue) == 2 {
				askParam = itemvalue[1]
			}
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
