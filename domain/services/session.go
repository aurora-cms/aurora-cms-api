package services

import "github.com/h4rdc0m/aurora-api/domain/entities"

// SessionService defines the interface for managing user authentication sessions.
type SessionService interface {
	// StoreSession saves the provided authentication session.
	StoreSession(session *entities.AuthSession) error

	// GetSession retrieves an authentication session by its state.
	GetSession(state string) (*entities.AuthSession, error)

	// DeleteSession removes an authentication session by its state.
	DeleteSession(state string) error

	// ValidateAndConsumeSession validates and consumes a session by its state, returning the session if valid.
	ValidateAndConsumeSession(state string) (*entities.AuthSession, error)

	// SerializeSession serializes the authentication session into a string.
	SerializeSession(session *entities.AuthSession) (string, error)

	// DeserializeSession deserializes the string data into an authentication session.
	DeserializeSession(data string) (*entities.AuthSession, error)
}
