package app

import (
	"context"
	"fmt"
	"log"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// appLogger adapts Wails runtime logging (and stdlog fallback) for app + Badger use.
type appLogger struct {
	ctx  context.Context
	name string
}

func newAppLogger(ctx context.Context, name string) *appLogger {
	return &appLogger{ctx: ctx, name: name}
}

func (l *appLogger) setContext(ctx context.Context) {
	l.ctx = ctx
}

func (l *appLogger) prefix(format string) string {
	if l.name == "" {
		return format
	}
	return "[" + l.name + "] " + format
}

func (l *appLogger) Errorf(format string, args ...interface{}) {
	format = l.prefix(format)
	if l.ctx != nil {
		runtime.LogErrorf(l.ctx, format, args...)
		return
	}
	log.Printf("ERROR: "+format, args...)
}

func (l *appLogger) Warnf(format string, args ...interface{}) {
	format = l.prefix(format)
	if l.ctx != nil {
		runtime.LogWarningf(l.ctx, format, args...)
		return
	}
	log.Printf("WARN: "+format, args...)
}

func (l *appLogger) Warningf(format string, args ...interface{}) {
	l.Warnf(format, args...)
}

func (l *appLogger) Infof(format string, args ...interface{}) {
	format = l.prefix(format)
	if l.ctx != nil {
		runtime.LogInfof(l.ctx, format, args...)
		return
	}
	log.Printf("INFO: "+format, args...)
}

func (l *appLogger) Info(message string) {
	l.Infof("%s", message)
}

func (l *appLogger) Debugf(format string, args ...interface{}) {
	format = l.prefix(format)
	if l.ctx != nil {
		runtime.LogDebugf(l.ctx, format, args...)
		return
	}
	log.Printf("DEBUG: "+format, args...)
}

func (l *appLogger) Debug(message string) {
	l.Debugf("%s", message)
}

func (l *appLogger) String() string {
	return fmt.Sprintf("appLogger(%s)", l.name)
}
