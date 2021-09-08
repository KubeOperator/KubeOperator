package polaris

import (
	"context"
	"fmt"
	"github.com/KubeOperator/KubeOperator/pkg/dto"
	conf "github.com/fairwindsops/polaris/pkg/config"
	"github.com/fairwindsops/polaris/pkg/kube"
	"github.com/fairwindsops/polaris/pkg/validator"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type Config struct {
	Host  string
	Token string
	Port  int
}

func NewResourceProvider(c *Config) (*kube.ResourceProvider, error) {

	kubeConf := &rest.Config{
		Host:        fmt.Sprintf("%s:%d", c.Host, c.Port),
		BearerToken: c.Token,
		TLSClientConfig: rest.TLSClientConfig{
			Insecure: true,
		},
	}
	api, err := kubernetes.NewForConfig(kubeConf)
	if err != nil {
		return nil, err
	}
	dynamicInterface, err := dynamic.NewForConfig(kubeConf)
	if err != nil {
		return nil, err
	}
	config := NewPolarisConfig()
	return kube.CreateResourceProviderFromAPI(context.Background(), api, "", &dynamicInterface, config)
}

func NewPolarisConfig() conf.Configuration {
	return conf.Configuration{
		DisplayName: "kube-operator",
		Checks: map[string]conf.Severity{
			//# reliability
			"multipleReplicasForDeployment": conf.SeverityIgnore,
			"priorityClassNotSet":           conf.SeverityIgnore,
			"tagNotSpecified":               conf.SeverityDanger,
			"pullPolicyNotAlways":           conf.SeverityWarning,
			"readinessProbeMissing":         conf.SeverityWarning,
			"livenessProbeMissing":          conf.SeverityWarning,
			//# efficiency
			"cpuRequestsMissing":    conf.SeverityWarning,
			"cpuLimitsMissing":      conf.SeverityWarning,
			"memoryRequestsMissing": conf.SeverityWarning,
			"memoryLimitsMissing":   conf.SeverityWarning,
			//# security
			"hostIPCSet":                 conf.SeverityDanger,
			"hostPIDSet":                 conf.SeverityDanger,
			"notReadOnlyRootFilesystem":  conf.SeverityWarning,
			"privilegeEscalationAllowed": conf.SeverityDanger,
			"runAsRootAllowed":           conf.SeverityWarning,
			"runAsPrivileged":            conf.SeverityDanger,
			"dangerousCapabilities":      conf.SeverityDanger,
			"insecureCapabilities":       conf.SeverityWarning,
			"hostNetworkSet":             conf.SeverityWarning,
			"hostPortSet":                conf.SeverityWarning,
			//# custom
			"resourceLimits":             conf.SeverityWarning,
			"imageRegistry":              conf.SeverityDanger,
			"metadataAndNameMismatched":  conf.SeverityIgnore,
			"pdbDisruptionsIsZero":       conf.SeverityWarning,
			"missingPodDisruptionBudget": conf.SeverityIgnore,
			"tlsSettingsMissing":         conf.SeverityWarning,
		},
		Exemptions: []conf.Exemption{
			{
				ControllerNames: []string{"kube-apiserver", "kube-proxy", "kube-scheduler", "etcd-manager-events", "kube-controller-manager", "kube-dns", "etcd-manager-main"},
				Rules:           []string{"hostPortSet", "hostNetworkSet", "readinessProbeMissing", "livenessProbeMissing", "cpuRequestsMissing", "cpuLimitsMissing", "memoryRequestsMissing", "memoryLimitsMissing", "runAsRootAllowed", "runAsPrivileged", "notReadOnlyRootFilesystem", "hostPIDSet"},
			},
			{
				ControllerNames: []string{"kube-flannel-ds"},
				Rules:           []string{"notReadOnlyRootFilesystem", "runAsRootAllowed", "notReadOnlyRootFilesystem", "readinessProbeMissing", "livenessProbeMissing", "cpuLimitsMissing"},
			},
			{
				ControllerNames: []string{"cert-manager"},
				Rules:           []string{"notReadOnlyRootFilesystem", "runAsRootAllowed", "readinessProbeMissing", "livenessProbeMissing"},
			},
			{
				ControllerNames: []string{"cluster-autoscaler"},
				Rules:           []string{"notReadOnlyRootFilesystem", "runAsRootAllowed", "readinessProbeMissing"},
			},
			{
				ControllerNames: []string{"vpa"},
				Rules:           []string{"runAsRootAllowed", "livenessProbeMissing", "readinessProbeMissing", "notReadOnlyRootFilesystem"},
			},
			{
				ControllerNames: []string{"datadog"},
				Rules:           []string{"runAsRootAllowed", "livenessProbeMissing", "readinessProbeMissing", "notReadOnlyRootFilesystem"},
			},
			{
				ControllerNames: []string{"nginx-ingress-controller"},
				Rules:           []string{"runAsRootAllowed", "privilegeEscalationAllowed", "insecureCapabilities"},
			},
			{
				ControllerNames: []string{"dns-controller", "datadog-datadog", "kube-flannel-ds", "kube2iam", "aws-iam-authenticator", "datadog", "kube2iam"},
				Rules:           []string{"hostNetworkSet"},
			},
			{
				ControllerNames: []string{"aws-iam-authenticator", "aws-cluster-autoscaler", "kube-state-metrics", "dns-controller", " external-dns", "dnsmasq", "autoscaler", "kubernetes-dashboard", "install-cni", "kube2iam"},
				Rules:           []string{"readinessProbeMissing", "livenessProbeMissing"},
			},
			{
				ControllerNames: []string{"aws-iam-authenticator", "nginx-ingress-default-backend", "aws-cluster-autoscaler", "kube-state-metrics", "dns-controller", "external-dns", "kubedns", "dnsmasq", "autoscaler", "tiller", "kube2iam"},
				Rules:           []string{"runAsRootAllowed"},
			},
			{
				ControllerNames: []string{"aws-iam-authenticator", "nginx-ingress-controller", "nginx-ingress-default-backend", "aws-cluster-autoscaler", "kube-state-metrics", "dns-controller", "external-dns", "kubedns", "dnsmasq", "autoscaler", "tiller", "kube2iam"},
				Rules:           []string{"notReadOnlyRootFilesystem"},
			},
			{
				ControllerNames: []string{"cert-manager", "dns-controller", "kubedns", "dnsmasq", "autoscaler", "insights-agent-goldilocks-vpa-install", "datadog"},
				Rules:           []string{"cpuRequestsMissing", "cpuLimitsMissing", "memoryRequestsMissing", "memoryLimitsMissing"},
			},
			{
				ControllerNames: []string{"kube2iam", "kube-flannel-ds"},
				Rules:           []string{"runAsPrivileged"},
			},
			{
				ControllerNames: []string{"kube-hunter"},
				Rules:           []string{"hostPIDSet"},
			},
			{
				ControllerNames: []string{"polaris", "kube-hunter", "goldilocks", "insights-agent-goldilocks-vpa-install"},
				Rules:           []string{"notReadOnlyRootFilesystem"},
			},
			{
				ControllerNames: []string{"insights-agent-goldilocks-controller"},
				Rules:           []string{"livenessProbeMissing", "readinessProbeMissing"},
			},
			{
				ControllerNames: []string{"insights-agent-goldilocks-vpa-install", "kube-hunter"},
				Rules:           []string{"runAsRootAllowed"},
			},
		},
	}
}

func RunGrade(c *Config) (*dto.ClusterGrade, error) {
	k, err := NewResourceProvider(c)
	if err != nil {
		return nil, err
	}

	config := NewPolarisConfig()
	auditData, err := validator.RunAudit(config, k)
	if err != nil {
		return nil, err
	}
	return GetCLusterGrade(auditData), nil
}

func GetCLusterGrade(data validator.AuditData) *dto.ClusterGrade {
	var clusterGrade dto.ClusterGrade
	score := data.GetSummary().GetScore()
	clusterGrade.Score = int(score)
	sums := data.GetSummaryByCategory()
	var total dto.Summary
	list := make(map[string]dto.Summary)
	for k, v := range sums {
		list[k] = dto.Summary{
			Danger:  int(v.Dangers),
			Success: int(v.Successes),
			Warning: int(v.Warnings),
		}
		total.Warning = total.Warning + int(v.Warnings)
		total.Success = total.Success + int(v.Successes)
		total.Danger = total.Danger + int(v.Dangers)
	}
	clusterGrade.TotalSum = total
	clusterGrade.ListSum = list
	namespaceResults := data.GetResultsByNamespace()
	for i := range namespaceResults {
		namespaceResult := dto.NamespaceResult{
			Namespace: i,
		}
		for _, v := range namespaceResults[i] {
			detail := dto.NamespaceResultDetail{
				Name: v.Name,
				Kind: v.Kind,
			}
			if v.PodResult != nil && v.PodResult.Results != nil {
				for _, pod := range v.PodResult.Results {
					podResult := dto.PodResult{
						ID:       pod.ID,
						Message:  pod.Message,
						Category: pod.Category,
						Success:  pod.Success,
						Severity: string(pod.Severity),
					}
					detail.PodResults = append(detail.PodResults, podResult)
				}
			}
			namespaceResult.Results = append(namespaceResult.Results, detail)
		}
		clusterGrade.Results = append(clusterGrade.Results, namespaceResult)
	}
	return &clusterGrade
}
