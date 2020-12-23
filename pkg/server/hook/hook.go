package hook

import "github.com/KubeOperator/KubeOperator/pkg/logger"

var log = logger.Default

type Func func() error

type Hook interface {
	Run() error
	AddFunc(f Func)
}

var BeforeApplicationStart = NewHook("BeforeApplicationStart")

type hook struct {
	fs   []Func
	Name string
}

func NewHook(Name string) Hook {
	return &hook{
		Name: Name,
	}
}

func (h *hook) Run() error {
	log.Infof("run hook: %s", h.Name)
	for _, f := range h.fs {
		if err := f(); err != nil {
			log.Errorf("run hook func error %s", err.Error())
			return err
		}
	}
	return nil
}

func (h *hook) AddFunc(f Func) {
	h.fs = append(h.fs, f)
}
