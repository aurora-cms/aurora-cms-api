package services

import (
	"encoding/json"
	"github.com/h4rdc0m/aurora-api/domain/entities"
	"github.com/h4rdc0m/aurora-api/tests/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

func TestSessionServiceImpl_StoreSession(t *testing.T) {
	logger := &mocks.Logger{}

	// Set up mock expectations to match actual calls: message, key, value
	logger.On("Info", "Stored session", "state", mock.AnythingOfType("string")).Return()

	service := NewSessionService(logger).(*SessionServiceImpl)

	tests := []struct {
		name    string
		session *entities.AuthSession
	}{
		{name: "valid session", session: &entities.AuthSession{State: "state1", Timestamp: time.Now().Unix()}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.StoreSession(tt.session)
			assert.NoError(t, err)
			assert.Contains(t, service.sessions, tt.session.State)
		})
	}

	// Verify all expectations were met
	logger.AssertExpectations(t)
}

func TestSessionServiceImpl_GetSession(t *testing.T) {
	logger := &mocks.Logger{}

	// Set up mock expectations for all possible log calls
	logger.On("Info", "Stored session", "state", mock.AnythingOfType("string")).Return()
	logger.On("Info", "Session not found", "state", mock.AnythingOfType("string")).Return()
	logger.On("Info", "Session expired", "state", mock.AnythingOfType("string")).Return()

	service := NewSessionService(logger).(*SessionServiceImpl)

	expiredSession := &entities.AuthSession{State: "expired", Timestamp: time.Now().Unix() - 301}
	validSession := &entities.AuthSession{State: "valid", Timestamp: time.Now().Unix()}

	_ = service.StoreSession(expiredSession)
	_ = service.StoreSession(validSession)

	tests := []struct {
		name      string
		state     string
		shouldErr bool
		errText   string
	}{
		{name: "valid session", state: "valid", shouldErr: false},
		{name: "expired session", state: "expired", shouldErr: true, errText: "session expired for state: expired"},
		{name: "nonexistent session", state: "missing", shouldErr: true, errText: "session not found for state: missing"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			session, err := service.GetSession(tt.state)
			if tt.shouldErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errText)
				assert.Nil(t, session)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, session)
				assert.Equal(t, tt.state, session.State)
			}
		})
	}
}

func TestSessionServiceImpl_DeleteSession(t *testing.T) {
	logger := &mocks.Logger{}

	// Set up mock expectations
	logger.On("Info", "Stored session", "state", mock.AnythingOfType("string")).Return()
	logger.On("Info", "Deleted session", "state", mock.AnythingOfType("string")).Return()
	logger.On("Warn", "Session not found for deletion", "state", mock.AnythingOfType("string")).Return()

	service := NewSessionService(logger).(*SessionServiceImpl)

	_ = service.StoreSession(&entities.AuthSession{State: "deletable"})

	tests := []struct {
		name      string
		state     string
		shouldErr bool
		errText   string
	}{
		{name: "delete existing session", state: "deletable", shouldErr: false},
		{name: "delete nonexisting session", state: "missing", shouldErr: true, errText: "session not found for state: missing"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.DeleteSession(tt.state)
			if tt.shouldErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errText)
			} else {
				assert.NoError(t, err)
				assert.NotContains(t, service.sessions, tt.state)
			}
		})
	}
}

func TestSessionServiceImpl_ValidateAndConsumeSession(t *testing.T) {
	logger := &mocks.Logger{}

	// Set up mock expectations
	logger.On("Info", "Stored session", "state", mock.AnythingOfType("string")).Return()
	logger.On("Info", "Deleted session", "state", mock.AnythingOfType("string")).Return()
	logger.On("Info", "Session not found", "state", mock.AnythingOfType("string")).Return()

	service := NewSessionService(logger).(*SessionServiceImpl)

	_ = service.StoreSession(&entities.AuthSession{State: "valid", Timestamp: time.Now().Unix()})

	tests := []struct {
		name      string
		state     string
		shouldErr bool
		errText   string
	}{
		{name: "validate and delete valid session", state: "valid", shouldErr: false},
		{name: "validate nonexisting session", state: "missing", shouldErr: true, errText: "session not found for state: missing"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			session, err := service.ValidateAndConsumeSession(tt.state)
			if tt.shouldErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errText)
				assert.Nil(t, session)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, session)
				assert.Equal(t, tt.state, session.State)
				assert.NotContains(t, service.sessions, tt.state)
			}
		})
	}
}

func TestSessionServiceImpl_SerializeSession(t *testing.T) {
	logger := &mocks.Logger{}

	// Set up mock expectation for potential error logging
	logger.On("Error", "Failed to serialize session", "error", mock.Anything).Maybe()

	service := NewSessionService(logger).(*SessionServiceImpl)

	session := &entities.AuthSession{State: "serializable"}

	data, err := service.SerializeSession(session)
	assert.NoError(t, err)

	var deserialized entities.AuthSession
	err = json.Unmarshal([]byte(data), &deserialized)
	assert.NoError(t, err)
	assert.Equal(t, session.State, deserialized.State)
}

func TestSessionServiceImpl_DeserializeSession(t *testing.T) {
	logger := &mocks.Logger{}

	// Set up mock expectation for error logging
	logger.On("Error", "Failed to deserialize session", "error", mock.Anything).Maybe()

	service := NewSessionService(logger).(*SessionServiceImpl)

	session := &entities.AuthSession{State: "deserializable"}
	data, _ := json.Marshal(session)

	tests := []struct {
		name      string
		input     string
		shouldErr bool
		errText   string
	}{
		{name: "valid data", input: string(data), shouldErr: false},
		{name: "invalid data", input: "{invalid json}", shouldErr: true, errText: "invalid character"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := service.DeserializeSession(tt.input)
			if tt.shouldErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errText)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, session.State, result.State)
			}
		})
	}
}
