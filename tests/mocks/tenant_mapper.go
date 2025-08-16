package mocks

import (
	"github.com/h4rdc0m/aurora-api/domain/entities"
	"github.com/h4rdc0m/aurora-api/infrastructure/persistence/models"
)

type MockTenantMapper struct {
	MockMapper[models.Tenant, entities.Tenant]
}

// ToModel converts a domain entity to a persistence model
func (m *MockTenantMapper) ToModel(entity *entities.Tenant) (*models.Tenant, error) {
	args := m.Called(entity)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Tenant), args.Error(1)
}

// ToDomain converts a persistence model to a domain entity
func (m *MockTenantMapper) ToDomain(model *models.Tenant) (*entities.Tenant, error) {
	args := m.Called(model)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Tenant), args.Error(1)
}

// ToModels converts a slice of domain entities to persistence models
func (m *MockTenantMapper) ToModels(entities []*entities.Tenant) ([]*models.Tenant, error) {
	args := m.Called(entities)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Tenant), args.Error(1)
}

// ToDomains converts a slice of persistence models to domain entities
func (m *MockTenantMapper) ToDomains(models []*models.Tenant) ([]*entities.Tenant, error) {
	args := m.Called(models)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entities.Tenant), args.Error(1)
}
