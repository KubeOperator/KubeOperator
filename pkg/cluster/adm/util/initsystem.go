package util

import (
	"fmt"
	"ko3-gin/pkg/util/ssh"
	"strings"
)

type InitSystem interface {
	EnableCommand(service string) string
	ServiceStart(service string) error
	ServiceStop(service string) error
	ServiceRestart(service string) error
	ServiceExists(service string) bool
	ServiceIsEnabled(service string) bool
	ServiceIsActive(service string) bool
}



// SystemdInitSystem defines systemd
type SystemdInitSystem struct {
	sshClient *ssh.SSH
}

// EnableCommand return a string describing how to enable a service
func (sysd SystemdInitSystem) EnableCommand(service string) string {
	return fmt.Sprintf("systemctl enable %s.service", service)
}

// reloadSystemd reloeads the systemd daemon
func (sysd SystemdInitSystem) reloadSystemd() error {
	if err := sysd.sshClient.Run("systemctl", "daemon-reload"); err != nil {
		return fmt.Errorf("failed to reload systemd: %v", err)
	}
	return nil
}

// ServiceStart tries to start a specific service
func (sysd SystemdInitSystem) ServiceStart(service string) error {
	// Before we try to start any service, make sure that systemd is ready
	if err := sysd.reloadSystemd(); err != nil {
		return err
	}
	args := []string{"systemctl", "start", service}
	return sysd.sshClient.Run(args...)
}

// ServiceRestart tries to reload the environment and restart the specific service
func (sysd SystemdInitSystem) ServiceRestart(service string) error {
	// Before we try to restart any service, make sure that systemd is ready
	if err := sysd.reloadSystemd(); err != nil {
		return err
	}
	args := []string{"systemctl", "restart", service}
	return sysd.sshClient.Run(args...)
}

// ServiceStop tries to stop a specific service
func (sysd SystemdInitSystem) ServiceStop(service string) error {
	args := []string{"systemctl", "stop", service}
	return sysd.sshClient.Run(args...)
}

// ServiceExists ensures the service is defined for this init system.
func (sysd SystemdInitSystem) ServiceExists(service string) bool {
	args := []string{"systemctl", "status", service}
	outBytes, _ := sysd.sshClient.CombinedOutput(args...)
	output := string(outBytes)
	return !strings.Contains(output, "Loaded: not-found")
}

// ServiceIsEnabled ensures the service is enabled to start on each boot.
func (sysd SystemdInitSystem) ServiceIsEnabled(service string) bool {
	args := []string{"systemctl", "is-enabled", service}
	err := sysd.sshClient.Run(args...)
	return err == nil
}

// ServiceIsActive will check is the service is "active". In the case of
// crash looping services (kubelet in our case) status will return as
// "activating", so we will consider this active as well.
func (sysd SystemdInitSystem) ServiceIsActive(service string) bool {
	args := []string{"systemctl", "is-active", service}
	// Ignoring error here, command returns non-0 if in "activating" status:
	outBytes, _ := sysd.sshClient.CombinedOutput(args...)
	output := strings.TrimSpace(string(outBytes))
	if output == "active" || output == "activating" {
		return true
	}
	return false
}

// GetInitSystem returns an InitSystem for the current system, or nil
// if we cannot detect a supported init system.
// This indicates we will skip init system checks, not an error.
func GetInitSystem(ssh *ssh.SSH) (InitSystem, error) {
	_, err := ssh.LookPath("systemctl")
	if err != nil {
		return nil, fmt.Errorf("no supported init system detected, skipping checking for services")
	}
	return &SystemdInitSystem{
		sshClient: ssh,
	}, nil
}
