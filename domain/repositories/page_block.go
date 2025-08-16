package repositories

import (
	"github.com/h4rdc0m/aurora-api/domain/entities"
)

// PageBlockRepository defines the interface for page block data operations
type PageBlockRepository interface {
	Save(block *entities.PageBlock) error
	FindByID(id entities.PageBlockID) (*entities.PageBlock, error)
	FindByPageVersionID(pageVersionID entities.PageVersionID) ([]*entities.PageBlock, error)
	FindByBlockKey(blockKey string, pageVersionID entities.PageVersionID) (*entities.PageBlock, error)
	Delete(id entities.PageBlockID) error
	DeleteByPageVersionID(pageVersionID entities.PageVersionID) error
}
