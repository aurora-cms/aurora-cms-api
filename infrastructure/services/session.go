package services

import (
	"encoding/json"
	"fmt"
	"github.com/h4rdc0m/aurora-api/domain/common"
	"github.com/h4rdc0m/aurora-api/domain/entities"
	domainServices "github.com/h4rdc0m/aurora-api/domain/services"
	"sync"
	"time"
)

// SessionServiceImpl is an implementation of SessionService for managing user authentication sessions.
// It uses an in-memory map to store sessions, protected by a read-write mutex for concurrent access.
// A logger is used to log events like creation, deletion, and expiration of sessions.
type SessionServiceImpl struct {
	logger   common.Logger
	sessions map[string]*entities.AuthSession
	mutex    sync.RWMutex
}

// NewSessionService initializes and returns a new instance of SessionService to manage user authentication sessions.
func NewSessionService(logger common.Logger) domainServices.SessionService {
	service := &SessionServiceImpl{
		logger:   logger,
		sessions: make(map[string]*entities.AuthSession),
	}

	// Start a cleanup routine to remove expired sessions
	go service.cleanupExpiredSessions()

	return service
}

// StoreSession stores an authentication session in memory using the session's state as the key. It is thread-safe.
func (s *SessionServiceImpl) StoreSession(session *entities.AuthSession) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Store the session using its state as the key
	s.sessions[session.State] = session
	s.logger.Info("Stored session", "state", session.State)

	return nil
}

// GetSession retrieves an AuthSession by its state. Returns an error if the session doesn't exist or has expired.
func (s *SessionServiceImpl) GetSession(state string) (*entities.AuthSession, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	session, exists := s.sessions[state]
	if !exists {
		s.logger.Info("Session not found", "state", state)
		return nil, fmt.Errorf("session not found for state: %s", state)
	}

	if time.Now().Unix()-session.Timestamp > 300 { // 5-minute expiration time
		s.logger.Info("Session expired", "state", state)
		return nil, fmt.Errorf("session expired for state: %s", state)
	}

	return session, nil
}

// DeleteSession removes a session associated with the given state from the session storage. Returns an error if not found.
func (s *SessionServiceImpl) DeleteSession(state string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if _, exists := s.sessions[state]; !exists {
		s.logger.Warn("Session not found for deletion", "state", state)
		return fmt.Errorf("session not found for state: %s", state)
	}

	delete(s.sessions, state)
	s.logger.Info("Deleted session", "state", state)
	return nil
}

// ValidateAndConsumeSession validates the session by its state and removes it from storage after successful validation.
func (s *SessionServiceImpl) ValidateAndConsumeSession(state string) (*entities.AuthSession, error) {
	session, err := s.GetSession(state)
	if err != nil {
		return nil, err
	}

	// Remove the session after validation
	if err := s.DeleteSession(state); err != nil {
		s.logger.Error("Failed to delete session after validation", "state", state, "error", err)
		return nil, err
	}
	return session, nil
}

// SerializeSession converts an AuthSession into its JSON string representation. Returns the serialized string or an error.
func (s *SessionServiceImpl) SerializeSession(session *entities.AuthSession) (string, error) {
	data, err := json.Marshal(session)
	if err != nil {
		s.logger.Error("Failed to serialize session", "error", err)
		return "", err
	}
	return string(data), nil
}

// DeserializeSession deserializes a JSON string into an AuthSession object and logs an error if deserialization fails.
func (s *SessionServiceImpl) DeserializeSession(data string) (*entities.AuthSession, error) {
	var session entities.AuthSession
	err := json.Unmarshal([]byte(data), &session)
	if err != nil {
		s.logger.Error("Failed to deserialize session", "error", err)
		return nil, err
	}
	return &session, nil
}

// cleanupExpiredSessions periodically removes expired authentication sessions from the session map.
// It runs in an infinite loop with a 1-minute interval, deleting sessions older than 5 minutes.
// Thread safety is maintained using a mutex to prevent race conditions during session removal.
func (s *SessionServiceImpl) cleanupExpiredSessions() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		s.mutex.Lock()
		currentTime := time.Now().Unix()
		for state, session := range s.sessions {
			if currentTime-session.Timestamp > 300 { // 5-minute expiration time
				delete(s.sessions, state)
				s.logger.Info("Removed expired session", "state", state)
			}
		}
		s.mutex.Unlock()
	}
}
