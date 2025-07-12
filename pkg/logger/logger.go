package logger

import (
	"go.uber.org/zap"
)

// Logger is a wrapper around the zap logger.
type Logger struct {
	*zap.Logger
}

// New creates a new Logger.
func New(config *Config) (*Logger, error) {
	if err := config.Validate(); err != nil {
		return nil, err
	}

	var logger *zap.Logger
	var err error

	if config.Production {
		logger, err = zap.NewProduction()
	} else {
		logger, err = zap.NewDevelopment()
	}

	if err != nil {
		return nil, err
	}

	return &Logger{logger}, nil
}

// Config is the logger configuration.
type Config struct {
	Production bool `mapstructure:"production"`
}

// Validate validates the configuration.
func (c *Config) Validate() error {
	return nil
}

// Sync flushes any buffered log entries.
func (l *Logger) Sync() error {
	return l.Logger.Sync()
}

// Debug logs a message at the debug level.
func (l *Logger) Debug(args ...interface{}) {
	l.Logger.Sugar().Debug(args...)
}

// Info logs a message at the info level.
func (l *Logger) Info(args ...interface{}) {
	l.Logger.Sugar().Info(args...)
}

// Warn logs a message at the warn level.
func (l *Logger) Warn(args ...interface{}) {
	l.Logger.Sugar().Warn(args...)
}

// Error logs a message at the error level.
func (l *Logger) Error(args ...interface{}) {
	l.Logger.Sugar().Error(args...)
}

// Fatal logs a message at the fatal level.
func (l *Logger) Fatal(args ...interface{}) {
	l.Logger.Sugar().Fatal(args...)
}

// Debugf logs a formatted message at the debug level.
func (l *Logger) Debugf(template string, args ...interface{}) {
	l.Logger.Sugar().Debugf(template, args...)
}

// Infof logs a formatted message at the info level.
func (l *Logger) Infof(template string, args ...interface{}) {
	l.Logger.Sugar().Infof(template, args...)
}

// Warnf logs a formatted message at the warn level.
func (l *Logger) Warnf(template string, args ...interface{}) {
	l.Logger.Sugar().Warnf(template, args...)
}

// Errorf logs a formatted message at the error level.
func (l *Logger) Errorf(template string, args ...interface{}) {
	l.Logger.Sugar().Errorf(template, args...)
}

// Fatalf logs a formatted message at the fatal level.
func (l *Logger) Fatalf(template string, args ...interface{}) {
	l.Logger.Sugar().Fatalf(template, args...)
}
