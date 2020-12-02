package helm

import (
	"fmt"
	"io"

	"github.com/go-kit/kit/log"
)

// logWriter wraps a `log.Logger` so it can be used as an `io.Writer`.
type logWriter struct {
	log.Logger
}

// NewLogWriter returns an `io.Writer` for the given logger.
func NewLogWriter(logger log.Logger) io.Writer {
	return &logWriter{logger}
}

// Write simply logs the given byes as a 'write' operation, the only
// modification it makes before logging the given bytes is the removal
// of a terminating newline if present.
func (l *logWriter) Write(p []byte) (n int, err error) {
	origLen := len(p)
	if len(p) > 0 && p[len(p)-1] == '\n' {
		p = p[:len(p)-1] // Cut terminating newline
	}
	l.Log("info", fmt.Sprintf("%s", p))
	return origLen, nil
}
