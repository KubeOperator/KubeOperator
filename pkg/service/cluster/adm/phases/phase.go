package phases

import (
	"github.com/KubeOperator/KubeOperator/pkg/util/kobe"
	"time"
)

const (
	PhaseInterval = 2 * time.Second
	PhaseTimeout  = 10 * time.Minute
)

type Interface interface {
	Name() string
	Run(p kobe.Interface) error
}
