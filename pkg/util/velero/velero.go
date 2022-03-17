package velero

import (
	"bytes"
	"context"
	"errors"
	"os/exec"
	"time"
)

const defaultVeleroPath = "/usr/local/bin/velero"

func Get(operate string, args []string) ([]byte, error) {
	operates := []string{"get", operate}
	args = append(operates, args...)
	args = append(args, "-o", "json")
	return ExecCommand(defaultVeleroPath, args)
}

func GetLogs(name, operate string, args []string) ([]byte, error) {
	logs := []string{operate, "logs", name}
	args = append(logs, args...)
	return ExecCommand(defaultVeleroPath, args)
}

func GetDescribe(name, operate string, args []string) ([]byte, error) {
	describes := []string{operate, "describe", name}
	args = append(describes, args...)
	return ExecCommand(defaultVeleroPath, args)
}

func Delete(name, operate string, args []string) ([]byte, error) {
	command := "echo y| /usr/local/bin/velero delete " + operate + " " + name
	del := []string{"-c", command}
	args = append(del, args...)
	return ExecCommand("/bin/bash", args)
}

func Create(name, operate string, args []string) ([]byte, error) {
	backups := []string{operate, "create", name}
	args = append(backups, args...)
	return ExecCommand(defaultVeleroPath, args)
}

func Restore(backupName string, args []string) ([]byte, error) {
	backups := []string{"restore", "create", "--from-backup", backupName}
	args = append(backups, args...)
	return ExecCommand(defaultVeleroPath, args)
}

func Install(args []string) ([]byte, error) {
	install := []string{"install"}
	args = append(install, args...)
	return ExecCommand(defaultVeleroPath, args)
}

func ExecCommand(command string, args []string) ([]byte, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	cmd := exec.CommandContext(ctx, command, args...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return []byte{}, err
	}
	cmd.Stderr = cmd.Stdout
	if err = cmd.Start(); err != nil {
		return []byte{}, err
	}

	var buffer bytes.Buffer
	done := make(chan bool, 1)
	go func() {
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
		done <- true
	}()

	select {
	case <-done:
		if err = cmd.Wait(); err != nil {
			return []byte{}, errors.New(buffer.String())
		}
		return buffer.Bytes(), nil
	case <-time.After(time.Second * 20):
		_ = stdout.Close()
		return []byte("time out"), errors.New("read log time out")
	}
}
