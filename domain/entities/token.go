package entities

import "github.com/golang-jwt/jwt/v5"

// KeycloakJWKS represents the JSON Web Key Set (JWKS) structure returned by Keycloak.
type KeycloakJWKS struct {
	Keys []KeycloakJWK `json:"keys"`
}

// KeycloakJWK represents a single JSON Web Key (JWK) used for verifying JWTs issued by Keycloak.
type KeycloakJWK struct {
	Kid string `json:"kid"`
	Kty string `json:"kty"`
	Alg string `json:"alg"`
	Use string `json:"use"`
	N   string `json:"n"`
	E   string `json:"e"`
}

// KeycloakClaims represents the claims structure of a JWT issued by Keycloak.
type KeycloakClaims struct {
	jwt.RegisteredClaims
	RealmAccess struct {
		Roles []string `json:"roles"`
	} `json:"realm_access"`
	ResourceAccess map[string]struct {
		Roles []string `json:"roles"`
	} `json:"resource_access"`
	PreferredUsername string `json:"preferred_username"`
	Email             string `json:"email"`
	EmailVerified     bool   `json:"email_verified"`
	Name              string `json:"name"`
	GivenName         string `json:"given_name"`
	FamilyName        string `json:"family_name"`
}
