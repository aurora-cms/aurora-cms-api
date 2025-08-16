package repositories

import (
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/h4rdc0m/aurora-api/domain/entities"
	"github.com/h4rdc0m/aurora-api/infrastructure/persistence/models"
	"github.com/h4rdc0m/aurora-api/tests/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestTenantRepository_Save(t *testing.T) {
	t.Run("insert success", func(t *testing.T) {
		mockDB := new(mocks.Database)
		mockLogger := new(mocks.Logger)
		repo := &TenantRepositoryImpl{db: mockDB, logger: mockLogger, mapper: &mocks.MockTenantMapper{}}
		tenant := &entities.Tenant{}
		model := &models.Tenant{Name: "Test", Base: models.Base{CreatedAt: time.Now(), UpdatedAt: time.Now()}}
		mapperMock := repo.mapper.(*mocks.MockTenantMapper)
		mapperMock.On("ToModel", tenant).Return(model, nil)
		mockResult := new(mocks.SqlResult)
		mockResult.On("LastInsertId").Return(int64(42), nil)
		mockDB.On("Exec", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(mockResult, nil)
		err := repo.Save(tenant)
		assert.NoError(t, err)
		assert.Equal(t, uint64(42), tenant.ID().Value())
		mockDB.AssertExpectations(t)
		mockLogger.AssertExpectations(t)
		mapperMock.AssertExpectations(t)
	})

	t.Run("update success", func(t *testing.T) {
		mockDB := new(mocks.Database)
		mockLogger := new(mocks.Logger)
		repo := &TenantRepositoryImpl{db: mockDB, logger: mockLogger, mapper: &mocks.MockTenantMapper{}}
		tenant := &entities.Tenant{}
		tenant.SetID(entities.NewTenantID(99))
		model := &models.Tenant{Name: "Test", Base: models.Base{ID: 99, CreatedAt: time.Now(), UpdatedAt: time.Now()}}
		mapperMock := repo.mapper.(*mocks.MockTenantMapper)
		mapperMock.On("ToModel", tenant).Return(model, nil)
		mockResult := new(mocks.SqlResult)
		mockDB.On("Exec", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(mockResult, nil)
		err := repo.Save(tenant)
		assert.NoError(t, err)
		mockDB.AssertExpectations(t)
		mockLogger.AssertExpectations(t)
		mapperMock.AssertExpectations(t)
	})

	t.Run("mapper error", func(t *testing.T) {
		mockDB := new(mocks.Database)
		mockLogger := new(mocks.Logger)
		repo := &TenantRepositoryImpl{db: mockDB, logger: mockLogger, mapper: &mocks.MockTenantMapper{}}
		tenant := &entities.Tenant{}
		mapperErr := errors.New("mapper error")
		mapperMock := repo.mapper.(*mocks.MockTenantMapper)
		mapperMock.On("ToModel", tenant).Return(nil, mapperErr)
		mockLogger.On("Error", "Failed to map tenant to model", "error", mapperErr).Return()
		err := repo.Save(tenant)
		assert.Error(t, err)
		assert.Equal(t, mapperErr, err)
		mockDB.AssertExpectations(t)
		mockLogger.AssertExpectations(t)
		mapperMock.AssertExpectations(t)
	})
}

func TestTenantRepository_Save_ErrorBranches(t *testing.T) {
	tenant := &entities.Tenant{}
	model := &models.Tenant{Name: "Test", Base: models.Base{CreatedAt: time.Now(), UpdatedAt: time.Now()}}

	t.Run("insert exec error", func(t *testing.T) {
		mockDB := new(mocks.Database)
		mockLogger := new(mocks.Logger)
		repo := &TenantRepositoryImpl{db: mockDB, logger: mockLogger, mapper: &mocks.MockTenantMapper{}}
		mapperMock := repo.mapper.(*mocks.MockTenantMapper)
		mapperMock.On("ToModel", tenant).Return(model, nil)
		mockDB.On("Exec", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(new(mocks.SqlResult), errors.New("exec error"))
		mockLogger.On("Error", "Failed to create tenant", "error", mock.Anything).Return()
		err := repo.Save(tenant)
		assert.Error(t, err)
		mockDB.AssertExpectations(t)
		mockLogger.AssertExpectations(t)
		mapperMock.AssertExpectations(t)
	})

	t.Run("lastInsertId error", func(t *testing.T) {
		mockDB := new(mocks.Database)
		mockLogger := new(mocks.Logger)
		repo := &TenantRepositoryImpl{db: mockDB, logger: mockLogger, mapper: &mocks.MockTenantMapper{}}
		mapperMock := repo.mapper.(*mocks.MockTenantMapper)
		mapperMock.On("ToModel", tenant).Return(model, nil)
		mockResult := new(mocks.SqlResult)
		mockResult.On("LastInsertId").Return(int64(0), errors.New("lastInsertId error"))
		mockDB.On("Exec", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(mockResult, nil)
		mockLogger.On("Error", "Failed to get last insert ID for tenant", "error", mock.Anything).Return()
		err := repo.Save(tenant)
		assert.Error(t, err)
		mockDB.AssertExpectations(t)
		mockLogger.AssertExpectations(t)
		mapperMock.AssertExpectations(t)
	})

	t.Run("update exec error", func(t *testing.T) {
		tenant.SetID(entities.NewTenantID(99))
		mockDB := new(mocks.Database)
		mockLogger := new(mocks.Logger)
		repo := &TenantRepositoryImpl{db: mockDB, logger: mockLogger, mapper: &mocks.MockTenantMapper{}}
		mapperMock := repo.mapper.(*mocks.MockTenantMapper)
		mapperMock.On("ToModel", tenant).Return(&models.Tenant{Name: "Test", Base: models.Base{ID: 99, CreatedAt: time.Now(), UpdatedAt: time.Now()}}, nil)
		mockDB.On("Exec", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(new(mocks.SqlResult), errors.New("exec error"))
		mockLogger.On("Error", "Failed to update tenant", "id", uint64(99), "error", mock.Anything).Return()
		err := repo.Save(tenant)
		assert.Error(t, err)
		mockDB.AssertExpectations(t)
		mockLogger.AssertExpectations(t)
		mapperMock.AssertExpectations(t)
	})
}

func TestTenantRepository_FindByID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockDB := new(mocks.Database)
		mockLogger := new(mocks.Logger)
		repo := &TenantRepositoryImpl{db: mockDB, logger: mockLogger, mapper: &mocks.MockTenantMapper{}}
		id := entities.NewTenantID(123)
		mockDB.On("Get", mock.AnythingOfType("*models.Tenant"), mock.Anything, id.Value()).Run(func(args mock.Arguments) {
			tenant := args.Get(0).(*models.Tenant)
			tenant.ID = 123
			tenant.Name = "Test"
			tenant.CreatedAt = time.Now()
			tenant.UpdatedAt = time.Now()
		}).Return(nil)
		expectedTenant := &entities.Tenant{}
		mapperMock := repo.mapper.(*mocks.MockTenantMapper)
		mapperMock.On("ToDomain", mock.AnythingOfType("*models.Tenant")).Return(expectedTenant, nil)
		result, err := repo.FindByID(id)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, expectedTenant, result)
		mockDB.AssertExpectations(t)
		mockLogger.AssertExpectations(t)
		mapperMock.AssertExpectations(t)
	})
	t.Run("not_found", func(t *testing.T) {
		mockDB := new(mocks.Database)
		mockLogger := new(mocks.Logger)
		repo := &TenantRepositoryImpl{db: mockDB, logger: mockLogger, mapper: &mocks.MockTenantMapper{}}
		id := entities.NewTenantID(999)
		mockDB.On("Get", mock.AnythingOfType("*models.Tenant"), mock.Anything, id.Value()).Return(sql.ErrNoRows)
		result, err := repo.FindByID(id)
		assert.NoError(t, err)
		assert.Nil(t, result)
		mockDB.AssertExpectations(t)
		mockLogger.AssertExpectations(t)
	})
	t.Run("db error", func(t *testing.T) {
		mockDB := new(mocks.Database)
		mockLogger := new(mocks.Logger)
		repo := &TenantRepositoryImpl{db: mockDB, logger: mockLogger, mapper: &mocks.MockTenantMapper{}}
		id := entities.NewTenantID(123)
		dbErr := errors.New("db error")
		mockDB.On("Get", mock.AnythingOfType("*models.Tenant"), mock.Anything, id.Value()).Return(dbErr)
		mockLogger.On("Error", "Failed to find tenant by ID", "id", id.Value(), "error", dbErr).Return()
		result, err := repo.FindByID(id)
		assert.Error(t, err)
		assert.Nil(t, result)
		mockDB.AssertExpectations(t)
		mockLogger.AssertExpectations(t)
	})
}

func TestTenantRepository_FindAll(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockDB := new(mocks.Database)
		repo := &TenantRepositoryImpl{db: mockDB, logger: new(mocks.Logger), mapper: &mocks.MockTenantMapper{}}
		modelList := []*models.Tenant{{Base: models.Base{ID: 1}, Name: "Tenant1"}, {Base: models.Base{ID: 2}, Name: "Tenant2"}}
		mockDB.On("Select", mock.AnythingOfType("*[]*models.Tenant"), mock.Anything).Run(func(args mock.Arguments) {
			tenants := args.Get(0).(*[]*models.Tenant)
			*tenants = modelList
		}).Return(nil)
		mapperMock := repo.mapper.(*mocks.MockTenantMapper)
		mapperMock.On("ToDomains", modelList).Return([]*entities.Tenant{{}, {}}, nil)
		result, err := repo.FindAll()
		assert.NoError(t, err)
		assert.Len(t, result, 2)
		mockDB.AssertExpectations(t)
		mapperMock.AssertExpectations(t)
	})
	t.Run("db error", func(t *testing.T) {
		mockDB := new(mocks.Database)
		mockLogger := new(mocks.Logger)
		repo := &TenantRepositoryImpl{db: mockDB, logger: mockLogger, mapper: &mocks.MockTenantMapper{}}
		dbErr := errors.New("db error")
		mockDB.On("Select", mock.AnythingOfType("*[]*models.Tenant"), mock.Anything).Return(dbErr)
		mockLogger.On("Error", "Failed to find all tenants", "error", dbErr).Return()
		result, err := repo.FindAll()
		assert.Error(t, err)
		assert.Nil(t, result)
		mockDB.AssertExpectations(t)
		mockLogger.AssertExpectations(t)
	})
}

func TestTenantRepository_FindActiveOnly(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockDB := new(mocks.Database)
		repo := &TenantRepositoryImpl{db: mockDB, logger: new(mocks.Logger), mapper: &mocks.MockTenantMapper{}}
		modelList := []*models.Tenant{{Base: models.Base{ID: 1}, Name: "Tenant1"}, {Base: models.Base{ID: 2}, Name: "Tenant2"}}
		mockDB.On("Select", mock.AnythingOfType("*[]*models.Tenant"), mock.Anything, true).Run(func(args mock.Arguments) {
			tenants := args.Get(0).(*[]*models.Tenant)
			*tenants = modelList
		}).Return(nil)
		mapperMock := repo.mapper.(*mocks.MockTenantMapper)
		mapperMock.On("ToDomains", modelList).Return([]*entities.Tenant{{}, {}}, nil)
		result, err := repo.FindActiveOnly()
		assert.NoError(t, err)
		assert.Len(t, result, 2)
		mockDB.AssertExpectations(t)
		mapperMock.AssertExpectations(t)
	})
	t.Run("none found", func(t *testing.T) {
		mockDB := new(mocks.Database)
		mockLogger := new(mocks.Logger)
		repo := &TenantRepositoryImpl{db: mockDB, logger: mockLogger, mapper: &mocks.MockTenantMapper{}}
		mockDB.On("Select", mock.AnythingOfType("*[]*models.Tenant"), mock.Anything, true).Run(func(args mock.Arguments) {
			tenants := args.Get(0).(*[]*models.Tenant)
			*tenants = []*models.Tenant{}
		}).Return(nil)
		mockLogger.On("Info", "No active tenants found").Return()
		result, err := repo.FindActiveOnly()
		assert.NoError(t, err)
		assert.Nil(t, result)
		mockDB.AssertExpectations(t)
		mockLogger.AssertExpectations(t)
	})
	t.Run("db error", func(t *testing.T) {
		mockDB := new(mocks.Database)
		mockLogger := new(mocks.Logger)
		repo := &TenantRepositoryImpl{db: mockDB, logger: mockLogger, mapper: &mocks.MockTenantMapper{}}
		mockDB.On("Select", mock.AnythingOfType("*[]*models.Tenant"), mock.Anything, true).Return(errors.New("db error"))
		mockLogger.On("Error", "Failed to find active tenants", "error", mock.Anything).Return()
		result, err := repo.FindActiveOnly()
		assert.Error(t, err)
		assert.Nil(t, result)
		mockDB.AssertExpectations(t)
		mockLogger.AssertExpectations(t)
	})
}

func TestTenantRepository_Delete(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockDB := new(mocks.Database)
		repo := &TenantRepositoryImpl{db: mockDB, logger: new(mocks.Logger), mapper: &mocks.MockTenantMapper{}}
		id := entities.NewTenantID(1)
		mockResult := new(mocks.SqlResult)
		mockDB.On("Exec", mock.Anything, id.Value()).Return(mockResult, nil)
		err := repo.Delete(id)
		assert.NoError(t, err)
		mockDB.AssertExpectations(t)
	})
	t.Run("query error", func(t *testing.T) {
		mockDB := new(mocks.Database)
		mockLogger := new(mocks.Logger)
		repo := &TenantRepositoryImpl{db: mockDB, logger: mockLogger, mapper: &mocks.MockTenantMapper{}}
		id := entities.NewTenantID(2)
		queryErr := errors.New("query error")
		mockResult := new(mocks.SqlResult)
		mockDB.On("Exec", mock.Anything, id.Value()).Return(mockResult, queryErr)
		mockLogger.On("Error", "Failed to delete tenant", "id", id.Value(), "error", queryErr).Return()
		err := repo.Delete(id)
		assert.Error(t, err)
		mockDB.AssertExpectations(t)
		mockLogger.AssertExpectations(t)
	})
}

func TestTenantRepository_ExistsByName(t *testing.T) {
	t.Run("exists", func(t *testing.T) {
		mockDB := new(mocks.Database)
		repo := &TenantRepositoryImpl{db: mockDB, logger: new(mocks.Logger), mapper: &mocks.MockTenantMapper{}}
		name := "exists"
		mockDB.On("Get", mock.AnythingOfType("*int64"), mock.Anything, name).Run(func(args mock.Arguments) {
			count := args.Get(0).(*int64)
			*count = 1
		}).Return(nil)
		exists, err := repo.ExistsByName(name)
		assert.NoError(t, err)
		assert.True(t, exists)
		mockDB.AssertExpectations(t)
	})
	t.Run("not exists", func(t *testing.T) {
		mockDB := new(mocks.Database)
		repo := &TenantRepositoryImpl{db: mockDB, logger: new(mocks.Logger), mapper: &mocks.MockTenantMapper{}}
		name := "notexists"
		mockDB.On("Get", mock.AnythingOfType("*int64"), mock.Anything, name).Run(func(args mock.Arguments) {
			count := args.Get(0).(*int64)
			*count = 0
		}).Return(nil)
		exists, err := repo.ExistsByName(name)
		assert.NoError(t, err)
		assert.False(t, exists)
		mockDB.AssertExpectations(t)
	})
	t.Run("db error", func(t *testing.T) {
		mockDB := new(mocks.Database)
		mockLogger := new(mocks.Logger)
		repo := &TenantRepositoryImpl{db: mockDB, logger: mockLogger, mapper: &mocks.MockTenantMapper{}}
		name := "error"
		dbErr := errors.New("db error")
		mockDB.On("Get", mock.AnythingOfType("*int64"), mock.Anything, name).Return(dbErr)
		mockLogger.On("Error", "Failed to check if tenant exists by name", "name", name, "error", dbErr).Return()
		exists, err := repo.ExistsByName(name)
		assert.Error(t, err)
		assert.False(t, exists)
		mockDB.AssertExpectations(t)
		mockLogger.AssertExpectations(t)
	})
}

func TestTenantRepository_FindByName(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockDB := new(mocks.Database)
		mockLogger := new(mocks.Logger)
		repo := &TenantRepositoryImpl{db: mockDB, logger: mockLogger, mapper: &mocks.MockTenantMapper{}}
		name := "TestTenant"
		mockDB.On("Get", mock.AnythingOfType("*models.Tenant"), mock.Anything, name).Run(func(args mock.Arguments) {
			tenant := args.Get(0).(*models.Tenant)
			tenant.ID = 1
			tenant.Name = name
			tenant.CreatedAt = time.Now()
			tenant.UpdatedAt = time.Now()
		}).Return(nil)
		expectedTenant := &entities.Tenant{}
		mapperMock := repo.mapper.(*mocks.MockTenantMapper)
		mapperMock.On("ToDomain", mock.AnythingOfType("*models.Tenant")).Return(expectedTenant, nil)
		result, err := repo.FindByName(name)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, expectedTenant, result)
		mockDB.AssertExpectations(t)
		mockLogger.AssertExpectations(t)
		mapperMock.AssertExpectations(t)
	})
	t.Run("not_found", func(t *testing.T) {
		mockDB := new(mocks.Database)
		mockLogger := new(mocks.Logger)
		repo := &TenantRepositoryImpl{db: mockDB, logger: mockLogger, mapper: &mocks.MockTenantMapper{}}
		name := "NotFound"
		mockDB.On("Get", mock.AnythingOfType("*models.Tenant"), mock.Anything, name).Return(sql.ErrNoRows)
		mockLogger.On("Error", "Failed to find tenant by name", "name", name, "error", mock.Anything).Return()
		result, err := repo.FindByName(name)
		assert.NoError(t, err)
		assert.Nil(t, result)
		mockDB.AssertExpectations(t)
		mockLogger.AssertExpectations(t)
	})
	t.Run("db error", func(t *testing.T) {
		mockDB := new(mocks.Database)
		mockLogger := new(mocks.Logger)
		repo := &TenantRepositoryImpl{db: mockDB, logger: mockLogger, mapper: &mocks.MockTenantMapper{}}
		name := "Error"
		dbErr := errors.New("db error")
		mockDB.On("Get", mock.AnythingOfType("*models.Tenant"), mock.Anything, name).Return(dbErr)
		mockLogger.On("Error", "Failed to find tenant by name", "name", name, "error", dbErr).Return()
		result, err := repo.FindByName(name)
		assert.Error(t, err)
		assert.Nil(t, result)
		mockDB.AssertExpectations(t)
		mockLogger.AssertExpectations(t)
	})
}
