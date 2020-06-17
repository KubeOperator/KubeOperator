package addon

const (
	addonEtcd = "06-etcd.yml"
)

//type EFKInstallPhase struct {
//	EtcdDataDir string
//}
//
//func (s EtcdPhase) Name() string {
//	return "InitEtcd"
//}
//
//func (s EtcdPhase) Run(b kobe.Interface) error {
//	if s.EtcdDataDir != "" {
//		b.SetVar(facts.EtcdDataDirFactName, s.EtcdDataDir)
//	}
//	return phases.RunPlaybookAndGetResult(b, initEtcd)
//}
