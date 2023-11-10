package log

import (
	"fmt"
	"io"
	"strings"

	"github.com/liamg/tml"
)

type Level int

const (
	LevelTrace Level = iota
	LevelDebug
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal
)

func ParseLevel(name string) Level {
	switch strings.ToLower(name) {
	case "trace":
		return LevelTrace
	case "debug":
		return LevelDebug
	case "info":
		return LevelInfo
	case "warn":
		return LevelWarn
	case "error":
		return LevelError
	case "fatal":
		return LevelFatal
	}
	return LevelInfo
}

func (l Level) String() string {
	switch l {
	case LevelTrace:
		return "<dim>TRC</dim>"
	case LevelDebug:
		return "DBG"
	case LevelInfo:
		return "<blue>INF</blue>"
	case LevelWarn:
		return "<yellow>WRN</yellow>"
	case LevelError:
		return "<red>ERR</red>"
	case LevelFatal:
		return "<red><bold>FTL</bold></red>"
	}
	return "<dim>UNK</dim>"
}

type Logger struct {
	prefix string
	w      io.Writer
	parent *Logger
	level  Level
}

func New(w io.Writer) *Logger {
	return &Logger{w: w}
}

func (l *Logger) SetLevel(level Level) {
	l.level = level
	if l.parent != nil {
		l.parent.SetLevel(level)
	}
}

func (l *Logger) WithPrefix(prefix string) *Logger {
	return &Logger{prefix: prefix, parent: l}
}

func (l *Logger) Printf(level Level, format string, v ...interface{}) {
	if level < l.level {
		return
	}
	line := strings.TrimSpace(tml.Sprintf(format, v...))
	if l.parent != nil {
		l.parent.Printf(level, "<dim>[</dim>%s<dim>]</dim> %s", l.prefix, line)
		return
	}
	_, _ = fmt.Fprintf(l.w, "%s %s\n", tml.Sprintf(level.String()), line)
}

func (l *Logger) Tracef(format string, v ...interface{}) {
	l.Printf(LevelTrace, format, v...)
}

func (l *Logger) Debugf(format string, v ...interface{}) {
	l.Printf(LevelDebug, format, v...)
}

func (l *Logger) Infof(format string, v ...interface{}) {
	l.Printf(LevelInfo, format, v...)
}

func (l *Logger) Warnf(format string, v ...interface{}) {
	l.Printf(LevelWarn, format, v...)
}

func (l *Logger) Errorf(format string, v ...interface{}) {
	l.Printf(LevelError, format, v...)
}

func (l *Logger) Fatalf(format string, v ...interface{}) {
	l.Printf(LevelFatal, format, v...)
}

func (l *Logger) Print(message string) {
	l.Printf(LevelDebug, "%s", message)
}

func (l *Logger) Trace(message string) {
	l.Printf(LevelDebug, "%s", message)
}

func (l *Logger) Debug(message string) {
	l.Printf(LevelDebug, "%s", message)
}

func (l *Logger) Info(message string) {
	l.Printf(LevelInfo, "%s", message)
}

func (l *Logger) Warning(message string) {
	l.Printf(LevelWarn, "%s", message)
}

func (l *Logger) Error(message string) {
	l.Printf(LevelError, "%s", message)
}

func (l *Logger) Fatal(message string) {
	l.Printf(LevelFatal, "%s", message)
}
