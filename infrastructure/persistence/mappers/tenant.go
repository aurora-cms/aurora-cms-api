package mappers

import (
	"github.com/h4rdc0m/aurora-api/domain/entities"
	"github.com/h4rdc0m/aurora-api/domain/errors"
	"github.com/h4rdc0m/aurora-api/infrastructure/persistence/models"
)

type TenantMapper struct{}

func NewTenantMapper() *TenantMapper {
	return &TenantMapper{}
}

func (m *TenantMapper) ToModel(tenant *entities.Tenant) (*models.Tenant, error) {
	if tenant == nil {
		return nil, nil
	}

	if tenant.Name() == "" {
		return nil, errors.ErrTenantNameEmpty
	}

	return &models.Tenant{
		Base: models.Base{
			ID:        tenant.ID().Value(),
			CreatedAt: tenant.CreatedAt(),
			UpdatedAt: tenant.UpdatedAt(),
		},
		Name:             tenant.Name(),
		Description:      tenant.Description(),
		IsActive:         tenant.IsActive(),
		IsBillingEnabled: tenant.IsBillingEnabled(),
	}, nil
}

func (m *TenantMapper) ToDomain(model *models.Tenant) (*entities.Tenant, error) {
	if model == nil {
		return nil, nil
	}
	if model.Name == "" {
		return nil, errors.ErrTenantNameEmpty
	}

	tenant, err := entities.NewTenant(model.Name, model.Description)
	if err != nil {
		return nil, err
	}

	tenant.SetID(entities.NewTenantID(model.ID))
	tenant.SetTimestamps(model.CreatedAt, model.UpdatedAt)

	if !model.IsActive {
		tenant.Deactivate()
	} else {
		tenant.Activate()
	}

	if model.IsBillingEnabled {
		tenant.EnableBilling()
	} else {
		tenant.DisableBilling()
	}

	return tenant, nil
}

func (m *TenantMapper) ToModels(tenants []*entities.Tenant) ([]*models.Tenant, error) {
	if tenants == nil {
		return nil, nil
	}

	result := make([]*models.Tenant, len(tenants))
	for i, tenant := range tenants {
		model, err := m.ToModel(tenant)
		if err != nil {
			return nil, err
		}
		result[i] = model
	}

	return result, nil
}

func (m *TenantMapper) ToDomains(modelsList []*models.Tenant) ([]*entities.Tenant, error) {
	if modelsList == nil {
		return nil, nil
	}

	result := make([]*entities.Tenant, len(modelsList))
	for i, model := range modelsList {
		domain, err := m.ToDomain(model)
		if err != nil {
			return nil, err
		}
		result[i] = domain
	}

	return result, nil
}
