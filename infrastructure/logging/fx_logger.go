package logging

import (
	"github.com/h4rdc0m/aurora-api/domain/common"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
)

// FxLogger wraps a common.Logger and provides logging integration for Fx events.
type FxLogger struct {
	logger    common.Logger
	zapLogger *ZapLogger
}

// NewFxLogger initializes and returns a new fxevent.Logger instance using the provided common.Logger interface.
func NewFxLogger(logger common.Logger) fxevent.Logger {
	zapLogger, ok := logger.(*ZapLogger)
	if !ok {
		panic("FxLogger requires a ZapLogger instance")
	}

	baseLogger := zapLogger.WithOptions(
		zap.WithCaller(false),
	)

	return &FxLogger{
		logger:    &ZapLogger{SugaredLogger: baseLogger},
		zapLogger: zapLogger,
	}
}

// LogEvent processes and logs events using the provided fxevent.Event, categorizing them based on their type.
func (l *FxLogger) LogEvent(event fxevent.Event) {
	switch e := event.(type) {
	case *fxevent.OnStartExecuting:
		l.logger.Debug("OnStart hook executing",
			zap.String("function", e.FunctionName),
			zap.String("caller", e.CallerName),
		)
	case *fxevent.OnStartExecuted:
		if e.Err != nil {
			l.logger.Error("OnStart hook failed",
				zap.String("function", e.FunctionName),
				zap.String("caller", e.CallerName),
				zap.Error(e.Err))
		} else {
			l.logger.Debug("OnStart hook executed successfully",
				zap.String("function", e.FunctionName),
				zap.String("caller", e.CallerName),
				zap.String("runtime", e.Runtime.String()),
			)
		}
	case *fxevent.OnStopExecuting:
		l.logger.Debug("OnStop hook executing",
			zap.String("function", e.FunctionName),
			zap.String("caller", e.CallerName),
		)
	case *fxevent.OnStopExecuted:
		if e.Err != nil {
			l.logger.Error("OnStop hook failed",
				zap.String("function", e.FunctionName),
				zap.String("caller", e.CallerName),
				zap.Error(e.Err))
		} else {
			l.logger.Debug("OnStop hook executed successfully",
				zap.String("function", e.FunctionName),
				zap.String("caller", e.CallerName),
				zap.String("runtime", e.Runtime.String()),
			)
		}
	case *fxevent.Supplied:
		l.logger.Debug("Supplied",
			zap.String("type", e.TypeName),
			zap.Error(e.Err),
		)
	case *fxevent.Provided:
		for _, rType := range e.OutputTypeNames {
			l.logger.Debug("Provided", zap.String("constructor", e.ConstructorName), zap.String("type", rType))
		}
	case *fxevent.Decorated:
		for _, rType := range e.OutputTypeNames {
			l.logger.Debug("Decorated", zap.String("decorator", e.DecoratorName), zap.String("type", rType))
		}
	case *fxevent.Invoking:
		l.logger.Debug("Invoking", zap.String("function", e.FunctionName))
	case *fxevent.Started:
		if e.Err == nil {
			l.logger.Debug("Started")
		}
	case *fxevent.Stopped:
		if e.Err == nil {
			l.logger.Debug("Stopped")
		}
	case *fxevent.LoggerInitialized:
		if e.Err == nil {
			l.logger.Debug("Initialized: custom fxevent.Logger", zap.String("constructor", e.ConstructorName))
		}
	default:
		l.logger.Warn("Unhandled fxevent type", e)
	}

}
