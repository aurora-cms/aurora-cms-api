// zap_logger_test.go
package logging

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestZapLogger_Debug(t *testing.T) {
	tests := []struct {
		name string
		msg  string
		args []interface{}
	}{
		{"simple_debug", "debug_message", nil},
		{"debug_with_args", "debug_message_with_args", []interface{}{"key", "value"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			logger := newTestLogger(&buf, zap.DebugLevel)
			logger.Debug(tt.msg, tt.args...)
			assert.Contains(t, buf.String(), tt.msg)
			if tt.args != nil {
				for _, arg := range tt.args {
					assert.Contains(t, buf.String(), arg.(string))
				}
			}
		})
	}
}

func TestZapLogger_Info(t *testing.T) {
	tests := []struct {
		name string
		msg  string
		args []interface{}
	}{
		{"simple_info", "info_message", nil},
		{"info_with_args", "info_message_with_args", []interface{}{"key", "value"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			logger := newTestLogger(&buf, zap.InfoLevel)
			logger.Info(tt.msg, tt.args...)
			assert.Contains(t, buf.String(), tt.msg)
			if tt.args != nil {
				for _, arg := range tt.args {
					assert.Contains(t, buf.String(), arg.(string))
				}
			}
		})
	}
}

func TestZapLogger_Error(t *testing.T) {
	tests := []struct {
		name string
		msg  string
		args []interface{}
	}{
		{"simple_error", "error_message", nil},
		{"error_with_args", "error_message_with_args", []interface{}{"key", "value"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			logger := newTestLogger(&buf, zap.ErrorLevel)
			logger.Error(tt.msg, tt.args...)
			assert.Contains(t, buf.String(), tt.msg)
			if tt.args != nil {
				for _, arg := range tt.args {
					assert.Contains(t, buf.String(), arg.(string))
				}
			}
		})
	}
}

func TestZapLogger_Debugf(t *testing.T) {
	tests := []struct {
		name   string
		format string
		args   []interface{}
	}{
		{"debugf", "debug %s", []interface{}{"log"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			logger := newTestLogger(&buf, zap.DebugLevel)
			logger.Debugf(tt.format, tt.args...)
			assert.Contains(t, buf.String(), "debug log")
		})
	}
}

func TestNewZapLogger_ConfiguresZapCorrectly(t *testing.T) {
	tests := []struct {
		name      string
		config    *LogConfig
		wantLevel zapcore.Level
	}{
		{"production_debug", &LogConfig{Environment: "production", LogLevel: "debug"}, zap.DebugLevel},
		{"production_error", &LogConfig{Environment: "production", LogLevel: "error"}, zap.ErrorLevel},
		{"development_info", &LogConfig{Environment: "development", LogLevel: "info"}, zap.InfoLevel},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			loggerIface := NewZapLogger(tt.config)
			logger, ok := loggerIface.(*ZapLogger)
			assert.True(t, ok)
			desugared := logger.Desugar()
			assert.True(t, desugared.Core().Enabled(tt.wantLevel))
		})
	}
}

func TestZapLogger_Integration(t *testing.T) {
	var buf bytes.Buffer
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		zapcore.AddSync(&buf),
		zap.DebugLevel,
	)
	logger := &ZapLogger{SugaredLogger: zap.New(core).Sugar()}

	logger.Debug("debug", "key", "val")
	logger.Info("info", "key", "val")
	logger.Warn("warn", "key", "val")
	logger.Errorf("errorf %s", "val")

	logs := buf.String()
	assert.Contains(t, logs, "debug")
	assert.Contains(t, logs, "info")
	assert.Contains(t, logs, "warn")
	assert.Contains(t, logs, "errorf val")
	assert.Contains(t, logs, "key")
	assert.Contains(t, logs, "val")
}

func newTestLogger(buf *bytes.Buffer, level zapcore.Level) *ZapLogger {
	writer := zapcore.AddSync(buf)
	encoder := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	core := zapcore.NewCore(encoder, writer, level)
	return &ZapLogger{zap.New(core).Sugar()}
}
