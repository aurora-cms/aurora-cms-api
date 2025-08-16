package repositories

import (
	"github.com/h4rdc0m/aurora-api/domain/entities"
)

// PageVersionRepository defines the interface for page version data operations
type PageVersionRepository interface {
	Save(version *entities.PageVersion) error
	FindByID(id entities.PageVersionID) (*entities.PageVersion, error)
	FindByPageID(pageID entities.PageID) ([]*entities.PageVersion, error)
	FindPublishedByPageID(pageID entities.PageID) (*entities.PageVersion, error)
	FindLatestByPageID(pageID entities.PageID) (*entities.PageVersion, error)
	Delete(id entities.PageVersionID) error
}
