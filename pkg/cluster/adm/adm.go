package adm

import (
	"fmt"
	"github.com/google/martian/log"
	"github.com/pkg/errors"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"ko3-gin/pkg/cluster"
	"ko3-gin/pkg/util/ssh"
	"reflect"
	"runtime"
	"strings"
	"time"
)

const (
	ReasonFailedProcess     = "FailedProcess"
	ReasonWaitingProcess    = "WaitingProcess"
	ReasonSuccessfulProcess = "SuccessfulProcess"
	ReasonSkipProcess       = "SkipProcess"

	ConditionTypeDone = "EnsureDone"
)

type Registry struct {
	Domain string
	Prefix string
	IP     string
}

type Credential struct {
	ClusterName       string
	ETCDCACert        []byte
	ETCDCAKey         []byte
	ETCDAPIClientCert []byte
	ETCDAPIClientKey  []byte
	CACert            []byte
	CAKey             []byte
	ClientCert        []byte
	ClientKey         []byte
	Token             *string
	BootstrapToken    *string
	CertificateKey    *string
}

type Cluster struct {
	cluster.Cluster
	Registry
	Credential Credential
	SSH        map[string]ssh.Interface
}

func (c *Cluster) SetCondition(newCondition cluster.Condition) {
	var conditions []cluster.Condition
	exist := false
	for _, condition := range c.Status.Conditions {
		if condition.Type == newCondition.Type {
			exist = true
			if newCondition.Status != condition.Status {
				condition.Status = newCondition.Status
			}
			if newCondition.Message != condition.Message {
				condition.Message = newCondition.Message
			}
			if newCondition.Reason != condition.Reason {
				condition.Reason = newCondition.Reason
			}
			if !newCondition.LastProbeTime.IsZero() && newCondition.LastProbeTime != condition.LastProbeTime {
				condition.LastProbeTime = newCondition.LastProbeTime
			}
			if !newCondition.LastTransitionTime.IsZero() && newCondition.LastTransitionTime != condition.LastTransitionTime {
				condition.LastTransitionTime = newCondition.LastTransitionTime
			}
		}
		conditions = append(conditions, condition)
	}
	if !exist {
		if newCondition.LastProbeTime.IsZero() {
			newCondition.LastProbeTime = time.Now()
		}
		if newCondition.LastTransitionTime.IsZero() {
			newCondition.LastTransitionTime = time.Now()
		}
		conditions = append(conditions, newCondition)
	}
	c.Status.Conditions = conditions

}
func (c *Cluster) Clientset() (*kubernetes.Clientset, error) {
	restConfig := &rest.Config{
		Host:        fmt.Sprintf("https://%s:6443", c.Spec.Machines[0].IP), // use the first host because the rest probably not join
		BearerToken: *c.Credential.Token,
		TLSClientConfig: rest.TLSClientConfig{
			CAData: c.Credential.CACert,
		},
		Timeout: 5 * time.Second,
	}

	return kubernetes.NewForConfig(restConfig)
}

func NewCluster(cluster cluster.Cluster) (*Cluster, error) {
	c := &Cluster{
		Cluster: cluster,
	}
	c.SSH = make(map[string]ssh.Interface)
	for _, m := range c.Spec.Machines {
		sshCfg := &ssh.Config{
			User:       m.Username,
			Host:       m.IP,
			Port:       m.Port,
			Password:   string(m.Password),
			PrivateKey: m.PrivateKey,
			PassPhrase: m.PassPhrase,
		}
		s, err := ssh.New(sshCfg)
		if err != nil {
			return nil, errors.Wrap(err, "Create ssh error")
		}
		c.SSH[m.IP] = s
	}
	return c, nil
}

type Handler func(*Cluster) error

func (h Handler) name() string {
	name := runtime.FuncForPC(reflect.ValueOf(h).Pointer()).Name()
	i := strings.Index(name, "Ensure")
	if i == -1 {
		return ""
	}
	return strings.TrimSuffix(name[i:], "-fm")
}

type ClusterAdm struct {
	createHandlers []Handler
}

func NewClusterAdm() (*ClusterAdm, error) {
	ca := new(ClusterAdm)
	ca.createHandlers = []Handler{
		ca.EnsureKernelModule,
		ca.EnsureSysctl,
		ca.EnsureDisableSwap,
		ca.EnsureClusterComplete,
		ca.EnsureDocker,
		ca.EnsureKubelet,
		ca.EnsureKubeadm,
		ca.EnsurePrepareForControlplane,
		ca.EnsureKubeadmInitKubeletStartPhase,
		ca.EnsureKubeadmInitCertsPhase,
		ca.EnsureStoreCredential,
		ca.EnsureKubeconfig,
		ca.EnsureKubeadmInitKubeConfigPhase,
		ca.EnsureKubeadmInitControlPlanePhase,
		ca.EnsureKubeadmInitEtcdPhase,
		ca.EnsureKubeadmInitWaitControlPlanePhase,
		ca.EnsureKubeadmInitUploadConfigPhase,
		ca.EnsureKubeadmInitUploadCertsPhase,
		ca.EnsureKubeadmInitBootstrapTokenPhase,
		ca.EnsureKubeadmInitAddonPhase,
		ca.EnsureJoinControlePlane,
	}
	return ca, nil
}

func (ca *ClusterAdm) OnInitialize(args cluster.Cluster) (cluster.Cluster, error) {
	c, err := NewCluster(args)
	if err != nil {
		return c.Cluster, err
	}
	err = ca.create(c)
	return c.Cluster, err
}

func (ca *ClusterAdm) create(c *Cluster) error {
	condition, err := ca.getCreateCurrentCondition(c)
	if err != nil {
		return err
	}
	now := time.Now()
	f := ca.getCreateHandler(condition.Type)
	if f == nil {
		return fmt.Errorf("can't get handler by %s", condition.Type)
	}
	err = f(c)
	if err != nil {
		c.SetCondition(cluster.Condition{
			Type:          condition.Type,
			Status:        cluster.ConditionFalse,
			LastProbeTime: now,
			Message:       err.Error(),
			Reason:        ReasonFailedProcess,
		})
		c.Status.Reason = ReasonFailedProcess
		c.Status.Message = err.Error()
		return nil
	}

	c.SetCondition(cluster.Condition{
		Type:               condition.Type,
		Status:             cluster.ConditionTrue,
		LastProbeTime:      now,
		LastTransitionTime: now,
		Reason:             ReasonSuccessfulProcess,
	})

	nextConditionType := ca.getNextConditionType(condition.Type)
	if nextConditionType == ConditionTypeDone {
		c.Status.Phase = cluster.ClusterRunning
	} else {
		c.SetCondition(cluster.Condition{
			Type:               nextConditionType,
			Status:             cluster.ConditionUnknown,
			LastProbeTime:      now,
			LastTransitionTime: now,
			Message:            "waiting process",
			Reason:             ReasonWaitingProcess,
		})

		log.Infof("%s is done, next is %s", condition.Type, nextConditionType)
	}
	return nil
}

func (ca *ClusterAdm) getCreateHandler(conditionType string) Handler {
	for _, f := range ca.createHandlers {
		if conditionType == f.name() {
			return f
		}
	}

	return nil
}

func (ca *ClusterAdm) getNextConditionType(conditionType string) string {
	var (
		i int
		f Handler
	)
	for i, f = range ca.createHandlers {
		name := runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
		if strings.Contains(name, conditionType) {
			break
		}
	}
	if i == len(ca.createHandlers)-1 {
		return ConditionTypeDone
	}
	next := ca.createHandlers[i+1]

	return next.name()
}

func (ca *ClusterAdm) getCreateCurrentCondition(c *Cluster) (*cluster.Condition, error) {
	if c.Status.Phase == cluster.ClusterRunning {
		return nil, errors.New("cluster phase is running now")
	}
	if len(ca.createHandlers) == 0 {
		return nil, errors.New("no create handlers")
	}

	if len(c.Status.Conditions) == 0 {
		return &cluster.Condition{
			Type:          ca.createHandlers[0].name(),
			Status:        cluster.ConditionUnknown,
			LastProbeTime: time.Now(),
			Message:       "waiting process",
			Reason:        ReasonWaitingProcess,
		}, nil
	}

	for _, condition := range c.Status.Conditions {
		if condition.Status == cluster.ConditionFalse || condition.Status == cluster.ConditionUnknown {
			return &condition, nil
		}
	}

	return nil, errors.New("no condition need process")
}
