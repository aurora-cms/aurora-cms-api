package mocks

import (
	"github.com/h4rdc0m/aurora-api/domain/common"
	"github.com/stretchr/testify/mock"
)

type Logger struct {
	mock.Mock
}

// Ensure Logger implements common.Logger
var _ common.Logger = (*Logger)(nil)

func (m *Logger) Debug(msg string, args ...interface{}) {
	callArgs := []interface{}{msg}
	callArgs = append(callArgs, args...)
	m.Called(callArgs...)
}

func (m *Logger) Info(msg string, args ...interface{}) {
	callArgs := []interface{}{msg}
	callArgs = append(callArgs, args...)
	m.Called(callArgs...)
}

func (m *Logger) Warn(msg string, args ...interface{}) {
	callArgs := []interface{}{msg}
	callArgs = append(callArgs, args...)
	m.Called(callArgs...)
}

func (m *Logger) Error(msg string, args ...interface{}) {
	callArgs := []interface{}{msg}
	callArgs = append(callArgs, args...)
	m.Called(callArgs...)
}

func (m *Logger) Fatal(msg string, args ...interface{}) {
	callArgs := []interface{}{msg}
	callArgs = append(callArgs, args...)
	m.Called(callArgs...)
}

func (m *Logger) Panic(msg string, args ...interface{}) {
	callArgs := []interface{}{msg}
	callArgs = append(callArgs, args...)
	m.Called(callArgs...)
}

func (m *Logger) Debugf(template string, args ...interface{}) {
	callArgs := []interface{}{template}
	callArgs = append(callArgs, args...)
	m.Called(callArgs...)
}

func (m *Logger) Infof(template string, args ...interface{}) {
	callArgs := []interface{}{template}
	callArgs = append(callArgs, args...)
	m.Called(callArgs...)
}

func (m *Logger) Warnf(template string, args ...interface{}) {
	callArgs := []interface{}{template}
	callArgs = append(callArgs, args...)
	m.Called(callArgs...)
}

func (m *Logger) Errorf(template string, args ...interface{}) {
	callArgs := []interface{}{template}
	callArgs = append(callArgs, args...)
	m.Called(callArgs...)
}

func (m *Logger) Fatalf(template string, args ...interface{}) {
	callArgs := []interface{}{template}
	callArgs = append(callArgs, args...)
	m.Called(callArgs...)
}

func (m *Logger) Panicf(template string, args ...interface{}) {
	callArgs := []interface{}{template}
	callArgs = append(callArgs, args...)
	m.Called(callArgs...)
}

func (m *Logger) Print(msg string, args ...interface{}) {
	callArgs := []interface{}{msg}
	callArgs = append(callArgs, args...)
	m.Called(callArgs...)
}

func (m *Logger) Printf(template string, args ...interface{}) {
	callArgs := []interface{}{template}
	callArgs = append(callArgs, args...)
	m.Called(callArgs...)
}
