package mappers

import (
	"github.com/h4rdc0m/aurora-api/domain/entities"
	"github.com/h4rdc0m/aurora-api/domain/errors"
	"github.com/h4rdc0m/aurora-api/infrastructure/persistence/models"
)

// PageVersionMapper handles conversion between domain entities and GORM models
type PageVersionMapper struct {
	blockMapper *PageBlockMapper
}

// NewPageVersionMapper creates a new PageVersionMapper
func NewPageVersionMapper() *PageVersionMapper {
	return &PageVersionMapper{
		blockMapper: NewPageBlockMapper(),
	}
}

// ToModel converts a domain PageVersion to a GORM models.PageVersion
func (m *PageVersionMapper) ToModel(version *entities.PageVersion) (*models.PageVersion, error) {
	if version == nil {
		return nil, nil
	}

	model := &models.PageVersion{
		Base: models.Base{
			ID:        version.ID().Value(),
			CreatedAt: version.CreatedAt(),
			UpdatedAt: version.UpdatedAt(),
		},
		PageID:      version.PageID().Value(),
		Version:     version.Version(),
		Title:       version.Title(),
		Description: version.Description(),
		IsPublished: version.IsPublished(),
	}

	return model, nil
}

// ToDomain converts a GORM models.PageVersion to a domain PageVersion
func (m *PageVersionMapper) ToDomain(model *models.PageVersion) (*entities.PageVersion, error) {
	if model == nil {
		return nil, nil
	}

	version, err := entities.NewPageVersion(
		entities.NewPageID(model.PageID),
		model.Version,
		model.Title,
		model.Description,
	)
	if err != nil {
		return nil, err
	}

	version.SetID(entities.NewPageVersionID(model.ID))
	version.SetTimestamps(model.CreatedAt, model.UpdatedAt)

	if model.IsPublished {
		version.Publish()
	}

	return version, nil
}

// ToModels converts a slice of domain PageVersions to GORM models
func (m *PageVersionMapper) ToModels(versions []*entities.PageVersion) ([]*models.PageVersion, error) {
	if versions == nil {
		return nil, nil
	}

	result := make([]*models.PageVersion, len(versions))
	for i, version := range versions {
		if version == nil {
			return nil, errors.ErrInvalidPageVersionModel
		}
		model, err := m.ToModel(version)
		if err != nil {
			return nil, err
		}
		result[i] = model
	}

	return result, nil
}

// ToDomains converts a slice of GORM models to domain PageVersions
func (m *PageVersionMapper) ToDomains(modelList []*models.PageVersion) ([]*entities.PageVersion, error) {
	if modelList == nil {
		return nil, nil
	}

	result := make([]*entities.PageVersion, len(modelList))
	for i, model := range modelList {
		if model == nil {
			return nil, errors.ErrInvalidPageVersionModel
		}
		version, err := m.ToDomain(model)
		if err != nil {
			return nil, err
		}
		result[i] = version
	}

	return result, nil
}
