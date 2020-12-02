package chartsync

// ChartUnavailableError is returned when the requested chart is
// unavailable, and the reason is known and finite.
type ChartUnavailableError struct {
	Err error
}

func (err ChartUnavailableError) Unwrap() error {
	return err.Err
}

func (err ChartUnavailableError) Error() string {
	return "chart unavailable: " + err.Err.Error()
}

// ChartNotReadyError is returned when the requested chart is
// unavailable at the moment, but may become at available a later stage
// without any interference from a human.
type ChartNotReadyError struct {
	Err error
}

func (err ChartNotReadyError) Unwrap() error {
	return err.Err
}

func (err ChartNotReadyError) Error() string {
	return "chart not ready: " + err.Err.Error()
}

// GitAuthError presents a error that has occured when handling
// the git auth details
type GitAuthError struct {
	Err error
}

func (err GitAuthError) Unwrap() error {
	return err.Err
}

func (err GitAuthError) Error() string {
	return "git auth error: " + err.Err.Error()
}
