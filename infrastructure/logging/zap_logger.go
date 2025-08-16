package logging

import (
	"fmt"
	"github.com/h4rdc0m/aurora-api/domain/common"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ZapLogger struct {
	*zap.SugaredLogger
}

func NewZapLogger(logConfig *LogConfig) common.Logger {
	var config zap.Config
	logOutput := logConfig.LogOutput
	if logConfig.Environment == "production" {
		config = zap.NewProductionConfig()
		if logOutput != "" {
			config.OutputPaths = []string{logOutput}
		}
	} else {
		// Development configuration
		config = zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	logLevel := logConfig.LogLevel
	level := zap.PanicLevel
	switch logLevel {
	case "debug":
		level = zap.DebugLevel
	case "info":
		level = zap.InfoLevel
	case "warn":
		level = zap.WarnLevel
	case "error":
		level = zap.ErrorLevel
	case "fatal":
		level = zap.FatalLevel
	default:
		level = zap.PanicLevel
	}
	config.Level.SetLevel(level)

	zapLogger, err := config.Build()
	if err != nil {
		panic(fmt.Sprintf("failed to initialize logger: %v", err))
	}
	return &ZapLogger{
		SugaredLogger: zapLogger.Sugar(),
	}
}

// Debug logs a debug message
func (l *ZapLogger) Debug(msg string, args ...interface{}) {
	l.SugaredLogger.Debugw(msg, args...)
}

// Info logs an info message
func (l *ZapLogger) Info(msg string, args ...interface{}) {
	l.SugaredLogger.Infow(msg, args...)
}

// Warn logs a warning message
func (l *ZapLogger) Warn(msg string, args ...interface{}) {
	l.SugaredLogger.Warnw(msg, args...)
}

// Error logs an error message
func (l *ZapLogger) Error(msg string, args ...interface{}) {
	l.SugaredLogger.Errorw(msg, args...)
}

// Fatal logs a fatal message and exits
func (l *ZapLogger) Fatal(msg string, args ...interface{}) {
	l.SugaredLogger.Fatalw(msg, args...)
}

// Panic logs a fatal message and exits
func (l *ZapLogger) Panic(msg string, args ...interface{}) {
	l.SugaredLogger.Panicw(msg, args...)
}

// Print logs a fatal message and exits
func (l *ZapLogger) Print(msg string, args ...interface{}) {
	l.SugaredLogger.Debugw(msg, args...)
}

// Debugf logs a formatted debug message
func (l *ZapLogger) Debugf(format string, args ...interface{}) {
	l.SugaredLogger.Debugf(format, args...)
}

// Infof logs a formatted info message
func (l *ZapLogger) Infof(format string, args ...interface{}) {
	l.SugaredLogger.Infof(format, args...)
}

// Warnf logs a formatted warning message
func (l *ZapLogger) Warnf(format string, args ...interface{}) {
	l.SugaredLogger.Warnf(format, args...)
}

// Errorf logs a formatted error message
func (l *ZapLogger) Errorf(format string, args ...interface{}) {
	l.SugaredLogger.Errorf(format, args...)
}

// Fatalf logs a formatted fatal message and exits
func (l *ZapLogger) Fatalf(format string, args ...interface{}) {
	l.SugaredLogger.Fatalf(format, args...)
}

// Panicf logs a formatted fatal message and exits
func (l *ZapLogger) Panicf(format string, args ...interface{}) {
	l.SugaredLogger.Panicf(format, args...)
}

// Printf logs a formatted fatal message and exits
func (l *ZapLogger) Printf(format string, args ...interface{}) {
	l.SugaredLogger.Debugf(format, args...)
}
