package common

// Logger defines a generic logging interface with methods for leveled and formatted logging.
type Logger interface {
	Debug(msg string, args ...interface{})
	Info(msg string, args ...interface{})
	Warn(msg string, args ...interface{})
	Error(msg string, args ...interface{})
	Fatal(msg string, args ...interface{})
	Panic(msg string, args ...interface{})
	Print(msg string, args ...interface{})

	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
	Panicf(format string, args ...interface{})
	Printf(format string, args ...interface{})
}

// GinLogger is an interface for implementing a logger compatible with the Write method.
// Write processes a byte slice and outputs a log message, returning the byte count and any errors encountered.
type GinLogger interface {
	Write(p []byte) (n int, err error)
}
