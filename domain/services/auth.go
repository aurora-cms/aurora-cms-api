package services

import "github.com/h4rdc0m/aurora-api/domain/entities"

type AuthService interface {
	// GetOIDCDiscovery retrieves the OpenID Connect discovery document.
	// Returns an OIDCDiscovery object or an error if retrieval fails.
	GetOIDCDiscovery() (*entities.OIDCDiscovery, error)

	// InitiateLogin starts the login process for a user.
	// Takes an LoginRequest and returns an LoginResponse, an AuthSession, or an error.
	InitiateLogin(request entities.LoginRequest) (*entities.LoginResponse, *entities.AuthSession, error)

	// ExchangeCodeForTokens exchanges an authorization code for access and refresh tokens.
	// Takes the authorization code and the current AuthSession, returns a TokenResponse or an error.
	ExchangeCodeForTokens(code string, session *entities.AuthSession) (*entities.TokenResponse, error)

	// RefreshToken refreshes the access token using a refresh token.
	// Takes a refresh token and returns a new TokenResponse or an error.
	RefreshToken(refreshToken string) (*entities.TokenResponse, error)

	// Logout terminates the user's session.
	// Takes a LogoutRequest and returns a LogoutResponse or an error.
	Logout(request entities.LogoutRequest) (*entities.LogoutResponse, error)

	// GetAuthConfig returns the authentication configuration as a map.
	GetAuthConfig() map[string]interface{}

	// GetDiscoveryConfig returns the OIDC discovery configuration as a map.
	GetDiscoveryConfig() map[string]interface{}
}
