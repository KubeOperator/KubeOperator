package host

import (
	"ko3-gin/pkg/util/ssh"
)

type Host struct {
	Client *ssh.SSH
}

