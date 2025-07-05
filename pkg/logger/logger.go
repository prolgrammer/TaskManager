package logger

import (
	"github.com/rs/zerolog"
)

func NewConsoleLogger(level int) Logger {
	return newConsoleZerolog(level)
}

type Logger interface {
	Debug() LogContext
	Info() LogContext
	Warn() LogContext
	Error() LogContext
	Fatal() LogContext
	Err(err error) LogContext
}

type LogContext interface {
	Msg(message string)
	Msgf(format string, args ...interface{})
	Err(err error) LogContext
	Error() LogContext
}

func levelSwitch(level int) zerolog.Level {
	switch level {
	case LevelInfo:
		return zerolog.InfoLevel
	case LevelWarn:
		return zerolog.WarnLevel
	case LevelError:
		return zerolog.ErrorLevel
	case LevelFatal:
		return zerolog.FatalLevel
	case LevelDebug:
		return zerolog.DebugLevel
	}
	return zerolog.InfoLevel
}

const (
	LevelInfo = iota
	LevelDebug
	LevelWarn
	LevelError
	LevelFatal
)

func LevelSwitch(level string) int {
	switch level {
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
