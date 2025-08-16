package repositories

import (
	"github.com/h4rdc0m/aurora-api/domain/entities"
)

// TemplateRepository defines the interface for template data operations
type TemplateRepository interface {
	Save(template *entities.Template) error
	FindByID(id entities.TemplateID) (*entities.Template, error)
	FindByName(name string) (*entities.Template, error)
	FindAll() ([]*entities.Template, error)
	FindEnabledOnly() ([]*entities.Template, error)
	Delete(id entities.TemplateID) error
	ExistsByName(name string) (bool, error)
	ExistsByFilePath(filePath string) (bool, error)
}
