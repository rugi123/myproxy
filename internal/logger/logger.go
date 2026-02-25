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

type Entry struct {
	Level   Level
	Message string
	Args    []interface{}
	Time    time.Time
}

type Logger struct {
	level Level
	out   io.Writer
	ch    chan Entry
}

func New(level Level, out io.Writer, ch chan Entry) *Logger {
	return &Logger{
		level: level,
		out:   out,
		ch:    ch,
	}
}

func (l *Logger) log(level Level, format string, args ...interface{}) {
	if l.level < level {
		return
	}
	l.ch <- Entry{
		Level:   level,
		Message: format,
		Args:    args,
		Time:    time.Now(),
	}
}

func (l *Logger) Debug(format string, args ...interface{}) {
	l.log(LevelDebug, format, args...)
}
func (l *Logger) Info(format string, args ...interface{}) {
	l.log(LevelInfo, format, args...)
}
func (l *Logger) Warn(format string, args ...interface{}) {
	l.log(LevelWarn, format, args...)
}
func (l *Logger) Error(format string, args ...interface{}) {
	l.log(LevelError, format, args...)
}
func (l *Logger) Fatal(format string, args ...interface{}) {
	l.log(LevelFatal, format, args...)
	os.Exit(1)
}

func (l *Logger) Close() {
	close(l.ch)
}

func (l *Logger) Run() {
	var entry Entry

	for {
		entry = <-l.ch

		var col func(...interface{}) string
		switch entry.Level {
		case LevelDebug:
			col = color.New(color.FgGreen).SprintFunc()
		case LevelInfo:
			col = color.New(color.FgYellow).SprintFunc()
		case LevelWarn:
			col = color.RGB(255, 127, 0).SprintFunc()
		case LevelError:
			col = color.New(color.FgRed).SprintFunc()
		case LevelFatal:
			col = color.New(color.FgRed).SprintFunc()
		}

		message := fmt.Sprintf(entry.Message, entry.Args...)
		timestamp := entry.Time.Format("2006/01/02 15:04:05.000")
		prefix := fmt.Sprintf("[%s]", entry.Level.String())
		log := fmt.Sprintf("%s %s %s", timestamp, col(prefix), message)

		fmt.Fprintln(l.out, log)
	}
}
