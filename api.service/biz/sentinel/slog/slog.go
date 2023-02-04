package slog

import (
	"github.com/cloudwego/hertz/pkg/common/hlog"
)

type Logger struct {
}

func (*Logger) Debug(msg string, keysAndValues ...interface{}) {
	hlog.Debugf(msg, keysAndValues)
}
func (*Logger) DebugEnabled() bool {
	return true
}

// Info logs a non-error message with the given key/value pairs as context.
//
// The msg argument should be used to add some constant description to
// the log line.  The key/value pairs can then be used to add additional
// variable information.  The key/value pairs should alternate string
// keys and arbitrary values.
func (*Logger) Info(msg string, keysAndValues ...interface{}) {
	hlog.Infof(msg, keysAndValues)
}
func (*Logger) InfoEnabled() bool {
	return true
}
func (*Logger) Warn(msg string, keysAndValues ...interface{}) {
	hlog.Warnf(msg, keysAndValues)
}
func (*Logger) WarnEnabled() bool {
	return true
}

func (*Logger) Error(err error, msg string, keysAndValues ...interface{}) {
	hlog.Errorf(err.Error()+msg, keysAndValues)
}
func (*Logger) ErrorEnabled() bool {
	return true
}
