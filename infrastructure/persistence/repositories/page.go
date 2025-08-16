package repositories

import (
	"database/sql"
	"errors"
	"github.com/Masterminds/squirrel"
	"github.com/h4rdc0m/aurora-api/domain/common"
	"github.com/h4rdc0m/aurora-api/domain/entities"
	"github.com/h4rdc0m/aurora-api/domain/repositories"
	"github.com/h4rdc0m/aurora-api/infrastructure/persistence/mappers"
	"github.com/h4rdc0m/aurora-api/infrastructure/persistence/models"
)

// PageRepositoryImpl provides the implementation of the PageRepository interface for interacting with page data.
// It uses GORM for database interactions and a PageMapper for entity-model conversion.
// The struct includes a logger for logging errors or operational information during repository operations.
type PageRepositoryImpl struct {
	db     common.Database
	logger common.Logger
	mapper common.Mapper[*entities.Page, *models.Page]
}

// NewPageRepository creates a new instance of PageRepository with the provided database connection and logger.
func NewPageRepository(db common.Database, logger common.Logger) repositories.PageRepository {
	return &PageRepositoryImpl{
		db:     db,
		logger: logger,
		mapper: mappers.NewPageMapper(),
	}
}

// Save saves a new page or updates an existing page in the database.
// Returns an error if the operation fails.
func (r *PageRepositoryImpl) Save(page *entities.Page) error {
	model, err := r.mapper.ToModel(page)
	if err != nil {
		r.logger.Error("Failed to map page entity to model", "error", err)
		return err
	}

	if model.ID == 0 {
		query, args, err := squirrel.Insert("pages").
			Columns("key", "path", "index", "site_id", "type", "link_url", "parent_id", "hard_link_page_id").
			Values(model.Key, model.Path, model.Index, model.SiteID, model.Type, model.LinkURL, model.ParentID, model.HardLinkPageID).
			PlaceholderFormat(squirrel.Question).
			ToSql()
		if err != nil {
			r.logger.Error("Failed to build insert query", "error", err)
			return err
		}
		result, err := r.db.Exec(query, args...)
		if err != nil {
			r.logger.Error("Failed to insert new page", "error", err)
			return err
		}
		id, err := result.LastInsertId()
		if err != nil {
			r.logger.Error("Failed to get last insert ID", "error", err)
			return err
		}
		page.SetID(entities.NewPageID(uint64(id)))
	} else {
		query, args, err := squirrel.Update("pages").
			Set("key", model.Key).
			Set("path", model.Path).
			Set("index", model.Index).
			Set("type", model.Type).
			Set("link_url", model.LinkURL).
			Set("parent_id", model.ParentID).
			Set("hard_link_page_id", model.HardLinkPageID).
			Where(squirrel.Eq{"id": model.ID}).
			PlaceholderFormat(squirrel.Question).
			ToSql()
		if err != nil {
			r.logger.Error("Failed to build update query", "error", err)
			return err
		}
		if _, err := r.db.Exec(query, args...); err != nil {
			r.logger.Error("Failed to update existing page", "id", model.ID, "error", err)
			return err
		}
	}
	return nil
}

// FindByID retrieves a page entity by its unique identifier.
// Returns the page or nil if not found, and an error if a failure occurs during the operation.
func (r *PageRepositoryImpl) FindByID(id entities.PageID) (*entities.Page, error) {
	var model models.Page
	query, args, err := squirrel.Select("*").From("pages").Where(squirrel.Eq{"id": id.Value()}).ToSql()
	if err != nil {
		r.logger.Error("Failed to build select query for FindByID", "error", err)
		return nil, err
	}
	if err := r.db.Get(&model, query, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			r.logger.Warn("Page not found", "id", id.Value())
			return nil, nil
		}
		r.logger.Error("Failed to find page by ID", "id", id.Value(), "error", err)
		return nil, err
	}
	return r.mapper.ToDomain(&model)
}

// FindByPath retrieves a page by its path and associated site ID from the database. Returns nil if no record is found.
func (r *PageRepositoryImpl) FindByPath(path string, siteID entities.SiteID) (*entities.Page, error) {
	var model models.Page
	query, args, err := squirrel.Select("*").From("pages").Where(squirrel.Eq{"path": path, "site_id": siteID.Value()}).ToSql()
	if err != nil {
		r.logger.Error("Failed to build select query for FindByPath", "error", err)
		return nil, err
	}
	if err := r.db.Get(&model, query, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			r.logger.Warn("Page not found by path", "path", path, "siteID", siteID)
			return nil, nil
		}
		r.logger.Error("Failed to find page by path", "path", path, "siteID", siteID, "error", err)
		return nil, err
	}
	return r.mapper.ToDomain(&model)
}

// FindBySiteID retrieves a list of pages associated with the given site ID, ordered by their index.
func (r *PageRepositoryImpl) FindBySiteID(siteID entities.SiteID) ([]*entities.Page, error) {
	var modelList []*models.Page
	query, args, err := squirrel.Select("*").From("pages").Where(squirrel.Eq{"site_id": siteID.Value()}).OrderBy("index ASC").ToSql()
	if err != nil {
		r.logger.Error("Failed to build select query for FindBySiteID", "error", err)
		return nil, err
	}
	if err := r.db.Select(&modelList, query, args...); err != nil {
		r.logger.Error("Failed to find pages by site ID", "siteID", siteID.Value(), "error", err)
		return nil, err
	}
	return r.mapper.ToDomains(modelList)
}

// FindRootPagesBySiteID retrieves root pages by site ID where parent ID is null, ordering them by index in ascending order.
func (r *PageRepositoryImpl) FindRootPagesBySiteID(siteID entities.SiteID) ([]*entities.Page, error) {
	var modelList []*models.Page
	query, args, err := squirrel.Select("*").From("pages").Where(squirrel.And{squirrel.Eq{"site_id": siteID.Value()}, squirrel.Expr("parent_id IS NULL")}).OrderBy("index ASC").ToSql()
	if err != nil {
		r.logger.Error("Failed to build select query for FindRootPagesBySiteID", "error", err)
		return nil, err
	}
	if err := r.db.Select(&modelList, query, args...); err != nil {
		r.logger.Error("Failed to find root pages by site ID", "siteID", siteID.Value(), "error", err)
		return nil, err
	}
	return r.mapper.ToDomains(modelList)
}

// FindChildrenByParentID retrieves all child pages associated with the given parent page ID, ordered by their index.
func (r *PageRepositoryImpl) FindChildrenByParentID(parentID entities.PageID) ([]*entities.Page, error) {
	var modelList []*models.Page
	query, args, err := squirrel.Select("*").From("pages").Where(squirrel.Eq{"parent_id": parentID.Value()}).OrderBy("index ASC").ToSql()
	if err != nil {
		r.logger.Error("Failed to build select query for FindChildrenByParentID", "error", err)
		return nil, err
	}
	if err := r.db.Select(&modelList, query, args...); err != nil {
		r.logger.Error("Failed to find children pages by parent ID", "parentID", parentID.Value(), "error", err)
		return nil, err
	}
	return r.mapper.ToDomains(modelList)
}

// Delete removes a page from the database using its unique identifier and returns an error if the operation fails.
func (r *PageRepositoryImpl) Delete(id entities.PageID) error {
	query, args, err := squirrel.Delete("pages").Where(squirrel.Eq{"id": id.Value()}).ToSql()
	if err != nil {
		r.logger.Error("Failed to build delete query", "id", id.Value(), "error", err)
		return err
	}
	if _, err := r.db.Exec(query, args...); err != nil {
		r.logger.Error("Failed to delete page", "id", id.Value(), "error", err)
		return err
	}
	return nil
}

// ExistsByPath checks if a page with the given path and site ID exists in the repository, returning a boolean result.
func (r *PageRepositoryImpl) ExistsByPath(path string, siteID entities.SiteID) (bool, error) {
	var count int64
	query, args, err := squirrel.Select("COUNT(*)").From("pages").Where(squirrel.Eq{"path": path, "site_id": siteID.Value()}).ToSql()
	if err != nil {
		r.logger.Error("Failed to build count query for ExistsByPath", "error", err)
		return false, err
	}
	if err := r.db.Get(&count, query, args...); err != nil {
		r.logger.Error("Failed to check if page exists by path", "path", path, "siteID", siteID.Value(), "error", err)
		return false, err
	}
	return count > 0, nil
}
