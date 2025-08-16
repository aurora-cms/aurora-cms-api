package mocks

import (
	"github.com/stretchr/testify/mock"
)

// MockMapper is a generic mock for the Mapper interface
type MockMapper[E any, M any] struct {
	mock.Mock
}

func (m *MockMapper[E, M]) ToModel(entity E) (M, error) {
	args := m.Called(entity)
	return args.Get(0).(M), args.Error(1)
}

func (m *MockMapper[E, M]) ToDomain(model M) (E, error) {
	args := m.Called(model)
	return args.Get(0).(E), args.Error(1)
}

func (m *MockMapper[E, M]) ToDomains(models []M) ([]E, error) {
	args := m.Called(models)
	return args.Get(0).([]E), args.Error(1)
}

func (m *MockMapper[E, M]) ToModels(entities []E) ([]M, error) {
	args := m.Called(entities)
	return args.Get(0).([]M), args.Error(1)
}
