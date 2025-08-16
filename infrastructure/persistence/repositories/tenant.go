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

type TenantRepositoryImpl struct {
	db     common.Database
	logger common.Logger
	mapper common.Mapper[*entities.Tenant, *models.Tenant]
}

func NewTenantRepository(db common.Database, logger common.Logger) repositories.TenantRepository {
	return &TenantRepositoryImpl{
		db:     db,
		logger: logger,
		mapper: mappers.NewTenantMapper(),
	}
}

func (r *TenantRepositoryImpl) Save(tenant *entities.Tenant) error {
	model, err := r.mapper.ToModel(tenant)
	if err != nil {
		r.logger.Error("Failed to map tenant to model", "error", err)
		return err
	}

	if model.ID == 0 {
		query, args, err := squirrel.Insert("tenants").
			Columns("name", "created_at", "updated_at").
			Values(model.Name, model.CreatedAt, model.UpdatedAt).
			PlaceholderFormat(squirrel.Question).
			ToSql()
		if err != nil {
			r.logger.Error("Failed to build insert query for tenant", "error", err)
			return err
		}
		result, err := r.db.Exec(query, args...)
		if err != nil {
			r.logger.Error("Failed to create tenant", "error", err)
			return err
		}
		id, err := result.LastInsertId()
		if err != nil {
			r.logger.Error("Failed to get last insert ID for tenant", "error", err)
			return err
		}
		tenant.SetID(entities.NewTenantID(uint64(id)))
	} else {
		query, args, err := squirrel.Update("tenants").
			Set("name", model.Name).
			Set("created_at", model.CreatedAt).
			Set("updated_at", model.UpdatedAt).
			Where(squirrel.Eq{"id": model.ID}).
			PlaceholderFormat(squirrel.Question).
			ToSql()
		if err != nil {
			r.logger.Error("Failed to build update query for tenant", "error", err)
			return err
		}
		if _, err := r.db.Exec(query, args...); err != nil {
			r.logger.Error("Failed to update tenant", "id", model.ID, "error", err)
			return err
		}
	}
	return nil
}

func (r *TenantRepositoryImpl) FindByID(id entities.TenantID) (*entities.Tenant, error) {
	var model models.Tenant
	query, args, err := squirrel.Select("*").From("tenants").Where(squirrel.Eq{"id": id.Value()}).ToSql()
	if err != nil {
		r.logger.Error("Failed to build select query for FindByID", "error", err)
		return nil, err
	}
	if err := r.db.Get(&model, query, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		r.logger.Error("Failed to find tenant by ID", "id", id.Value(), "error", err)
		return nil, err
	}
	return r.mapper.ToDomain(&model)
}

func (r *TenantRepositoryImpl) FindByName(name string) (*entities.Tenant, error) {
	var model models.Tenant
	query, args, err := squirrel.Select("*").From("tenants").Where(squirrel.Eq{"name": name}).ToSql()
	if err != nil {
		r.logger.Error("Failed to build select query for FindByName", "error", err)
		return nil, err
	}

	if err := r.db.Get(&model, query, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			r.logger.Error("Failed to find tenant by name", "name", name, "error", err)
			return nil, nil
		}
		r.logger.Error("Failed to find tenant by name", "name", name, "error", err)
		return nil, err
	}

	if model.ID == 0 {
		r.logger.Error("Tenant found by name has no ID", "name", name)
		return nil, errors.New("tenant not found")
	}

	return r.mapper.ToDomain(&model)
}

func (r *TenantRepositoryImpl) FindAll() ([]*entities.Tenant, error) {
	var modelList []*models.Tenant
	query, args, err := squirrel.Select("*").From("tenants").ToSql()
	if err != nil {
		r.logger.Error("Failed to build select query for FindAll", "error", err)
		return nil, err
	}
	if err := r.db.Select(&modelList, query, args...); err != nil {
		r.logger.Error("Failed to find all tenants", "error", err)
		return nil, err
	}
	return r.mapper.ToDomains(modelList)
}

func (r *TenantRepositoryImpl) FindActiveOnly() ([]*entities.Tenant, error) {
	var modelList []*models.Tenant
	query, args, err := squirrel.Select("*").From("tenants").Where(squirrel.Eq{"active": true}).ToSql()
	if err != nil {
		r.logger.Error("Failed to build select query for FindActiveOnly", "error", err)
		return nil, err
	}

	if err := r.db.Select(&modelList, query, args...); err != nil {
		r.logger.Error("Failed to find active tenants", "error", err)
		return nil, err
	}
	if len(modelList) == 0 {
		r.logger.Info("No active tenants found")
		return nil, nil
	}

	return r.mapper.ToDomains(modelList)
}

func (r *TenantRepositoryImpl) Delete(id entities.TenantID) error {
	query, args, err := squirrel.Delete("tenants").Where(squirrel.Eq{"id": id.Value()}).ToSql()
	if err != nil {
		r.logger.Error("Failed to build delete query for tenant", "id", id.Value(), "error", err)
		return err
	}
	if _, err := r.db.Exec(query, args...); err != nil {
		r.logger.Error("Failed to delete tenant", "id", id.Value(), "error", err)
		return err
	}
	return nil
}

func (r *TenantRepositoryImpl) ExistsByName(name string) (bool, error) {
	var count int64
	query, args, err := squirrel.Select("COUNT(*)").From("tenants").Where(squirrel.Eq{"name": name}).ToSql()
	if err != nil {
		r.logger.Error("Failed to build count query for ExistsByName", "name", name, "error", err)
		return false, err
	}
	if err := r.db.Get(&count, query, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		r.logger.Error("Failed to check if tenant exists by name", "name", name, "error", err)
		return false, err
	}
	return count > 0, nil
}
