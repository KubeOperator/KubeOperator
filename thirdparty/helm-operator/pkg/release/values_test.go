package release

import (
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"

	helmfluxv1 "github.com/fluxcd/helm-operator/pkg/apis/helm.fluxcd.io/v1"
	"github.com/fluxcd/helm-operator/pkg/helm"
)

func TestComposeValues(t *testing.T) {
	defaultNamespace := "flux"
	otherNamespace := "other-namespace"
	falseVal := false

	client := fake.NewSimpleClientset(
		&corev1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "release-configmap",
				Namespace: defaultNamespace,
			},
			Data: map[string]string{
				"values.yaml": `valuesDict:
  same-namespace-configmap: true`,
			},
		},
		&corev1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "release-configmap",
				Namespace: otherNamespace,
			},
			Data: map[string]string{
				"values.yaml": `valuesDict:
  cross-namespace-configmap: true`,
			},
		},
		&corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "release-secret",
				Namespace: defaultNamespace,
			},
			Data: map[string][]byte{
				"values.yaml": []byte(`valuesDict:
  same-namespace-secret: true`),
			},
		},
		&corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "release-secret",
				Namespace: otherNamespace,
			},
			Data: map[string][]byte{
				"values.yaml": []byte(`valuesDict:
  cross-namespace-secret: true`),
			},
		},
	)

	cases := []struct {
		description      string
		releaseNamespace string
		valuesFromSource []helmfluxv1.ValuesFromSource
		assertions       []func(*testing.T, helm.Values)
	}{
		{
			description:      "simple same-namespace test",
			releaseNamespace: defaultNamespace,
			valuesFromSource: []helmfluxv1.ValuesFromSource{
				{
					ConfigMapKeyRef: &helmfluxv1.NamespacedConfigMapKeySelector{
						Name:     "release-configmap",
						Key:      "values.yaml",
						Optional: &falseVal,
					},
					SecretKeyRef:      nil,
					ExternalSourceRef: nil,
					ChartFileRef:      nil,
				},
				{
					ConfigMapKeyRef: nil,
					SecretKeyRef: &helmfluxv1.NamespacedSecretKeySelector{
						Name:     "release-secret",
						Key:      "values.yaml",
						Optional: &falseVal,
					},
					ExternalSourceRef: nil,
					ChartFileRef:      nil,
				},
			},
			assertions: []func(*testing.T, helm.Values){
				func(t *testing.T, values helm.Values) {
					assert.NotNil(t, values["valuesDict"].(map[string]interface{})["same-namespace-configmap"])
				},
				func(t *testing.T, values helm.Values) {
					assert.NotNil(t, values["valuesDict"].(map[string]interface{})["same-namespace-secret"])
				},
			},
		},
		{
			description:      "simple cross-namespace test",
			releaseNamespace: defaultNamespace,
			valuesFromSource: []helmfluxv1.ValuesFromSource{
				{
					ConfigMapKeyRef: &helmfluxv1.NamespacedConfigMapKeySelector{
						Name:      "release-configmap",
						Namespace: otherNamespace,
						Key:       "values.yaml",
						Optional:  &falseVal,
					},
					SecretKeyRef:      nil,
					ExternalSourceRef: nil,
					ChartFileRef:      nil,
				},
				{
					ConfigMapKeyRef: nil,
					SecretKeyRef: &helmfluxv1.NamespacedSecretKeySelector{
						Name:      "release-secret",
						Namespace: otherNamespace,
						Key:       "values.yaml",
						Optional:  &falseVal,
					},
					ExternalSourceRef: nil,
					ChartFileRef:      nil,
				},
			},
			assertions: []func(*testing.T, helm.Values){
				func(t *testing.T, values helm.Values) {
					assert.NotNil(t, values["valuesDict"].(map[string]interface{})["cross-namespace-configmap"])
				},
				func(t *testing.T, values helm.Values) {
					assert.NotNil(t, values["valuesDict"].(map[string]interface{})["cross-namespace-secret"])
				},
			},
		},
		{
			description:      "same and cross-namespace test",
			releaseNamespace: defaultNamespace,
			valuesFromSource: []helmfluxv1.ValuesFromSource{
				{
					ConfigMapKeyRef: &helmfluxv1.NamespacedConfigMapKeySelector{
						Name:     "release-configmap",
						Key:      "values.yaml",
						Optional: &falseVal,
					},
					SecretKeyRef:      nil,
					ExternalSourceRef: nil,
					ChartFileRef:      nil,
				},
				{
					ConfigMapKeyRef: nil,
					SecretKeyRef: &helmfluxv1.NamespacedSecretKeySelector{
						Name:     "release-secret",
						Key:      "values.yaml",
						Optional: &falseVal,
					},
					ExternalSourceRef: nil,
					ChartFileRef:      nil,
				},
				{
					ConfigMapKeyRef: &helmfluxv1.NamespacedConfigMapKeySelector{
						Name:      "release-configmap",
						Namespace: otherNamespace,
						Key:       "values.yaml",
						Optional:  &falseVal,
					},
					SecretKeyRef:      nil,
					ExternalSourceRef: nil,
					ChartFileRef:      nil,
				},
				{
					ConfigMapKeyRef: nil,
					SecretKeyRef: &helmfluxv1.NamespacedSecretKeySelector{
						Name:      "release-secret",
						Namespace: otherNamespace,
						Key:       "values.yaml",
						Optional:  &falseVal,
					},
					ExternalSourceRef: nil,
					ChartFileRef:      nil,
				},
			},
			assertions: []func(*testing.T, helm.Values){
				func(t *testing.T, values helm.Values) {
					assert.NotNil(t, values["valuesDict"].(map[string]interface{})["cross-namespace-configmap"])
				},
				func(t *testing.T, values helm.Values) {
					assert.NotNil(t, values["valuesDict"].(map[string]interface{})["cross-namespace-secret"])
				},
				func(t *testing.T, values helm.Values) {
					assert.NotNil(t, values["valuesDict"].(map[string]interface{})["cross-namespace-configmap"])
				},
				func(t *testing.T, values helm.Values) {
					assert.NotNil(t, values["valuesDict"].(map[string]interface{})["cross-namespace-secret"])
				},
			},
		},
	}
	for _, c := range cases {
		t.Run(c.description, func(t *testing.T) {
			hr := &helmfluxv1.HelmRelease{
				Spec: helmfluxv1.HelmReleaseSpec{
					ValuesFrom: c.valuesFromSource,
				},
			}
			hr.Namespace = c.releaseNamespace

			values, err := composeValues(client.CoreV1(), hr, "")
			t.Log(values)
			assert.NoError(t, err)
			for _, assertion := range c.assertions {
				assertion(t, values)
			}
		})
	}
}
