package logging

import (
	"github.com/h4rdc0m/aurora-api/tests/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
)

func TestNewFxLogger(t *testing.T) {
	t.Run("valid ZapLogger creation", func(t *testing.T) {
		zapLogger := &ZapLogger{SugaredLogger: zap.NewExample().Sugar()}

		logger := NewFxLogger(zapLogger)
		assert.NotNil(t, logger)
		assert.IsType(t, &FxLogger{}, logger)
	})

	t.Run("panic on invalid Logger type", func(t *testing.T) {
		mockLogger := &mocks.Logger{}

		assert.PanicsWithValue(t, "FxLogger requires a ZapLogger instance", func() {
			NewFxLogger(mockLogger)
		})
	})
}

func TestFxLogger_LogEvent(t *testing.T) {
	t.Run("log fxevent.Started", func(t *testing.T) {
		mockLogger := &mocks.Logger{}
		//mockLogger.On("Info", "Started", mock.Anything).Return()

		fxLogger := &FxLogger{
			logger: mockLogger,
		}

		mockLogger.On("Debug", "Started").Return()

		event := &fxevent.Started{}
		fxLogger.LogEvent(event)

		mockLogger.AssertCalled(t, "Debug", "Started", mock.Anything)
		mockLogger.AssertExpectations(t)
	})

	t.Run("log fxevent.Stopped", func(t *testing.T) {
		mockLogger := &mocks.Logger{}
		mockLogger.On("Debug", "Stopped", mock.Anything).Return()

		fxLogger := &FxLogger{
			logger: mockLogger,
		}

		event := &fxevent.Stopped{}
		fxLogger.LogEvent(event)

		mockLogger.AssertCalled(t, "Debug", "Stopped", mock.Anything)
	})

	t.Run("log fxevent.LoggerInitialized", func(t *testing.T) {
		mockLogger := &mocks.Logger{}
		mockLogger.On("Debug", "Initialized: custom fxevent.Logger", mock.Anything).Return()

		fxLogger := &FxLogger{
			logger: mockLogger,
		}

		event := &fxevent.LoggerInitialized{ConstructorName: "example"}
		fxLogger.LogEvent(event)

		mockLogger.AssertCalled(t, "Debug", "Initialized: custom fxevent.Logger", mock.Anything)
	})

	t.Run("log unknown event", func(t *testing.T) {
		mockLogger := &mocks.Logger{}
		mockLogger.On("Warn", mock.Anything, mock.Anything).Return()

		fxLogger := &FxLogger{
			logger: mockLogger,
		}

		event := struct{ fxevent.Event }{}
		fxLogger.LogEvent(event)

		mockLogger.AssertCalled(t, "Warn", mock.Anything, mock.Anything)
	})
}
