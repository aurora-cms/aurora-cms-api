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

func TestPageVersionRepository_Save(t *testing.T) {
	t.Run("insert success", func(t *testing.T) {
		mockDB := new(mocks.Database)
		mockLogger := new(mocks.Logger)
		repo := &PageVersionRepositoryImpl{
			db:     mockDB,
			logger: mockLogger,
			mapper: &mocks.MockPageVersionMapper{},
		}

		version := &entities.PageVersion{}
		model := &models.PageVersion{PageID: 1, Version: 1, IsPublished: true, Base: models.Base{CreatedAt: time.Now(), UpdatedAt: time.Now()}}
		mapperMock := repo.mapper.(*mocks.MockPageVersionMapper)
		mapperMock.On("ToModel", version).Return(model, nil)

		mockResult := new(mocks.SqlResult)
		mockResult.On("LastInsertId").Return(int64(42), nil)
		mockDB.On("Exec", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(mockResult, nil)

		err := repo.Save(version)
		assert.NoError(t, err)
		assert.Equal(t, uint64(42), version.ID().Value())
		mockDB.AssertExpectations(t)
		mockLogger.AssertExpectations(t)
		mapperMock.AssertExpectations(t)
	})

	t.Run("update success", func(t *testing.T) {
		mockDB := new(mocks.Database)
		mockLogger := new(mocks.Logger)
		repo := &PageVersionRepositoryImpl{
			db:     mockDB,
			logger: mockLogger,
			mapper: &mocks.MockPageVersionMapper{},
		}

		version := &entities.PageVersion{}
		version.SetID(entities.NewPageVersionID(99))
		model := &models.PageVersion{Base: models.Base{ID: 99, CreatedAt: time.Now(), UpdatedAt: time.Now()}, PageID: 1, Version: 2, IsPublished: false}
		mapperMock := repo.mapper.(*mocks.MockPageVersionMapper)
		mapperMock.On("ToModel", version).Return(model, nil)

		mockResult := new(mocks.SqlResult)
		mockDB.On("Exec", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(mockResult, nil)

		err := repo.Save(version)
		assert.NoError(t, err)
		mockDB.AssertExpectations(t)
		mockLogger.AssertExpectations(t)
		mapperMock.AssertExpectations(t)
	})

	t.Run("mapper error", func(t *testing.T) {
		mockDB := new(mocks.Database)
		mockLogger := new(mocks.Logger)
		repo := &PageVersionRepositoryImpl{
			db:     mockDB,
			logger: mockLogger,
			mapper: &mocks.MockPageVersionMapper{},
		}

		version := &entities.PageVersion{}
		mapperErr := errors.New("mapper error")
		mapperMock := repo.mapper.(*mocks.MockPageVersionMapper)
		mapperMock.On("ToModel", version).Return(nil, mapperErr)
		mockLogger.On("Error", "Failed to convert page version to model", "error", mapperErr).Return()

		err := repo.Save(version)
		assert.Error(t, err)
		assert.Equal(t, mapperErr, err)
		mockDB.AssertExpectations(t)
		mockLogger.AssertExpectations(t)
		mapperMock.AssertExpectations(t)
	})
}

func TestPageVersionRepository_Save_ErrorBranches(t *testing.T) {
	version := &entities.PageVersion{}
	model := &models.PageVersion{PageID: 1, Version: 1, IsPublished: true, Base: models.Base{CreatedAt: time.Now(), UpdatedAt: time.Now()}}

	t.Run("insert exec error", func(t *testing.T) {
		mockDB := new(mocks.Database)
		mockLogger := new(mocks.Logger)
		repo := &PageVersionRepositoryImpl{db: mockDB, logger: mockLogger, mapper: &mocks.MockPageVersionMapper{}}
		mapperMock := repo.mapper.(*mocks.MockPageVersionMapper)
		mapperMock.On("ToModel", version).Return(model, nil)
		mockDB.On("Exec", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(new(mocks.SqlResult), errors.New("exec error"))
		mockLogger.On("Error", "Failed to create page version", "error", mock.Anything).Return()
		err := repo.Save(version)
		assert.Error(t, err)
		mockDB.AssertExpectations(t)
		mockLogger.AssertExpectations(t)
		mapperMock.AssertExpectations(t)
	})

	t.Run("lastInsertId error", func(t *testing.T) {
		mockDB := new(mocks.Database)
		mockLogger := new(mocks.Logger)
		repo := &PageVersionRepositoryImpl{db: mockDB, logger: mockLogger, mapper: &mocks.MockPageVersionMapper{}}
		mapperMock := repo.mapper.(*mocks.MockPageVersionMapper)
		mapperMock.On("ToModel", version).Return(model, nil)
		mockResult := new(mocks.SqlResult)
		mockResult.On("LastInsertId").Return(int64(0), errors.New("lastInsertId error"))
		mockDB.On("Exec", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(mockResult, nil)
		mockLogger.On("Error", "Failed to get last insert ID for page version", "error", mock.Anything).Return()
		err := repo.Save(version)
		assert.Error(t, err)
		mockDB.AssertExpectations(t)
		mockLogger.AssertExpectations(t)
		mapperMock.AssertExpectations(t)
	})

	t.Run("update exec error", func(t *testing.T) {
		version.SetID(entities.NewPageVersionID(99))
		mockDB := new(mocks.Database)
		mockLogger := new(mocks.Logger)
		repo := &PageVersionRepositoryImpl{db: mockDB, logger: mockLogger, mapper: &mocks.MockPageVersionMapper{}}
		mapperMock := repo.mapper.(*mocks.MockPageVersionMapper)
		mapperMock.On("ToModel", version).Return(&models.PageVersion{Base: models.Base{ID: 99, CreatedAt: time.Now(), UpdatedAt: time.Now()}, PageID: 1, Version: 2, IsPublished: false}, nil)
		mockDB.On("Exec", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(new(mocks.SqlResult), errors.New("exec error"))
		mockLogger.On("Error", "Failed to update page version", "id", uint64(99), "error", mock.Anything).Return()
		err := repo.Save(version)
		assert.Error(t, err)
		mockDB.AssertExpectations(t)
		mockLogger.AssertExpectations(t)
		mapperMock.AssertExpectations(t)
	})
}

func TestPageVersionRepository_FindByID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockDB := new(mocks.Database)
		mockLogger := new(mocks.Logger)
		repo := &PageVersionRepositoryImpl{
			db:     mockDB,
			logger: mockLogger,
			mapper: &mocks.MockPageVersionMapper{},
		}

		id := entities.NewPageVersionID(123)
		mockDB.On("Get", mock.AnythingOfType("*models.PageVersion"), mock.Anything, id.Value()).Run(func(args mock.Arguments) {
			version := args.Get(0).(*models.PageVersion)
			version.ID = 123
			version.PageID = 1
			version.Version = 1
			version.IsPublished = true
			version.CreatedAt = time.Now()
			version.UpdatedAt = time.Now()
		}).Return(nil)

		expectedVersion := &entities.PageVersion{}
		mapperMock := repo.mapper.(*mocks.MockPageVersionMapper)
		mapperMock.On("ToDomain", mock.AnythingOfType("*models.PageVersion")).Return(expectedVersion, nil)

		result, err := repo.FindByID(id)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, expectedVersion, result)
		mockDB.AssertExpectations(t)
		mockLogger.AssertExpectations(t)
		mapperMock.AssertExpectations(t)
	})

	t.Run("not_found", func(t *testing.T) {
		mockDB := new(mocks.Database)
		mockLogger := new(mocks.Logger)
		repo := &PageVersionRepositoryImpl{
			db:     mockDB,
			logger: mockLogger,
			mapper: &mocks.MockPageVersionMapper{},
		}
		id := entities.NewPageVersionID(999)
		mockDB.On("Get", mock.AnythingOfType("*models.PageVersion"), mock.Anything, id.Value()).Return(sql.ErrNoRows)
		result, err := repo.FindByID(id)
		assert.NoError(t, err)
		assert.Nil(t, result)
		mockDB.AssertExpectations(t)
		mockLogger.AssertExpectations(t)
	})

	t.Run("db error", func(t *testing.T) {
		mockDB := new(mocks.Database)
		mockLogger := new(mocks.Logger)
		repo := &PageVersionRepositoryImpl{
			db:     mockDB,
			logger: mockLogger,
			mapper: &mocks.MockPageVersionMapper{},
		}
		id := entities.NewPageVersionID(123)
		dbErr := errors.New("db error")
		mockDB.On("Get", mock.AnythingOfType("*models.PageVersion"), mock.Anything, id.Value()).Return(dbErr)
		mockLogger.On("Error", "Failed to find page version by ID", "id", id.Value(), "error", dbErr).Return()
		result, err := repo.FindByID(id)
		assert.Error(t, err)
		assert.Nil(t, result)
		mockDB.AssertExpectations(t)
		mockLogger.AssertExpectations(t)
	})
}

func TestPageVersionRepository_FindByPageID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockDB := new(mocks.Database)
		mockLogger := new(mocks.Logger)
		repo := &PageVersionRepositoryImpl{db: mockDB, logger: mockLogger, mapper: &mocks.MockPageVersionMapper{}}
		pageID := entities.NewPageID(1)
		modelList := []*models.PageVersion{{Base: models.Base{ID: 1}, PageID: 1, Version: 1}, {Base: models.Base{ID: 2}, PageID: 1, Version: 2}}
		mockDB.On("Select", mock.AnythingOfType("*[]*models.PageVersion"), mock.Anything, pageID.Value()).Run(func(args mock.Arguments) {
			versions := args.Get(0).(*[]*models.PageVersion)
			*versions = modelList
		}).Return(nil)
		mapperMock := repo.mapper.(*mocks.MockPageVersionMapper)
		mapperMock.On("ToDomains", modelList).Return([]*entities.PageVersion{{}, {}}, nil)
		result, err := repo.FindByPageID(pageID)
		assert.NoError(t, err)
		assert.Len(t, result, 2)
		mockDB.AssertExpectations(t)
		mapperMock.AssertExpectations(t)
	})
	t.Run("db error", func(t *testing.T) {
		mockDB := new(mocks.Database)
		mockLogger := new(mocks.Logger)
		repo := &PageVersionRepositoryImpl{db: mockDB, logger: mockLogger, mapper: &mocks.MockPageVersionMapper{}}
		pageID := entities.NewPageID(2)
		dbErr := errors.New("db error")
		mockDB.On("Select", mock.AnythingOfType("*[]*models.PageVersion"), mock.Anything, pageID.Value()).Return(dbErr)
		mockLogger.On("Error", "Failed to find page versions by page ID", "page_id", pageID.Value(), "error", dbErr).Return()
		result, err := repo.FindByPageID(pageID)
		assert.Error(t, err)
		assert.Nil(t, result)
		mockDB.AssertExpectations(t)
		mockLogger.AssertExpectations(t)
	})
}

func TestPageVersionRepository_FindPublishedByPageID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockDB := new(mocks.Database)
		mockLogger := new(mocks.Logger)
		repo := &PageVersionRepositoryImpl{db: mockDB, logger: mockLogger, mapper: &mocks.MockPageVersionMapper{}}
		pageID := entities.NewPageID(1)
		mockDB.On("Get", mock.AnythingOfType("*models.PageVersion"), mock.Anything, true, pageID.Value()).Run(func(args mock.Arguments) {
			version := args.Get(0).(*models.PageVersion)
			version.ID = 1
			version.PageID = pageID.Value()
			version.IsPublished = true
		}).Return(nil)
		expectedVersion := &entities.PageVersion{}
		mapperMock := repo.mapper.(*mocks.MockPageVersionMapper)
		mapperMock.On("ToDomain", mock.AnythingOfType("*models.PageVersion")).Return(expectedVersion, nil)
		result, err := repo.FindPublishedByPageID(pageID)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, expectedVersion, result)
		mockDB.AssertExpectations(t)
		mockLogger.AssertExpectations(t)
		mapperMock.AssertExpectations(t)
	})
	t.Run("not_found", func(t *testing.T) {
		mockDB := new(mocks.Database)
		mockLogger := new(mocks.Logger)
		repo := &PageVersionRepositoryImpl{db: mockDB, logger: mockLogger, mapper: &mocks.MockPageVersionMapper{}}
		pageID := entities.NewPageID(2)
		mockDB.On("Get", mock.AnythingOfType("*models.PageVersion"), mock.Anything, true, pageID.Value()).Return(sql.ErrNoRows)
		result, err := repo.FindPublishedByPageID(pageID)
		assert.NoError(t, err)
		assert.Nil(t, result)
		mockDB.AssertExpectations(t)
		mockLogger.AssertExpectations(t)
	})
	t.Run("db error", func(t *testing.T) {
		mockDB := new(mocks.Database)
		mockLogger := new(mocks.Logger)
		repo := &PageVersionRepositoryImpl{db: mockDB, logger: mockLogger, mapper: &mocks.MockPageVersionMapper{}}
		pageID := entities.NewPageID(3)
		dbErr := errors.New("db error")
		mockDB.On("Get", mock.AnythingOfType("*models.PageVersion"), mock.Anything, true, pageID.Value()).Return(dbErr)
		mockLogger.On("Error", "Failed to find published page version", "page_id", pageID.Value(), "error", dbErr).Return()
		result, err := repo.FindPublishedByPageID(pageID)
		assert.Error(t, err)
		assert.Nil(t, result)
		mockDB.AssertExpectations(t)
		mockLogger.AssertExpectations(t)
	})
}

func TestPageVersionRepository_FindLatestByPageID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockDB := new(mocks.Database)
		mockLogger := new(mocks.Logger)
		repo := &PageVersionRepositoryImpl{db: mockDB, logger: mockLogger, mapper: &mocks.MockPageVersionMapper{}}
		pageID := entities.NewPageID(1)
		mockDB.On("Get", mock.AnythingOfType("*models.PageVersion"), mock.Anything, pageID.Value()).Run(func(args mock.Arguments) {
			version := args.Get(0).(*models.PageVersion)
			version.ID = 2
			version.PageID = pageID.Value()
			version.Version = 2
		}).Return(nil)
		expectedVersion := &entities.PageVersion{}
		mapperMock := repo.mapper.(*mocks.MockPageVersionMapper)
		mapperMock.On("ToDomain", mock.AnythingOfType("*models.PageVersion")).Return(expectedVersion, nil)
		result, err := repo.FindLatestByPageID(pageID)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, expectedVersion, result)
		mockDB.AssertExpectations(t)
		mockLogger.AssertExpectations(t)
		mapperMock.AssertExpectations(t)
	})
	t.Run("not_found", func(t *testing.T) {
		mockDB := new(mocks.Database)
		mockLogger := new(mocks.Logger)
		repo := &PageVersionRepositoryImpl{db: mockDB, logger: mockLogger, mapper: &mocks.MockPageVersionMapper{}}
		pageID := entities.NewPageID(2)
		mockDB.On("Get", mock.AnythingOfType("*models.PageVersion"), mock.Anything, pageID.Value()).Return(sql.ErrNoRows)
		result, err := repo.FindLatestByPageID(pageID)
		assert.NoError(t, err)
		assert.Nil(t, result)
		mockDB.AssertExpectations(t)
		mockLogger.AssertExpectations(t)
	})
	t.Run("db error", func(t *testing.T) {
		mockDB := new(mocks.Database)
		mockLogger := new(mocks.Logger)
		repo := &PageVersionRepositoryImpl{db: mockDB, logger: mockLogger, mapper: &mocks.MockPageVersionMapper{}}
		pageID := entities.NewPageID(3)
		dbErr := errors.New("db error")
		mockDB.On("Get", mock.AnythingOfType("*models.PageVersion"), mock.Anything, pageID.Value()).Return(dbErr)
		mockLogger.On("Error", "Failed to find latest page version", "page_id", pageID.Value(), "error", dbErr).Return()
		result, err := repo.FindLatestByPageID(pageID)
		assert.Error(t, err)
		assert.Nil(t, result)
		mockDB.AssertExpectations(t)
		mockLogger.AssertExpectations(t)
	})
}

func TestPageVersionRepository_Delete(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockDB := new(mocks.Database)
		mockLogger := new(mocks.Logger)
		repo := &PageVersionRepositoryImpl{db: mockDB, logger: mockLogger, mapper: &mocks.MockPageVersionMapper{}}
		id := entities.NewPageVersionID(1)
		mockResult := new(mocks.SqlResult)
		mockDB.On("Exec", mock.Anything, id.Value()).Return(mockResult, nil)
		err := repo.Delete(id)
		assert.NoError(t, err)
		mockDB.AssertExpectations(t)
	})
	t.Run("query error", func(t *testing.T) {
		mockDB := new(mocks.Database)
		mockLogger := new(mocks.Logger)
		repo := &PageVersionRepositoryImpl{db: mockDB, logger: mockLogger, mapper: &mocks.MockPageVersionMapper{}}
		id := entities.NewPageVersionID(2)
		queryErr := errors.New("query error")
		mockResult := new(mocks.SqlResult)
		mockDB.On("Exec", mock.Anything, id.Value()).Return(mockResult, queryErr)
		mockLogger.On("Error", "Failed to delete page version", "id", id.Value(), "error", queryErr).Return()
		err := repo.Delete(id)
		assert.Error(t, err)
		mockDB.AssertExpectations(t)
		mockLogger.AssertExpectations(t)
	})
}
