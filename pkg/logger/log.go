package logger

import (
	"io"
	"os"
	"strings"

	kitlog "github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

type Level int

const (
	DebugLevel Level = iota
	InfoLevel
	WarnLevel
	ErrorLevel
)

var (
	mlog mLog
)

type Mlogger interface {
	Log(keyvals ...interface{}) error
	Print(keyvals ...interface{})
	Info(keyvals ...interface{})
	Error(keyvals ...interface{})
	Debug(keyvals ...interface{})
	Warn(keyvals ...interface{})
	WithPrefix(keyvals ...interface{}) Mlogger
}

type mLog struct {
	lg     kitlog.Logger
	allowd Level
}

func getLevel(lev string) Level {
	var allowd Level
	switch strings.ToLower(lev) {
	case "info":
		allowd = InfoLevel
	case "debug":
		allowd = DebugLevel
	case "warn":
		allowd = WarnLevel
	case "error":
		allowd = ErrorLevel
	default:
		allowd = InfoLevel
	}
	return allowd
}

func NewStrLogger(out io.Writer, lev string) Mlogger {
	var l mLog
	l.lg = kitlog.With(kitlog.NewLogfmtLogger(out), "ts", kitlog.DefaultTimestampUTC)
	l.allowd = getLevel(lev)
	return &l
}

func NewJsonLogger(out io.Writer, lev string) Mlogger {
	var l mLog
	l.lg = kitlog.With(kitlog.NewJSONLogger(out), "ts", kitlog.DefaultTimestampUTC)
	l.allowd = getLevel(lev)
	return &l
}

func (l *mLog) Print(keyvals ...interface{}) {
	l.lg.Log(keyvals...)
}

func (l *mLog) Log(keyvals ...interface{}) error {
	return l.lg.Log(keyvals...)
}

func (l *mLog) Info(keyvals ...interface{}) {
	if l.allowd <= InfoLevel {
		level.Info(l.lg).Log(keyvals...)
	}
}

func (l *mLog) Error(keyvals ...interface{}) {
	if l.allowd <= ErrorLevel {
		level.Error(l.lg).Log(keyvals...)
	}
}
func (l *mLog) Debug(keyvals ...interface{}) {
	if l.allowd <= DebugLevel {
		level.Debug(l.lg).Log(keyvals...)
	}
}

func (l *mLog) Warn(keyvals ...interface{}) {
	if l.allowd <= WarnLevel {
		level.Warn(l.lg).Log(keyvals...)
	}
}

func (l *mLog) WithPrefix(keyvals ...interface{}) Mlogger {
	l.lg = kitlog.With(l.lg, keyvals...)
	return l
}

func DefaultLog() Mlogger {
	return NewJsonLogger(os.Stdout, "info")
}
