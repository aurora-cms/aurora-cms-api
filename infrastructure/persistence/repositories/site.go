package repositories

import (
	"database/sql"
	"errors"
	"github.com/Masterminds/squirrel"
	"github.com/h4rdc0m/aurora-api/domain/common"
	"github.com/h4rdc0m/aurora-api/domain/entities"
	"github.com/h4rdc0m/aurora-api/domain/repositories"
	"github.com/h4rdc0m/aurora-api/domain/value_objects"
	"github.com/h4rdc0m/aurora-api/infrastructure/persistence/mappers"
	"github.com/h4rdc0m/aurora-api/infrastructure/persistence/models"
)

// SiteRepositoryImpl implements SiteRepository using sqlx and squirrel
type SiteRepositoryImpl struct {
	db     common.Database
	logger common.Logger
	mapper common.Mapper[*entities.Site, *models.Site]
}

// NewSiteRepository creates a new SiteRepository implementation
func NewSiteRepository(db common.Database, logger common.Logger) repositories.SiteRepository {
	return &SiteRepositoryImpl{
		db:     db,
		logger: logger,
		mapper: mappers.NewSiteMapper(),
	}
}

// Save saves a site (create or update)
func (r *SiteRepositoryImpl) Save(site *entities.Site) error {
	model, err := r.mapper.ToModel(site)
	if err != nil {
		r.logger.Error("Failed to convert site to model", "error", err)
		return err
	}

	if model.ID == 0 {
		query, args, err := squirrel.Insert("sites").
			Columns("domain", "name", "tenant_id", "enabled", "created_at", "updated_at").
			Values(model.Domain, model.Name, model.TenantID, model.Enabled, model.CreatedAt, model.UpdatedAt).
			PlaceholderFormat(squirrel.Question).
			ToSql()
		if err != nil {
			r.logger.Error("Failed to build insert query for site", "error", err)
			return err
		}
		result, err := r.db.Exec(query, args...)
		if err != nil {
			r.logger.Error("Failed to create site", "error", err)
			return err
		}
		id, err := result.LastInsertId()
		if err != nil {
			r.logger.Error("Failed to get last insert ID for site", "error", err)
			return err
		}
		err = site.SetID(entities.NewSiteID(uint64(id)))
		if err != nil {
			return err
		}
	} else {
		query, args, err := squirrel.Update("sites").
			Set("domain", model.Domain).
			Set("name", model.Name).
			Set("tenant_id", model.TenantID).
			Set("enabled", model.Enabled).
			Set("created_at", model.CreatedAt).
			Set("updated_at", model.UpdatedAt).
			Where(squirrel.Eq{"id": model.ID}).
			PlaceholderFormat(squirrel.Question).
			ToSql()
		if err != nil {
			r.logger.Error("Failed to build update query for site", "error", err)
			return err
		}
		if _, err := r.db.Exec(query, args...); err != nil {
			r.logger.Error("Failed to update site", "id", model.ID, "error", err)
			return err
		}
	}
	return nil
}

// FindByID retrieves a site by ID
func (r *SiteRepositoryImpl) FindByID(id entities.SiteID) (*entities.Site, error) {
	var model models.Site
	query, args, err := squirrel.Select("*").From("sites").Where(squirrel.Eq{"id": id.Value()}).ToSql()
	if err != nil {
		r.logger.Error("Failed to build select query for FindByID", "error", err)
		return nil, err
	}
	if err := r.db.Get(&model, query, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		r.logger.Error("Failed to find site by ID", "id", id.Value(), "error", err)
		return nil, err
	}
	return r.mapper.ToDomain(&model)
}

// FindByDomain retrieves a site by domain
func (r *SiteRepositoryImpl) FindByDomain(domain *value_objects.DomainName) (*entities.Site, error) {
	var model models.Site
	query, args, err := squirrel.Select("*").From("sites").Where(squirrel.Eq{"domain": domain.Value()}).ToSql()
	if err != nil {
		r.logger.Error("Failed to build select query for FindByDomain", "error", err)
		return nil, err
	}
	if err := r.db.Get(&model, query, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		r.logger.Error("Failed to find site by domain", "domain", domain.Value(), "error", err)
		return nil, err
	}
	return r.mapper.ToDomain(&model)
}

// FindByTenantID retrieves all sites for a specific tenant
func (r *SiteRepositoryImpl) FindByTenantID(tenantID entities.TenantID) ([]*entities.Site, error) {
	var modelList []*models.Site
	query, args, err := squirrel.Select("*").From("sites").Where(squirrel.Eq{"tenant_id": tenantID.Value()}).ToSql()
	if err != nil {
		r.logger.Error("Failed to build select query for FindByTenantID", "error", err)
		return nil, err
	}
	if err := r.db.Select(&modelList, query, args...); err != nil {
		r.logger.Error("Failed to find sites by tenant ID", "tenant_id", tenantID.Value(), "error", err)
		return nil, err
	}
	return r.mapper.ToDomains(modelList)
}

// FindAll retrieves all sites
func (r *SiteRepositoryImpl) FindAll() ([]*entities.Site, error) {
	var modelList []*models.Site
	query, args, err := squirrel.Select("*").From("sites").ToSql()
	if err != nil {
		r.logger.Error("Failed to build select query for FindAll", "error", err)
		return nil, err
	}
	if err := r.db.Select(&modelList, query, args...); err != nil {
		r.logger.Error("Failed to find all sites", "error", err)
		return nil, err
	}
	return r.mapper.ToDomains(modelList)
}

// FindEnabledByTenantID retrieves only enabled sites for a tenant
func (r *SiteRepositoryImpl) FindEnabledByTenantID(tenantID entities.TenantID) ([]*entities.Site, error) {
	var modelList []*models.Site
	query, args, err := squirrel.Select("*").From("sites").Where(squirrel.Eq{"tenant_id": tenantID.Value(), "enabled": true}).ToSql()
	if err != nil {
		r.logger.Error("Failed to build select query for FindEnabledByTenantID", "error", err)
		return nil, err
	}
	if err := r.db.Select(&modelList, query, args...); err != nil {
		r.logger.Error("Failed to find enabled sites by tenant", "tenant_id", tenantID.Value(), "error", err)
		return nil, err
	}
	return r.mapper.ToDomains(modelList)
}

// Delete deletes a site (soft delete)
func (r *SiteRepositoryImpl) Delete(id entities.SiteID) error {
	query, args, err := squirrel.Delete("sites").Where(squirrel.Eq{"id": id.Value()}).ToSql()
	if err != nil {
		r.logger.Error("Failed to build delete query for site", "id", id.Value(), "error", err)
		return err
	}
	if _, err := r.db.Exec(query, args...); err != nil {
		r.logger.Error("Failed to delete site", "id", id.Value(), "error", err)
		return err
	}
	return nil
}

// ExistsByDomain checks if a site with the given domain exists
func (r *SiteRepositoryImpl) ExistsByDomain(domain *value_objects.DomainName) (bool, error) {
	var count int64
	query, args, err := squirrel.Select("COUNT(*)").From("sites").Where(squirrel.Eq{"domain": domain.Value()}).ToSql()
	if err != nil {
		r.logger.Error("Failed to build select query for ExistsByDomain", "domain", domain.Value(), "error", err)
		return false, err
	}
	if err := r.db.Get(&count, query, args...); err != nil {
		r.logger.Error("Failed to check site existence by domain", "domain", domain.Value(), "error", err)
		return false, err
	}
	return count > 0, nil
}
