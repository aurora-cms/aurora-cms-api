package repositories

import (
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"github.com/h4rdc0m/aurora-api/domain/entities"
	"github.com/h4rdc0m/aurora-api/domain/value_objects"
	"github.com/h4rdc0m/aurora-api/infrastructure/persistence/models"
	"github.com/h4rdc0m/aurora-api/tests/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

// Test constructor
func TestNewUserRepository(t *testing.T) {
	mockDB := new(mocks.Database)
	mockLogger := new(mocks.Logger)

	repo := NewUserRepository(mockDB, mockLogger)

	assert.NotNil(t, repo)
	assert.IsType(t, &UserRepositoryImpl{}, repo)
}

// Enhanced Save tests with all error conditions
func TestUserRepositoryImpl_Save(t *testing.T) {
	t.Run("create_new_user_success", func(t *testing.T) {
		mockDB := new(mocks.Database)
		mockLogger := new(mocks.Logger)
		mockMapper := new(mocks.MockUserMapper)
		repo := &UserRepositoryImpl{
			db:     mockDB,
			logger: mockLogger,
			mapper: mockMapper,
		}

		role, _ := value_objects.NewUserRole("user")
		domainUser, _ := entities.NewUser(
			value_objects.NewKeycloakIDFromUUID(uuid.MustParse("123e4567-e89b-12d3-a456-426614174000")),
			role,
		)
		model := &models.User{
			Role:       models.RoleUser,
			KeycloakID: uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
			Base:       models.Base{CreatedAt: time.Now(), UpdatedAt: time.Now()},
		}

		mockMapper.On("ToModel", domainUser).Return(model, nil)
		mockResult := new(mocks.SqlResult)
		mockResult.On("LastInsertId").Return(int64(1), nil)
		mockDB.On("Exec", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(mockResult, nil)

		err := repo.Save(domainUser)
		assert.NoError(t, err)
		assert.Equal(t, uint64(1), domainUser.ID().Value())

		mockDB.AssertExpectations(t)
		mockLogger.AssertExpectations(t)
		mockMapper.AssertExpectations(t)
	})

	t.Run("update_existing_user_success", func(t *testing.T) {
		mockDB := new(mocks.Database)
		mockLogger := new(mocks.Logger)
		mockMapper := new(mocks.MockUserMapper)
		repo := &UserRepositoryImpl{
			db:     mockDB,
			logger: mockLogger,
			mapper: mockMapper,
		}

		role, _ := value_objects.NewUserRole("user")
		domainUser, _ := entities.NewUser(
			value_objects.NewKeycloakIDFromUUID(uuid.MustParse("123e4567-e89b-12d3-a456-426614174000")),
			role,
		)
		domainUser.SetID(entities.NewUserID(1))

		model := &models.User{
			Base:       models.Base{ID: 1, CreatedAt: time.Now(), UpdatedAt: time.Now()},
			Role:       models.RoleUser,
			KeycloakID: uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
		}

		mockMapper.On("ToModel", domainUser).Return(model, nil)
		mockResult := new(mocks.SqlResult)
		mockDB.On("Exec", mock.Anything, model.KeycloakID, model.ID).Return(mockResult, nil)

		err := repo.Save(domainUser)
		assert.NoError(t, err)

		mockDB.AssertExpectations(t)
		mockLogger.AssertExpectations(t)
		mockMapper.AssertExpectations(t)
	})

	t.Run("mapper_error", func(t *testing.T) {
		mockDB := new(mocks.Database)
		mockLogger := new(mocks.Logger)
		mockMapper := new(mocks.MockUserMapper)
		repo := &UserRepositoryImpl{
			db:     mockDB,
			logger: mockLogger,
			mapper: mockMapper,
		}

		role, _ := value_objects.NewUserRole("user")
		domainUser, _ := entities.NewUser(
			value_objects.NewKeycloakIDFromUUID(uuid.MustParse("123e4567-e89b-12d3-a456-426614174000")),
			role,
		)
		mapperErr := errors.New("mapper error")

		mockMapper.On("ToModel", domainUser).Return(nil, mapperErr)
		mockLogger.On("Error", "Failed to map user entity to model", "error", mapperErr).Return()

		err := repo.Save(domainUser)
		assert.Error(t, err)
		assert.Equal(t, mapperErr, err)

		mockDB.AssertExpectations(t)
		mockLogger.AssertExpectations(t)
		mockMapper.AssertExpectations(t)
	})

	t.Run("insert_exec_error", func(t *testing.T) {
		mockDB := new(mocks.Database)
		mockLogger := new(mocks.Logger)
		mockMapper := new(mocks.MockUserMapper)
		repo := &UserRepositoryImpl{
			db:     mockDB,
			logger: mockLogger,
			mapper: mockMapper,
		}

		role, _ := value_objects.NewUserRole("user")
		domainUser, _ := entities.NewUser(
			value_objects.NewKeycloakIDFromUUID(uuid.MustParse("123e4567-e89b-12d3-a456-426614174000")),
			role,
		)
		model := &models.User{
			Role:       models.RoleUser,
			KeycloakID: uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
			Base:       models.Base{CreatedAt: time.Now(), UpdatedAt: time.Now()},
		}
		execErr := errors.New("exec error")

		mockMapper.On("ToModel", domainUser).Return(model, nil)
		mockResult := new(mocks.SqlResult)
		mockDB.On("Exec", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(mockResult, execErr)
		mockLogger.On("Error", "Failed to create user", "error", execErr).Return()

		err := repo.Save(domainUser)
		assert.Error(t, err)
		assert.Equal(t, execErr, err)

		mockDB.AssertExpectations(t)
		mockLogger.AssertExpectations(t)
		mockMapper.AssertExpectations(t)
	})

	t.Run("lastInsertId_error", func(t *testing.T) {
		mockDB := new(mocks.Database)
		mockLogger := new(mocks.Logger)
		mockMapper := new(mocks.MockUserMapper)
		repo := &UserRepositoryImpl{
			db:     mockDB,
			logger: mockLogger,
			mapper: mockMapper,
		}

		role, _ := value_objects.NewUserRole("user")
		domainUser, _ := entities.NewUser(
			value_objects.NewKeycloakIDFromUUID(uuid.MustParse("123e4567-e89b-12d3-a456-426614174000")),
			role,
		)
		model := &models.User{
			Role:       models.RoleUser,
			KeycloakID: uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
			Base:       models.Base{CreatedAt: time.Now(), UpdatedAt: time.Now()},
		}
		lastInsertErr := errors.New("lastInsertId error")

		mockMapper.On("ToModel", domainUser).Return(model, nil)
		mockResult := new(mocks.SqlResult)
		mockResult.On("LastInsertId").Return(int64(0), lastInsertErr)
		mockDB.On("Exec", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(mockResult, nil)
		mockLogger.On("Error", "Failed to get last insert ID for user", "error", lastInsertErr).Return()

		err := repo.Save(domainUser)
		assert.Error(t, err)
		assert.Equal(t, lastInsertErr, err)

		mockDB.AssertExpectations(t)
		mockLogger.AssertExpectations(t)
		mockMapper.AssertExpectations(t)
	})

	t.Run("update_exec_error", func(t *testing.T) {
		mockDB := new(mocks.Database)
		mockLogger := new(mocks.Logger)
		mockMapper := new(mocks.MockUserMapper)
		repo := &UserRepositoryImpl{
			db:     mockDB,
			logger: mockLogger,
			mapper: mockMapper,
		}

		role, _ := value_objects.NewUserRole("user")
		domainUser, _ := entities.NewUser(
			value_objects.NewKeycloakIDFromUUID(uuid.MustParse("123e4567-e89b-12d3-a456-426614174000")),
			role,
		)
		domainUser.SetID(entities.NewUserID(1))
		model := &models.User{
			Base:       models.Base{ID: 1, CreatedAt: time.Now(), UpdatedAt: time.Now()},
			Role:       models.RoleUser,
			KeycloakID: uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
		}
		updateErr := errors.New("update error")

		mockMapper.On("ToModel", domainUser).Return(model, nil)
		mockResult := new(mocks.SqlResult)
		mockDB.On("Exec", mock.Anything, model.KeycloakID, model.ID).Return(mockResult, updateErr)
		mockLogger.On("Error", "Failed to update user", "error", updateErr).Return()

		err := repo.Save(domainUser)
		assert.Error(t, err)
		assert.Equal(t, updateErr, err)

		mockDB.AssertExpectations(t)
		mockLogger.AssertExpectations(t)
		mockMapper.AssertExpectations(t)
	})
}

// Enhanced FindByID tests
func TestUserRepositoryImpl_FindByID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockDB := new(mocks.Database)
		mockLogger := new(mocks.Logger)
		mockMapper := new(mocks.MockUserMapper)
		repo := &UserRepositoryImpl{
			db:     mockDB,
			logger: mockLogger,
			mapper: mockMapper,
		}

		userID := entities.NewUserID(1)
		model := &models.User{Base: models.Base{ID: 1}}
		domainUser := &entities.User{}

		mockDB.On("Get", mock.Anything, mock.Anything, userID.Value()).Run(func(args mock.Arguments) {
			arg := args.Get(0).(*models.User)
			*arg = *model
		}).Return(nil)
		mockMapper.On("ToDomain", model).Return(domainUser, nil)

		result, err := repo.FindByID(userID)
		assert.NoError(t, err)
		assert.Equal(t, domainUser, result)

		mockDB.AssertExpectations(t)
		mockLogger.AssertExpectations(t)
		mockMapper.AssertExpectations(t)
	})

	t.Run("not_found", func(t *testing.T) {
		mockDB := new(mocks.Database)
		mockLogger := new(mocks.Logger)
		mockMapper := new(mocks.MockUserMapper)
		repo := &UserRepositoryImpl{
			db:     mockDB,
			logger: mockLogger,
			mapper: mockMapper,
		}

		userID := entities.NewUserID(999)

		mockDB.On("Get", mock.Anything, mock.Anything, userID.Value()).Return(sql.ErrNoRows)
		mockLogger.On("Warn", "User not found by ID", "id", userID).Return()

		result, err := repo.FindByID(userID)
		assert.NoError(t, err)
		assert.Nil(t, result)

		mockDB.AssertExpectations(t)
		mockLogger.AssertExpectations(t)
		mockMapper.AssertExpectations(t)
	})

	t.Run("database_error", func(t *testing.T) {
		mockDB := new(mocks.Database)
		mockLogger := new(mocks.Logger)
		mockMapper := new(mocks.MockUserMapper)
		repo := &UserRepositoryImpl{
			db:     mockDB,
			logger: mockLogger,
			mapper: mockMapper,
		}

		userID := entities.NewUserID(1)
		dbErr := errors.New("database error")

		mockDB.On("Get", mock.Anything, mock.Anything, userID.Value()).Return(dbErr)
		mockLogger.On("Error", "Failed to find user by ID", "id", userID.Value(), "error", dbErr).Return()

		result, err := repo.FindByID(userID)
		assert.Error(t, err)
		assert.Equal(t, dbErr, err)
		assert.Nil(t, result)

		mockDB.AssertExpectations(t)
		mockLogger.AssertExpectations(t)
		mockMapper.AssertExpectations(t)
	})

	t.Run("mapper_error", func(t *testing.T) {
		mockDB := new(mocks.Database)
		mockLogger := new(mocks.Logger)
		mockMapper := new(mocks.MockUserMapper)
		repo := &UserRepositoryImpl{
			db:     mockDB,
			logger: mockLogger,
			mapper: mockMapper,
		}

		userID := entities.NewUserID(1)
		model := &models.User{Base: models.Base{ID: 1}}
		mapperErr := errors.New("mapper error")

		mockDB.On("Get", mock.Anything, mock.Anything, userID.Value()).Run(func(args mock.Arguments) {
			arg := args.Get(0).(*models.User)
			*arg = *model
		}).Return(nil)
		mockMapper.On("ToDomain", model).Return(nil, mapperErr)

		result, err := repo.FindByID(userID)
		assert.Error(t, err)
		assert.Equal(t, mapperErr, err)
		assert.Nil(t, result)

		mockDB.AssertExpectations(t)
		mockLogger.AssertExpectations(t)
		mockMapper.AssertExpectations(t)
	})
}

// Enhanced FindByKeycloakID tests
func TestUserRepositoryImpl_FindByKeycloakID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockDB := new(mocks.Database)
		mockLogger := new(mocks.Logger)
		mockMapper := new(mocks.MockUserMapper)
		repo := &UserRepositoryImpl{
			db:     mockDB,
			logger: mockLogger,
			mapper: mockMapper,
		}

		keycloakID := value_objects.NewKeycloakIDFromUUID(uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"))
		model := &models.User{Base: models.Base{ID: 1}, KeycloakID: keycloakID.Value()}
		domainUser := &entities.User{}

		mockDB.On("Get", mock.Anything, mock.Anything, *keycloakID).Run(func(args mock.Arguments) {
			arg := args.Get(0).(*models.User)
			*arg = *model
		}).Return(nil)
		mockMapper.On("ToDomain", model).Return(domainUser, nil)

		result, err := repo.FindByKeycloakID(*keycloakID)
		assert.NoError(t, err)
		assert.Equal(t, domainUser, result)

		mockDB.AssertExpectations(t)
		mockLogger.AssertExpectations(t)
		mockMapper.AssertExpectations(t)
	})

	t.Run("not_found", func(t *testing.T) {
		mockDB := new(mocks.Database)
		mockLogger := new(mocks.Logger)
		mockMapper := new(mocks.MockUserMapper)
		repo := &UserRepositoryImpl{
			db:     mockDB,
			logger: mockLogger,
			mapper: mockMapper,
		}

		keycloakID := value_objects.NewKeycloakIDFromUUID(uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"))

		mockDB.On("Get", mock.Anything, mock.Anything, *keycloakID).Return(sql.ErrNoRows)
		mockLogger.On("Warn", "User not found by Keycloak ID", "keycloakID", *keycloakID).Return()

		result, err := repo.FindByKeycloakID(*keycloakID)
		assert.NoError(t, err)
		assert.Nil(t, result)

		mockDB.AssertExpectations(t)
		mockLogger.AssertExpectations(t)
		mockMapper.AssertExpectations(t)
	})

	t.Run("database_error", func(t *testing.T) {
		mockDB := new(mocks.Database)
		mockLogger := new(mocks.Logger)
		mockMapper := new(mocks.MockUserMapper)
		repo := &UserRepositoryImpl{
			db:     mockDB,
			logger: mockLogger,
			mapper: mockMapper,
		}

		keycloakID := value_objects.NewKeycloakIDFromUUID(uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"))
		dbErr := errors.New("database error")

		mockDB.On("Get", mock.Anything, mock.Anything, *keycloakID).Return(dbErr)
		mockLogger.On("Error", "Failed to find user by Keycloak ID", "keycloakID", *keycloakID, "error", dbErr).Return()

		result, err := repo.FindByKeycloakID(*keycloakID)
		assert.Error(t, err)
		assert.Equal(t, dbErr, err)
		assert.Nil(t, result)

		mockDB.AssertExpectations(t)
		mockLogger.AssertExpectations(t)
		mockMapper.AssertExpectations(t)
	})

	t.Run("mapper_error", func(t *testing.T) {
		mockDB := new(mocks.Database)
		mockLogger := new(mocks.Logger)
		mockMapper := new(mocks.MockUserMapper)
		repo := &UserRepositoryImpl{
			db:     mockDB,
			logger: mockLogger,
			mapper: mockMapper,
		}

		keycloakID := value_objects.NewKeycloakIDFromUUID(uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"))
		model := &models.User{Base: models.Base{ID: 1}, KeycloakID: keycloakID.Value()}
		mapperErr := errors.New("mapper error")

		mockDB.On("Get", mock.Anything, mock.Anything, *keycloakID).Run(func(args mock.Arguments) {
			arg := args.Get(0).(*models.User)
			*arg = *model
		}).Return(nil)
		mockMapper.On("ToDomain", model).Return(nil, mapperErr)

		result, err := repo.FindByKeycloakID(*keycloakID)
		assert.Error(t, err)
		assert.Equal(t, mapperErr, err)
		assert.Nil(t, result)

		mockDB.AssertExpectations(t)
		mockLogger.AssertExpectations(t)
		mockMapper.AssertExpectations(t)
	})
}

// Enhanced FindAll tests
func TestUserRepositoryImpl_FindAll(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockDB := new(mocks.Database)
		mockLogger := new(mocks.Logger)
		mockMapper := new(mocks.MockUserMapper)
		repo := &UserRepositoryImpl{
			db:     mockDB,
			logger: mockLogger,
			mapper: mockMapper,
		}

		modelsList := []*models.User{{Base: models.Base{ID: 1}}}
		domainList := []*entities.User{{}}

		mockDB.On("Select", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
			arg := args.Get(0).(*[]*models.User)
			*arg = modelsList
		}).Return(nil)
		mockMapper.On("ToDomains", modelsList).Return(domainList, nil)

		result, err := repo.FindAll()
		assert.NoError(t, err)
		assert.Equal(t, domainList, result)

		mockDB.AssertExpectations(t)
		mockLogger.AssertExpectations(t)
		mockMapper.AssertExpectations(t)
	})

	t.Run("no_rows_found", func(t *testing.T) {
		mockDB := new(mocks.Database)
		mockLogger := new(mocks.Logger)
		mockMapper := new(mocks.MockUserMapper)
		repo := &UserRepositoryImpl{
			db:     mockDB,
			logger: mockLogger,
			mapper: mockMapper,
		}

		mockDB.On("Select", mock.Anything, mock.Anything).Return(sql.ErrNoRows)
		mockLogger.On("Warn", "No users found", "error", sql.ErrNoRows).Return()

		result, err := repo.FindAll()
		assert.NoError(t, err)
		assert.Nil(t, result)

		mockDB.AssertExpectations(t)
		mockLogger.AssertExpectations(t)
		mockMapper.AssertExpectations(t)
	})

	t.Run("database_error", func(t *testing.T) {
		mockDB := new(mocks.Database)
		mockLogger := new(mocks.Logger)
		mockMapper := new(mocks.MockUserMapper)
		repo := &UserRepositoryImpl{
			db:     mockDB,
			logger: mockLogger,
			mapper: mockMapper,
		}

		dbErr := errors.New("database error")

		mockDB.On("Select", mock.Anything, mock.Anything).Return(dbErr)
		mockLogger.On("Error", "Failed to find all users", "error", dbErr).Return()

		result, err := repo.FindAll()
		assert.Error(t, err)
		assert.Equal(t, dbErr, err)
		assert.Nil(t, result)

		mockDB.AssertExpectations(t)
		mockLogger.AssertExpectations(t)
		mockMapper.AssertExpectations(t)
	})

	t.Run("mapper_error", func(t *testing.T) {
		mockDB := new(mocks.Database)
		mockLogger := new(mocks.Logger)
		mockMapper := new(mocks.MockUserMapper)
		repo := &UserRepositoryImpl{
			db:     mockDB,
			logger: mockLogger,
			mapper: mockMapper,
		}

		modelsList := []*models.User{{Base: models.Base{ID: 1}}}
		mapperErr := errors.New("mapper error")

		mockDB.On("Select", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
			arg := args.Get(0).(*[]*models.User)
			*arg = modelsList
		}).Return(nil)
		mockMapper.On("ToDomains", modelsList).Return(nil, mapperErr)

		result, err := repo.FindAll()
		assert.Error(t, err)
		assert.Equal(t, mapperErr, err)
		assert.Nil(t, result)

		mockDB.AssertExpectations(t)
		mockLogger.AssertExpectations(t)
		mockMapper.AssertExpectations(t)
	})
}

// Enhanced FindAllByTenantID tests
func TestUserRepositoryImpl_FindAllByTenantID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockDB := new(mocks.Database)
		mockLogger := new(mocks.Logger)
		mockMapper := new(mocks.MockUserMapper)
		repo := &UserRepositoryImpl{
			db:     mockDB,
			logger: mockLogger,
			mapper: mockMapper,
		}

		tenantID := entities.NewTenantID(1)
		modelsList := []*models.User{{Base: models.Base{ID: 1}}}
		domainList := []*entities.User{{}}

		mockDB.On("Select", mock.Anything, mock.Anything, tenantID.Value()).Run(func(args mock.Arguments) {
			arg := args.Get(0).(*[]*models.User)
			*arg = modelsList
		}).Return(nil)
		mockMapper.On("ToDomains", modelsList).Return(domainList, nil)

		result, err := repo.FindAllByTenantID(tenantID)
		assert.NoError(t, err)
		assert.Equal(t, domainList, result)

		mockDB.AssertExpectations(t)
		mockLogger.AssertExpectations(t)
		mockMapper.AssertExpectations(t)
	})

	t.Run("no_rows_found", func(t *testing.T) {
		mockDB := new(mocks.Database)
		mockLogger := new(mocks.Logger)
		mockMapper := new(mocks.MockUserMapper)
		repo := &UserRepositoryImpl{
			db:     mockDB,
			logger: mockLogger,
			mapper: mockMapper,
		}

		tenantID := entities.NewTenantID(999)

		mockDB.On("Select", mock.Anything, mock.Anything, tenantID.Value()).Return(sql.ErrNoRows)
		mockLogger.On("Warn", "No users found for tenant", "tenantID", tenantID).Return()

		result, err := repo.FindAllByTenantID(tenantID)
		assert.NoError(t, err)
		assert.Nil(t, result)

		mockDB.AssertExpectations(t)
		mockLogger.AssertExpectations(t)
		mockMapper.AssertExpectations(t)
	})

	t.Run("database_error", func(t *testing.T) {
		mockDB := new(mocks.Database)
		mockLogger := new(mocks.Logger)
		mockMapper := new(mocks.MockUserMapper)
		repo := &UserRepositoryImpl{
			db:     mockDB,
			logger: mockLogger,
			mapper: mockMapper,
		}

		tenantID := entities.NewTenantID(1)
		dbErr := errors.New("database error")

		mockDB.On("Select", mock.Anything, mock.Anything, tenantID.Value()).Return(dbErr)
		mockLogger.On("Error", "Failed to find users by tenant ID", "tenantID", tenantID, "error", dbErr).Return()

		result, err := repo.FindAllByTenantID(tenantID)
		assert.Error(t, err)
		assert.Equal(t, dbErr, err)
		assert.Nil(t, result)

		mockDB.AssertExpectations(t)
		mockLogger.AssertExpectations(t)
		mockMapper.AssertExpectations(t)
	})

	t.Run("mapper_error", func(t *testing.T) {
		mockDB := new(mocks.Database)
		mockLogger := new(mocks.Logger)
		mockMapper := new(mocks.MockUserMapper)
		repo := &UserRepositoryImpl{
			db:     mockDB,
			logger: mockLogger,
			mapper: mockMapper,
		}

		tenantID := entities.NewTenantID(1)
		modelsList := []*models.User{{Base: models.Base{ID: 1}}}
		mapperErr := errors.New("mapper error")

		mockDB.On("Select", mock.Anything, mock.Anything, tenantID.Value()).Run(func(args mock.Arguments) {
			arg := args.Get(0).(*[]*models.User)
			*arg = modelsList
		}).Return(nil)
		mockMapper.On("ToDomains", modelsList).Return(nil, mapperErr)

		result, err := repo.FindAllByTenantID(tenantID)
		assert.Error(t, err)
		assert.Equal(t, mapperErr, err)
		assert.Nil(t, result)

		mockDB.AssertExpectations(t)
		mockLogger.AssertExpectations(t)
		mockMapper.AssertExpectations(t)
	})
}

// Enhanced Delete tests
func TestUserRepositoryImpl_Delete(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockDB := new(mocks.Database)
		mockLogger := new(mocks.Logger)
		mockMapper := new(mocks.MockUserMapper)
		repo := &UserRepositoryImpl{
			db:     mockDB,
			logger: mockLogger,
			mapper: mockMapper,
		}

		userID := entities.NewUserID(1)
		mockResult := new(mocks.SqlResult)

		mockDB.On("Exec", mock.Anything, userID.Value()).Return(mockResult, nil)

		err := repo.Delete(userID)
		assert.NoError(t, err)

		mockDB.AssertExpectations(t)
		mockLogger.AssertExpectations(t)
		mockMapper.AssertExpectations(t)
	})

	t.Run("database_error", func(t *testing.T) {
		mockDB := new(mocks.Database)
		mockLogger := new(mocks.Logger)
		mockMapper := new(mocks.MockUserMapper)
		repo := &UserRepositoryImpl{
			db:     mockDB,
			logger: mockLogger,
			mapper: mockMapper,
		}

		userID := entities.NewUserID(1)
		dbErr := errors.New("database error")
		mockResult := new(mocks.SqlResult)

		mockDB.On("Exec", mock.Anything, userID.Value()).Return(mockResult, dbErr)
		mockLogger.On("Error", "Failed to delete user", "id", userID.Value(), "error", dbErr).Return()

		err := repo.Delete(userID)
		assert.Error(t, err)
		assert.Equal(t, dbErr, err)

		mockDB.AssertExpectations(t)
		mockLogger.AssertExpectations(t)
		mockMapper.AssertExpectations(t)
	})
}

// Complete ExistsByKeycloakID tests (from previous response)
func TestUserRepositoryImpl_ExistsByKeycloakID(t *testing.T) {
	t.Run("user_exists", func(t *testing.T) {
		mockDB := new(mocks.Database)
		mockLogger := new(mocks.Logger)
		mockMapper := new(mocks.MockUserMapper)
		repo := &UserRepositoryImpl{
			db:     mockDB,
			logger: mockLogger,
			mapper: mockMapper,
		}

		keycloakID := value_objects.NewKeycloakIDFromUUID(uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"))

		mockDB.On("Get", mock.Anything, mock.Anything, *keycloakID).Run(func(args mock.Arguments) {
			arg := args.Get(0).(*int64)
			*arg = 1
		}).Return(nil)

		result, err := repo.ExistsByKeycloakID(*keycloakID)
		assert.NoError(t, err)
		assert.True(t, result)

		mockDB.AssertExpectations(t)
		mockLogger.AssertExpectations(t)
		mockMapper.AssertExpectations(t)
	})

	t.Run("user_does_not_exist", func(t *testing.T) {
		mockDB := new(mocks.Database)
		mockLogger := new(mocks.Logger)
		mockMapper := new(mocks.MockUserMapper)
		repo := &UserRepositoryImpl{
			db:     mockDB,
			logger: mockLogger,
			mapper: mockMapper,
		}

		keycloakID := value_objects.NewKeycloakIDFromUUID(uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"))

		mockDB.On("Get", mock.Anything, mock.Anything, *keycloakID).Run(func(args mock.Arguments) {
			arg := args.Get(0).(*int64)
			*arg = 0
		}).Return(nil)

		result, err := repo.ExistsByKeycloakID(*keycloakID)
		assert.NoError(t, err)
		assert.False(t, result)

		mockDB.AssertExpectations(t)
		mockLogger.AssertExpectations(t)
		mockMapper.AssertExpectations(t)
	})

	t.Run("database_error_not_no_rows", func(t *testing.T) {
		mockDB := new(mocks.Database)
		mockLogger := new(mocks.Logger)
		mockMapper := new(mocks.MockUserMapper)
		repo := &UserRepositoryImpl{
			db:     mockDB,
			logger: mockLogger,
			mapper: mockMapper,
		}

		keycloakID := value_objects.NewKeycloakIDFromUUID(uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"))
		dbError := errors.New("database connection error")

		mockDB.On("Get", mock.Anything, mock.Anything, *keycloakID).Return(dbError)
		mockLogger.On("Error", "Failed to check existence of user by Keycloak ID", "keycloakID", *keycloakID, "error", dbError).Return()

		result, err := repo.ExistsByKeycloakID(*keycloakID)
		assert.Error(t, err)
		assert.Equal(t, dbError, err)
		assert.False(t, result)

		mockDB.AssertExpectations(t)
		mockLogger.AssertExpectations(t)
		mockMapper.AssertExpectations(t)
	})

	t.Run("database_error_sql_no_rows", func(t *testing.T) {
		mockDB := new(mocks.Database)
		mockLogger := new(mocks.Logger)
		mockMapper := new(mocks.MockUserMapper)
		repo := &UserRepositoryImpl{
			db:     mockDB,
			logger: mockLogger,
			mapper: mockMapper,
		}

		keycloakID := value_objects.NewKeycloakIDFromUUID(uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"))

		mockDB.On("Get", mock.Anything, mock.Anything, *keycloakID).Return(sql.ErrNoRows)
		mockLogger.On("Warn", "No users found with Keycloak ID", "keycloakID", *keycloakID).Return()

		result, err := repo.ExistsByKeycloakID(*keycloakID)
		assert.NoError(t, err)
		assert.False(t, result)

		mockDB.AssertExpectations(t)
		mockLogger.AssertExpectations(t)
		mockMapper.AssertExpectations(t)
	})

	t.Run("count_greater_than_one", func(t *testing.T) {
		mockDB := new(mocks.Database)
		mockLogger := new(mocks.Logger)
		mockMapper := new(mocks.MockUserMapper)
		repo := &UserRepositoryImpl{
			db:     mockDB,
			logger: mockLogger,
			mapper: mockMapper,
		}

		keycloakID := value_objects.NewKeycloakIDFromUUID(uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"))

		mockDB.On("Get", mock.Anything, mock.Anything, *keycloakID).Run(func(args mock.Arguments) {
			arg := args.Get(0).(*int64)
			*arg = 5 // Multiple users found
		}).Return(nil)

		result, err := repo.ExistsByKeycloakID(*keycloakID)
		assert.NoError(t, err)
		assert.True(t, result) // count > 0, so exists should return true

		mockDB.AssertExpectations(t)
		mockLogger.AssertExpectations(t)
		mockMapper.AssertExpectations(t)
	})
}
