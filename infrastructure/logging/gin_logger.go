package logging

import "github.com/h4rdc0m/aurora-api/domain/common"

// GinLogger is a logger adapter for the Gin framework that utilizes a common logging interface and a Zap logger.
// It integrates ZapLogger for structured and leveled logging and is designed as a Write interface implementation.
// GinLogger supports seamless logging within middleware and various functional layers in the application.
type GinLogger struct {
	logger    common.Logger
	zapLogger *ZapLogger
}

// NewGinLogger initializes and returns a new GinLogger instance using the provided common.Logger interface.
func NewGinLogger(logger common.Logger) common.GinLogger {
	zapLogger, ok := logger.(*ZapLogger)
	if !ok {
		panic("GinLogger requires a ZapLogger instance")
	}

	return &GinLogger{
		logger:    logger,
		zapLogger: zapLogger,
	}
}

// Write writes the given byte slice to the logger as an info message and returns the number of bytes written and any error.
func (l *GinLogger) Write(p []byte) (n int, err error) {
	l.logger.Info(string(p))
	return len(p), nil
}
