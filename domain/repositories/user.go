package repositories

import (
	"github.com/h4rdc0m/aurora-api/domain/entities"
	"github.com/h4rdc0m/aurora-api/domain/value_objects"
)

type UserRepository interface {
	Save(user *entities.User) error
	FindByID(id entities.UserID) (*entities.User, error)
	FindByKeycloakID(keycloakID value_objects.KeycloakID) (*entities.User, error)
	FindAll() ([]*entities.User, error)
	FindAllByTenantID(tenantID entities.TenantID) ([]*entities.User, error)
	Delete(id entities.UserID) error
	ExistsByKeycloakID(keycloakID value_objects.KeycloakID) (bool, error)
}
