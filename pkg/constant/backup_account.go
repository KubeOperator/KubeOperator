package constant

import "path"

const (
	Valid           = "VALID"
	DisConnect      = "DISCONNECT"
	VerifyFailed    = "VERIFYFAILED"
	Azure           = "AZURE"
	S3              = "S3"
	OSS             = "OSS"
	DefaultFireName = "./ko-backup-test.json"
)

var (
	DefaultDataDir    = "/var/ko/data"
	DefaultBackupDir  = path.Join(DefaultDataDir, "backup")
	DefaultRestoreDir = path.Join(DefaultDataDir, "restore")
)
