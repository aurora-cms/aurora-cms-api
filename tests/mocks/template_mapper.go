package mocks

import (
	"github.com/h4rdc0m/aurora-api/domain/entities"
	"github.com/h4rdc0m/aurora-api/infrastructure/persistence/models"
)

type MockTemplateMapper struct {
	MockMapper[models.Template, entities.Template]
}

// ToModel converts a domain entity to a persistence model
func (m *MockTemplateMapper) ToModel(entity *entities.Template) (*models.Template, error) {
	args := m.Called(entity)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Template), args.Error(1)
}

// ToDomain converts a persistence model to a domain entity
func (m *MockTemplateMapper) ToDomain(model *models.Template) (*entities.Template, error) {
	args := m.Called(model)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Template), args.Error(1)
}

// ToModels converts a slice of domain entities to persistence models
func (m *MockTemplateMapper) ToModels(entities []*entities.Template) ([]*models.Template, error) {
	args := m.Called(entities)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Template), args.Error(1)
}

// ToDomains converts a slice of persistence models to domain entities
func (m *MockTemplateMapper) ToDomains(models []*models.Template) ([]*entities.Template, error) {
	args := m.Called(models)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entities.Template), args.Error(1)
}
