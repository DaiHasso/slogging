package slogging

import (
	"log"
)

// Logger is the base interface for all loggers.
type Logger interface {
	Debug(string) LogInstance
	Warn(string) LogInstance
	Error(string) LogInstance
	Info(string) LogInstance
	Log(LogLevel, []byte)
	SetPretty(bool)
	SetInternalLogger(*log.Logger)
	GetInternalLogger() *log.Logger
	GetPseudoWriter(LogLevel) PseudoWriter
	GetStdLogger(LogLevel) *log.Logger
	SetLogLevel(string) error
}
