package config

import (
	"errors"

	"github.com/go-ozzo/ozzo-validation/v4"
)

const (
	LoggerLevelError   LoggerLevel = "error"
	LoggerLevelWarning LoggerLevel = "warning"
	LoggerLevelInfo    LoggerLevel = "info"
	LoggerLevelDebug   LoggerLevel = "debug"
	LoggerLevelNone    LoggerLevel = "none"
)

type LoggerLevel = string

type Logger struct {
	Level LoggerLevel       `yaml:"level" default:"debug"`
	Tags  map[string]string `yaml:"tags"`
}

func (l *Logger) Validate() (err error) {
	return validation.ValidateStruct(l,
		validation.Field(&l.Level, validation.By(l.checkLogLevel)))
}

func (l *Logger) checkLogLevel(v interface{}) error {
	level := v.(LoggerLevel)
	if !(level == LoggerLevelError ||
		level == LoggerLevelWarning ||
		level == LoggerLevelInfo ||
		level == LoggerLevelDebug ||
		level == LoggerLevelNone) {
		return errors.New("must be one of 'error', 'info', 'warning', or 'debug")
	}
	return nil
}
