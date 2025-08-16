package mappers

import (
	"github.com/h4rdc0m/aurora-api/domain/common"
	"github.com/h4rdc0m/aurora-api/domain/entities"
	"github.com/h4rdc0m/aurora-api/domain/errors"
	"github.com/h4rdc0m/aurora-api/domain/value_objects"
	"github.com/h4rdc0m/aurora-api/infrastructure/persistence/models"
)

// PageMapper handles conversion between domain entities and GORM models
type PageMapper struct {
	pageTypeMapper common.Mapper[entities.PageType, models.PageType]
}

// NewPageMapper creates a new PageMapper
func NewPageMapper() common.Mapper[*entities.Page, *models.Page] {
	return &PageMapper{
		pageTypeMapper: NewPageTypeMapper(),
	}
}

// ToModel converts a entities.Page to models.Page
func (m *PageMapper) ToModel(page *entities.Page) (*models.Page, error) {
	if page == nil {
		return nil, nil
	}

	ptype, err := m.pageTypeMapper.ToModel(page.Type())
	if err != nil {
		return nil, err
	}
	model := &models.Page{
		Base: models.Base{
			ID:        page.ID().Value(),
			CreatedAt: page.CreatedAt(),
			UpdatedAt: page.UpdatedAt(),
		},
		Key:     page.Key().Value(),
		Path:    page.Path(),
		Index:   page.Index(),
		SiteID:  page.SiteID().Value(),
		Type:    ptype,
		LinkURL: page.LinkURL(),
	}

	if page.ParentID() != nil {
		parentID := page.ParentID().Value()
		model.ParentID = &parentID
	}

	if page.HardLinkPageID() != nil {
		hardLinkPageID := page.HardLinkPageID().Value()
		model.HardLinkPageID = &hardLinkPageID
	}

	return model, nil
}

// ToDomain converts a models.Page to a entities.Page
func (m *PageMapper) ToDomain(model *models.Page) (*entities.Page, error) {
	if model == nil {
		return nil, nil
	}

	key, err := value_objects.NewPageKey(model.Key)
	if err != nil {
		return nil, err
	}

	ptype, err := m.pageTypeMapper.ToDomain(model.Type)
	if err != nil {
		return nil, err
	}
	page, err := entities.NewPage(
		key,
		model.Path,
		entities.NewSiteID(model.SiteID),
		ptype,
	)
	if err != nil {
		return nil, err
	}

	page.SetID(entities.NewPageID(model.ID))
	page.SetTimestamps(model.CreatedAt, model.UpdatedAt)
	page.UpdateIndex(model.Index)

	if model.ParentID != nil {
		parentID := entities.NewPageID(*model.ParentID)
		page.SetParent(&parentID)
	}

	if model.LinkURL != nil {
		err := page.SetLinkURL(model.LinkURL)
		if err != nil {
			return nil, err
		}
	}

	if model.HardLinkPageID != nil {
		hardLinkPageID := entities.NewPageID(*model.HardLinkPageID)
		err := page.SetHardLinkPageID(&hardLinkPageID)
		if err != nil {
			return nil, err
		}
	}

	return page, nil
}

// ToModels converts []entities.Page to []models.Page
func (m *PageMapper) ToModels(pages []*entities.Page) ([]*models.Page, error) {
	if pages == nil {
		return nil, nil
	}

	result := make([]*models.Page, len(pages))
	for i, page := range pages {
		model, err := m.ToModel(page)
		if err != nil {
			return nil, err
		}
		result[i] = model
	}

	return result, nil
}

// ToDomains converts []models.Page to []entities.Page
func (m *PageMapper) ToDomains(modelList []*models.Page) ([]*entities.Page, error) {
	if modelList == nil {
		return nil, nil
	}

	result := make([]*entities.Page, len(modelList))
	for i, model := range modelList {
		page, err := m.ToDomain(model)
		if err != nil {
			return nil, err
		}
		result[i] = page
	}

	return result, nil
}

type PageTypeMapper struct {
}

func (p PageTypeMapper) ToModel(entity entities.PageType) (models.PageType, error) {
	switch entity {
	case entities.PageTypeContent:
		return models.PageTypeContent, nil
	case entities.PageTypeLink:
		return models.PageTypeLink, nil
	case entities.PageTypeHardLink:
		return models.PageTypeHardLink, nil
	case entities.PageTypeSnippet:
		return models.PageTypeSnippet, nil
	default:
		return "", errors.ErrInvalidPageType
	}
}

func (p PageTypeMapper) ToDomain(model models.PageType) (entities.PageType, error) {
	switch model {
	case models.PageTypeContent:
		return entities.PageTypeContent, nil
	case models.PageTypeLink:
		return entities.PageTypeLink, nil
	case models.PageTypeHardLink:
		return entities.PageTypeHardLink, nil
	case models.PageTypeSnippet:
		return entities.PageTypeSnippet, nil
	default:
		return "", errors.ErrInvalidPageType
	}
}

func (p PageTypeMapper) ToDomains(models []models.PageType) ([]entities.PageType, error) {
	if models == nil {
		return nil, nil
	}

	result := make([]entities.PageType, len(models))
	for i, model := range models {
		entity, err := p.ToDomain(model)
		if err != nil {
			return nil, err
		}
		result[i] = entity
	}

	return result, nil
}

func (p PageTypeMapper) ToModels(entities []entities.PageType) ([]models.PageType, error) {
	if entities == nil {
		return nil, nil
	}

	result := make([]models.PageType, len(entities))
	for i, entity := range entities {
		model, err := p.ToModel(entity)
		if err != nil {
			return nil, err
		}
		result[i] = model
	}

	return result, nil
}

func NewPageTypeMapper() common.Mapper[entities.PageType, models.PageType] {
	return &PageTypeMapper{}
}
