package kubernetes

import (
	"context"
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"log"
	"testing"
)

func TestNewKubernetesClient(t *testing.T) {
	c, err := NewKubernetesExtensionClient(&Config{
		Host:  "172.16.10.210",
		Token: "eyJhbGciOiJSUzI1NiIsImtpZCI6IllzTEdnV2VCamtlRTlyTkFpWmNuamdITWljRGNzc0lvS2p0X0JCbEtETjQifQ.eyJpc3MiOiJrdWJlcm5ldGVzL3NlcnZpY2VhY2NvdW50Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9uYW1lc3BhY2UiOiJrdWJlLXN5c3RlbSIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VjcmV0Lm5hbWUiOiJrby1hZG1pbi10b2tlbi1idnA5cSIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VydmljZS1hY2NvdW50Lm5hbWUiOiJrby1hZG1pbiIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VydmljZS1hY2NvdW50LnVpZCI6IjU1MmZiMTg0LTdiOGMtNDRkNi1iNzI0LTZhMGY5YzdhNTdlNyIsInN1YiI6InN5c3RlbTpzZXJ2aWNlYWNjb3VudDprdWJlLXN5c3RlbTprby1hZG1pbiJ9.Hd0I9mjMAw0mgkhoKiG6_zlJ4S0sPLfTtSOTa7cg0VIyordqycQZ2pUpEEFXImGFrt3-EUcWJ_kmQbgCVRtJgw-FHEe84yKJx3UiPwu_gSC3ecx6lALoSJeOVC7AD3M4zwCvY2xhle6Kp6V3V4Asrp-m_uNj-ojJlimCbgevNqItMNncl4kS8qtEE3S3fKEi7Y3k73J06j_g5LFxN86QbVaNcxntu9xqqAC7ouxW_mMEcgQTT7mLzaR6aGrhyimxyMZCT0VzqrxXLXYLo7DI21zgLNhHWJ_EsKxapIK9kZmz_hT8sTf2MATSqH_Dw2oLn2ptwSjI9qNdMvdlkxvRRQ",
		Port:  8443,
	})
	if err != nil {
		log.Fatal(err)
	}
	l, err := c.ApiextensionsV1beta1().CustomResourceDefinitions().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(l)
}
