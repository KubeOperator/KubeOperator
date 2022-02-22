package velero

import (
	"bytes"
	"errors"
	"os/exec"
)

const defaultVeleroPath = "/usr/local/bin"
const bashPath = "/bin/bash"

func Get(operate string, args []string) ([]byte, error) {
	operates := []string{"get", operate}
	args = append(operates, args...)
	args = append(args, "-o", "json")
	return ExecCommand(defaultVeleroPath, "velero", args)
}

func GetLogs(name, operate string, args []string) ([]byte, error) {
	logs := []string{operate, "logs", name}
	args = append(logs, args...)
	return ExecCommand(defaultVeleroPath, "velero", args)
}

func GetDescribe(name, operate string, args []string) ([]byte, error) {
	describes := []string{operate, "describe", name}
	args = append(describes, args...)
	return ExecCommand(defaultVeleroPath, "velero", args)
}

func Delete(name, operate string, args []string) ([]byte, error) {
	command := "echo y| /usr/local/bin/velero delete " + operate + " " + name
	del := []string{"-c", command}
	args = append(del, args...)
	return ExecCommand("", bashPath, args)
}

func Create(name, operate string, args []string) ([]byte, error) {
	backups := []string{operate, "create", name}
	args = append(backups, args...)
	return ExecCommand(defaultVeleroPath, "velero", args)
}

func Restore(backupName string, args []string) ([]byte, error) {
	backups := []string{"restore", "create", "--from-backup", backupName}
	args = append(backups, args...)

	return ExecCommand(defaultVeleroPath, "velero", args)
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
		return []byte{}, errors.New(buffer.String())
	}
	return buffer.Bytes(), nil
}
