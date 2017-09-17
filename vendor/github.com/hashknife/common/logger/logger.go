package logger

import (
	"io"
	stdlog "log"
	"os"

	"github.com/go-kit/kit/log"
)

// SetupJSONLogger sets up a new go-kit JSON logger with some useful
// values to share across each log line Redirects the stdlib logging
// output to go through this logger
func SetupJSONLogger(serviceName string) log.Logger {
	logger := log.NewJSONLogger(os.Stdout)
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)
	logger = log.With(logger, "service", serviceName)
	stdlog.SetOutput(NewStdlibAdapter(logger)) // redirect stdlib logging to us
	stdlog.SetFlags(0)
	return logger
}

// StdlibAdapter is a copy of Go_kits logger
type StdlibAdapter struct {
	log.Logger
}

// NewStdlibAdapter returns a new StdlibAdapter wrapper around the passed
// logger. It's designed to be passed to log.SetOutput
func NewStdlibAdapter(logger log.Logger) io.Writer {
	return StdlibAdapter{
		Logger: logger,
	}
}

// Write
func (a StdlibAdapter) Write(p []byte) (int, error) {
	if err := a.Logger.Log("msg", string(p)); err != nil {
		return 0, err
	}
	return len(p), nil
}
