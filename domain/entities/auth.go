package entities

// OIDCDiscovery represents the OpenID Connect discovery document structure.
type OIDCDiscovery struct {
	AuthorizationEndpoint string `json:"authorization_endpoint"`
	TokenEndpoint         string `json:"token_endpoint"`
	UserInfoEndpoint      string `json:"userinfo_endpoint"`
	EndSessionEndpoint    string `json:"end_session_endpoint"`
	RevocationEndpoint    string `json:"revocation_endpoint"`
	JwksURI               string `json:"jwks_uri"`
	Issuer                string `json:"issuer"`
}

// TokenResponse represents the structure of the token response from the Keycloak server.
type TokenResponse struct {
	AccessToken      string `json:"access_token"`
	RefreshToken     string `json:"refresh_token"`
	ExpiresIn        int    `json:"expires_in"`
	RefreshExpiresIn int    `json:"refresh_expires_in"`
	TokenType        string `json:"token_type"`
	IDToken          string `json:"id_token,omitempty"`
	Scope            string `json:"scope"`
}

// LoginRequest represents the request structure for initiating a login.
type LoginRequest struct {
	RedirectURI string `json:"redirect_uri"`
}

// LoginResponse represents the response structure for a login request.
type LoginResponse struct {
	AuthURL string `json:"auth_url"`
	State   string `json:"state"`
}

// AuthSession represents the session data for an authentication request.
type AuthSession struct {
	State         string `json:"state"`
	CodeVerifier  string `json:"code_verifier"`
	CodeChallenge string `json:"code_challenge"`
	RedirectURI   string `json:"redirect_uri"`
	Timestamp     int64  `json:"timestamp"`
}

// CallbackRequest represents the request structure for handling the callback from the authentication provider.
type CallbackRequest struct {
	Code  string `json:"code"`
	State string `json:"state"`
}

// RefreshTokenRequest represents the request structure for refreshing an access token.
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

// LogoutRequest represents the request structure for logging out a user.
type LogoutRequest struct {
	RefreshToken string `json:"refresh_token"`
	RedirectURI  string `json:"redirect_uri"`
}

// LogoutResponse represents the response structure for a logout request.
type LogoutResponse struct {
	LogoutURL string `json:"logout_url"`
	Message   string `json:"message"`
}
