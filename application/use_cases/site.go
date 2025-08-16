package use_cases

import (
	"errors"
	"github.com/h4rdc0m/aurora-api/domain/common"
	"github.com/h4rdc0m/aurora-api/domain/entities"
	"github.com/h4rdc0m/aurora-api/domain/repositories"
	"github.com/h4rdc0m/aurora-api/domain/value_objects"
)

// SiteUseCase handles site business logic
type SiteUseCase struct {
	siteRepo   repositories.SiteRepository
	tenantRepo repositories.TenantRepository
	logger     common.Logger
}

// NewSiteUseCase creates a new SiteUseCase
func NewSiteUseCase(
	siteRepo repositories.SiteRepository,
	tenantRepo repositories.TenantRepository,
	logger common.Logger,
) *SiteUseCase {
	return &SiteUseCase{
		siteRepo:   siteRepo,
		tenantRepo: tenantRepo,
		logger:     logger,
	}
}

// CreateSite creates a new site
func (u *SiteUseCase) CreateSite(name string, description *string, domainStr string, templateID uint64, tenantID uint64) (*entities.Site, error) {
	// Validate tenant exists
	tenant, err := u.tenantRepo.FindByID(entities.NewTenantID(tenantID))
	if err != nil {
		return nil, err
	}
	if tenant == nil {
		return nil, errors.New("tenant not found")
	}

	// Create domain value object
	domain, err := value_objects.NewDomainName(domainStr)
	if err != nil {
		return nil, err
	}

	// Check if domain already exists
	exists, err := u.siteRepo.ExistsByDomain(domain)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("site with this domain already exists")
	}

	// Create new site entity
	site, err := entities.NewSite(
		name,
		description,
		domain,
		entities.NewTemplateID(templateID),
		entities.NewTenantID(tenantID),
	)
	if err != nil {
		return nil, err
	}

	// Save site
	if err := u.siteRepo.Save(site); err != nil {
		return nil, err
	}

	return site, nil
}

// GetSite retrieves a site by ID
func (u *SiteUseCase) GetSite(id uint64) (*entities.Site, error) {
	return u.siteRepo.FindByID(entities.NewSiteID(id))
}

// GetSiteByDomain retrieves a site by domain
func (u *SiteUseCase) GetSiteByDomain(domainStr string) (*entities.Site, error) {
	domain, err := value_objects.NewDomainName(domainStr)
	if err != nil {
		return nil, err
	}

	return u.siteRepo.FindByDomain(domain)
}

// GetSitesByTenant retrieves all sites for a tenant
func (u *SiteUseCase) GetSitesByTenant(tenantID uint64) ([]*entities.Site, error) {
	return u.siteRepo.FindByTenantID(entities.NewTenantID(tenantID))
}

// GetEnabledSitesByTenant retrieves only enabled sites for a tenant
func (u *SiteUseCase) GetEnabledSitesByTenant(tenantID uint64) ([]*entities.Site, error) {
	return u.siteRepo.FindEnabledByTenantID(entities.NewTenantID(tenantID))
}

// UpdateSite updates a site
func (u *SiteUseCase) UpdateSite(id uint64, name string, description *string, domainStr string, templateID uint64) (*entities.Site, error) {
	// Get existing site
	site, err := u.siteRepo.FindByID(entities.NewSiteID(id))
	if err != nil {
		return nil, err
	}
	if site == nil {
		return nil, errors.New("site not found")
	}

	// Update name if provided
	if name != "" {
		if err := site.UpdateName(name); err != nil {
			return nil, err
		}
	}

	// Update description
	err = site.UpdateDescription(description)
	if err != nil {
		return nil, err
	}

	// Update domain if provided
	if domainStr != "" {
		domain, err := value_objects.NewDomainName(domainStr)
		if err != nil {
			return nil, err
		}

		// Check if new domain already exists (and it's not the current site)
		existingSite, err := u.siteRepo.FindByDomain(domain)
		if err != nil {
			return nil, err
		}
		if existingSite != nil && existingSite.ID().Value() != site.ID().Value() {
			return nil, errors.New("site with this domain already exists")
		}

		if err := site.UpdateDomain(domain); err != nil {
			return nil, err
		}
	}

	// Update template if provided
	if templateID != 0 {
		site.UpdateTemplate(entities.NewTemplateID(templateID))
	}

	// Save updated site
	if err := u.siteRepo.Save(site); err != nil {
		return nil, err
	}

	return site, nil
}

// DeleteSite deletes a site
func (u *SiteUseCase) DeleteSite(id uint64) error {
	siteID := entities.NewSiteID(id)

	// Check if site exists
	site, err := u.siteRepo.FindByID(siteID)
	if err != nil {
		return err
	}
	if site == nil {
		return errors.New("site not found")
	}

	return u.siteRepo.Delete(siteID)
}

// EnableSite enables a site
func (u *SiteUseCase) EnableSite(id uint64) error {
	site, err := u.siteRepo.FindByID(entities.NewSiteID(id))
	if err != nil {
		return err
	}
	if site == nil {
		return errors.New("site not found")
	}

	site.Enable()
	return u.siteRepo.Save(site)
}

// DisableSite disables a site
func (u *SiteUseCase) DisableSite(id uint64) error {
	site, err := u.siteRepo.FindByID(entities.NewSiteID(id))
	if err != nil {
		return err
	}
	if site == nil {
		return errors.New("site not found")
	}

	site.Disable()
	return u.siteRepo.Save(site)
}
