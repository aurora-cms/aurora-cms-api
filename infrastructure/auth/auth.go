package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/h4rdc0m/aurora-api/domain/common"
	"github.com/h4rdc0m/aurora-api/domain/entities"
	"github.com/h4rdc0m/aurora-api/domain/services"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Service provides methods for authentication and authorization using OIDC and Keycloak.
// It handles login initiation, token exchange, token refresh, and OIDC discovery.
// Dependencies include a logger, configuration, and an HTTP client.
type Service struct {
	logger     common.Logger
	config     Config
	httpClient common.HTTPClient
}

// NewService initializes and returns a new instance of the AuthService.
// It requires a logger, configuration details, and an HTTP client as dependencies.
func NewService(
	logger common.Logger,
	config Config,
	httpClient common.HTTPClient,
) services.AuthService {
	return &Service{
		logger:     logger,
		config:     config,
		httpClient: httpClient,
	}
}

// GetOIDCDiscovery retrieves the OpenID Connect discovery document for the configured Keycloak realm.
func (s Service) GetOIDCDiscovery() (*entities.OIDCDiscovery, error) {
	discoveryURL := fmt.Sprintf("%s/realms/%s/.well-known/openid-configuration", s.config.KeycloakURL, s.config.KeycloakRealm)

	resp, err := s.httpClient.Get(discoveryURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch OIDC discovery document: %w", err)
	}
	defer func(body io.ReadCloser) {
		err := body.Close()
		if err != nil {
			s.logger.Error("Failed to close response body", "error", err)
		}
	}(resp.Body)

	var discovery entities.OIDCDiscovery
	if err := json.NewDecoder(resp.Body).Decode(&discovery); err != nil {
		return nil, fmt.Errorf("failed to decode OIDC discovery document: %w", err)
	}

	return &discovery, nil
}

// InitiateLogin starts the login process by generating state, code verifier, and code challenge for PKCE flow.
// It returns the login URL, state, and the authentication session needed for completing the flow.
// An error is returned if any step in initiating the login process fails.
func (s Service) InitiateLogin(request entities.LoginRequest) (*entities.LoginResponse, *entities.AuthSession, error) {
	redirectURI := request.RedirectURI
	if redirectURI == "" {
		redirectURI = s.config.KeycloakDefaultRedirectURI
	} else {
		redirectURI = fmt.Sprintf("%s/auth/callback", s.config.BaseURL)
	}

	state, err := s.generateRandomString(32)
	if err != nil {
		s.logger.Error("Failed to generate state", "error", err)
		return nil, nil, fmt.Errorf("failed to generate state: %w", err)
	}

	codeVerifier, err := s.generateRandomString(64)
	if err != nil {
		s.logger.Error("Failed to generate code verifier", "error", err)
		return nil, nil, fmt.Errorf("failed to generate code verifier: %w", err)
	}

	codeChallenge := s.generateCodeChallenge(codeVerifier)

	session := &entities.AuthSession{
		State:         state,
		CodeVerifier:  codeVerifier,
		CodeChallenge: codeChallenge,
		RedirectURI:   redirectURI,
		Timestamp:     time.Now().Unix(),
	}

	discovery, err := s.GetOIDCDiscovery()
	if err != nil {
		s.logger.Error("Failed to get OIDC discovery", "error", err)
		return nil, nil, fmt.Errorf("failed to get OIDC discovery: %w", err)
	}

	authURL, err := url.Parse(discovery.AuthorizationEndpoint)
	if err != nil {
		s.logger.Error("Failed to parse authorization endpoint URL", "error", err)
		return nil, nil, fmt.Errorf("failed to parse authorization endpoint URL: %w", err)
	}

	params := url.Values{}
	params.Add("client_id", s.config.KeycloakClientID)
	params.Add("response_type", "code")
	params.Add("scope", "openid profile email roles")
	params.Add("redirect_uri", redirectURI)
	params.Add("state", state)
	params.Add("code_challenge", codeChallenge)
	params.Add("code_challenge_method", "S256")

	authURL.RawQuery = params.Encode()

	response := &entities.LoginResponse{
		AuthURL: authURL.String(),
		State:   state,
	}

	return response, session, nil
}

// ExchangeCodeForTokens exchanges an authorization code for tokens using the OIDC token endpoint. Returns a TokenResponse or an error.
func (s Service) ExchangeCodeForTokens(code string, session *entities.AuthSession) (*entities.TokenResponse, error) {
	discovery, err := s.GetOIDCDiscovery()
	if err != nil {
		s.logger.Error("Failed to get OIDC discovery", "error", err)
		return nil, fmt.Errorf("failed to get OIDC discovery: %w", err)
	}

	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("client_id", s.config.KeycloakClientID)
	data.Set("client_secret", s.config.KeycloakClientSecret)
	data.Set("code", code)
	data.Set("redirect_uri", session.RedirectURI)
	data.Set("code_verifier", session.CodeVerifier)

	req, err := http.NewRequest("POST", discovery.TokenEndpoint, strings.NewReader(data.Encode()))
	if err != nil {
		s.logger.Error("Failed to create token exchange request", "error", err)
		return nil, fmt.Errorf("failed to create token exchange request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		s.logger.Error("Failed to exchange code for tokens", "error", err)
		return nil, fmt.Errorf("failed to exchange code for tokens: %w", err)
	}
	defer func(body io.ReadCloser) {
		err := body.Close()
		if err != nil {
			s.logger.Error("Failed to close response body", "error", err)
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		s.logger.Error("Token exchange failed", "status_code", resp.StatusCode)
		return nil, fmt.Errorf("token exchange failed with status code: %d", resp.StatusCode)
	}

	var tokenResponse entities.TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResponse); err != nil {
		s.logger.Error("Failed to decode token response", "error", err)
		return nil, fmt.Errorf("failed to decode token response: %w", err)
	}

	return &tokenResponse, nil
}

// RefreshToken exchanges a refresh token for a new access token and refresh token by interacting with the Keycloak token endpoint.
func (s Service) RefreshToken(refreshToken string) (*entities.TokenResponse, error) {
	discovery, err := s.GetOIDCDiscovery()
	if err != nil {
		s.logger.Error("Failed to get OIDC discovery", "error", err)
		return nil, fmt.Errorf("failed to get OIDC discovery: %w", err)
	}

	data := url.Values{}
	data.Set("grant_type", "refresh_token")
	data.Set("client_id", s.config.KeycloakClientID)
	data.Set("client_secret", s.config.KeycloakClientSecret)
	data.Set("refresh_token", refreshToken)

	req, err := http.NewRequest("POST", discovery.TokenEndpoint, strings.NewReader(data.Encode()))
	if err != nil {
		s.logger.Error("Failed to create refresh token request", "error", err)
		return nil, fmt.Errorf("failed to create refresh token request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		s.logger.Error("Failed to refresh token", "error", err)
		return nil, fmt.Errorf("failed to refresh token: %w", err)
	}
	defer func(body io.ReadCloser) {
		err := body.Close()
		if err != nil {
			s.logger.Error("Failed to close response body", "error", err)
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		s.logger.Error("Token refresh failed", "status_code", resp.StatusCode)
		return nil, fmt.Errorf("token refresh failed with status code: %d", resp.StatusCode)
	}

	var tokenResponse entities.TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResponse); err != nil {
		s.logger.Error("Failed to decode token response", "error", err)
		return nil, fmt.Errorf("failed to decode token response: %w", err)
	}

	return &tokenResponse, nil
}

// Logout handles user logout by revoking tokens and generating an end-session URL for redirection.
// It uses OIDC discovery endpoints for revocation and end-session operations.
// Returns a LogoutResponse containing a logout URL and a message upon successful completion or an error for failures.
func (s Service) Logout(request entities.LogoutRequest) (*entities.LogoutResponse, error) {
	discovery, err := s.GetOIDCDiscovery()
	if err != nil {
		s.logger.Error("Failed to get OIDC discovery", "error", err)
		return nil, fmt.Errorf("failed to get OIDC discovery: %w", err)
	}

	data := url.Values{}
	data.Set("client_id", s.config.KeycloakClientID)
	data.Set("client_secret", s.config.KeycloakClientSecret)
	data.Set("refresh_token", request.RefreshToken)
	data.Set("id_token_hint", "refresh_token")

	req, err := http.NewRequest("POST", discovery.RevocationEndpoint, strings.NewReader(data.Encode()))
	if err == nil {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		_, err = s.httpClient.Do(req)
		if err != nil {
			s.logger.Warn("Failed to logout", "error", err)
		}
	}

	logoutURL, err := url.Parse(discovery.EndSessionEndpoint)
	if err != nil {
		s.logger.Error("Failed to parse end session endpoint URL", "error", err)
		return nil, fmt.Errorf("failed to parse end session endpoint URL: %w", err)
	}

	params := url.Values{}
	params.Add("redirect_uri", request.RedirectURI)
	logoutURL.RawQuery = params.Encode()

	return &entities.LogoutResponse{
		LogoutURL: logoutURL.String(),
		Message:   "User logged out successfully",
	}, nil
}

// GetAuthConfig retrieves authentication-related configuration values as a map for constructing URLs and settings.
func (s Service) GetAuthConfig() map[string]interface{} {
	return map[string]interface{}{
		"keycloak_url":       s.config.KeycloakURL,
		"keycloak_realm":     s.config.KeycloakRealm,
		"keycloak_client_id": s.config.KeycloakClientID,
		"base_url":           s.config.BaseURL,
		"login_url":          fmt.Sprintf("%s/auth/login", s.config.BaseURL),
		"callback_url":       fmt.Sprintf("%s/auth/callback", s.config.BaseURL),
		"refresh_url":        fmt.Sprintf("%s/auth/refresh", s.config.BaseURL),
		"logout_url":         fmt.Sprintf("%s/auth/logout", s.config.BaseURL),
		"userinfo_url":       fmt.Sprintf("%s/auth/userinfo", s.config.BaseURL),
	}

}

// GetDiscoveryConfig returns a map containing the OpenID Connect discovery configuration for the service.
func (s Service) GetDiscoveryConfig() map[string]interface{} {
	return map[string]interface{}{
		"issuer":                                fmt.Sprintf("%s/realms/%s", s.config.KeycloakURL, s.config.KeycloakRealm),
		"authorization_endpoint":                fmt.Sprintf("%s/auth/login", s.config.BaseURL),
		"token_endpoint":                        fmt.Sprintf("%s/auth/callback", s.config.BaseURL),
		"userinfo_endpoint":                     fmt.Sprintf("%s/auth/userinfo", s.config.BaseURL),
		"end_session_endpoint":                  fmt.Sprintf("%s/auth/logout", s.config.BaseURL),
		"jwks_uri":                              fmt.Sprintf("%s/realms/%s/protocol/openid-connect/certs", s.config.KeycloakURL, s.config.KeycloakRealm),
		"response_types_supported":              []string{"code"},
		"grant_types_supported":                 []string{"authorization_code", "refresh_token"},
		"subject_types_supported":               []string{"public"},
		"id_token_signing_alg_values_supported": []string{"RS256"},
		"scopes_supported":                      []string{"openid", "profile", "email", "roles"},
		"code_challenge_methods_supported":      []string{"S256"},
	}
}

// generateRandomString generates a random URL-safe string of the specified length and returns it or an error if it fails.
func (s Service) generateRandomString(length int) (string, error) {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		s.logger.Error("Failed to generate random bytes", "error", err)
		return "", fmt.Errorf("failed to generate random string: %w", err)
	}

	return base64.RawURLEncoding.EncodeToString(bytes), nil
}

// generateCodeChallenge generates a code challenge string by hashing the given code verifier using SHA-256 and encoding it in Base64.
func (s Service) generateCodeChallenge(codeVerifier string) string {
	hash := sha256.Sum256([]byte(codeVerifier))
	return base64.RawURLEncoding.EncodeToString(hash[:])
}
