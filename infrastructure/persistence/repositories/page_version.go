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

// PageVersionRepositoryImpl implements PageVersionRepository using sqlx and squirrel
type PageVersionRepositoryImpl struct {
	db     common.Database
	logger common.Logger
	mapper common.Mapper[*entities.PageVersion, *models.PageVersion]
}

// NewPageVersionRepository creates a new PageVersionRepository implementation
func NewPageVersionRepository(db common.Database, logger common.Logger) repositories.PageVersionRepository {
	return &PageVersionRepositoryImpl{
		db:     db,
		logger: logger,
		mapper: mappers.NewPageVersionMapper(),
	}
}

// Save saves a page version (create or update)
func (r *PageVersionRepositoryImpl) Save(version *entities.PageVersion) error {
	model, err := r.mapper.ToModel(version)
	if err != nil {
		r.logger.Error("Failed to convert page version to model", "error", err)
		return err
	}

	if model.ID == 0 {
		query, args, err := squirrel.Insert("page_versions").
			Columns("page_id", "version", "is_published").
			Values(model.PageID, model.Version, model.IsPublished).
			PlaceholderFormat(squirrel.Question).
			ToSql()
		if err != nil {
			r.logger.Error("Failed to build insert query for page version", "error", err)
			return err
		}
		result, err := r.db.Exec(query, args...)
		if err != nil {
			r.logger.Error("Failed to create page version", "error", err)
			return err
		}
		id, err := result.LastInsertId()
		if err != nil {
			r.logger.Error("Failed to get last insert ID for page version", "error", err)
			return err
		}
		version.SetID(entities.NewPageVersionID(uint64(id)))
	} else {
		query, args, err := squirrel.Update("page_versions").
			Set("page_id", model.PageID).
			Set("version", model.Version).
			Set("is_published", model.IsPublished).
			Set("created_at", model.CreatedAt).
			Where(squirrel.Eq{"id": model.ID}).
			PlaceholderFormat(squirrel.Question).
			ToSql()
		if err != nil {
			r.logger.Error("Failed to build update query for page version", "error", err)
			return err
		}
		if _, err := r.db.Exec(query, args...); err != nil {
			r.logger.Error("Failed to update page version", "id", model.ID, "error", err)
			return err
		}
	}
	return nil
}

// FindByID retrieves a page version by ID
func (r *PageVersionRepositoryImpl) FindByID(id entities.PageVersionID) (*entities.PageVersion, error) {
	var model models.PageVersion
	query, args, err := squirrel.Select("*").From("page_versions").Where(squirrel.Eq{"id": id.Value()}).ToSql()
	if err != nil {
		r.logger.Error("Failed to build select query for FindByID", "error", err)
		return nil, err
	}
	if err := r.db.Get(&model, query, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		r.logger.Error("Failed to find page version by ID", "id", id.Value(), "error", err)
		return nil, err
	}
	return r.mapper.ToDomain(&model)
}

// FindByPageID retrieves all versions for a specific page
func (r *PageVersionRepositoryImpl) FindByPageID(pageID entities.PageID) ([]*entities.PageVersion, error) {
	var modelList []*models.PageVersion
	query, args, err := squirrel.Select("*").From("page_versions").Where(squirrel.Eq{"page_id": pageID.Value()}).OrderBy("version DESC").ToSql()
	if err != nil {
		r.logger.Error("Failed to build select query for FindByPageID", "error", err)
		return nil, err
	}
	if err := r.db.Select(&modelList, query, args...); err != nil {
		r.logger.Error("Failed to find page versions by page ID", "page_id", pageID.Value(), "error", err)
		return nil, err
	}
	return r.mapper.ToDomains(modelList)
}

// FindPublishedByPageID retrieves the published version for a page
func (r *PageVersionRepositoryImpl) FindPublishedByPageID(pageID entities.PageID) (*entities.PageVersion, error) {
	var model models.PageVersion
	query, args, err := squirrel.Select("*").From("page_versions").Where(squirrel.Eq{"page_id": pageID.Value(), "is_published": true}).ToSql()
	if err != nil {
		r.logger.Error("Failed to build select query for FindPublishedByPageID", "error", err)
		return nil, err
	}
	if err := r.db.Get(&model, query, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		r.logger.Error("Failed to find published page version", "page_id", pageID.Value(), "error", err)
		return nil, err
	}
	return r.mapper.ToDomain(&model)
}

// FindLatestByPageID retrieves the latest version for a page
func (r *PageVersionRepositoryImpl) FindLatestByPageID(pageID entities.PageID) (*entities.PageVersion, error) {
	var model models.PageVersion
	query, args, err := squirrel.Select("*").From("page_versions").Where(squirrel.Eq{"page_id": pageID.Value()}).OrderBy("version DESC").Limit(1).ToSql()
	if err != nil {
		r.logger.Error("Failed to build select query for FindLatestByPageID", "error", err)
		return nil, err
	}
	if err := r.db.Get(&model, query, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		r.logger.Error("Failed to find latest page version", "page_id", pageID.Value(), "error", err)
		return nil, err
	}
	return r.mapper.ToDomain(&model)
}

// Delete deletes a page version (soft delete)
func (r *PageVersionRepositoryImpl) Delete(id entities.PageVersionID) error {
	query, args, err := squirrel.Delete("page_versions").Where(squirrel.Eq{"id": id.Value()}).ToSql()
	if err != nil {
		r.logger.Error("Failed to build delete query for page version", "id", id.Value(), "error", err)
		return err
	}
	if _, err := r.db.Exec(query, args...); err != nil {
		r.logger.Error("Failed to delete page version", "id", id.Value(), "error", err)
		return err
	}
	return nil
}
