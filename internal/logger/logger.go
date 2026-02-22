package logger

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/fatih/color"
)

type Level int

const (
	LevelFatal Level = iota
	LevelError
	LevelWarn
	LevelInfo
	LevelDebug
)

func (l Level) String() string {
	return [...]string{"FATAL", "ERROR", "WARN", "INFO", "DEBUG"}[l]
}

type Logger struct {
	level  Level
	out    io.Writer
	errOut io.Writer
}

type Options func(*Logger)

func WithOutput(w io.Writer) Options {
	return func(l *Logger) {
		l.out = w
	}
}

func WithErrOutput(w io.Writer) Options {
	return func(l *Logger) {
		l.errOut = w
	}
}

func New(level Level, opts ...Options) *Logger {
	l := &Logger{ // значения по умолчанию
		level:  level,
		out:    os.Stdout,
		errOut: os.Stderr,
	}

	for _, opt := range opts { // применяем все опции
		opt(l)
	}

	return l
}

func (l *Logger) Log(level Level, format string, args ...interface{}) {
	if level > l.level {
		return
	}

	timestamp := time.Now().Format("2006/01/02 15:04:05.000")

	message := fmt.Sprintf(format, args...)

	var col *color.Color

	switch level {
	case LevelDebug:
		col = color.New(color.FgGreen)
	case LevelInfo:
		col = color.New(color.FgYellow)
	case LevelWarn:
		col = color.RGB(255, 165, 0) //оранжевый
	case LevelError:
		col = color.New(color.FgRed)
	case LevelFatal:
		col = color.New(color.FgRed)
	}

	levelLine := col.SprintfFunc()("[%s]", level.String())

	log := fmt.Sprintf("%s %s %s", timestamp, levelLine, message)

	var out io.Writer
	if level >= LevelError {
		out = l.errOut
	} else {
		out = l.out
	}

	fmt.Fprintln(out, log)
}

func (l *Logger) Debug(format string, args ...interface{}) {
	l.Log(LevelDebug, format, args...)
}
func (l *Logger) Info(format string, args ...interface{}) {
	l.Log(LevelInfo, format, args...)
}
func (l *Logger) Warn(format string, args ...interface{}) {
	l.Log(LevelWarn, format, args...)
}
func (l *Logger) Error(format string, args ...interface{}) {
	l.Log(LevelError, format, args...)
}
func (l *Logger) Fatal(format string, args ...interface{}) {
	l.Log(LevelFatal, format, args...)
	os.Exit(1)
}
