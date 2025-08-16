package use_cases

import (
	"github.com/h4rdc0m/aurora-api/domain/common"
	"github.com/h4rdc0m/aurora-api/domain/entities"
	"github.com/h4rdc0m/aurora-api/domain/errors"
	"github.com/h4rdc0m/aurora-api/domain/repositories"
)

type TenantUseCase struct {
	tenantRepo repositories.TenantRepository
	siteRepo   repositories.SiteRepository
	logger     common.Logger
}

// NewTenantUseCase creates a new TenantUseCase
func NewTenantUseCase(
	tenantRepo repositories.TenantRepository,
	siteRepo repositories.SiteRepository,
	logger common.Logger,
) *TenantUseCase {
	return &TenantUseCase{
		tenantRepo: tenantRepo,
		siteRepo:   siteRepo,
		logger:     logger,
	}
}

// GetTenant retrieves a tenant by its ID
func (u *TenantUseCase) GetTenant(id uint64) (*entities.Tenant, error) {
	tenant, err := u.tenantRepo.FindByID(entities.NewTenantID(id))
	if err != nil {
		u.logger.Error("Failed to get tenant", "id", id, "error", err)
		return nil, err
	}
	if tenant == nil {
		u.logger.Warn("Tenant not found", "id", id)
		return nil, nil // or return an error if preferred
	}
	return tenant, nil
}

// GetAllTenants retrieves all tenants
func (u *TenantUseCase) GetAllTenants() ([]*entities.Tenant, error) {
	tenants, err := u.tenantRepo.FindAll()
	if err != nil {
		u.logger.Error("Failed to get all tenants", "error", err)
		return nil, err
	}
	return tenants, nil
}

// CreateTenant creates a new tenant
func (u *TenantUseCase) CreateTenant(name string, description *string) (*entities.Tenant, error) {
	if name == "" {
		return nil, errors.ErrTenantNameEmpty
	}

	tenant, err := entities.NewTenant(name, description)
	if err != nil {
		u.logger.Error("Failed to create tenant", "name", name, "error", err)
		return nil, err
	}
	tenant.Activate()

	err = u.tenantRepo.Save(tenant)
	if err != nil {
		u.logger.Error("Failed to save tenant", "name", name, "error", err)
		return nil, err
	}

	return tenant, nil
}

// GetActiveTenants retrieves only active tenants
func (u *TenantUseCase) GetActiveTenants() ([]*entities.Tenant, error) {
	tenants, err := u.tenantRepo.FindActiveOnly()
	if err != nil {
		u.logger.Error("Failed to get active tenants", "error", err)
		return nil, err
	}
	return tenants, nil
}

// UpdateTenant updates an existing tenant
func (u *TenantUseCase) UpdateTenant(id uint64, name string, description *string) (*entities.Tenant, error) {
	if name == "" {
		return nil, errors.ErrTenantNameEmpty
	}

	tenant, err := u.tenantRepo.FindByID(entities.NewTenantID(id))
	if err != nil {
		u.logger.Error("Failed to find tenant for update", "id", id, "error", err)
		return nil, err
	}
	if tenant == nil {
		u.logger.Warn("Tenant not found for update", "id", id)
		return nil, errors.ErrTenantNotFound
	}

	err = tenant.UpdateName(name)
	if err != nil {
		u.logger.Error("Failed to update tenant name", "id", id, "error", err)
		return nil, err
	}

	if description != nil {
		err = tenant.UpdateDescription(description)
		if err != nil {
			u.logger.Error("Failed to update tenant description", "id", id, "error", err)
			return nil, err
		}
	}

	err = u.tenantRepo.Save(tenant)
	if err != nil {
		u.logger.Error("Failed to save updated tenant", "id", id, "error", err)
		return nil, err
	}

	return tenant, nil
}

// DeleteTenant deletes a tenant by its ID
func (u *TenantUseCase) DeleteTenant(id uint64) error {
	tenant, err := u.tenantRepo.FindByID(entities.NewTenantID(id))
	if err != nil {
		u.logger.Error("Failed to find tenant for deletion", "id", id, "error", err)
		return err
	}
	if tenant == nil {
		u.logger.Warn("Tenant not found for deletion", "id", id)
		return errors.ErrTenantNotFound
	}

	err = u.tenantRepo.Delete(tenant.ID())
	if err != nil {
		u.logger.Error("Failed to delete tenant", "id", id, "error", err)
		return err
	}

	return nil
}
