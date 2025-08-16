package mappers

import (
	"github.com/h4rdc0m/aurora-api/domain/entities"
	"github.com/h4rdc0m/aurora-api/domain/errors"
	"github.com/h4rdc0m/aurora-api/domain/value_objects"
	"github.com/h4rdc0m/aurora-api/infrastructure/persistence/models"
)

// SiteMapper handles conversion between domain entities and GORM models
type SiteMapper struct{}

// NewSiteMapper creates a new SiteMapper
func NewSiteMapper() *SiteMapper {
	return &SiteMapper{}
}

// ToModel converts a domain Site to a GORM models.Site
func (m *SiteMapper) ToModel(site *entities.Site) (*models.Site, error) {
	if site == nil {
		return nil, nil
	}

	return &models.Site{
		Base: models.Base{
			ID:        site.ID().Value(),
			CreatedAt: site.CreatedAt(),
			UpdatedAt: site.UpdatedAt(),
		},
		Name:          site.Name(),
		Description:   site.Description(),
		Domain:        site.Domain().Value(),
		TitleTemplate: site.TitleTemplate(),
		Enabled:       site.IsEnabled(),
		TemplateID:    site.TemplateID().Value(),
		TenantID:      site.TenantID().Value(),
	}, nil
}

// ToDomain converts a GORM models.Site to a domain Site
func (m *SiteMapper) ToDomain(model *models.Site) (*entities.Site, error) {
	if model == nil {
		return nil, nil
	}

	domain, err := value_objects.NewDomainName(model.Domain)
	if err != nil {
		return nil, err
	}

	site, err := entities.NewSite(
		model.Name,
		model.Description,
		domain,
		entities.NewTemplateID(model.TemplateID),
		entities.NewTenantID(model.TenantID),
	)
	if err != nil {
		return nil, err
	}

	err = site.SetID(entities.NewSiteID(model.ID))
	if err != nil {
		return nil, err
	}
	site.SetTimestamps(model.CreatedAt, model.UpdatedAt)

	if model.TitleTemplate != nil {
		site.UpdateTitleTemplate(model.TitleTemplate)
	}

	if !model.Enabled {
		site.Disable()
	}

	return site, nil
}

// ToModels converts a slice of domain Sites to GORM models
func (m *SiteMapper) ToModels(sites []*entities.Site) ([]*models.Site, error) {
	if sites == nil {
		return nil, nil
	}

	result := make([]*models.Site, len(sites))
	for i, site := range sites {
		if site == nil {
			return nil, errors.ErrSiteEmpty
		}
		model, err := m.ToModel(site)
		if err != nil {
			return nil, err
		}
		result[i] = model
	}

	return result, nil
}

// ToDomains converts a slice of GORM models to domain Sites
func (m *SiteMapper) ToDomains(modelList []*models.Site) ([]*entities.Site, error) {
	if modelList == nil {
		return nil, nil
	}

	result := make([]*entities.Site, len(modelList))
	for i, model := range modelList {
		site, err := m.ToDomain(model)
		if err != nil {
			return nil, err
		}
		result[i] = site
	}

	return result, nil
}
