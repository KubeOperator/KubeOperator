package sync

import (
	"context"
	"encoding/json"
	"fmt"

	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	kubernetes "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/rest"
)

const syncMarkerKey = "flux.weave.works/sync-hwm"

// NativeSyncProvider keeps information related to the native state of a sync marker stored in a "native" kubernetes resource.
type NativeSyncProvider struct {
	namespace    string
	revision     string
	resourceName string
	resourceAPI  v1.SecretInterface
}

// NewNativeSyncProvider creates a new NativeSyncProvider
func NewNativeSyncProvider(namespace string, resourceName string) (NativeSyncProvider, error) {
	//config, err := rest.InClusterConfig()
	//if err != nil {
	//	return NativeSyncProvider{}, err
	//}

	kubeConf := &rest.Config{
		Host:        fmt.Sprintf("%s:%d", "https://172.16.10.210", 8443),
		BearerToken: "eyJhbGciOiJSUzI1NiIsImtpZCI6IllzTEdnV2VCamtlRTlyTkFpWmNuamdITWljRGNzc0lvS2p0X0JCbEtETjQifQ.eyJpc3MiOiJrdWJlcm5ldGVzL3NlcnZpY2VhY2NvdW50Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9uYW1lc3BhY2UiOiJrdWJlLXN5c3RlbSIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VjcmV0Lm5hbWUiOiJrby1hZG1pbi10b2tlbi1idnA5cSIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VydmljZS1hY2NvdW50Lm5hbWUiOiJrby1hZG1pbiIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VydmljZS1hY2NvdW50LnVpZCI6IjU1MmZiMTg0LTdiOGMtNDRkNi1iNzI0LTZhMGY5YzdhNTdlNyIsInN1YiI6InN5c3RlbTpzZXJ2aWNlYWNjb3VudDprdWJlLXN5c3RlbTprby1hZG1pbiJ9.Hd0I9mjMAw0mgkhoKiG6_zlJ4S0sPLfTtSOTa7cg0VIyordqycQZ2pUpEEFXImGFrt3-EUcWJ_kmQbgCVRtJgw-FHEe84yKJx3UiPwu_gSC3ecx6lALoSJeOVC7AD3M4zwCvY2xhle6Kp6V3V4Asrp-m_uNj-ojJlimCbgevNqItMNncl4kS8qtEE3S3fKEi7Y3k73J06j_g5LFxN86QbVaNcxntu9xqqAC7ouxW_mMEcgQTT7mLzaR6aGrhyimxyMZCT0VzqrxXLXYLo7DI21zgLNhHWJ_EsKxapIK9kZmz_hT8sTf2MATSqH_Dw2oLn2ptwSjI9qNdMvdlkxvRRQ",
		TLSClientConfig: rest.TLSClientConfig{
			Insecure: true,
		},
	}
	clientset, err := kubernetes.NewForConfig(kubeConf)
	if err != nil {
		return NativeSyncProvider{}, err
	}

	return NativeSyncProvider{
		resourceAPI:  clientset.CoreV1().Secrets(namespace),
		namespace:    namespace,
		resourceName: resourceName,
	}, nil
}

func (p NativeSyncProvider) String() string {
	return "kubernetes " + p.namespace + ":secret/" + p.resourceName
}

// GetRevision gets the revision of the current sync marker (representing the place flux has synced to).
func (p NativeSyncProvider) GetRevision(ctx context.Context) (string, error) {
	resource, err := p.resourceAPI.Get(context.TODO(), p.resourceName, meta_v1.GetOptions{})
	if err != nil {
		return "", err
	}
	revision, exists := resource.Annotations[syncMarkerKey]
	if !exists {
		return "", p.setRevision("")
	}
	return revision, nil
}

// UpdateMarker updates the revision the sync marker points to.
func (p NativeSyncProvider) UpdateMarker(ctx context.Context, revision string) error {
	return p.setRevision(revision)
}

// DeleteMarker resets the state of the object.
func (p NativeSyncProvider) DeleteMarker(ctx context.Context) error {
	return p.setRevision("")
}

func (p NativeSyncProvider) setRevision(revision string) error {
	jsonPatch, err := json.Marshal(patch(revision))
	if err != nil {
		return err
	}

	_, err = p.resourceAPI.Patch(
		context.TODO(),
		p.resourceName,
		types.StrategicMergePatchType,
		jsonPatch,
		meta_v1.PatchOptions{},
	)
	return err
}

func patch(revision string) map[string]map[string]map[string]string {
	return map[string]map[string]map[string]string{
		"metadata": map[string]map[string]string{
			"annotations": map[string]string{
				syncMarkerKey: revision,
			},
		},
	}
}
