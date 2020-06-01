package constant

var (
	PlaybookBase                   = "01-base.yml"
	PlaybookContainerd             = "02-containerd.yml"
	PlaybookDocker                 = "02-runtime.yml"
	PlaybookKubernetesComponent    = "03-kubernetes-component.yml"
	PlaybookLoadBalancer           = "04-load-balancer.yml"
	PlaybookETCD                   = "05-etcd.yml"
	PlaybookKubernetesCertificates = "06-kubernetes-certificates.yml"
	PlaybookKubernetesMaster       = "07-kubernetes-master.yml"
	PlaybookKubernetesWorker       = "08-kubernetes-worker.yml"
	PlaybookPost                   = "09-post.yml"
	PlaybookNetworkPlugin          = "10-network-plugin.yml"
	PlaybookInitCluster            = "90-initial-cluster.yml"
	PlaybookAddWorker              = "91-add-worker.yml"
	PlaybookUpgradeCluster         = "92-upgrade-cluster.yml"
	PlaybookCertificatesRenew      = "93-certificates-renew.yml"
	PlaybookBackupCluster          = "94-backup-cluster.yml"
	PlaybookRestoreCluster         = "95-restore-cluster.yml"
	PlaybookResetCluster           = "99-reset-cluster.yml"
)
