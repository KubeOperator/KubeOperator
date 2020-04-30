package constants

const (
	KubernetesDir         = "/etc/kubernetes/"
	KubeletPodManifestDir = KubernetesDir + "manifests/"
	TokenFile             = KubernetesDir + "known_tokens.csv"
	KubectlConfigFile     = "/root/.kube/config"

	CertificatesDir             = KubernetesDir + "pki/"
	CACertName                  = CertificatesDir + "ca.crt"
	CAKeyName                   = CertificatesDir + "ca.key"
	EtcdCACertName              = CertificatesDir + "etcd/ca.crt"
	EtcdCAKeyName               = CertificatesDir + "etcd/ca.key"
	EtcdListenClientPort        = 2379
	EtcdListenPeerPort          = 2380
	APIServerEtcdClientCertName = CertificatesDir + "apiserver-etcd-client.crt"
	APIServerEtcdClientKeyName  = CertificatesDir + "apiserver-etcd-client.key"

	EtcdVersion    = "3.4.3-0"
	CoreDNSVersion = "1.6.7"
	PauseVersion   = "3.2"

	BinDir = "/usr/bin/"
	TmpDir = "/tmp/ko/"

	ResourceDir = "resource/"
	ConfDir     = ResourceDir + "conf/"
	SrcDir      = ResourceDir + "static/"

	DefaultSystemdUnitFilePath = "/etc/systemd/system"
)
