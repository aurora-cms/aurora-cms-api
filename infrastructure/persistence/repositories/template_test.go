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

func TestTemplateRepository_Save(t *testing.T) {
	t.Run("insert success", func(t *testing.T) {
		mockDB := new(mocks.Database)
		mockLogger := new(mocks.Logger)
		repo := &TemplateRepositoryImpl{db: mockDB, logger: mockLogger, mapper: &mocks.MockTemplateMapper{}}
		template := &entities.Template{}
		model := &models.Template{Name: "Test", Description: nil, FilePath: "content", Base: models.Base{CreatedAt: time.Now(), UpdatedAt: time.Now()}}
		mapperMock := repo.mapper.(*mocks.MockTemplateMapper)
		mapperMock.On("ToModel", template).Return(model, nil)
		mockResult := new(mocks.SqlResult)
		mockResult.On("LastInsertId").Return(int64(42), nil)
		mockDB.On("Exec", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(mockResult, nil)
		err := repo.Save(template)
		assert.NoError(t, err)
		assert.Equal(t, uint64(42), template.ID().Value())
		mockDB.AssertExpectations(t)
		mockLogger.AssertExpectations(t)
		mapperMock.AssertExpectations(t)
	})

	t.Run("update success", func(t *testing.T) {
		mockDB := new(mocks.Database)
		mockLogger := new(mocks.Logger)
		repo := &TemplateRepositoryImpl{db: mockDB, logger: mockLogger, mapper: &mocks.MockTemplateMapper{}}
		template := &entities.Template{}
		template.SetID(entities.NewTemplateID(99))
		model := &models.Template{Name: "Test", Description: nil, FilePath: "content", Base: models.Base{ID: 99, CreatedAt: time.Now(), UpdatedAt: time.Now()}}
		mapperMock := repo.mapper.(*mocks.MockTemplateMapper)
		mapperMock.On("ToModel", template).Return(model, nil)
		mockResult := new(mocks.SqlResult)
		mockDB.On("Exec", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(mockResult, nil)
		err := repo.Save(template)
		assert.NoError(t, err)
		mockDB.AssertExpectations(t)
		mockLogger.AssertExpectations(t)
		mapperMock.AssertExpectations(t)
	})

	t.Run("mapper error", func(t *testing.T) {
		mockDB := new(mocks.Database)
		mockLogger := new(mocks.Logger)
		repo := &TemplateRepositoryImpl{db: mockDB, logger: mockLogger, mapper: &mocks.MockTemplateMapper{}}
		template := &entities.Template{}
		mapperErr := errors.New("mapper error")
		mapperMock := repo.mapper.(*mocks.MockTemplateMapper)
		mapperMock.On("ToModel", template).Return(nil, mapperErr)
		mockLogger.On("Error", "Failed to convert template to model", "error", mapperErr).Return()
		err := repo.Save(template)
		assert.Error(t, err)
		assert.Equal(t, mapperErr, err)
		mockDB.AssertExpectations(t)
		mockLogger.AssertExpectations(t)
		mapperMock.AssertExpectations(t)
	})
}

func TestTemplateRepository_Save_ErrorBranches(t *testing.T) {
	template := &entities.Template{}
	model := &models.Template{Name: "Test", Description: nil, FilePath: "content", Base: models.Base{CreatedAt: time.Now(), UpdatedAt: time.Now()}}

	t.Run("insert exec error", func(t *testing.T) {
		mockDB := new(mocks.Database)
		mockLogger := new(mocks.Logger)
		repo := &TemplateRepositoryImpl{db: mockDB, logger: mockLogger, mapper: &mocks.MockTemplateMapper{}}
		mapperMock := repo.mapper.(*mocks.MockTemplateMapper)
		mapperMock.On("ToModel", template).Return(model, nil)
		mockDB.On("Exec", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(new(mocks.SqlResult), errors.New("exec error"))
		mockLogger.On("Error", "Failed to create template", "error", mock.Anything).Return()
		err := repo.Save(template)
		assert.Error(t, err)
		mockDB.AssertExpectations(t)
		mockLogger.AssertExpectations(t)
		mapperMock.AssertExpectations(t)
	})

	t.Run("lastInsertId error", func(t *testing.T) {
		mockDB := new(mocks.Database)
		mockLogger := new(mocks.Logger)
		repo := &TemplateRepositoryImpl{db: mockDB, logger: mockLogger, mapper: &mocks.MockTemplateMapper{}}
		mapperMock := repo.mapper.(*mocks.MockTemplateMapper)
		mapperMock.On("ToModel", template).Return(model, nil)
		mockResult := new(mocks.SqlResult)
		mockResult.On("LastInsertId").Return(int64(0), errors.New("lastInsertId error"))
		mockDB.On("Exec", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(mockResult, nil)
		mockLogger.On("Error", "Failed to get last insert ID for template", "error", mock.Anything).Return()
		err := repo.Save(template)
		assert.Error(t, err)
		mockDB.AssertExpectations(t)
		mockLogger.AssertExpectations(t)
		mapperMock.AssertExpectations(t)
	})

	t.Run("update exec error", func(t *testing.T) {
		template.SetID(entities.NewTemplateID(99))
		mockDB := new(mocks.Database)
		mockLogger := new(mocks.Logger)
		repo := &TemplateRepositoryImpl{db: mockDB, logger: mockLogger, mapper: &mocks.MockTemplateMapper{}}
		mapperMock := repo.mapper.(*mocks.MockTemplateMapper)
		mapperMock.On("ToModel", template).Return(&models.Template{Name: "Test", Description: nil, FilePath: "content", Base: models.Base{ID: 99, CreatedAt: time.Now(), UpdatedAt: time.Now()}}, nil)
		mockDB.On("Exec", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(new(mocks.SqlResult), errors.New("exec error"))
		mockLogger.On("Error", "Failed to update template", "id", uint64(99), "error", mock.Anything).Return()
		err := repo.Save(template)
		assert.Error(t, err)
		mockDB.AssertExpectations(t)
		mockLogger.AssertExpectations(t)
		mapperMock.AssertExpectations(t)
	})
}

func TestTemplateRepository_FindByID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockDB := new(mocks.Database)
		mockLogger := new(mocks.Logger)
		repo := &TemplateRepositoryImpl{db: mockDB, logger: mockLogger, mapper: &mocks.MockTemplateMapper{}}
		id := entities.NewTemplateID(123)
		mockDB.On("Get", mock.AnythingOfType("*models.Template"), mock.Anything, id.Value()).Run(func(args mock.Arguments) {
			template := args.Get(0).(*models.Template)
			template.ID = 123
			template.Name = "Test"
			template.Description = nil
			template.FilePath = "content"
			template.CreatedAt = time.Now()
			template.UpdatedAt = time.Now()
		}).Return(nil)
		expectedTemplate := &entities.Template{}
		mapperMock := repo.mapper.(*mocks.MockTemplateMapper)
		mapperMock.On("ToDomain", mock.AnythingOfType("*models.Template")).Return(expectedTemplate, nil)
		result, err := repo.FindByID(id)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, expectedTemplate, result)
		mockDB.AssertExpectations(t)
		mockLogger.AssertExpectations(t)
		mapperMock.AssertExpectations(t)
	})
	t.Run("not_found", func(t *testing.T) {
		mockDB := new(mocks.Database)
		mockLogger := new(mocks.Logger)
		repo := &TemplateRepositoryImpl{db: mockDB, logger: mockLogger, mapper: &mocks.MockTemplateMapper{}}
		id := entities.NewTemplateID(999)
		mockDB.On("Get", mock.AnythingOfType("*models.Template"), mock.Anything, id.Value()).Return(sql.ErrNoRows)
		result, err := repo.FindByID(id)
		assert.NoError(t, err)
		assert.Nil(t, result)
		mockDB.AssertExpectations(t)
		mockLogger.AssertExpectations(t)
	})
	t.Run("db error", func(t *testing.T) {
		mockDB := new(mocks.Database)
		mockLogger := new(mocks.Logger)
		repo := &TemplateRepositoryImpl{db: mockDB, logger: mockLogger, mapper: &mocks.MockTemplateMapper{}}
		id := entities.NewTemplateID(123)
		dbErr := errors.New("db error")
		mockDB.On("Get", mock.AnythingOfType("*models.Template"), mock.Anything, id.Value()).Return(dbErr)
		mockLogger.On("Error", "Failed to find template by ID", "id", id.Value(), "error", dbErr).Return()
		result, err := repo.FindByID(id)
		assert.Error(t, err)
		assert.Nil(t, result)
		mockDB.AssertExpectations(t)
		mockLogger.AssertExpectations(t)
	})
}

func TestTemplateRepository_FindAll(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockDB := new(mocks.Database)
		repo := &TemplateRepositoryImpl{db: mockDB, logger: new(mocks.Logger), mapper: &mocks.MockTemplateMapper{}}
		modelList := []*models.Template{{Base: models.Base{ID: 1}, Name: "Test", Description: nil, FilePath: "content"}, {Base: models.Base{ID: 2}, Name: "Test2", Description: nil, FilePath: "content2"}}
		mockDB.On("Select", mock.AnythingOfType("*[]*models.Template"), mock.Anything).Run(func(args mock.Arguments) {
			templates := args.Get(0).(*[]*models.Template)
			*templates = modelList
		}).Return(nil)
		mapperMock := repo.mapper.(*mocks.MockTemplateMapper)
		mapperMock.On("ToDomains", modelList).Return([]*entities.Template{{}, {}}, nil)
		result, err := repo.FindAll()
		assert.NoError(t, err)
		assert.Len(t, result, 2)
		mockDB.AssertExpectations(t)
		mapperMock.AssertExpectations(t)
	})
	t.Run("db error", func(t *testing.T) {
		mockDB := new(mocks.Database)
		mockLogger := new(mocks.Logger)
		repo := &TemplateRepositoryImpl{db: mockDB, logger: mockLogger, mapper: &mocks.MockTemplateMapper{}}
		dbErr := errors.New("db error")
		mockDB.On("Select", mock.AnythingOfType("*[]*models.Template"), mock.Anything).Return(dbErr)
		mockLogger.On("Error", "Failed to find all templates", "error", dbErr).Return()
		result, err := repo.FindAll()
		assert.Error(t, err)
		assert.Nil(t, result)
		mockDB.AssertExpectations(t)
		mockLogger.AssertExpectations(t)
	})
}

func TestTemplateRepository_FindEnabledOnly(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockDB := new(mocks.Database)
		repo := &TemplateRepositoryImpl{db: mockDB, logger: new(mocks.Logger), mapper: &mocks.MockTemplateMapper{}}
		modelList := []*models.Template{{Base: models.Base{ID: 1}, Name: "Test", Description: nil, FilePath: "content"}, {Base: models.Base{ID: 2}, Name: "Test2", Description: nil, FilePath: "content2"}}
		mockDB.On("Select", mock.AnythingOfType("*[]*models.Template"), mock.Anything, true).Run(func(args mock.Arguments) {
			templates := args.Get(0).(*[]*models.Template)
			*templates = modelList
		}).Return(nil)
		mapperMock := repo.mapper.(*mocks.MockTemplateMapper)
		mapperMock.On("ToDomains", modelList).Return([]*entities.Template{{}, {}}, nil)
		result, err := repo.FindEnabledOnly()
		assert.NoError(t, err)
		assert.Len(t, result, 2)
		mockDB.AssertExpectations(t)
		mapperMock.AssertExpectations(t)
	})
	t.Run("none found", func(t *testing.T) {
		mockDB := new(mocks.Database)
		mockLogger := new(mocks.Logger)
		repo := &TemplateRepositoryImpl{db: mockDB, logger: mockLogger, mapper: &mocks.MockTemplateMapper{}}
		mockDB.On("Select", mock.AnythingOfType("*[]*models.Template"), mock.Anything, true).Run(func(args mock.Arguments) {
			templates := args.Get(0).(*[]*models.Template)
			*templates = []*models.Template{}
		}).Return(nil)
		mockLogger.On("Info", "No enabled templates found").Return()
		result, err := repo.FindEnabledOnly()
		assert.NoError(t, err)
		assert.Nil(t, result)
		mockDB.AssertExpectations(t)
		mockLogger.AssertExpectations(t)
	})
	t.Run("db error", func(t *testing.T) {
		mockDB := new(mocks.Database)
		mockLogger := new(mocks.Logger)
		repo := &TemplateRepositoryImpl{db: mockDB, logger: mockLogger, mapper: &mocks.MockTemplateMapper{}}
		mockDB.On("Select", mock.AnythingOfType("*[]*models.Template"), mock.Anything, true).Return(errors.New("db error"))
		mockLogger.On("Error", "Failed to find enabled templates", "error", mock.Anything).Return()
		result, err := repo.FindEnabledOnly()
		assert.Error(t, err)
		assert.Nil(t, result)
		mockDB.AssertExpectations(t)
		mockLogger.AssertExpectations(t)
	})
}

func TestTemplateRepository_Delete(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockDB := new(mocks.Database)
		repo := &TemplateRepositoryImpl{db: mockDB, logger: new(mocks.Logger), mapper: &mocks.MockTemplateMapper{}}
		id := entities.NewTemplateID(1)
		mockResult := new(mocks.SqlResult)
		mockDB.On("Exec", mock.Anything, id.Value()).Return(mockResult, nil)
		err := repo.Delete(id)
		assert.NoError(t, err)
		mockDB.AssertExpectations(t)
	})
	t.Run("query error", func(t *testing.T) {
		mockDB := new(mocks.Database)
		mockLogger := new(mocks.Logger)
		repo := &TemplateRepositoryImpl{db: mockDB, logger: mockLogger, mapper: &mocks.MockTemplateMapper{}}
		id := entities.NewTemplateID(2)
		queryErr := errors.New("query error")
		mockResult := new(mocks.SqlResult)
		mockDB.On("Exec", mock.Anything, id.Value()).Return(mockResult, queryErr)
		mockLogger.On("Error", "Failed to delete template", "id", id.Value(), "error", queryErr).Return()
		err := repo.Delete(id)
		assert.Error(t, err)
		mockDB.AssertExpectations(t)
		mockLogger.AssertExpectations(t)
	})
}

func TestTemplateRepository_ExistsByName(t *testing.T) {
	t.Run("exists", func(t *testing.T) {
		mockDB := new(mocks.Database)
		repo := &TemplateRepositoryImpl{db: mockDB, logger: new(mocks.Logger), mapper: &mocks.MockTemplateMapper{}}
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
		repo := &TemplateRepositoryImpl{db: mockDB, logger: new(mocks.Logger), mapper: &mocks.MockTemplateMapper{}}
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
		repo := &TemplateRepositoryImpl{db: mockDB, logger: mockLogger, mapper: &mocks.MockTemplateMapper{}}
		name := "error"
		dbErr := errors.New("db error")
		mockDB.On("Get", mock.AnythingOfType("*int64"), mock.Anything, name).Return(dbErr)
		mockLogger.On("Error", "Failed to check template existence by name", "name", name, "error", dbErr).Return()
		exists, err := repo.ExistsByName(name)
		assert.Error(t, err)
		assert.False(t, exists)
		mockDB.AssertExpectations(t)
		mockLogger.AssertExpectations(t)
	})
}

func TestTemplateRepository_ExistsByFilePath(t *testing.T) {
	t.Run("exists", func(t *testing.T) {
		mockDB := new(mocks.Database)
		repo := &TemplateRepositoryImpl{db: mockDB, logger: new(mocks.Logger), mapper: &mocks.MockTemplateMapper{}}
		filePath := "/exists/path"
		mockDB.On("Get", mock.AnythingOfType("*int64"), mock.Anything, filePath).Run(func(args mock.Arguments) {
			count := args.Get(0).(*int64)
			*count = 1
		}).Return(nil)
		exists, err := repo.ExistsByFilePath(filePath)
		assert.NoError(t, err)
		assert.True(t, exists)
		mockDB.AssertExpectations(t)
	})
	t.Run("not exists", func(t *testing.T) {
		mockDB := new(mocks.Database)
		repo := &TemplateRepositoryImpl{db: mockDB, logger: new(mocks.Logger), mapper: &mocks.MockTemplateMapper{}}
		filePath := "/notexists/path"
		mockDB.On("Get", mock.AnythingOfType("*int64"), mock.Anything, filePath).Run(func(args mock.Arguments) {
			count := args.Get(0).(*int64)
			*count = 0
		}).Return(nil)
		exists, err := repo.ExistsByFilePath(filePath)
		assert.NoError(t, err)
		assert.False(t, exists)
		mockDB.AssertExpectations(t)
	})
	t.Run("db error", func(t *testing.T) {
		mockDB := new(mocks.Database)
		mockLogger := new(mocks.Logger)
		repo := &TemplateRepositoryImpl{db: mockDB, logger: mockLogger, mapper: &mocks.MockTemplateMapper{}}
		filePath := "/error/path"
		dbErr := errors.New("db error")
		mockDB.On("Get", mock.AnythingOfType("*int64"), mock.Anything, filePath).Return(dbErr)
		mockLogger.On("Error", "Failed to check template existence by file path", "file_path", filePath, "error", dbErr).Return()
		exists, err := repo.ExistsByFilePath(filePath)
		assert.Error(t, err)
		assert.False(t, exists)
		mockDB.AssertExpectations(t)
		mockLogger.AssertExpectations(t)
	})
}

func TestTemplateRepository_FindByName(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockDB := new(mocks.Database)
		mockLogger := new(mocks.Logger)
		repo := &TemplateRepositoryImpl{db: mockDB, logger: mockLogger, mapper: &mocks.MockTemplateMapper{}}
		name := "Test"
		mockDB.On("Get", mock.AnythingOfType("*models.Template"), mock.Anything, name).Run(func(args mock.Arguments) {
			template := args.Get(0).(*models.Template)
			template.ID = 1
			template.Name = name
			template.Description = nil
			template.FilePath = "content"
			template.CreatedAt = time.Now()
			template.UpdatedAt = time.Now()
		}).Return(nil)
		expectedTemplate := &entities.Template{}
		mapperMock := repo.mapper.(*mocks.MockTemplateMapper)
		mapperMock.On("ToDomain", mock.AnythingOfType("*models.Template")).Return(expectedTemplate, nil)
		result, err := repo.FindByName(name)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, expectedTemplate, result)
		mockDB.AssertExpectations(t)
		mockLogger.AssertExpectations(t)
		mapperMock.AssertExpectations(t)
	})
	t.Run("not_found", func(t *testing.T) {
		mockDB := new(mocks.Database)
		mockLogger := new(mocks.Logger)
		repo := &TemplateRepositoryImpl{db: mockDB, logger: mockLogger, mapper: &mocks.MockTemplateMapper{}}
		name := "NotFound"
		mockDB.On("Get", mock.AnythingOfType("*models.Template"), mock.Anything, name).Return(sql.ErrNoRows)
		result, err := repo.FindByName(name)
		assert.NoError(t, err)
		assert.Nil(t, result)
		mockDB.AssertExpectations(t)
		mockLogger.AssertExpectations(t)
	})
	t.Run("db error", func(t *testing.T) {
		mockDB := new(mocks.Database)
		mockLogger := new(mocks.Logger)
		repo := &TemplateRepositoryImpl{db: mockDB, logger: mockLogger, mapper: &mocks.MockTemplateMapper{}}
		name := "Error"
		dbErr := errors.New("db error")
		mockDB.On("Get", mock.AnythingOfType("*models.Template"), mock.Anything, name).Return(dbErr)
		mockLogger.On("Error", "Failed to find template by name", "name", name, "error", dbErr).Return()
		result, err := repo.FindByName(name)
		assert.Error(t, err)
		assert.Nil(t, result)
		mockDB.AssertExpectations(t)
		mockLogger.AssertExpectations(t)
	})
}
