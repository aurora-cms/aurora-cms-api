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

func TestPageBlockRepository_Save(t *testing.T) {
	t.Run("insert success", func(t *testing.T) {
		mockDB := new(mocks.Database)
		mockLogger := new(mocks.Logger)
		repo := &PageBlockRepositoryImpl{
			db:     mockDB,
			logger: mockLogger,
			mapper: &mocks.MockPageBlockMapper{},
		}

		block := &entities.PageBlock{}
		model := &models.PageBlock{BlockKey: "block1", PageVersionID: 1, Index: 0, ContentType: "type", Content: "content"}
		mapperMock := repo.mapper.(*mocks.MockPageBlockMapper)
		mapperMock.On("ToModel", block).Return(model, nil)

		mockResult := new(mocks.SqlResult)
		mockResult.On("LastInsertId").Return(int64(42), nil)
		mockDB.On("Exec", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(mockResult, nil)

		err := repo.Save(block)
		assert.NoError(t, err)
		assert.Equal(t, uint64(42), block.ID().Value())
		mockDB.AssertExpectations(t)
		mockLogger.AssertExpectations(t)
		mapperMock.AssertExpectations(t)
	})

	t.Run("update success", func(t *testing.T) {
		mockDB := new(mocks.Database)
		mockLogger := new(mocks.Logger)
		repo := &PageBlockRepositoryImpl{
			db:     mockDB,
			logger: mockLogger,
			mapper: &mocks.MockPageBlockMapper{},
		}

		block := &entities.PageBlock{}
		block.SetID(entities.NewPageBlockID(99))
		model := &models.PageBlock{Base: models.Base{ID: 99}, BlockKey: "block1", PageVersionID: 1, Index: 0, ContentType: "type", Content: "content"}
		mapperMock := repo.mapper.(*mocks.MockPageBlockMapper)
		mapperMock.On("ToModel", block).Return(model, nil)

		mockResult := new(mocks.SqlResult)
		mockDB.On("Exec", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(mockResult, nil)

		err := repo.Save(block)
		assert.NoError(t, err)
		mockDB.AssertExpectations(t)
		mockLogger.AssertExpectations(t)
		mapperMock.AssertExpectations(t)
	})

	t.Run("mapper error", func(t *testing.T) {
		mockDB := new(mocks.Database)
		mockLogger := new(mocks.Logger)
		repo := &PageBlockRepositoryImpl{
			db:     mockDB,
			logger: mockLogger,
			mapper: &mocks.MockPageBlockMapper{},
		}

		block := &entities.PageBlock{}
		mapperErr := errors.New("mapper error")
		mapperMock := repo.mapper.(*mocks.MockPageBlockMapper)
		mapperMock.On("ToModel", block).Return(nil, mapperErr)
		mockLogger.On("Error", "Failed to convert page block to model", "error", mapperErr).Return()

		err := repo.Save(block)
		assert.Error(t, err)
		assert.Equal(t, mapperErr, err)
		mockDB.AssertExpectations(t)
		mockLogger.AssertExpectations(t)
		mapperMock.AssertExpectations(t)
	})
}

func TestPageBlockRepository_Save_ErrorBranches(t *testing.T) {
	block := &entities.PageBlock{}
	model := &models.PageBlock{BlockKey: "block1", PageVersionID: 1, Index: 0, ContentType: "type", Content: "content"}

	t.Run("insert exec error", func(t *testing.T) {
		mockDB := new(mocks.Database)
		mockLogger := new(mocks.Logger)
		repo := &PageBlockRepositoryImpl{db: mockDB, logger: mockLogger, mapper: &mocks.MockPageBlockMapper{}}
		mapperMock := repo.mapper.(*mocks.MockPageBlockMapper)
		mapperMock.On("ToModel", block).Return(model, nil)
		mockDB.On("Exec", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(new(mocks.SqlResult), errors.New("exec error"))
		mockLogger.On("Error", "Failed to insert new page block", "error", mock.Anything).Return()
		err := repo.Save(block)
		assert.Error(t, err)
		mockDB.AssertExpectations(t)
		mockLogger.AssertExpectations(t)
		mapperMock.AssertExpectations(t)
	})

	t.Run("lastInsertId error", func(t *testing.T) {
		mockDB := new(mocks.Database)
		mockLogger := new(mocks.Logger)
		repo := &PageBlockRepositoryImpl{db: mockDB, logger: mockLogger, mapper: &mocks.MockPageBlockMapper{}}
		mapperMock := repo.mapper.(*mocks.MockPageBlockMapper)
		mapperMock.On("ToModel", block).Return(model, nil)
		mockResult := new(mocks.SqlResult)
		mockResult.On("LastInsertId").Return(int64(0), errors.New("lastInsertId error"))
		mockDB.On("Exec", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(mockResult, nil)
		mockLogger.On("Error", "Failed to get last insert ID for page block", "error", mock.Anything).Return()
		err := repo.Save(block)
		assert.Error(t, err)
		mockDB.AssertExpectations(t)
		mockLogger.AssertExpectations(t)
		mapperMock.AssertExpectations(t)
	})
}

func TestPageBlockRepository_Save_UpdateErrorBranch(t *testing.T) {
	block := &entities.PageBlock{}
	block.SetID(entities.NewPageBlockID(99))
	model := &models.PageBlock{Base: models.Base{ID: 99}, BlockKey: "block1", PageVersionID: 1, Index: 0, ContentType: "type", Content: "content"}

	mockDB := new(mocks.Database)
	mockLogger := new(mocks.Logger)
	repo := &PageBlockRepositoryImpl{db: mockDB, logger: mockLogger, mapper: &mocks.MockPageBlockMapper{}}
	mapperMock := repo.mapper.(*mocks.MockPageBlockMapper)
	mapperMock.On("ToModel", block).Return(model, nil)
	mockDB.On("Exec", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(new(mocks.SqlResult), errors.New("exec error"))
	mockLogger.On("Error", "Failed to update page block", "id", model.ID, "error", mock.Anything).Return()

	err := repo.Save(block)
	assert.Error(t, err)
	mockDB.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
	mapperMock.AssertExpectations(t)
}

func TestPageBlockRepository_FindByID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockDB := new(mocks.Database)
		mockLogger := new(mocks.Logger)
		repo := &PageBlockRepositoryImpl{
			db:     mockDB,
			logger: mockLogger,
			mapper: &mocks.MockPageBlockMapper{},
		}

		id := entities.NewPageBlockID(123)
		mockDB.On("Get", mock.AnythingOfType("*models.PageBlock"), mock.Anything, id.Value()).Run(func(args mock.Arguments) {
			block := args.Get(0).(*models.PageBlock)
			block.ID = 123
			block.BlockKey = "block1"
			block.PageVersionID = 1
			block.Index = 0
			block.ContentType = "type"
			block.Content = "content"
			block.CreatedAt = time.Now()
			block.UpdatedAt = time.Now()
		}).Return(nil)

		expectedBlock := &entities.PageBlock{}
		mapperMock := repo.mapper.(*mocks.MockPageBlockMapper)
		mapperMock.On("ToDomain", mock.AnythingOfType("*models.PageBlock")).Return(expectedBlock, nil)

		result, err := repo.FindByID(id)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, expectedBlock, result)
		mockDB.AssertExpectations(t)
		mockLogger.AssertExpectations(t)
		mapperMock.AssertExpectations(t)
	})

	t.Run("not_found", func(t *testing.T) {
		mockDB := new(mocks.Database)
		mockLogger := new(mocks.Logger)
		repo := &PageBlockRepositoryImpl{
			db:     mockDB,
			logger: mockLogger,
			mapper: &mocks.MockPageBlockMapper{},
		}
		id := entities.NewPageBlockID(999)
		mockDB.On("Get", mock.AnythingOfType("*models.PageBlock"), mock.Anything, id.Value()).Return(sql.ErrNoRows)
		result, err := repo.FindByID(id)
		assert.NoError(t, err)
		assert.Nil(t, result)
		mockDB.AssertExpectations(t)
		mockLogger.AssertExpectations(t)
	})

	t.Run("db error", func(t *testing.T) {
		mockDB := new(mocks.Database)
		mockLogger := new(mocks.Logger)
		repo := &PageBlockRepositoryImpl{
			db:     mockDB,
			logger: mockLogger,
			mapper: &mocks.MockPageBlockMapper{},
		}
		id := entities.NewPageBlockID(123)
		dbErr := errors.New("db error")
		mockDB.On("Get", mock.AnythingOfType("*models.PageBlock"), mock.Anything, id.Value()).Return(dbErr)
		mockLogger.On("Error", "Failed to find page block by ID", "id", id.Value(), "error", dbErr).Return()
		result, err := repo.FindByID(id)
		assert.Error(t, err)
		assert.Nil(t, result)
		mockDB.AssertExpectations(t)
		mockLogger.AssertExpectations(t)
	})
}

func TestPageBlockRepository_FindByPageVersionID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockDB := new(mocks.Database)
		mockLogger := new(mocks.Logger)
		repo := &PageBlockRepositoryImpl{db: mockDB, logger: mockLogger, mapper: &mocks.MockPageBlockMapper{}}
		pageVersionID := entities.NewPageVersionID(1)
		modelList := []*models.PageBlock{{Base: models.Base{ID: 1}, BlockKey: "block1", PageVersionID: 1}, {Base: models.Base{ID: 2}, BlockKey: "block2", PageVersionID: 1}}
		mockDB.On("Select", mock.AnythingOfType("*[]*models.PageBlock"), mock.Anything, pageVersionID.Value()).Run(func(args mock.Arguments) {
			blocks := args.Get(0).(*[]*models.PageBlock)
			*blocks = modelList
		}).Return(nil)
		mapperMock := repo.mapper.(*mocks.MockPageBlockMapper)
		mapperMock.On("ToDomains", modelList).Return([]*entities.PageBlock{{}, {}}, nil)
		result, err := repo.FindByPageVersionID(pageVersionID)
		assert.NoError(t, err)
		assert.Len(t, result, 2)
		mockDB.AssertExpectations(t)
		mapperMock.AssertExpectations(t)
	})
	t.Run("db error", func(t *testing.T) {
		mockDB := new(mocks.Database)
		mockLogger := new(mocks.Logger)
		repo := &PageBlockRepositoryImpl{db: mockDB, logger: mockLogger, mapper: &mocks.MockPageBlockMapper{}}
		pageVersionID := entities.NewPageVersionID(2)
		dbErr := errors.New("db error")
		mockDB.On("Select", mock.AnythingOfType("*[]*models.PageBlock"), mock.Anything, pageVersionID.Value()).Return(dbErr)
		mockLogger.On("Error", "Failed to find page blocks by page version ID", "page_version_id", pageVersionID.Value(), "error", dbErr).Return()
		result, err := repo.FindByPageVersionID(pageVersionID)
		assert.Error(t, err)
		assert.Nil(t, result)
		mockDB.AssertExpectations(t)
		mockLogger.AssertExpectations(t)
	})
}

func TestPageBlockRepository_FindByBlockKey(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockDB := new(mocks.Database)
		mockLogger := new(mocks.Logger)
		repo := &PageBlockRepositoryImpl{db: mockDB, logger: mockLogger, mapper: &mocks.MockPageBlockMapper{}}
		blockKey := "block1"
		pageVersionID := entities.NewPageVersionID(1)
		mockDB.On("Get", mock.AnythingOfType("*models.PageBlock"), mock.Anything, blockKey, pageVersionID.Value()).Run(func(args mock.Arguments) {
			block := args.Get(0).(*models.PageBlock)
			block.ID = 1
			block.BlockKey = blockKey
			block.PageVersionID = pageVersionID.Value()
		}).Return(nil)
		expectedBlock := &entities.PageBlock{}
		mapperMock := repo.mapper.(*mocks.MockPageBlockMapper)
		mapperMock.On("ToDomain", mock.AnythingOfType("*models.PageBlock")).Return(expectedBlock, nil)
		result, err := repo.FindByBlockKey(blockKey, pageVersionID)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, expectedBlock, result)
		mockDB.AssertExpectations(t)
		mockLogger.AssertExpectations(t)
		mapperMock.AssertExpectations(t)
	})
	t.Run("not_found", func(t *testing.T) {
		mockDB := new(mocks.Database)
		mockLogger := new(mocks.Logger)
		repo := &PageBlockRepositoryImpl{db: mockDB, logger: mockLogger, mapper: &mocks.MockPageBlockMapper{}}
		blockKey := "block2"
		pageVersionID := entities.NewPageVersionID(2)
		mockDB.On("Get", mock.AnythingOfType("*models.PageBlock"), mock.Anything, blockKey, pageVersionID.Value()).Return(sql.ErrNoRows)
		result, err := repo.FindByBlockKey(blockKey, pageVersionID)
		assert.NoError(t, err)
		assert.Nil(t, result)
		mockDB.AssertExpectations(t)
		mockLogger.AssertExpectations(t)
	})
	t.Run("db error", func(t *testing.T) {
		mockDB := new(mocks.Database)
		mockLogger := new(mocks.Logger)
		repo := &PageBlockRepositoryImpl{db: mockDB, logger: mockLogger, mapper: &mocks.MockPageBlockMapper{}}
		blockKey := "block3"
		pageVersionID := entities.NewPageVersionID(3)
		dbErr := errors.New("db error")
		mockDB.On("Get", mock.AnythingOfType("*models.PageBlock"), mock.Anything, blockKey, pageVersionID.Value()).Return(dbErr)
		mockLogger.On("Error", "Failed to find page block by key", "block_key", blockKey, "page_version_id", pageVersionID.Value(), "error", dbErr).Return()
		result, err := repo.FindByBlockKey(blockKey, pageVersionID)
		assert.Error(t, err)
		assert.Nil(t, result)
		mockDB.AssertExpectations(t)
		mockLogger.AssertExpectations(t)
	})
}

func TestPageBlockRepository_Delete(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockDB := new(mocks.Database)
		mockLogger := new(mocks.Logger)
		repo := &PageBlockRepositoryImpl{db: mockDB, logger: mockLogger, mapper: &mocks.MockPageBlockMapper{}}
		id := entities.NewPageBlockID(1)
		mockResult := new(mocks.SqlResult)
		mockDB.On("Exec", mock.Anything, id.Value()).Return(mockResult, nil)
		err := repo.Delete(id)
		assert.NoError(t, err)
		mockDB.AssertExpectations(t)
	})
	t.Run("query error", func(t *testing.T) {
		mockDB := new(mocks.Database)
		mockLogger := new(mocks.Logger)
		repo := &PageBlockRepositoryImpl{db: mockDB, logger: mockLogger, mapper: &mocks.MockPageBlockMapper{}}
		id := entities.NewPageBlockID(2)
		queryErr := errors.New("query error")
		mockResult := new(mocks.SqlResult)
		mockDB.On("Exec", mock.Anything, id.Value()).Return(mockResult, queryErr)
		mockLogger.On("Error", "Failed to delete page block", "id", id.Value(), "error", queryErr).Return()
		err := repo.Delete(id)
		assert.Error(t, err)
		mockDB.AssertExpectations(t)
		mockLogger.AssertExpectations(t)
	})
}

func TestPageBlockRepository_DeleteByPageVersionID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockDB := new(mocks.Database)
		mockLogger := new(mocks.Logger)
		repo := &PageBlockRepositoryImpl{db: mockDB, logger: mockLogger, mapper: &mocks.MockPageBlockMapper{}}
		pageVersionID := entities.NewPageVersionID(1)
		mockResult := new(mocks.SqlResult)
		mockDB.On("Exec", mock.Anything, pageVersionID.Value()).Return(mockResult, nil)
		err := repo.DeleteByPageVersionID(pageVersionID)
		assert.NoError(t, err)
		mockDB.AssertExpectations(t)
	})
	t.Run("query error", func(t *testing.T) {
		mockDB := new(mocks.Database)
		mockLogger := new(mocks.Logger)
		repo := &PageBlockRepositoryImpl{db: mockDB, logger: mockLogger, mapper: &mocks.MockPageBlockMapper{}}
		pageVersionID := entities.NewPageVersionID(2)
		queryErr := errors.New("query error")
		mockResult := new(mocks.SqlResult)
		mockDB.On("Exec", mock.Anything, pageVersionID.Value()).Return(mockResult, queryErr)
		mockLogger.On("Error", "Failed to delete page blocks by page version ID", "page_version_id", pageVersionID.Value(), "error", queryErr).Return()
		err := repo.DeleteByPageVersionID(pageVersionID)
		assert.Error(t, err)
		mockDB.AssertExpectations(t)
		mockLogger.AssertExpectations(t)
	})
}
