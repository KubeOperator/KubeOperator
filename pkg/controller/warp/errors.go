package warp

import (
	"encoding/json"
	"github.com/KubeOperator/KubeOperator/pkg/logger"
)

var log = logger.Default

type ControllerError struct {
	err error
}

func NewControllerError(err error) error {
	return &ControllerError{err: err}
}

func (e *ControllerError) Error() string {
	warp := struct {
		Msg string `json:"msg"`
	}{e.err.Error()}
	b, err := json.Marshal(&warp)
	if err != nil {
		log.Errorf("json marshal failed, %v", warp)
	}
	return string(b)
}

func (e *ControllerError) Unwrap() error {
	return e.err
}
