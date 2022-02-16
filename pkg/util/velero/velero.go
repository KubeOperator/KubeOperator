package velero

import (
	"bytes"
	"os/exec"
	"time"
)

const defaultVeleroPath = "/usr/local/bin"

func Backup(cluster string, args []string) ([]byte, error) {

	now := time.Now()
	day := now.Format("2006-01-02-15-04")
	backupName := cluster + "-" + day
	backups := []string{"backup", "create", backupName}
	args = append(backups, args...)

	result, err := ExecCommand(defaultVeleroPath, "velero", args)
	if err != nil {
		return result, err
	}
	return result, nil
}

func GetBackups(args []string) ([]byte, error) {
	backups := []string{"get", "backups"}
	args = append(backups, args...)
	args = append(args, "-o", "json")

	result, err := ExecCommand(defaultVeleroPath, "velero", args)
	if err != nil {
		return result, err
	}
	return result, nil
}

func GetBackupDescribe(backupName string, args []string) ([]byte, error) {
	describes := []string{"backup", "describe", backupName}
	args = append(describes, args...)

	result, err := ExecCommand(defaultVeleroPath, "velero", args)
	if err != nil {
		return result, err
	}
	return result, nil
}

func GetBackupLogs(backupName string, args []string) ([]byte, error) {
	logs := []string{"backup", "logs", backupName}
	args = append(logs, args...)

	result, err := ExecCommand(defaultVeleroPath, "velero", args)
	if err != nil {
		return result, err
	}
	return result, nil
}

func Restore(backupName string, args []string) ([]byte, error) {
	backups := []string{"restore", "create", "--from-backup", backupName}
	args = append(backups, args...)

	result, err := ExecCommand(defaultVeleroPath, "velero", args)
	if err != nil {
		return result, err
	}
	return result, nil
}

func GetRestores(args []string) ([]byte, error) {
	backups := []string{"get", "backups"}
	args = append(backups, args...)

	result, err := ExecCommand(defaultVeleroPath, "velero", args)
	if err != nil {
		return result, err
	}
	return result, nil
}

func GetRestoreDescribe(restoreName string, args []string) ([]byte, error) {
	describes := []string{"restore", "describe", restoreName}
	args = append(describes, args...)

	result, err := ExecCommand(defaultVeleroPath, "velero", args)
	if err != nil {
		return result, err
	}
	return result, nil
}

func ExecCommand(path string, command string, args []string) ([]byte, error) {
	cmd := exec.Command(command, args...)
	cmd.Dir = path
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return []byte{}, err
	}
	cmd.Stderr = cmd.Stdout
	if err = cmd.Start(); err != nil {
		return []byte{}, err
	}

	var buffer bytes.Buffer
	for {
		out := make([]byte, 1024)
		length, err := stdout.Read(out)
		if err != nil {
			break
		}
		if length > 0 {
			buffer.Write(out[:length])
		}
	}

	if err = cmd.Wait(); err != nil {
		return buffer.Bytes(), err
	}
	return buffer.Bytes(), nil
}
