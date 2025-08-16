package mappers

import (
	"github.com/h4rdc0m/aurora-api/domain/entities"
	"github.com/h4rdc0m/aurora-api/infrastructure/persistence/models"
)

type TemplateMapper struct{}

// NewTemplateMapper creates a new instance of TemplateMapper.
func NewTemplateMapper() *TemplateMapper {
	return &TemplateMapper{}
}

// ToModel converts a domain Template to a GORM models.Template.
func (m *TemplateMapper) ToModel(template *entities.Template) (*models.Template, error) {
	if template == nil {
		return nil, nil
	}

	model := &models.Template{
		Base: models.Base{
			ID:        template.ID().Value(),
			CreatedAt: template.CreatedAt(),
			UpdatedAt: template.UpdatedAt(),
		},
		Name:        template.Name(),
		Description: template.Description(),
		FilePath:    template.FilePath(),
		Enabled:     template.IsEnabled(),
	}

	return model, nil
}

// ToDomain converts a GORM models.Template to a domain Template entity.
func (m *TemplateMapper) ToDomain(model *models.Template) (*entities.Template, error) {
	if model == nil {
		return nil, nil
	}

	template, err := entities.NewTemplate(
		model.Name,
		model.FilePath,
		model.Description,
	)
	if err != nil {
		return nil, err
	}

	template.SetID(entities.NewTemplateID(model.ID))
	template.SetTimestamps(model.CreatedAt, model.UpdatedAt)
	if model.Enabled {
		template.Enable()
	} else {
		template.Disable()
	}
	return template, nil
}

func (m *TemplateMapper) ToModels(templates []*entities.Template) ([]*models.Template, error) {
	if templates == nil {
		return nil, nil
	}

	// Convert a slice of domain Template entities to a slice of GORM modelList.Template.
	modelList := make([]*models.Template, len(templates))
	for i, template := range templates {
		model, err := m.ToModel(template)
		if err != nil {
			return nil, err
		}
		modelList[i] = model
	}
	return modelList, nil
}

// ToDomains converts a slice of GORM models.Template to a slice of domain Template entities.
func (m *TemplateMapper) ToDomains(models []*models.Template) ([]*entities.Template, error) {
	if models == nil {
		return nil, nil
	}

	templates := make([]*entities.Template, len(models))
	for i, model := range models {
		template, err := m.ToDomain(model)
		if err != nil {
			return nil, err
		}
		templates[i] = template
	}
	return templates, nil
}
