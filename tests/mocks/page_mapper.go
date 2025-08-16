package mocks

import (
	"github.com/h4rdc0m/aurora-api/domain/entities"
	"github.com/h4rdc0m/aurora-api/infrastructure/persistence/models"
	"github.com/stretchr/testify/mock"
)

// MockPageMapper is a mock implementation of the Mapper interface for Page entities
type MockPageMapper struct {
	mock.Mock
}

// ToModel converts a domain entity to a persistence model
func (m *MockPageMapper) ToModel(entity *entities.Page) (*models.Page, error) {
	args := m.Called(entity)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Page), args.Error(1)
}

// ToDomain converts a persistence model to a domain entity
func (m *MockPageMapper) ToDomain(model *models.Page) (*entities.Page, error) {
	args := m.Called(model)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Page), args.Error(1)
}

// ToModels converts a slice of domain entities to persistence models
func (m *MockPageMapper) ToModels(entities []*entities.Page) ([]*models.Page, error) {
	args := m.Called(entities)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Page), args.Error(1)
}

// ToDomains converts a slice of persistence models to domain entities
func (m *MockPageMapper) ToDomains(models []*models.Page) ([]*entities.Page, error) {
	args := m.Called(models)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entities.Page), args.Error(1)
}

// MockPageBlockMapper is a mock implementation of the Mapper interface for Page entities
type MockPageBlockMapper struct {
	mock.Mock
}

// ToModel converts a domain entity to a persistence model
func (m *MockPageBlockMapper) ToModel(entity *entities.PageBlock) (*models.PageBlock, error) {
	args := m.Called(entity)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.PageBlock), args.Error(1)
}

// ToDomain converts a persistence model to a domain entity
func (m *MockPageBlockMapper) ToDomain(model *models.PageBlock) (*entities.PageBlock, error) {
	args := m.Called(model)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.PageBlock), args.Error(1)
}

// ToModels converts a slice of domain entities to persistence models
func (m *MockPageBlockMapper) ToModels(entities []*entities.PageBlock) ([]*models.PageBlock, error) {
	args := m.Called(entities)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.PageBlock), args.Error(1)
}

// ToDomains converts a slice of persistence models to domain entities
func (m *MockPageBlockMapper) ToDomains(models []*models.PageBlock) ([]*entities.PageBlock, error) {
	args := m.Called(models)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entities.PageBlock), args.Error(1)
}

// MockPageVersionMapper is a mock implementation of the Mapper interface for Page entities
type MockPageVersionMapper struct {
	mock.Mock
}

// ToModel converts a domain entity to a persistence model
func (m *MockPageVersionMapper) ToModel(entity *entities.PageVersion) (*models.PageVersion, error) {
	args := m.Called(entity)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.PageVersion), args.Error(1)
}

// ToDomain converts a persistence model to a domain entity
func (m *MockPageVersionMapper) ToDomain(model *models.PageVersion) (*entities.PageVersion, error) {
	args := m.Called(model)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.PageVersion), args.Error(1)
}

// ToModels converts a slice of domain entities to persistence models
func (m *MockPageVersionMapper) ToModels(entities []*entities.PageVersion) ([]*models.PageVersion, error) {
	args := m.Called(entities)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.PageVersion), args.Error(1)
}

// ToDomains converts a slice of persistence models to domain entities
func (m *MockPageVersionMapper) ToDomains(models []*models.PageVersion) ([]*entities.PageVersion, error) {
	args := m.Called(models)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entities.PageVersion), args.Error(1)
}
