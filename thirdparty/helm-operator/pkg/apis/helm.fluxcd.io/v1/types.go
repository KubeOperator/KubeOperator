package v1

import (
	"fmt"
	"strings"
	"time"

	"github.com/fluxcd/flux/pkg/resource"
	"github.com/ghodss/yaml"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/fluxcd/helm-operator/pkg/helm"
	helmv2 "github.com/fluxcd/helm-operator/pkg/helm/v2"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// HelmRelease represents custom resource associated with a Helm Chart
type HelmRelease struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`

	Spec   HelmReleaseSpec   `json:"spec"`
	Status HelmReleaseStatus `json:"status"`
}

// ResourceID returns an ID made from the identifying parts of the
// resource, as a convenience for Flux, which uses them
// everywhere.
func (hr HelmRelease) ResourceID() resource.ID {
	return resource.MakeID(hr.Namespace, "HelmRelease", hr.Name)
}

// GetReleaseName returns the configured release name, or constructs and
// returns one based on the namespace and name of the HelmRelease.
// When the HelmRelease's metadata.namespace and spec.targetNamespace
// differ, both are used in the generated name.
// This name is used for naming and operating on the release in Helm.
func (hr HelmRelease) GetReleaseName() string {
	if hr.Spec.ReleaseName == "" {
		namespace := hr.GetDefaultedNamespace()
		targetNamespace := hr.GetTargetNamespace()

		if namespace != targetNamespace {
			// prefix the releaseName with the administering HelmRelease namespace as well
			return fmt.Sprintf("%s-%s-%s", namespace, targetNamespace, hr.Name)
		}
		return fmt.Sprintf("%s-%s", targetNamespace, hr.Name)
	}

	return hr.Spec.ReleaseName
}

// GetDefaultedNamespace returns the HelmRelease's namespace
// defaulting to the "default" if not set.
func (hr HelmRelease) GetDefaultedNamespace() string {
	if hr.GetNamespace() == "" {
		return "default"
	}
	return hr.Namespace
}

// GetTargetNamespace returns the configured release targetNamespace
// defaulting to the namespace of the HelmRelease if not set.
func (hr HelmRelease) GetTargetNamespace() string {
	if hr.Spec.TargetNamespace == "" {
		return hr.GetDefaultedNamespace()
	}
	return hr.Spec.TargetNamespace
}

// ValuesFromSource represents a source of values.
// Only one of its fields may be set.
type ValuesFromSource struct {
	// Selects a key of a ConfigMap.
	// +optional
	ConfigMapKeyRef *NamespacedConfigMapKeySelector `json:"configMapKeyRef,omitempty"`
	// Selects a key of a Secret.
	// +optional
	SecretKeyRef *NamespacedSecretKeySelector `json:"secretKeyRef,omitempty"`
	// Selects an URL.
	// +optional
	ExternalSourceRef *ExternalSourceSelector `json:"externalSourceRef,omitempty"`
	// Selects a file from git source helm chart.
	// +optional
	ChartFileRef *ChartFileSelector `json:"chartFileRef,omitempty"`
}

type NamespacedConfigMapKeySelector struct {
	Name string `json:"name"`
	Key  string `json:"key"`
	// +optional
	Optional *bool `json:"optional,omitempty"`
	// +optional
	Namespace string `json:"namespace,omitempty"`
}
type NamespacedSecretKeySelector struct {
	Name string `json:"name"`
	Key  string `json:"key"`
	// +optional
	Optional *bool `json:"optional,omitempty"`
	// +optional
	Namespace string `json:"namespace,omitempty"`
}

type ChartFileSelector struct {
	Path string `json:"path"`
	// Do not fail if chart file could not be retrieved
	// +optional
	Optional *bool `json:"optional,omitempty"`
}

type ExternalSourceSelector struct {
	URL string `json:"url"`
	// Do not fail if external source could not be retrieved
	// +optional
	Optional *bool `json:"optional,omitempty"`
}

type ChartSource struct {
	// one of the following...
	// +optional
	*GitChartSource
	// +optional
	*RepoChartSource
}

type GitChartSource struct {
	GitURL    string                       `json:"git"`
	Ref       string                       `json:"ref"`
	Path      string                       `json:"path"`
	SecretRef *corev1.LocalObjectReference `json:"secretRef,omitempty"`
	// Do not run 'dep' update (assume requirements.yaml is already fulfilled)
	// +optional
	SkipDepUpdate bool `json:"skipDepUpdate,omitempty"`
}

// RefOrDefault returns the configured ref of the chart source. If the chart source
// does not specify a ref, the provided default is used instead.
func (s GitChartSource) RefOrDefault(defaultGitRef string) string {
	if s.Ref == "" {
		return defaultGitRef
	}
	return s.Ref
}

type RepoChartSource struct {
	RepoURL string `json:"repository"`
	Name    string `json:"name"`
	Version string `json:"version"`
	// An authentication secret for accessing the chart repo
	// +optional
	ChartPullSecret *corev1.LocalObjectReference `json:"chartPullSecret,omitempty"`
}

// CleanRepoURL returns the RepoURL but ensures it ends with a trailing slash
func (s RepoChartSource) CleanRepoURL() string {
	cleanURL := strings.TrimRight(s.RepoURL, "/")
	return cleanURL + "/"
}

type Rollback struct {
	Enable       bool   `json:"enable,omitempty"`
	Retry        bool   `json:"retry,omitempty"`
	MaxRetries   *int64 `json:"maxRetries,omitempty"`
	Force        bool   `json:"force,omitempty"`
	Recreate     bool   `json:"recreate,omitempty"`
	DisableHooks bool   `json:"disableHooks,omitempty"`
	Timeout      *int64 `json:"timeout,omitempty"`
	Wait         bool   `json:"wait,omitempty"`
}

func (r Rollback) GetTimeout() time.Duration {
	if r.Timeout == nil {
		return 300 * time.Second
	}
	return time.Duration(*r.Timeout) * time.Second
}

func (r Rollback) GetMaxRetries() int64 {
	if r.MaxRetries == nil {
		return 5
	}
	return *r.MaxRetries
}

// HelmReleaseSpec is the spec for a HelmRelease resource
type HelmReleaseSpec struct {
	ChartSource      `json:"chart"`
	HelmVersion      string                        `json:"helmVersion,omitempty"`
	ReleaseName      string                        `json:"releaseName,omitempty"`
	MaxHistory       *int                          `json:"maxHistory,omitempty"`
	ValueFileSecrets []corev1.LocalObjectReference `json:"valueFileSecrets,omitempty"`
	ValuesFrom       []ValuesFromSource            `json:"valuesFrom,omitempty"`
	HelmValues       `json:",inline"`
	// Override the target namespace, defaults to metadata.namespace
	// +optional
	TargetNamespace string `json:"targetNamespace,omitempty"`
	// Install or upgrade timeout in seconds
	// +optional
	Timeout *int64 `json:"timeout,omitempty"`
	// Reset values on helm upgrade
	// +optional
	ResetValues bool `json:"resetValues,omitempty"`
	// Skip helm3 CRD installation
	// +optional
	SkipCRDs bool `json:"skipCRDs,omitempty"`
	// Wait for the install or upgrade to complete before marking release as successful
	// +optional
	Wait bool `json:"wait,omitempty"`
	// Force resource update through delete/recreate, allows recovery from a failed state
	// +optional
	ForceUpgrade bool `json:"forceUpgrade,omitempty"`
	// Enable rollback and configure options
	// +optional
	Rollback Rollback `json:"rollback,omitempty"`
}

func (hr HelmRelease) GetHelmVersion(defaultVersion string) string {
	if hr.Spec.HelmVersion != "" {
		return hr.Spec.HelmVersion
	}
	if defaultVersion != "" {
		return defaultVersion
	}
	return helmv2.VERSION
}

// GetTimeout returns the install or upgrade timeout (defaults to 300s)
func (hr HelmRelease) GetTimeout() time.Duration {
	if hr.Spec.Timeout == nil {
		return 300 * time.Second
	}
	return time.Duration(*hr.Spec.Timeout) * time.Second
}

// GetMaxHistory returns the maximum number of release
// revisions to keep (defaults to 10)
func (hr HelmRelease) GetMaxHistory() int {
	if hr.Spec.MaxHistory == nil {
		return 10
	}
	return *hr.Spec.MaxHistory
}

// GetValuesFromSources maintains backwards compatibility with
// ValueFileSecrets by merging them into the ValuesFrom array.
func (hr HelmRelease) GetValuesFromSources() []ValuesFromSource {
	valuesFrom := hr.Spec.ValuesFrom
	// Maintain backwards compatibility with ValueFileSecrets
	if hr.Spec.ValueFileSecrets != nil {
		var secretKeyRefs []ValuesFromSource
		for _, ref := range hr.Spec.ValueFileSecrets {
			s := &NamespacedSecretKeySelector{Name: ref.Name}
			secretKeyRefs = append(secretKeyRefs, ValuesFromSource{SecretKeyRef: s})
		}
		valuesFrom = append(secretKeyRefs, valuesFrom...)
	}
	return valuesFrom
}

type HelmReleaseStatus struct {
	// ReleaseName is the name as either supplied or generated.
	// +optional
	ReleaseName string `json:"releaseName"`

	// ReleaseStatus is the status as given by Helm for the release
	// managed by this resource.
	ReleaseStatus string `json:"releaseStatus"`

	// ObservedGeneration is the most recent generation observed by
	// the controller.
	ObservedGeneration int64 `json:"observedGeneration"`

	// Revision would define what Git hash or Chart version has currently
	// been deployed.
	// +optional
	Revision string `json:"revision,omitempty"`

	// RollbackCount defines the amount of rollback attempts made,
	// it is incremented after a rollback failure and reset after a
	// successful upgrade or revision change.
	RollbackCount int64 `json:"rollbackCount,omitempty"`

	// Conditions contains observations of the resource's state, e.g.,
	// has the chart which it refers to been fetched.
	// +optional
	// +patchMergeKey=type
	// +patchStrategy=merge
	Conditions []HelmReleaseCondition `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type"`
}

type HelmReleaseCondition struct {
	Type   HelmReleaseConditionType `json:"type"`
	Status corev1.ConditionStatus   `json:"status"`
	// +optional
	LastUpdateTime metav1.Time `json:"lastUpdateTime,omitempty"`
	// +optional
	LastTransitionTime metav1.Time `json:"lastTransitionTime,omitempty"`
	// +optional
	Reason string `json:"reason,omitempty"`
	// +optional
	Message string `json:"message,omitempty"`
}

type HelmReleaseConditionType string

const (
	// ChartFetched means the chart to which the HelmRelease refers
	// has been fetched successfully
	HelmReleaseChartFetched HelmReleaseConditionType = "ChartFetched"
	// Released means the chart release, as specified in this
	// HelmRelease, has been processed by Helm.
	HelmReleaseReleased HelmReleaseConditionType = "Released"
	// RolledBack means the chart to which the HelmRelease refers
	// has been rolled back
	HelmReleaseRolledBack HelmReleaseConditionType = "RolledBack"
)

// FluxHelmValues embeds chartutil.Values so we can implement deepcopy on map[string]interface{}
// +k8s:deepcopy-gen=false
type HelmValues struct {
	helm.Values `json:"values,omitempty"`
}

// DeepCopyInto implements deepcopy-gen method for use in generated code
func (in *HelmValues) DeepCopyInto(out *HelmValues) {
	if in == nil {
		return
	}

	b, err := yaml.Marshal(in.Values)
	if err != nil {
		return
	}
	var values helm.Values
	err = yaml.Unmarshal(b, &values)
	if err != nil {
		return
	}
	out.Values = values
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// HelmReleaseList is a list of HelmRelease resources
type HelmReleaseList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []HelmRelease `json:"items"`
}
