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

// TemplateRepositoryImpl implements TemplateRepository using sqlx and squirrel
type TemplateRepositoryImpl struct {
	db     common.Database
	logger common.Logger
	mapper common.Mapper[*entities.Template, *models.Template]
}

// NewTemplateRepository creates a new TemplateRepository implementation
func NewTemplateRepository(db common.Database, logger common.Logger) repositories.TemplateRepository {
	return &TemplateRepositoryImpl{
		db:     db,
		logger: logger,
		mapper: mappers.NewTemplateMapper(),
	}
}

// Save saves a template (create or update)
func (r *TemplateRepositoryImpl) Save(template *entities.Template) error {
	model, err := r.mapper.ToModel(template)
	if err != nil {
		r.logger.Error("Failed to convert template to model", "error", err)
		return err
	}

	if model.ID == 0 {
		query, args, err := squirrel.Insert("templates").
			Columns("name", "description", "file_path", "created_at", "updated_at").
			Values(model.Name, model.Description, model.FilePath, model.CreatedAt, model.UpdatedAt).
			PlaceholderFormat(squirrel.Question).
			ToSql()
		if err != nil {
			r.logger.Error("Failed to build insert query for template", "error", err)
			return err
		}
		result, err := r.db.Exec(query, args...)
		if err != nil {
			r.logger.Error("Failed to create template", "error", err)
			return err
		}
		id, err := result.LastInsertId()
		if err != nil {
			r.logger.Error("Failed to get last insert ID for template", "error", err)
			return err
		}
		template.SetID(entities.NewTemplateID(uint64(id)))

	} else {
		query, args, err := squirrel.Update("templates").
			Set("name", model.Name).
			Set("description", model.Description).
			Set("file_path", model.FilePath).
			Set("created_at", model.CreatedAt).
			Set("updated_at", model.UpdatedAt).
			Where(squirrel.Eq{"id": model.ID}).
			PlaceholderFormat(squirrel.Question).
			ToSql()
		if err != nil {
			r.logger.Error("Failed to build update query for template", "error", err)
			return err
		}
		if _, err := r.db.Exec(query, args...); err != nil {
			r.logger.Error("Failed to update template", "id", model.ID, "error", err)
			return err
		}
	}
	return nil
}

// FindByID retrieves a template by ID
func (r *TemplateRepositoryImpl) FindByID(id entities.TemplateID) (*entities.Template, error) {
	var model models.Template
	query, args, err := squirrel.Select("*").From("templates").Where(squirrel.Eq{"id": id.Value()}).ToSql()
	if err != nil {
		r.logger.Error("Failed to build select query for FindByID", "error", err)
		return nil, err
	}
	if err := r.db.Get(&model, query, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		r.logger.Error("Failed to find template by ID", "id", id.Value(), "error", err)
		return nil, err
	}
	return r.mapper.ToDomain(&model)
}

// FindByName retrieves a template by name
func (r *TemplateRepositoryImpl) FindByName(name string) (*entities.Template, error) {
	var model models.Template
	query, args, err := squirrel.Select("*").From("templates").Where(squirrel.Eq{"name": name}).ToSql()
	if err != nil {
		r.logger.Error("Failed to build select query for FindByName", "error", err)
		return nil, err
	}
	if err := r.db.Get(&model, query, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		r.logger.Error("Failed to find template by name", "name", name, "error", err)
		return nil, err
	}

	return r.mapper.ToDomain(&model)
}

// FindAll retrieves all templates
func (r *TemplateRepositoryImpl) FindAll() ([]*entities.Template, error) {
	var modelList []*models.Template
	query, args, err := squirrel.Select("*").From("templates").ToSql()
	if err != nil {
		r.logger.Error("Failed to build select query for FindAll", "error", err)
		return nil, err
	}
	if err := r.db.Select(&modelList, query, args...); err != nil {
		r.logger.Error("Failed to find all templates", "error", err)
		return nil, err
	}
	return r.mapper.ToDomains(modelList)
}

// FindEnabledOnly retrieves only enabled templates
func (r *TemplateRepositoryImpl) FindEnabledOnly() ([]*entities.Template, error) {
	var modelList []*models.Template
	query, args, err := squirrel.Select("*").From("templates").Where(squirrel.Eq{"enabled": true}).ToSql()
	if err != nil {
		r.logger.Error("Failed to build select query for FindEnabledOnly", "error", err)
		return nil, err
	}
	if err := r.db.Select(&modelList, query, args...); err != nil {
		r.logger.Error("Failed to find enabled templates", "error", err)
		return nil, err
	}

	if len(modelList) == 0 {
		r.logger.Info("No enabled templates found")
		return nil, nil
	}

	return r.mapper.ToDomains(modelList)
}

// Delete deletes a template
func (r *TemplateRepositoryImpl) Delete(id entities.TemplateID) error {
	query, args, err := squirrel.Delete("templates").Where(squirrel.Eq{"id": id.Value()}).ToSql()
	if err != nil {
		r.logger.Error("Failed to build delete query for template", "id", id.Value(), "error", err)
		return err
	}
	if _, err := r.db.Exec(query, args...); err != nil {
		r.logger.Error("Failed to delete template", "id", id.Value(), "error", err)
		return err
	}
	return nil
}

// ExistsByName checks if a template with the given name exists
func (r *TemplateRepositoryImpl) ExistsByName(name string) (bool, error) {
	var count int64
	// Use squirrel to build the count query
	query, args, err := squirrel.Select("COUNT(*)").From("templates").Where(squirrel.Eq{"name": name}).ToSql()
	if err != nil {
		r.logger.Error("Failed to build count query for template existence by name", "name", name, "error", err)
		return false, err
	}

	if err := r.db.Get(&count, query, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil // No rows means the template does not exist
		}
		r.logger.Error("Failed to check template existence by name", "name", name, "error", err)
		return false, err
	}
	return count > 0, nil
}

// ExistsByFilePath checks if a template with the given file path exists
func (r *TemplateRepositoryImpl) ExistsByFilePath(filePath string) (bool, error) {
	var count int64
	// Use squirrel to build the count query
	query, args, err := squirrel.Select("COUNT(*)").From("templates").Where(squirrel.Eq{"file_path": filePath}).ToSql()
	if err != nil {
		r.logger.Error("Failed to build count query for template existence by file path", "file_path", filePath, "error", err)
		return false, err
	}

	if err := r.db.Get(&count, query, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil // No rows means the template does not exist
		}
		r.logger.Error("Failed to check template existence by file path", "file_path", filePath, "error", err)
		return false, err
	}

	return count > 0, nil
}
