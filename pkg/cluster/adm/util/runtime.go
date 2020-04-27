package util

import (
	"github.com/pkg/errors"
	errorsutil "k8s.io/apimachinery/pkg/util/errors"
	"ko3-gin/pkg/cluster/adm/constants"
	"ko3-gin/pkg/util/ssh"
	"strings"
)

type ContainerRuntime interface {
	IsDocker() bool
	IsRunning() error
	ListKubeContainers() ([]string, error)
	RemoveContainers(containers []string) error
	PullImage(image string) error
	ImageExists(image string) (bool, error)
}

type DockerRuntime struct {
	sshClient *ssh.SSH
}

func NewContainerRuntime(ssh *ssh.SSH) (ContainerRuntime, error) {
	toolName := "docker"
	runtime := &DockerRuntime{ssh}
	if _, err := ssh.LookPath(toolName); err != nil {
		return nil, errors.Wrapf(err, "%s is required for container runtime", toolName)
	}

	return runtime, nil
}

func (runtime *DockerRuntime) IsDocker() bool {
	return true
}

func (runtime *DockerRuntime) IsRunning() error {
	if out, err := runtime.sshClient.CombinedOutput("docker info"); err != nil {
		return errors.Wrapf(err, "container runtime is not running: output: %s, error", string(out))
	}
	return nil
}

func (runtime *DockerRuntime) ListKubeContainers() ([]string, error) {
	output, err := runtime.sshClient.CombinedOutput("docker", "ps", "-a", "--filter", "name=k8s_", "-q")
	return strings.Fields(string(output)), err
}

// RemoveContainers removes running containers
func (runtime *DockerRuntime) RemoveContainers(containers []string) error {
	errs := []error{}
	for _, container := range containers {
		out, err := runtime.sshClient.CombinedOutput("docker", "rm", "--force", "--volumes", container)
		if err != nil {
			// don't stop on errors, try to remove as many containers as possible
			errs = append(errs, errors.Wrapf(err, "failed to remove running container %s: output: %s, error", container, string(out)))
		}
	}
	return errorsutil.NewAggregate(errs)
}

// PullImage pulls the image
func (runtime *DockerRuntime) PullImage(image string) error {
	var err error
	var out []byte
	for i := 0; i < constants.PullImageRetry; i++ {
		out, err = runtime.sshClient.CombinedOutput("docker", "pull", image)
		if err == nil {
			return nil
		}
	}
	return errors.Wrapf(err, "output: %s, error", out)
}

// ImageExists checks to see if the image exists on the system
func (runtime *DockerRuntime) ImageExists(image string) (bool, error) {
	_, err := runtime.sshClient.CombinedOutput("docker", "inspect", image)
	return err == nil, nil
}
