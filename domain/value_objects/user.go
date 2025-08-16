package value_objects

import (
	"github.com/google/uuid"
	"github.com/h4rdc0m/aurora-api/domain/errors"
)

type KeycloakID struct {
	value uuid.UUID
}

func NewKeycloakID(id string) (*KeycloakID, error) {
	if id == "" {
		return nil, errors.ErrKeycloakIDEmpty
	}

	parsedUUID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.ErrKeycloakIDInvalid
	}

	return &KeycloakID{value: parsedUUID}, nil
}

func NewKeycloakIDFromUUID(id uuid.UUID) *KeycloakID {
	return &KeycloakID{value: id}
}

func (k KeycloakID) Value() uuid.UUID {
	return k.value
}

func (k KeycloakID) String() string {
	return k.value.String()
}

func (k KeycloakID) Equals(other KeycloakID) bool {
	return k.value == other.value
}
