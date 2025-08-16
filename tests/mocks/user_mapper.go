package mocks

import (
	"github.com/h4rdc0m/aurora-api/domain/entities"
	"github.com/h4rdc0m/aurora-api/infrastructure/persistence/models"
)

type MockUserMapper struct {
	MockMapper[models.User, entities.User]
}

// ToModel converts a domain entity to a persistence model
func (m *MockUserMapper) ToModel(entity *entities.User) (*models.User, error) {
	args := m.Called(entity)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

// ToDomain converts a persistence model to a domain entity
func (m *MockUserMapper) ToDomain(model *models.User) (*entities.User, error) {
	args := m.Called(model)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.User), args.Error(1)
}

// ToModels converts a slice of domain entities to persistence models
func (m *MockUserMapper) ToModels(entities []*entities.User) ([]*models.User, error) {
	args := m.Called(entities)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.User), args.Error(1)
}

// ToDomains converts a slice of persistence models to domain entities
func (m *MockUserMapper) ToDomains(models []*models.User) ([]*entities.User, error) {
	args := m.Called(models)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entities.User), args.Error(1)
}
