package mocks

import (
	"github.com/h4rdc0m/aurora-api/domain/entities"
	"github.com/h4rdc0m/aurora-api/infrastructure/persistence/models"
)

type MockSiteMapper struct {
	MockMapper[models.Site, entities.Site]
}

// ToModel converts a domain entity to a persistence model
func (m *MockSiteMapper) ToModel(entity *entities.Site) (*models.Site, error) {
	args := m.Called(entity)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Site), args.Error(1)
}

// ToDomain converts a persistence model to a domain entity
func (m *MockSiteMapper) ToDomain(model *models.Site) (*entities.Site, error) {
	args := m.Called(model)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Site), args.Error(1)
}

// ToModels converts a slice of domain entities to persistence models
func (m *MockSiteMapper) ToModels(entities []*entities.Site) ([]*models.Site, error) {
	args := m.Called(entities)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Site), args.Error(1)
}

// ToDomains converts a slice of persistence models to domain entities
func (m *MockSiteMapper) ToDomains(models []*models.Site) ([]*entities.Site, error) {
	args := m.Called(models)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entities.Site), args.Error(1)
}
