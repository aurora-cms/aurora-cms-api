package repositories

import "github.com/h4rdc0m/aurora-api/domain/entities"

type TenantRepository interface {
	Save(tenant *entities.Tenant) error
	FindByID(id entities.TenantID) (*entities.Tenant, error)
	FindByName(name string) (*entities.Tenant, error)
	FindAll() ([]*entities.Tenant, error)
	FindActiveOnly() ([]*entities.Tenant, error)
	Delete(id entities.TenantID) error
	ExistsByName(name string) (bool, error)
}
