package logger

import (
	"log/slog"
	"os"
	"path"
)

var (
	// Logger is the shared logger instance accessible by other packages
	Logger *slog.Logger
)

// Init initializes the Logger with the specified log level and output destination.
func Init(logLevel string, output *os.File) {
	Logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.SourceKey {
				s := a.Value.Any().(*slog.Source)
				s.File = path.Base(s.File)
			}
			return a
		},
	}))
}
