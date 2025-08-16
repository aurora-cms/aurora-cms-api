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

// PageBlockRepositoryImpl implements PageBlockRepository using sqlx and squirrel
type PageBlockRepositoryImpl struct {
	db     common.Database
	logger common.Logger
	mapper common.Mapper[*entities.PageBlock, *models.PageBlock]
}

// NewPageBlockRepository creates a new PageBlockRepository implementation
func NewPageBlockRepository(db common.Database, logger common.Logger) repositories.PageBlockRepository {
	return &PageBlockRepositoryImpl{
		db:     db,
		logger: logger,
		mapper: mappers.NewPageBlockMapper(),
	}
}

// Save saves a page block (create or update)
func (r *PageBlockRepositoryImpl) Save(block *entities.PageBlock) error {
	model, err := r.mapper.ToModel(block)
	if err != nil {
		r.logger.Error("Failed to convert page block to model", "error", err)
		return err
	}

	if model.ID == 0 {
		query, args, err := squirrel.Insert("page_blocks").
			Columns("block_key", "page_version_id", "index", "content_type", "content").
			Values(model.BlockKey, model.PageVersionID, model.Index, model.ContentType, model.Content).
			PlaceholderFormat(squirrel.Question).
			ToSql()
		if err != nil {
			r.logger.Error("Failed to build insert query for page block", "error", err)
			return err
		}
		result, err := r.db.Exec(query, args...)
		if err != nil {
			r.logger.Error("Failed to insert new page block", "error", err)
			return err
		}
		id, err := result.LastInsertId()
		if err != nil {
			r.logger.Error("Failed to get last insert ID for page block", "error", err)
			return err
		}
		block.SetID(entities.NewPageBlockID(uint64(id)))
	} else {
		query, args, err := squirrel.Update("page_blocks").
			Set("block_key", model.BlockKey).
			Set("page_version_id", model.PageVersionID).
			Set("index", model.Index).
			Set("content_type", model.ContentType).
			Set("content", model.Content).
			Where(squirrel.Eq{"id": model.ID}).
			PlaceholderFormat(squirrel.Question).
			ToSql()
		if err != nil {
			r.logger.Error("Failed to build update query for page block", "error", err)
			return err
		}
		if _, err := r.db.Exec(query, args...); err != nil {
			r.logger.Error("Failed to update page block", "id", model.ID, "error", err)
			return err
		}
	}
	return nil
}

// FindByID retrieves a page block by ID
func (r *PageBlockRepositoryImpl) FindByID(id entities.PageBlockID) (*entities.PageBlock, error) {
	var model models.PageBlock
	query, args, err := squirrel.Select("*").From("page_blocks").Where(squirrel.Eq{"id": id.Value()}).ToSql()
	if err != nil {
		r.logger.Error("Failed to build select query for FindByID", "error", err)
		return nil, err
	}
	if err := r.db.Get(&model, query, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		r.logger.Error("Failed to find page block by ID", "id", id.Value(), "error", err)
		return nil, err
	}
	return r.mapper.ToDomain(&model)
}

// FindByPageVersionID retrieves all blocks for a specific page version
func (r *PageBlockRepositoryImpl) FindByPageVersionID(pageVersionID entities.PageVersionID) ([]*entities.PageBlock, error) {
	var modelList []*models.PageBlock
	query, args, err := squirrel.Select("*").From("page_blocks").Where(squirrel.Eq{"page_version_id": pageVersionID.Value()}).OrderBy("index ASC").ToSql()
	if err != nil {
		r.logger.Error("Failed to build select query for FindByPageVersionID", "error", err)
		return nil, err
	}
	if err := r.db.Select(&modelList, query, args...); err != nil {
		r.logger.Error("Failed to find page blocks by page version ID", "page_version_id", pageVersionID.Value(), "error", err)
		return nil, err
	}
	return r.mapper.ToDomains(modelList)
}

// FindByBlockKey retrieves a block by key and page version ID
func (r *PageBlockRepositoryImpl) FindByBlockKey(blockKey string, pageVersionID entities.PageVersionID) (*entities.PageBlock, error) {
	var model models.PageBlock
	query, args, err := squirrel.Select("*").From("page_blocks").Where(squirrel.Eq{"block_key": blockKey, "page_version_id": pageVersionID.Value()}).ToSql()
	if err != nil {
		r.logger.Error("Failed to build select query for FindByBlockKey", "error", err)
		return nil, err
	}
	if err := r.db.Get(&model, query, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		r.logger.Error("Failed to find page block by key", "block_key", blockKey, "page_version_id", pageVersionID.Value(), "error", err)
		return nil, err
	}
	return r.mapper.ToDomain(&model)
}

// Delete deletes a page block (soft delete)
func (r *PageBlockRepositoryImpl) Delete(id entities.PageBlockID) error {
	query, args, err := squirrel.Delete("page_blocks").Where(squirrel.Eq{"id": id.Value()}).ToSql()
	if err != nil {
		r.logger.Error("Failed to build delete query for page block", "id", id.Value(), "error", err)
		return err
	}
	if _, err := r.db.Exec(query, args...); err != nil {
		r.logger.Error("Failed to delete page block", "id", id.Value(), "error", err)
		return err
	}
	return nil
}

// DeleteByPageVersionID deletes all blocks for a page version
func (r *PageBlockRepositoryImpl) DeleteByPageVersionID(pageVersionID entities.PageVersionID) error {
	query, args, err := squirrel.Delete("page_blocks").Where(squirrel.Eq{"page_version_id": pageVersionID.Value()}).ToSql()
	if err != nil {
		r.logger.Error("Failed to build delete query for page blocks", "page_version_id", pageVersionID.Value(), "error", err)
		return err
	}
	if _, err := r.db.Exec(query, args...); err != nil {
		r.logger.Error("Failed to delete page blocks by page version ID", "page_version_id", pageVersionID.Value(), "error", err)
		return err
	}
	return nil
}
