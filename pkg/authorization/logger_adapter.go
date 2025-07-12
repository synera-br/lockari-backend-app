package authorization

import (
	"log/slog"
)

// SlogAdapter adapta slog.Logger para a interface Logger
type SlogAdapter struct {
	logger *slog.Logger
}

// NewSlogAdapter cria um novo adapter para slog.Logger
func NewSlogAdapter(logger *slog.Logger) Logger {
	return &SlogAdapter{logger: logger}
}

// Debug registra uma mensagem de debug
func (s *SlogAdapter) Debug(msg string, fields ...interface{}) {
	s.logger.Debug(msg, fields...)
}

// Info registra uma mensagem informativa
func (s *SlogAdapter) Info(msg string, fields ...interface{}) {
	s.logger.Info(msg, fields...)
}

// Warn registra uma mensagem de aviso
func (s *SlogAdapter) Warn(msg string, fields ...interface{}) {
	s.logger.Warn(msg, fields...)
}

// Error registra uma mensagem de erro
func (s *SlogAdapter) Error(msg string, fields ...interface{}) {
	s.logger.Error(msg, fields...)
}

// With adiciona campos ao logger
func (s *SlogAdapter) With(fields ...interface{}) Logger {
	return &SlogAdapter{logger: s.logger.With(fields...)}
}
