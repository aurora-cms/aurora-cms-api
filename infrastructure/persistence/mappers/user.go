package mappers

import (
	"github.com/h4rdc0m/aurora-api/domain/entities"
	"github.com/h4rdc0m/aurora-api/domain/value_objects"
	"github.com/h4rdc0m/aurora-api/infrastructure/persistence/models"
)

type UserMapper struct{}

func NewUserMapper() *UserMapper {
	return &UserMapper{}
}

func (m *UserMapper) ToModel(user *entities.User) (*models.User, error) {
	if user == nil {
		return nil, nil
	}

	return &models.User{
		Base: models.Base{
			ID:        user.ID().Value(),
			CreatedAt: user.CreatedAt(),
			UpdatedAt: user.UpdatedAt(),
		},
		KeycloakID: user.KeycloakID().Value(),
		Role:       models.UserRole(user.Role().Value()),
	}, nil
}

func (m *UserMapper) ToDomain(model *models.User) (*entities.User, error) {
	if model == nil {
		return nil, nil
	}

	keycloakID := value_objects.NewKeycloakIDFromUUID(model.KeycloakID)

	role, err := value_objects.NewUserRole(string(model.Role))
	if err != nil {
		return nil, err
	}

	user, err := entities.NewUser(keycloakID, role)
	if err != nil {
		return nil, err
	}

	user.SetID(entities.NewUserID(model.ID))
	user.SetTimestamps(model.CreatedAt, model.UpdatedAt)

	return user, nil
}

func (m *UserMapper) ToModels(users []*entities.User) ([]*models.User, error) {
	if users == nil {
		return nil, nil
	}

	result := make([]*models.User, len(users))
	for i, user := range users {
		model, err := m.ToModel(user)
		if err != nil {
			return nil, err
		}
		result[i] = model
	}

	return result, nil
}

func (m *UserMapper) ToDomains(modelsList []*models.User) ([]*entities.User, error) {
	if modelsList == nil {
		return nil, nil
	}

	result := make([]*entities.User, len(modelsList))
	for i, model := range modelsList {
		domain, err := m.ToDomain(model)
		if err != nil {
			return nil, err
		}
		result[i] = domain
	}

	return result, nil
}
