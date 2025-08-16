package repositories

import (
	"github.com/h4rdc0m/aurora-api/domain/entities"
	"github.com/h4rdc0m/aurora-api/domain/value_objects"
)

// SiteRepository defines the interface for site data operations
type SiteRepository interface {
	Save(site *entities.Site) error
	FindByID(id entities.SiteID) (*entities.Site, error)
	FindByDomain(domain *value_objects.DomainName) (*entities.Site, error)
	FindByTenantID(tenantID entities.TenantID) ([]*entities.Site, error)
	FindAll() ([]*entities.Site, error)
	FindEnabledByTenantID(tenantID entities.TenantID) ([]*entities.Site, error)
	Delete(id entities.SiteID) error
	ExistsByDomain(domain *value_objects.DomainName) (bool, error)
}
