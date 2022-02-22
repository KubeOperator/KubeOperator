package velero

import (
	"bytes"
	"os/exec"
)

const defaultVeleroPath = "/usr/local/bin"

func Backup(backupName string, args []string) ([]byte, error) {
	backups := []string{"backup", "create", backupName}
	args = append(backups, args...)

	return ExecCommand(defaultVeleroPath, "velero", args)
}

func GetBackups(args []string) ([]byte, error) {
	backups := []string{"get", "backups"}
	args = append(backups, args...)
	args = append(args, "-o", "json")

	return ExecCommand(defaultVeleroPath, "velero", args)
}

func GetBackupDescribe(backupName string, args []string) ([]byte, error) {
	describes := []string{"backup", "describe", backupName}
	args = append(describes, args...)

	return ExecCommand(defaultVeleroPath, "velero", args)
}

func GetBackupLogs(backupName string, args []string) ([]byte, error) {
	logs := []string{"backup", "logs", backupName}
	args = append(logs, args...)

	return ExecCommand(defaultVeleroPath, "velero", args)
}

func DeleteBackup(backupName string) ([]byte, error) {
	del := []string{"backup", "delete", backupName}
	return ExecCommand(defaultVeleroPath, "velero", del)
}

func Restore(backupName string, args []string) ([]byte, error) {
	backups := []string{"restore", "create", "--from-backup", backupName}
	args = append(backups, args...)

	return ExecCommand(defaultVeleroPath, "velero", args)
}

func GetRestores(args []string) ([]byte, error) {
	backups := []string{"get", "restores"}
	args = append(backups, args...)
	args = append(args, "-o", "json")

	return ExecCommand(defaultVeleroPath, "velero", args)
}

func GetRestoreDescribe(restoreName string, args []string) ([]byte, error) {
	describes := []string{"restore", "describe", restoreName}
	args = append(describes, args...)

	return ExecCommand(defaultVeleroPath, "velero", args)
}

func GetRestoreLogs(backupName string, args []string) ([]byte, error) {
	logs := []string{"restore", "logs", backupName}
	args = append(logs, args...)

	return ExecCommand(defaultVeleroPath, "velero", args)
}

func Schedule(scheduleName string, args []string) ([]byte, error) {
	backups := []string{"schedule", "create", scheduleName}
	args = append(backups, args...)

	return ExecCommand(defaultVeleroPath, "velero", args)
}

func GetSchedules(args []string) ([]byte, error) {
	backups := []string{"get", "schedule"}
	args = append(backups, args...)
	args = append(args, "-o", "json")

	return ExecCommand(defaultVeleroPath, "velero", args)
}

func GetScheduleDescribe(scheduleName string, args []string) ([]byte, error) {
	describes := []string{"schedule", "describe", scheduleName}
	args = append(describes, args...)

	return ExecCommand(defaultVeleroPath, "velero", args)
}

func DeleteSchedule(scheduleName string) ([]byte, error) {
	del := []string{"schedule", "delete", scheduleName}
	return ExecCommand(defaultVeleroPath, "velero", del)
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
