package main

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/spf13/pflag"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/klog"

	"github.com/fluxcd/flux/pkg/checkpoint"

	"github.com/fluxcd/helm-operator/pkg/chartsync"
	clientset "github.com/fluxcd/helm-operator/pkg/client/clientset/versioned"
	ifinformers "github.com/fluxcd/helm-operator/pkg/client/informers/externalversions"
	"github.com/fluxcd/helm-operator/pkg/helm"
	"github.com/fluxcd/helm-operator/pkg/helm/v2"
	"github.com/fluxcd/helm-operator/pkg/helm/v3"
	daemonhttp "github.com/fluxcd/helm-operator/pkg/http/daemon"
	"github.com/fluxcd/helm-operator/pkg/operator"
	"github.com/fluxcd/helm-operator/pkg/release"
	"github.com/fluxcd/helm-operator/pkg/status"
)

var (
	fs     *pflag.FlagSet
	logger log.Logger

	versionFlag *bool

	logFormat *string

	kubeconfig *string
	master     *string
	namespace  *string

	workers *int

	tillerIP        *string
	tillerPort      *string
	tillerNamespace *string

	tillerTLSVerify   *bool
	tillerTLSEnable   *bool
	tillerTLSKey      *string
	tillerTLSCert     *string
	tillerTLSCACert   *string
	tillerTLSHostname *string

	chartsSyncInterval   *time.Duration
	statusUpdateInterval *time.Duration
	logReleaseDiffs      *bool
	updateDependencies   *bool

	gitTimeout      *time.Duration
	gitPollInterval *time.Duration
	gitDefaultRef   *string

	listenAddr *string

	versionedHelmRepositoryIndexes *[]string

	enabledHelmVersions *[]string
	defaultHelmVersion  *string
)

const (
	product = "weave-flux-helm"
)

var version = "unversioned"

func init() {
	// Flags processing
	fs = pflag.NewFlagSet("default", pflag.ExitOnError)
	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "DESCRIPTION\n")
		fmt.Fprintf(os.Stderr, "  helm-operator releases Helm charts.\n")
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "FLAGS\n")
		fs.PrintDefaults()
	}

	versionFlag = fs.Bool("version", false, "print version and exit")

	logFormat = fs.String("log-format", "fmt", "change the log format.")

	kubeconfig = fs.String("kubeconfig", "", "path to a kubeconfig; required if out-of-cluster")
	master = fs.String("master", "", "address of the Kubernetes API server; overrides any value in kubeconfig; required if out-of-cluster")
	namespace = fs.String("allow-namespace", "", "if set, this limits the scope to a single namespace; if not specified, all namespaces will be watched")

	workers = fs.Int("workers", 2, "amount of workers processing releases")

	listenAddr = fs.StringP("listen", "l", ":3030", "Listen address where /metrics and API will be served")

	tillerIP = fs.String("tiller-ip", "", "Tiller IP address; required if run out-of-cluster")
	tillerPort = fs.String("tiller-port", "", "Tiller port; required if run out-of-cluster")
	tillerNamespace = fs.String("tiller-namespace", "kube-system", "Tiller namespace")

	tillerTLSVerify = fs.Bool("tiller-tls-verify", false, "verify TLS certificate from Tiller; will enable TLS communication when provided")
	tillerTLSEnable = fs.Bool("tiller-tls-enable", false, "enable TLS communication with Tiller; if provided, requires TLSKey and TLSCert to be provided as well")
	tillerTLSKey = fs.String("tiller-tls-key-path", "/etc/fluxd/helm/tls.key", "path to private key file used to communicate with the Tiller server")
	tillerTLSCert = fs.String("tiller-tls-cert-path", "/etc/fluxd/helm/tls.crt", "path to certificate file used to communicate with the Tiller server")
	tillerTLSCACert = fs.String("tiller-tls-ca-cert-path", "", "path to CA certificate file used to validate the Tiller server; required if tiller-tls-verify is enabled")
	tillerTLSHostname = fs.String("tiller-tls-hostname", "", "server name used to verify the hostname on the returned certificates from the server")

	chartsSyncInterval = fs.Duration("charts-sync-interval", 3*time.Minute, "period on which to reconcile the Helm releases with HelmRelease resources")
	statusUpdateInterval = fs.Duration("status-update-interval", 10*time.Second, "period on which to update the Helm release status in HelmRelease resources")
	logReleaseDiffs = fs.Bool("log-release-diffs", false, "log the diff when a chart release diverges; potentially insecure")
	updateDependencies = fs.Bool("update-chart-deps", true, "update chart dependencies before installing/upgrading a release")

	gitTimeout = fs.Duration("git-timeout", 20*time.Second, "duration after which git operations time out")
	gitPollInterval = fs.Duration("git-poll-interval", 5*time.Minute, "period on which to poll git chart sources for changes")
	gitDefaultRef = fs.String("git-default-ref", "master", "ref to clone chart from if ref is unspecified in a HelmRelease")

	versionedHelmRepositoryIndexes = fs.StringSlice("helm-repository-import", nil, "Targeted version and the path of the Helm repository index to import, i.e. v3:/tmp/v3/index.yaml,v2:/tmp/v2/index.yaml")

	enabledHelmVersions = fs.StringSlice("enabled-helm-versions", []string{v2.VERSION, v3.VERSION}, "Helm versions supported by this operator instance")
}

func main() {
	// explicitly initialize klog to enable stderr logging,
	// and parse our own flags.
	klog.InitFlags(nil)
	fs.Parse(os.Args)

	if *versionFlag {
		println(version)
		os.Exit(0)
	}

	// support enabling the Helm supported versions through an
	// environment variable.
	helmVersionEnv := getEnvAsSlice("HELM_VERSION", []string{})
	if len(helmVersionEnv) > 0 && !fs.Changed("enabled-helm-versions") {
		enabledHelmVersions = &helmVersionEnv
	}

	// init go-kit log
	{
		switch *logFormat {
		case "json":
			logger = log.NewJSONLogger(log.NewSyncWriter(os.Stderr))
		case "fmt":
			logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
		default:
			logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
		}
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	// error channel
	errc := make(chan error)

	// shutdown triggers
	shutdown := make(chan struct{})
	shutdownWg := &sync.WaitGroup{}

	// wait for SIGTERM
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errc <- fmt.Errorf("%s", <-c)
	}()

	mainLogger := log.With(logger, "component", "helm-operator")

	// build Kubernetes clients
	cfg, err := clientcmd.BuildConfigFromFlags(*master, *kubeconfig)
	if err != nil {
		mainLogger.Log("error", fmt.Sprintf("error building kubeconfig: %v", err))
		os.Exit(1)
	}

	kubeClient, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		mainLogger.Log("error", fmt.Sprintf("error building kubernetes clientset: %v", err))
		os.Exit(1)
	}

	ifClient, err := clientset.NewForConfig(cfg)
	if err != nil {
		mainLogger.Log("error", fmt.Sprintf("error building integrations clientset: %v", err))
		os.Exit(1)
	}

	// initialize versioned Helm clients
	helmClients := &helm.Clients{}
	for _, v := range *enabledHelmVersions {
		versionedLogger := log.With(logger, "component", "helm", "version", v)

		switch v {
		case v2.VERSION:
			helmClients.Add(v2.VERSION, v2.New(versionedLogger, kubeClient, v2.TillerOptions{
				Host:        *tillerIP,
				Port:        *tillerPort,
				Namespace:   *tillerNamespace,
				TLSVerify:   *tillerTLSVerify,
				TLSEnable:   *tillerTLSEnable,
				TLSKey:      *tillerTLSKey,
				TLSCert:     *tillerTLSCert,
				TLSCACert:   *tillerTLSCACert,
				TLSHostname: *tillerTLSHostname,
			}))
		case v3.VERSION:
			client := v3.New(versionedLogger, cfg)
			helmClients.Add(v3.VERSION, client)
		default:
			mainLogger.Log("error", fmt.Sprintf("unsupported Helm version: %s", v))
			continue
		}

		if defaultHelmVersion == nil {
			defaultVersion := v
			defaultHelmVersion = &defaultVersion
		}
	}

	// import Helm chart repositories from provided indexes
	for _, i := range *versionedHelmRepositoryIndexes {
		parts := strings.Split(i, ":")
		if len(parts) != 2 {
			mainLogger.Log("error", fmt.Sprintf("invalid version/path pair: %s, expected format is [version]:[path]", i))
			continue
		}
		v, p := parts[0], parts[1]
		client, ok := helmClients.Load(v)
		if !ok {
			mainLogger.Log("error", fmt.Sprintf("no Helm client found for version: %s", v))
			continue
		}
		if err := client.RepositoryImport(p); err != nil {
			mainLogger.Log("error", fmt.Sprintf("failed to import Helm chart repositories for %s from %s: %v", v, p, err))
		}
	}

	// setup shared informer for HelmReleases
	nsOpt := ifinformers.WithNamespace(*namespace)
	ifInformerFactory := ifinformers.NewSharedInformerFactoryWithOptions(ifClient, *chartsSyncInterval, nsOpt)
	hrInformer := ifInformerFactory.Helm().V1().HelmReleases()

	// setup workqueue for HelmReleases
	queue := workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "ChartRelease")

	gitChartSync := chartsync.NewGitChartSync(
		log.With(logger, "component", "gitchartsync"),
		kubeClient.CoreV1(),
		hrInformer.Lister(),
		chartsync.GitConfig{GitTimeout: *gitTimeout, GitPollInterval: *gitPollInterval, GitDefaultRef: *gitDefaultRef},
		queue,
	)

	rel := release.New(
		log.With(logger, "component", "release"),
		kubeClient.CoreV1(),
		ifClient.HelmV1(),
		gitChartSync,
		release.Config{LogDiffs: *logReleaseDiffs, UpdateDeps: *updateDependencies},
	)

	// prepare operator and start FluxRelease informer
	// NB: the operator needs to do its magic with the informer
	// _before_ starting it or else the cache sync seems to hang at
	// random
	opr := operator.New(log.With(logger, "component", "operator"),
		*logReleaseDiffs, kubeClient, hrInformer, queue, rel, helmClients, *defaultHelmVersion)
	go ifInformerFactory.Start(shutdown)

	// wait for the caches to be synced before starting _any_ workers
	mainLogger.Log("info", "waiting for informer caches to sync")
	if ok := cache.WaitForCacheSync(shutdown, hrInformer.Informer().HasSynced); !ok {
		mainLogger.Log("error", "failed to wait for caches to sync")
		os.Exit(1)
	}
	mainLogger.Log("info", "informer caches synced")

	// start operator
	go opr.Run(*workers, shutdown, shutdownWg)

	// start git chart sync loop
	go gitChartSync.Run(shutdown, errc, shutdownWg)

	// the status updater, to keep track of the release status for
	// every HelmRelease
	statusUpdater := status.New(ifClient, hrInformer.Lister(), helmClients, *defaultHelmVersion)
	go statusUpdater.Loop(shutdown, *statusUpdateInterval, log.With(logger, "component", "statusupdater"))

	// start HTTP server
	go daemonhttp.ListenAndServe(*listenAddr, gitChartSync, log.With(logger, "component", "daemonhttp"), shutdown)

	checkpoint.CheckForUpdates(product, version, nil, log.With(logger, "component", "checkpoint"))

	shutdownErr := <-errc
	logger.Log("exiting...", shutdownErr)
	close(shutdown)
	shutdownWg.Wait()
}

func getEnv(key string, defaultValue string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	return defaultValue
}

func getEnvAsSlice(name string, defaultValue []string) []string {
	v := getEnv(name, "")
	if v == "" {
		return defaultValue
	}
	return strings.Split(v, ",")
}
