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

// UserRepositoryImpl is an implementation of the UserRepository interface for user persistence operations.
// It uses a UserDB for database interactions, a Logger for logging, and a Mapper for entity-model conversions.
type UserRepositoryImpl struct {
	db     common.Database
	logger common.Logger
	mapper common.Mapper[*entities.User, *models.User]
}

// NewUserRepository creates a new instance of UserRepository using the provided database and logger dependencies.
func NewUserRepository(db common.Database, logger common.Logger) repositories.UserRepository {
	return &UserRepositoryImpl{
		db:     db,
		logger: logger,
		mapper: mappers.NewUserMapper(),
	}
}

// Save persists a new or existing user entity to the database. Returns an error if the operation fails.
func (r *UserRepositoryImpl) Save(user *entities.User) error {
	model, err := r.mapper.ToModel(user)
	if err != nil {
		r.logger.Error("Failed to map user entity to model", "error", err)
		return err
	}

	if model.ID == 0 {
		// Create new user
		query, args, err := squirrel.Insert("users").
			Columns("keycloak_id", "created_at", "updated_at").
			Values(model.KeycloakID, model.CreatedAt, model.UpdatedAt).
			PlaceholderFormat(squirrel.Question).
			ToSql()
		if err != nil {
			r.logger.Error("Failed to build insert query for user", "error", err)
			return err
		}

		result, err := r.db.Exec(query, args...)
		if err != nil {
			r.logger.Error("Failed to create user", "error", err)
			return err
		}
		id, err := result.LastInsertId()
		if err != nil {
			r.logger.Error("Failed to get last insert ID for user", "error", err)
			return err
		}
		user.SetID(entities.NewUserID(uint64(id)))
	} else {
		// Update existing user
		query, args, err := squirrel.Update("users").
			Set("keycloak_id", model.KeycloakID).Where("id = ?", model.ID).
			PlaceholderFormat(squirrel.Question).ToSql()
		if err != nil {
			r.logger.Error("Failed to build update query for user", "error", err)
			return err
		}

		if _, err := r.db.Exec(query, args...); err != nil {
			r.logger.Error("Failed to update user", "error", err)
			return err
		}
	}

	return nil
}

// FindByID retrieves a user from the database by their unique identifier and returns the domain entity or an error.
func (r *UserRepositoryImpl) FindByID(id entities.UserID) (*entities.User, error) {
	var model models.User
	query, args, err := squirrel.Select("*").From("users").Where("id = ?", id.Value()).ToSql()
	if err != nil {
		r.logger.Error("Failed to build select query for FindByID", "error", err)
		return nil, err
	}

	if err := r.db.Get(&model, query, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			r.logger.Warn("User not found by ID", "id", id)
			return nil, nil
		}
		r.logger.Error("Failed to find user by ID", "id", id.Value(), "error", err)
		return nil, err
	}

	return r.mapper.ToDomain(&model)
}

// FindByKeycloakID retrieves a user entity by its Keycloak ID from the data source. Returns nil if no user is found.
func (r *UserRepositoryImpl) FindByKeycloakID(keycloakID value_objects.KeycloakID) (*entities.User, error) {
	var model models.User
	query, args, err := squirrel.Select("*").From("users").Where("keycloak_id = ?", keycloakID).ToSql()
	if err != nil {
		r.logger.Error("Failed to build select query for FindByKeycloakID", "error", err)
		return nil, err
	}

	if err := r.db.Get(&model, query, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			r.logger.Warn("User not found by Keycloak ID", "keycloakID", keycloakID)
			return nil, nil
		}
		r.logger.Error("Failed to find user by Keycloak ID", "keycloakID", keycloakID, "error", err)
		return nil, err
	}

	return r.mapper.ToDomain(&model)
}

// FindAll retrieves all users from the database and maps them to domain entities. Returns an error if the operation fails.
func (r *UserRepositoryImpl) FindAll() ([]*entities.User, error) {
	var modelList []*models.User
	query, args, err := squirrel.Select("*").From("users").ToSql()
	if err != nil {
		r.logger.Error("Failed to build select query for FindAll", "error", err)
		return nil, err
	}
	if err := r.db.Select(&modelList, query, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			r.logger.Warn("No users found", "error", err)
			return nil, nil
		}
		r.logger.Error("Failed to find all users", "error", err)
		return nil, err
	}

	return r.mapper.ToDomains(modelList)
}

// FindAllByTenantID retrieves all users associated with the specified tenant ID. Returns a list of users or an error.
func (r *UserRepositoryImpl) FindAllByTenantID(tenantID entities.TenantID) ([]*entities.User, error) {
	var modelList []*models.User
	query, args, err := squirrel.Select("*").From("users").Where(squirrel.Eq{"tenant_id": tenantID.Value()}).ToSql()
	if err != nil {
		r.logger.Error("Failed to build select query for FindAllByTenantID", "error", err)
		return nil, err
	}

	if err := r.db.Select(&modelList, query, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			r.logger.Warn("No users found for tenant", "tenantID", tenantID)
			return nil, nil
		}
		r.logger.Error("Failed to find users by tenant ID", "tenantID", tenantID, "error", err)
		return nil, err
	}

	return r.mapper.ToDomains(modelList)
}

// Delete removes a user record from the database identified by the provided UserID. Returns an error if the operation fails.
func (r *UserRepositoryImpl) Delete(id entities.UserID) error {
	query, args, err := squirrel.Delete("users").Where(squirrel.Eq{"id": id.Value()}).ToSql()
	if err != nil {
		r.logger.Error("Failed to build delete query for user", "error", err)
		return err
	}
	if _, err := r.db.Exec(query, args...); err != nil {
		r.logger.Error("Failed to delete user", "id", id.Value(), "error", err)
		return err
	}
	return nil
}

// ExistsByKeycloakID checks if a user with the given Keycloak ID exists in the database and returns a boolean result.
func (r *UserRepositoryImpl) ExistsByKeycloakID(keycloakID value_objects.KeycloakID) (bool, error) {
	var count int64
	query, args, err := squirrel.Select("COUNT(*)").From("users").Where(squirrel.Eq{"keycloak_id": keycloakID}).ToSql()
	if err != nil {
		r.logger.Error("Failed to build count query for ExistsByKeycloakID", "keycloakID", keycloakID, "error", err)
		return false, err
	}
	if err := r.db.Get(&count, query, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			r.logger.Warn("No users found with Keycloak ID", "keycloakID", keycloakID)
			return false, nil
		}
		r.logger.Error("Failed to check existence of user by Keycloak ID", "keycloakID", keycloakID, "error", err)
		return false, err
	}

	return count > 0, nil
}
