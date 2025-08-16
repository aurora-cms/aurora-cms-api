package logging

import (
	"github.com/h4rdc0m/aurora-api/infrastructure/config"
	"testing"

	"github.com/h4rdc0m/aurora-api/domain/common"
	"github.com/stretchr/testify/assert"
)

func TestWrite(t *testing.T) {
	tests := []struct {
		name        string
		input       []byte
		expectedLen int
		expectedErr error
	}{
		{"valid input", []byte("test message"), len("test message"), nil},
		{"empty input", []byte(""), 0, nil},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockConfig := NewLogConfig(&config.Env{})

			// Create GinLogger using the actual NewGinLogger function
			ginLogger := NewGinLogger(NewZapLogger(mockConfig))

			n, err := ginLogger.Write(tc.input)

			assert.Equal(t, tc.expectedLen, n)
			assert.Equal(t, tc.expectedErr, err)

		})
	}
}

func TestNewGinLogger(t *testing.T) {
	tests := []struct {
		name   string
		logger common.Logger
	}{
		{"valid logger", &ZapLogger{}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.NotPanics(t, func() {
				logger := NewGinLogger(tc.logger)
				assert.NotNil(t, logger)
			})
		})
	}
}
