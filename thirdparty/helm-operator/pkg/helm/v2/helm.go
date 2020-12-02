package v2

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/go-kit/kit/log"
	"google.golang.org/grpc/status"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/helm/pkg/getter"
	helmv2 "k8s.io/helm/pkg/helm"
	"k8s.io/helm/pkg/helm/environment"
	"k8s.io/helm/pkg/helm/helmpath"
	"k8s.io/helm/pkg/tlsutil"

	"github.com/fluxcd/helm-operator/pkg/helm"
)

const VERSION = "v2"

var (
	repositoryConfig = helmHome().RepositoryFile()
	repositoryCache  = helmHome().Cache()
)

// TillerOptions holds configuration options for tiller
type TillerOptions struct {
	Host        string
	Port        string
	Namespace   string
	TLSVerify   bool
	TLSEnable   bool
	TLSKey      string
	TLSCert     string
	TLSCACert   string
	TLSHostname string
}

// HelmV2 provides access to the Helm v2 client, while adhering
// to the generic Client interface
type HelmV2 struct {
	client *helmv2.Client
	logger log.Logger
}

func (h *HelmV2) Version() string {
	return VERSION
}

// New attempts to setup a Helm client
func New(logger log.Logger, kubeClient *kubernetes.Clientset, opts TillerOptions) helm.Client {
	var h *HelmV2
	for {
		client, host, err := newHelmClient(kubeClient, opts)
		if err != nil {
			logger.Log("error", fmt.Sprintf("error creating Client (v2) client: %s", err.Error()))
			time.Sleep(20 * time.Second)
			continue
		}
		h = &HelmV2{client: client, logger: logger}
		version, err := h.getVersion()
		if err != nil {
			logger.Log("warning", "unable to connect to Tiller", "err", err, "host", host, "options", fmt.Sprintf("%+v", opts))
			time.Sleep(20 * time.Second)
			continue
		}
		logger.Log("info", "connected to Tiller", "version", version, "host", host, "options", fmt.Sprintf("%+v", opts))
		break
	}
	return h
}

// getVersion retrieves the Tiller version. This is a _V2 only_  method
// and used internally during the setup of the client.
func (h *HelmV2) getVersion() (string, error) {
	v, err := h.client.GetVersion()
	if err != nil {
		return "", fmt.Errorf("error getting tiller version: %v", err)
	}
	return v.GetVersion().String(), nil
}

// newHelmClient creates a new Helm v2 client
func newHelmClient(kubeClient *kubernetes.Clientset, opts TillerOptions) (*helmv2.Client, string, error) {
	host, err := tillerHost(kubeClient, opts)
	if err != nil {
		return nil, "", err
	}

	// host = "tiller-deploy.kube-system:44134"
	options := []helmv2.Option{helmv2.Host(host)}
	if opts.TLSVerify || opts.TLSEnable {
		tlsopts := tlsutil.Options{
			KeyFile:            opts.TLSKey,
			CertFile:           opts.TLSCert,
			InsecureSkipVerify: true,
		}
		if opts.TLSVerify {
			tlsopts.CaCertFile = opts.TLSCACert
			tlsopts.InsecureSkipVerify = false
		}
		if opts.TLSHostname != "" {
			tlsopts.ServerName = opts.TLSHostname
		}
		tlsCfg, err := tlsutil.ClientConfig(tlsopts)
		if err != nil {
			return nil, "", err
		}
		options = append(options, helmv2.WithTLS(tlsCfg))
	}

	return helmv2.NewClient(options...), host, nil
}

// tillerHost either constructs the host string based on the configured options
// or attempts to resolve the `tiller-deploy` service in the configured namespace,
// in case of any failure during the resolving of the service it returns the error.
func tillerHost(kubeClient *kubernetes.Clientset, opts TillerOptions) (string, error) {
	if opts.Host == "" || opts.Port == "" {
		ts, err := kubeClient.CoreV1().Services(opts.Namespace).Get("tiller-deploy", metav1.GetOptions{})
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("%s.%s:%v", ts.Name, ts.Namespace, ts.Spec.Ports[0].Port), nil
	}

	return fmt.Sprintf("%s:%s", opts.Host, opts.Port), nil
}

func getterProviders() getter.Providers {
	return getter.All(environment.EnvSettings{
		Home: helmHome(),
	})
}

func helmHome() helmpath.Home {
	if v, ok := os.LookupEnv("HELM_HOME"); ok {
		return helmpath.Home(v)
	}
	return helmpath.Home(environment.DefaultHelmHome)
}

func statusMessageErr(err error) error {
	if s, ok := status.FromError(err); ok {
		return errors.New(s.Message())
	}
	return err
}
