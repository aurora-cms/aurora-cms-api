package services

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"github.com/golang-jwt/jwt/v5"
	"github.com/h4rdc0m/aurora-api/domain/entities"
	"github.com/h4rdc0m/aurora-api/tests/mocks"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestValidateToken(t *testing.T) {
	privateKey, _ := rsa.GenerateKey(rand.Reader, 2048)
	publicKey := &privateKey.PublicKey
	service := &TokenServiceImpl{
		logger:    &mocks.Logger{},
		config:    TokenServiceConfig{},
		publicKey: publicKey,
	}

	tests := []struct {
		name      string
		setup     func() (string, error)
		expectErr bool
		expectMsg string
	}{
		{
			"valid token",
			func() (string, error) {
				claims := &entities.KeycloakClaims{
					RegisteredClaims: jwt.RegisteredClaims{
						ExpiresAt: jwt.NewNumericDate(time.Now().Add(10 * time.Minute)),
					},
				}
				token := jwt.NewWithClaims(jwt.SigningMethodRS512, claims)
				return token.SignedString(privateKey)
			},
			false,
			"",
		},
		{
			"expired token",
			func() (string, error) {
				claims := &entities.KeycloakClaims{
					RegisteredClaims: jwt.RegisteredClaims{
						ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Minute)),
					},
				}
				token := jwt.NewWithClaims(jwt.SigningMethodRS512, claims)
				return token.SignedString(privateKey)
			},
			true,
			"token is expired",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			token, _ := tc.setup()
			claims, err := service.ValidateToken(token)
			if tc.expectErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectMsg)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, claims)
			}
		})
	}
}

func TestExtractTokenFromHeader(t *testing.T) {
	service := &TokenServiceImpl{}

	tests := []struct {
		name     string
		header   string
		expected string
	}{
		{"valid header", "Bearer valid.token.string", "valid.token.string"},
		{"empty header", "", ""},
		{"missing bearer prefix", "valid.token.string", ""},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			token := service.ExtractTokenFromHeader(tc.header)
			assert.Equal(t, tc.expected, token)
		})
	}
}

func TestGetUserInfo(t *testing.T) {
	service := &TokenServiceImpl{}

	// Create proper time values for the claims
	now := time.Now()
	expiresAt := jwt.NewNumericDate(now.Add(1 * time.Hour))
	issuedAt := jwt.NewNumericDate(now)

	claims := &entities.KeycloakClaims{
		PreferredUsername: "john.doe",
		Email:             "john.doe@example.com",
		EmailVerified:     true,
		Name:              "John Doe",
		GivenName:         "John",
		FamilyName:        "Doe",
		RealmAccess: struct {
			Roles []string `json:"roles"`
		}{
			Roles: []string{"admin", "user"},
		},
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   "1",
			ExpiresAt: expiresAt,
			IssuedAt:  issuedAt,
		},
	}

	expected := map[string]interface{}{
		"id":             "1",
		"username":       "john.doe",
		"email":          "john.doe@example.com",
		"email_verified": true,
		"name":           "John Doe",
		"given_name":     "John",
		"family_name":    "Doe",
		"roles":          []string{"admin", "user"},
		"exp":            expiresAt.Unix(),
		"iat":            issuedAt.Unix(),
	}

	userInfo := service.GetUserInfo(claims)
	assert.Equal(t, expected, userInfo)
}

func TestHasRole(t *testing.T) {
	service := &TokenServiceImpl{}
	claims := &entities.KeycloakClaims{
		RealmAccess: struct {
			Roles []string `json:"roles"`
		}{
			Roles: []string{"admin", "user"},
		},
	}

	tests := []struct {
		name        string
		role        string
		expectation bool
	}{
		{"role present", "admin", true},
		{"role absent", "viewer", false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			hasRole := service.HasRole(claims, tc.role)
			assert.Equal(t, tc.expectation, hasRole)
		})
	}
}

func TestHasAnyRole(t *testing.T) {
	service := &TokenServiceImpl{}
	claims := &entities.KeycloakClaims{
		RealmAccess: struct {
			Roles []string `json:"roles"`
		}{
			Roles: []string{"admin", "user"},
		},
	}

	tests := []struct {
		name        string
		roles       []string
		expectation bool
	}{
		{"at least one role matches", []string{"admin", "viewer"}, true},
		{"no roles match", []string{"viewer", "editor"}, false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			hasAnyRole := service.HasAnyRole(claims, tc.roles)
			assert.Equal(t, tc.expectation, hasAnyRole)
		})
	}
}

func TestLoadPublicKey(t *testing.T) {
	// Generate a real RSA key pair for testing
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	assert.NoError(t, err)

	publicKey := &privateKey.PublicKey

	// Convert the real public key to JWK format
	nBytes := publicKey.N.Bytes()
	eBytes := make([]byte, 4)
	eBytes[0] = byte(publicKey.E >> 24)
	eBytes[1] = byte(publicKey.E >> 16)
	eBytes[2] = byte(publicKey.E >> 8)
	eBytes[3] = byte(publicKey.E)

	// Remove leading zeros from exponent
	for len(eBytes) > 1 && eBytes[0] == 0 {
		eBytes = eBytes[1:]
	}

	validKey := entities.KeycloakJWK{
		N:   base64.RawURLEncoding.EncodeToString(nBytes),
		E:   base64.RawURLEncoding.EncodeToString(eBytes),
		Alg: "RS256",
	}

	jwks := entities.KeycloakJWKS{Keys: []entities.KeycloakJWK{validKey}}

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(jwks)
	}))
	defer mockServer.Close()

	logger := &mocks.Logger{}

	// Set up mock expectation for the Debug call that happens when the response body is closed successfully
	logger.On("Debug", "Closed JWKS response body successfully").Return()

	service := &TokenServiceImpl{
		logger: logger,
		config: TokenServiceConfig{KeycloakURL: mockServer.URL, Realm: "realm"},
	}

	err = service.loadPublicKey()
	assert.NoError(t, err)
	assert.NotNil(t, service.publicKey)

	// Verify the loaded key matches our original key
	assert.Equal(t, publicKey.N, service.publicKey.N)
	assert.Equal(t, publicKey.E, service.publicKey.E)

	// Verify all expectations were met
	logger.AssertExpectations(t)
}
