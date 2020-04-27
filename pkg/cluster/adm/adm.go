package adm

type Interface interface {
	Upgrade()
	Install()
	Reset()
	JoinMaster()
	JoinWorker()
}
