package repositories

import (
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/h4rdc0m/aurora-api/domain/entities"
	"github.com/h4rdc0m/aurora-api/domain/value_objects"
	"github.com/h4rdc0m/aurora-api/infrastructure/persistence/models"
	"github.com/h4rdc0m/aurora-api/tests/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSiteRepository_Save(t *testing.T) {
	t.Run("insert success", func(t *testing.T) {
		mockDB := new(mocks.Database)
		mockLogger := new(mocks.Logger)
		repo := &SiteRepositoryImpl{db: mockDB, logger: mockLogger, mapper: &mocks.MockSiteMapper{}}
		domain, _ := value_objects.NewDomainName("example.com")
		site := &entities.Site{}
		model := &models.Site{Domain: domain.Value(), Name: "Example", TenantID: 1, Enabled: true, Base: models.Base{CreatedAt: time.Now(), UpdatedAt: time.Now()}}
		mapperMock := repo.mapper.(*mocks.MockSiteMapper)
		mapperMock.On("ToModel", site).Return(model, nil)
		mockResult := new(mocks.SqlResult)
		mockResult.On("LastInsertId").Return(int64(42), nil)
		mockDB.On("Exec", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(mockResult, nil)
		err := repo.Save(site)
		assert.NoError(t, err)
		assert.Equal(t, uint64(42), site.ID().Value())
		mockDB.AssertExpectations(t)
		mockLogger.AssertExpectations(t)
		mapperMock.AssertExpectations(t)
	})

	t.Run("update success", func(t *testing.T) {
		mockDB := new(mocks.Database)
		mockLogger := new(mocks.Logger)
		repo := &SiteRepositoryImpl{db: mockDB, logger: mockLogger, mapper: &mocks.MockSiteMapper{}}
		domain, _ := value_objects.NewDomainName("example.com")
		site := &entities.Site{}
		err := site.SetID(entities.NewSiteID(99))
		assert.NoError(t, err)
		model := &models.Site{Base: models.Base{ID: 99, CreatedAt: time.Now(), UpdatedAt: time.Now()}, Domain: domain.Value(), Name: "Example", TenantID: 1, Enabled: true}
		mapperMock := repo.mapper.(*mocks.MockSiteMapper)
		mapperMock.On("ToModel", site).Return(model, nil)
		mockResult := new(mocks.SqlResult)
		mockDB.On("Exec", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(mockResult, nil)
		err = repo.Save(site)
		assert.NoError(t, err)
		mockDB.AssertExpectations(t)
		mockLogger.AssertExpectations(t)
		mapperMock.AssertExpectations(t)
	})

	t.Run("mapper error", func(t *testing.T) {
		mockDB := new(mocks.Database)
		mockLogger := new(mocks.Logger)
		repo := &SiteRepositoryImpl{db: mockDB, logger: mockLogger, mapper: &mocks.MockSiteMapper{}}
		site := &entities.Site{}
		mapperErr := errors.New("mapper error")
		mapperMock := repo.mapper.(*mocks.MockSiteMapper)
		mapperMock.On("ToModel", site).Return(nil, mapperErr)
		mockLogger.On("Error", "Failed to convert site to model", "error", mapperErr).Return()
		err := repo.Save(site)
		assert.Error(t, err)
		assert.Equal(t, mapperErr, err)
		mockDB.AssertExpectations(t)
		mockLogger.AssertExpectations(t)
		mapperMock.AssertExpectations(t)
	})
}

func TestSiteRepository_Save_ErrorBranches(t *testing.T) {
	site := &entities.Site{}
	model := &models.Site{Domain: "example.com", Name: "Example", TenantID: 1, Enabled: true, Base: models.Base{CreatedAt: time.Now(), UpdatedAt: time.Now()}}

	t.Run("insert exec error", func(t *testing.T) {
		mockDB := new(mocks.Database)
		mockLogger := new(mocks.Logger)
		repo := &SiteRepositoryImpl{db: mockDB, logger: mockLogger, mapper: &mocks.MockSiteMapper{}}
		mapperMock := repo.mapper.(*mocks.MockSiteMapper)
		mapperMock.On("ToModel", site).Return(model, nil)
		mockDB.On("Exec", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(new(mocks.SqlResult), errors.New("exec error"))
		mockLogger.On("Error", "Failed to create site", "error", mock.Anything).Return()
		err := repo.Save(site)
		assert.Error(t, err)
		mockDB.AssertExpectations(t)
		mockLogger.AssertExpectations(t)
		mapperMock.AssertExpectations(t)
	})

	t.Run("lastInsertId error", func(t *testing.T) {
		mockDB := new(mocks.Database)
		mockLogger := new(mocks.Logger)
		repo := &SiteRepositoryImpl{db: mockDB, logger: mockLogger, mapper: &mocks.MockSiteMapper{}}
		mapperMock := repo.mapper.(*mocks.MockSiteMapper)
		mapperMock.On("ToModel", site).Return(model, nil)
		mockResult := new(mocks.SqlResult)
		mockResult.On("LastInsertId").Return(int64(0), errors.New("lastInsertId error"))
		mockDB.On("Exec", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(mockResult, nil)
		mockLogger.On("Error", "Failed to get last insert ID for site", "error", mock.Anything).Return()
		err := repo.Save(site)
		assert.Error(t, err)
		mockDB.AssertExpectations(t)
		mockLogger.AssertExpectations(t)
		mapperMock.AssertExpectations(t)
	})

	t.Run("update exec error", func(t *testing.T) {
		_ = site.SetID(entities.NewSiteID(99))
		mockDB := new(mocks.Database)
		mockLogger := new(mocks.Logger)
		repo := &SiteRepositoryImpl{db: mockDB, logger: mockLogger, mapper: &mocks.MockSiteMapper{}}
		mapperMock := repo.mapper.(*mocks.MockSiteMapper)
		mapperMock.On("ToModel", site).Return(&models.Site{Base: models.Base{ID: 99, CreatedAt: time.Now(), UpdatedAt: time.Now()}, Domain: "example.com", Name: "Example", TenantID: 1, Enabled: true}, nil)
		mockDB.On("Exec", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(new(mocks.SqlResult), errors.New("exec error"))
		mockLogger.On("Error", "Failed to update site", "id", uint64(99), "error", mock.Anything).Return()
		err := repo.Save(site)
		assert.Error(t, err)
		mockDB.AssertExpectations(t)
		mockLogger.AssertExpectations(t)
		mapperMock.AssertExpectations(t)
	})
}

func TestSiteRepository_FindByID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockDB := new(mocks.Database)
		mockLogger := new(mocks.Logger)
		repo := &SiteRepositoryImpl{db: mockDB, logger: mockLogger, mapper: &mocks.MockSiteMapper{}}
		id := entities.NewSiteID(123)
		mockDB.On("Get", mock.AnythingOfType("*models.Site"), mock.Anything, id.Value()).Run(func(args mock.Arguments) {
			site := args.Get(0).(*models.Site)
			site.ID = 123
			site.Domain = "example.com"
			site.Name = "Example"
			site.TenantID = 1
			site.Enabled = true
			site.CreatedAt = time.Now()
			site.UpdatedAt = time.Now()
		}).Return(nil)
		expectedSite := &entities.Site{}
		mapperMock := repo.mapper.(*mocks.MockSiteMapper)
		mapperMock.On("ToDomain", mock.AnythingOfType("*models.Site")).Return(expectedSite, nil)
		result, err := repo.FindByID(id)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, expectedSite, result)
		mockDB.AssertExpectations(t)
		mockLogger.AssertExpectations(t)
		mapperMock.AssertExpectations(t)
	})
	t.Run("not_found", func(t *testing.T) {
		mockDB := new(mocks.Database)
		mockLogger := new(mocks.Logger)
		repo := &SiteRepositoryImpl{db: mockDB, logger: mockLogger, mapper: &mocks.MockSiteMapper{}}
		id := entities.NewSiteID(999)
		mockDB.On("Get", mock.AnythingOfType("*models.Site"), mock.Anything, id.Value()).Return(sql.ErrNoRows)
		result, err := repo.FindByID(id)
		assert.NoError(t, err)
		assert.Nil(t, result)
		mockDB.AssertExpectations(t)
		mockLogger.AssertExpectations(t)
	})
	t.Run("db error", func(t *testing.T) {
		mockDB := new(mocks.Database)
		mockLogger := new(mocks.Logger)
		repo := &SiteRepositoryImpl{db: mockDB, logger: mockLogger, mapper: &mocks.MockSiteMapper{}}
		id := entities.NewSiteID(123)
		dbErr := errors.New("db error")
		mockDB.On("Get", mock.AnythingOfType("*models.Site"), mock.Anything, id.Value()).Return(dbErr)
		mockLogger.On("Error", "Failed to find site by ID", "id", id.Value(), "error", dbErr).Return()
		result, err := repo.FindByID(id)
		assert.Error(t, err)
		assert.Nil(t, result)
		mockDB.AssertExpectations(t)
		mockLogger.AssertExpectations(t)
	})
}

func TestSiteRepository_FindByDomain(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockDB := new(mocks.Database)
		mockLogger := new(mocks.Logger)
		repo := &SiteRepositoryImpl{db: mockDB, logger: mockLogger, mapper: &mocks.MockSiteMapper{}}
		domain, _ := value_objects.NewDomainName("example.com")
		mockDB.On("Get", mock.AnythingOfType("*models.Site"), mock.Anything, domain.Value()).Run(func(args mock.Arguments) {
			site := args.Get(0).(*models.Site)
			site.ID = 1
			site.Domain = domain.Value()
			site.Name = "Example"
			site.TenantID = 1
			site.Enabled = true
		}).Return(nil)
		expectedSite := &entities.Site{}
		mapperMock := repo.mapper.(*mocks.MockSiteMapper)
		mapperMock.On("ToDomain", mock.AnythingOfType("*models.Site")).Return(expectedSite, nil)
		result, err := repo.FindByDomain(domain)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, expectedSite, result)
		mockDB.AssertExpectations(t)
		mockLogger.AssertExpectations(t)
		mapperMock.AssertExpectations(t)
	})
	t.Run("not_found", func(t *testing.T) {
		mockDB := new(mocks.Database)
		mockLogger := new(mocks.Logger)
		repo := &SiteRepositoryImpl{db: mockDB, logger: mockLogger, mapper: &mocks.MockSiteMapper{}}
		domain, _ := value_objects.NewDomainName("notfound.com")
		mockDB.On("Get", mock.AnythingOfType("*models.Site"), mock.Anything, domain.Value()).Return(sql.ErrNoRows)
		result, err := repo.FindByDomain(domain)
		assert.NoError(t, err)
		assert.Nil(t, result)
		mockDB.AssertExpectations(t)
		mockLogger.AssertExpectations(t)
	})

	t.Run("db error", func(t *testing.T) {
		mockDB := new(mocks.Database)
		mockLogger := new(mocks.Logger)
		repo := &SiteRepositoryImpl{db: mockDB, logger: mockLogger, mapper: &mocks.MockSiteMapper{}}
		domain, _ := value_objects.NewDomainName("error.com")
		dbErr := errors.New("db error")
		mockDB.On("Get", mock.AnythingOfType("*models.Site"), mock.Anything, domain.Value()).Return(dbErr)
		mockLogger.On("Error", "Failed to find site by domain", "domain", domain.Value(), "error", dbErr).Return()
		result, err := repo.FindByDomain(domain)
		assert.Error(t, err)
		assert.Nil(t, result)
		mockDB.AssertExpectations(t)
		mockLogger.AssertExpectations(t)
	})
}

func TestSiteRepository_FindByTenantID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockDB := new(mocks.Database)
		mockLogger := new(mocks.Logger)
		repo := &SiteRepositoryImpl{db: mockDB, logger: mockLogger, mapper: &mocks.MockSiteMapper{}}
		tenantID := entities.NewTenantID(1)
		modelList := []*models.Site{{Base: models.Base{ID: 1}, Domain: "example.com", Name: "Example", TenantID: 1, Enabled: true}, {Base: models.Base{ID: 2}, Domain: "test.com", Name: "Test", TenantID: 1, Enabled: true}}
		mockDB.On("Select", mock.AnythingOfType("*[]*models.Site"), mock.Anything, tenantID.Value()).Run(func(args mock.Arguments) {
			sites := args.Get(0).(*[]*models.Site)
			*sites = modelList
		}).Return(nil)
		mapperMock := repo.mapper.(*mocks.MockSiteMapper)
		mapperMock.On("ToDomains", modelList).Return([]*entities.Site{{}, {}}, nil)
		result, err := repo.FindByTenantID(tenantID)
		assert.NoError(t, err)
		assert.Len(t, result, 2)
		mockDB.AssertExpectations(t)
		mapperMock.AssertExpectations(t)
	})
	t.Run("db error", func(t *testing.T) {
		mockDB := new(mocks.Database)
		mockLogger := new(mocks.Logger)
		repo := &SiteRepositoryImpl{db: mockDB, logger: mockLogger, mapper: &mocks.MockSiteMapper{}}
		tenantID := entities.NewTenantID(2)
		dbErr := errors.New("db error")
		mockDB.On("Select", mock.AnythingOfType("*[]*models.Site"), mock.Anything, tenantID.Value()).Return(dbErr)
		mockLogger.On("Error", "Failed to find sites by tenant ID", "tenant_id", tenantID.Value(), "error", dbErr).Return()
		result, err := repo.FindByTenantID(tenantID)
		assert.Error(t, err)
		assert.Nil(t, result)
		mockDB.AssertExpectations(t)
		mockLogger.AssertExpectations(t)
	})
}

func TestSiteRepository_FindAll(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockDB := new(mocks.Database)
		repo := &SiteRepositoryImpl{db: mockDB, logger: new(mocks.Logger), mapper: &mocks.MockSiteMapper{}}
		modelList := []*models.Site{{Base: models.Base{ID: 1}, Domain: "example.com", Name: "Example", TenantID: 1, Enabled: true}, {Base: models.Base{ID: 2}, Domain: "test.com", Name: "Test", TenantID: 1, Enabled: true}}
		mockDB.On("Select", mock.AnythingOfType("*[]*models.Site"), mock.Anything).Run(func(args mock.Arguments) {
			sites := args.Get(0).(*[]*models.Site)
			*sites = modelList
		}).Return(nil)
		mapperMock := repo.mapper.(*mocks.MockSiteMapper)
		mapperMock.On("ToDomains", modelList).Return([]*entities.Site{{}, {}}, nil)
		result, err := repo.FindAll()
		assert.NoError(t, err)
		assert.Len(t, result, 2)
		mockDB.AssertExpectations(t)
		mapperMock.AssertExpectations(t)
	})
	t.Run("db error", func(t *testing.T) {
		mockDB := new(mocks.Database)
		mockLogger := new(mocks.Logger)
		repo := &SiteRepositoryImpl{db: mockDB, logger: mockLogger, mapper: &mocks.MockSiteMapper{}}
		dbErr := errors.New("db error")
		mockDB.On("Select", mock.AnythingOfType("*[]*models.Site"), mock.Anything).Return(dbErr)
		mockLogger.On("Error", "Failed to find all sites", "error", dbErr).Return()
		result, err := repo.FindAll()
		assert.Error(t, err)
		assert.Nil(t, result)
		mockDB.AssertExpectations(t)
		mockLogger.AssertExpectations(t)
	})
}

func TestSiteRepository_FindEnabledByTenantID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockDB := new(mocks.Database)
		repo := &SiteRepositoryImpl{db: mockDB, logger: new(mocks.Logger), mapper: &mocks.MockSiteMapper{}}
		tenantID := entities.NewTenantID(1)
		modelList := []*models.Site{{Base: models.Base{ID: 1}, Domain: "example.com", Name: "Example", TenantID: 1, Enabled: true}, {Base: models.Base{ID: 2}, Domain: "test.com", Name: "Test", TenantID: 1, Enabled: true}}
		mockDB.On("Select", mock.AnythingOfType("*[]*models.Site"), mock.Anything, true, tenantID.Value()).Run(func(args mock.Arguments) {
			sites := args.Get(0).(*[]*models.Site)
			*sites = modelList
		}).Return(nil)
		mapperMock := repo.mapper.(*mocks.MockSiteMapper)
		mapperMock.On("ToDomains", modelList).Return([]*entities.Site{{}, {}}, nil)
		result, err := repo.FindEnabledByTenantID(tenantID)
		assert.NoError(t, err)
		assert.Len(t, result, 2)
		mockDB.AssertExpectations(t)
		mapperMock.AssertExpectations(t)
	})
	t.Run("db error", func(t *testing.T) {
		mockDB := new(mocks.Database)
		mockLogger := new(mocks.Logger)
		repo := &SiteRepositoryImpl{db: mockDB, logger: mockLogger, mapper: &mocks.MockSiteMapper{}}
		tenantID := entities.NewTenantID(2)
		dbErr := errors.New("db error")
		mockDB.On("Select", mock.AnythingOfType("*[]*models.Site"), mock.Anything, true, tenantID.Value()).Return(dbErr)
		mockLogger.On("Error", "Failed to find enabled sites by tenant", "tenant_id", tenantID.Value(), "error", dbErr).Return()
		result, err := repo.FindEnabledByTenantID(tenantID)
		assert.Error(t, err)
		assert.Nil(t, result)
		mockDB.AssertExpectations(t)
		mockLogger.AssertExpectations(t)
	})
}

func TestSiteRepository_Delete(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockDB := new(mocks.Database)
		repo := &SiteRepositoryImpl{db: mockDB, logger: new(mocks.Logger), mapper: &mocks.MockSiteMapper{}}
		id := entities.NewSiteID(1)
		mockResult := new(mocks.SqlResult)
		mockDB.On("Exec", mock.Anything, id.Value()).Return(mockResult, nil)
		err := repo.Delete(id)
		assert.NoError(t, err)
		mockDB.AssertExpectations(t)
	})
	t.Run("query error", func(t *testing.T) {
		mockDB := new(mocks.Database)
		mockLogger := new(mocks.Logger)
		repo := &SiteRepositoryImpl{db: mockDB, logger: mockLogger, mapper: &mocks.MockSiteMapper{}}
		id := entities.NewSiteID(2)
		queryErr := errors.New("query error")
		mockResult := new(mocks.SqlResult)
		mockDB.On("Exec", mock.Anything, id.Value()).Return(mockResult, queryErr)
		mockLogger.On("Error", "Failed to delete site", "id", id.Value(), "error", queryErr).Return()
		err := repo.Delete(id)
		assert.Error(t, err)
		mockDB.AssertExpectations(t)
		mockLogger.AssertExpectations(t)
	})
}

func TestSiteRepository_ExistsByDomain(t *testing.T) {
	t.Run("exists", func(t *testing.T) {
		mockDB := new(mocks.Database)
		repo := &SiteRepositoryImpl{db: mockDB, logger: new(mocks.Logger), mapper: &mocks.MockSiteMapper{}}
		domain, _ := value_objects.NewDomainName("exists.com")
		mockDB.On("Get", mock.AnythingOfType("*int64"), mock.Anything, domain.Value()).Run(func(args mock.Arguments) {
			count := args.Get(0).(*int64)
			*count = 1
		}).Return(nil)
		exists, err := repo.ExistsByDomain(domain)
		assert.NoError(t, err)
		assert.True(t, exists)
		mockDB.AssertExpectations(t)
	})
	t.Run("not exists", func(t *testing.T) {
		mockDB := new(mocks.Database)
		repo := &SiteRepositoryImpl{db: mockDB, logger: new(mocks.Logger), mapper: &mocks.MockSiteMapper{}}
		domain, _ := value_objects.NewDomainName("notexists.com")
		mockDB.On("Get", mock.AnythingOfType("*int64"), mock.Anything, domain.Value()).Run(func(args mock.Arguments) {
			count := args.Get(0).(*int64)
			*count = 0
		}).Return(nil)
		exists, err := repo.ExistsByDomain(domain)
		assert.NoError(t, err)
		assert.False(t, exists)
		mockDB.AssertExpectations(t)
	})
	t.Run("db error", func(t *testing.T) {
		mockDB := new(mocks.Database)
		mockLogger := new(mocks.Logger)
		repo := &SiteRepositoryImpl{db: mockDB, logger: mockLogger, mapper: &mocks.MockSiteMapper{}}
		domain, _ := value_objects.NewDomainName("error.com")
		dbErr := errors.New("db error")
		mockDB.On("Get", mock.AnythingOfType("*int64"), mock.Anything, domain.Value()).Return(dbErr)
		mockLogger.On("Error", "Failed to check site existence by domain", "domain", domain.Value(), "error", dbErr).Return()
		exists, err := repo.ExistsByDomain(domain)
		assert.Error(t, err)
		assert.False(t, exists)
		mockDB.AssertExpectations(t)
		mockLogger.AssertExpectations(t)
	})
}
