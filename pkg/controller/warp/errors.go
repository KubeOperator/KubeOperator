package warp

import "encoding/json"

type ControllerError struct {
	err error
}

func NewControllerError(err error) error {
	return &ControllerError{err: err}
}

func (e *ControllerError) Error() string {
	warp := struct {
		Msg string `json:"msg"`
	}{e.err.Error(),}
	b, _ := json.Marshal(&warp)
	return string(b)
}

func (e *ControllerError) Unwrap() error {
	return e.err
}
