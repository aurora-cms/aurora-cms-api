package repositories

import (
	"github.com/h4rdc0m/aurora-api/domain/entities"
)

// PageRepository defines the interface for page data operations
type PageRepository interface {
	Save(page *entities.Page) error
	FindByID(id entities.PageID) (*entities.Page, error)
	FindByPath(path string, siteID entities.SiteID) (*entities.Page, error)
	FindBySiteID(siteID entities.SiteID) ([]*entities.Page, error)
	FindRootPagesBySiteID(siteID entities.SiteID) ([]*entities.Page, error)
	FindChildrenByParentID(parentID entities.PageID) ([]*entities.Page, error)
	Delete(id entities.PageID) error
	ExistsByPath(path string, siteID entities.SiteID) (bool, error)
}
