package kubernetes

import (
	"context"
	"fmt"
	"testing"

	"github.com/KubeOperator/KubeOperator/pkg/logger"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestNewKubernetesClient(t *testing.T) {
	c, err := NewKubernetesExtensionClient(&Config{
		Hosts: []Host{"172.16.10.228:8443", "172.16.10.225:8443", "172.16.10.226:8443"},
		Token: "eyJhbGciOiJSUzI1NiIsImtpZCI6IlphWEJqNEppc1NUdVJqU0x3eTVBMThWYjZoUmpmdFNiWVBaXzc0eno2RjQifQ.eyJpc3MiOiJrdWJlcm5ldGVzL3NlcnZpY2VhY2NvdW50Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9uYW1lc3BhY2UiOiJrdWJlLXN5c3RlbSIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VjcmV0Lm5hbWUiOiJrby1hZG1pbi10b2tlbi1zYjJkZCIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VydmljZS1hY2NvdW50Lm5hbWUiOiJrby1hZG1pbiIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VydmljZS1hY2NvdW50LnVpZCI6IjA4MjRhMzc2LTNhYTUtNGZlMy1iMjg1LTNmMDJkZGVkMDk3NiIsInN1YiI6InN5c3RlbTpzZXJ2aWNlYWNjb3VudDprdWJlLXN5c3RlbTprby1hZG1pbiJ9.d-uHAwfT4UzgsVAkff71bKKMigpYaKVLb_JXMEQT_FL-QgRpuxpM912pklfg7308FvbFEezGY9QgWAnOHe-11_fCNPgx_MdbLwoCQPz_jwGNT-luSEdprVslVqyrHJh66oe-w-oEP-GydanOC6M2L079gdJxgVTq_GN3laeyWvEl3tUyE9bP8zFayk0ae7BZXV7bCe-hmpd3pcO5Z_Gtnrg8SYoIjdBhxQGYzBqrKZIIYBLCIxik44smjAfHpQFhzDO7Xv9z2W-hv6aTSDzA8rKRh7TpLUBKjtbrlNsixbSX8HE7I-mTmf4CV7FC2fZ-JikA362nMxq3_dogXGxy-w",
	})
	if err != nil {
		logger.Log.Fatal(err)
	}
	l, err := c.ApiextensionsV1beta1().CustomResourceDefinitions().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logger.Log.Fatal(err)
	}
	fmt.Println(l)
}
