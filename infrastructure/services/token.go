package services

import (
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/h4rdc0m/aurora-api/domain/common"
	"github.com/h4rdc0m/aurora-api/domain/entities"
	domainServices "github.com/h4rdc0m/aurora-api/domain/services"
	"io"
	"math/big"
	"net/http"
	"time"
)

// TokenServiceImpl is a concrete implementation of the TokenService interface for handling token-related operations.
// It uses a logger for logging and configuration details for connecting to Keycloak, managing public keys for token validation.
type TokenServiceImpl struct {
	logger    common.Logger
	config    TokenServiceConfig
	publicKey *rsa.PublicKey
}

// NewTokenService creates and returns a new instance of the TokenService implementation with the provided logger and config.
func NewTokenService(
	logger common.Logger,
	config TokenServiceConfig,
) domainServices.TokenService {
	return &TokenServiceImpl{
		logger: logger,
		config: config,
	}
}

// ValidateToken verifies the authenticity and validity of a JWT and extracts claims if successful. Returns error if invalid.
func (t *TokenServiceImpl) ValidateToken(tokenString string) (*entities.KeycloakClaims, error) {
	if err := t.loadPublicKey(); err != nil {
		return nil, err
	}

	token, err := jwt.ParseWithClaims(tokenString, &entities.KeycloakClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return t.publicKey, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(*entities.KeycloakClaims)
	if !ok {
		return nil, fmt.Errorf("invalid claims")
	}

	if time.Now().Unix() > claims.ExpiresAt.Unix() {
		return nil, fmt.Errorf("token expired")
	}

	return claims, nil
}

// ExtractTokenFromHeader extracts the token string from an authorization header with the "Bearer " prefix. Returns an empty string if invalid.
func (t *TokenServiceImpl) ExtractTokenFromHeader(authHeader string) string {
	if authHeader == "" {
		return ""
	}

	const bearerPrefix = "Bearer "
	if len(authHeader) > len(bearerPrefix) && authHeader[:len(bearerPrefix)] == bearerPrefix {
		return authHeader[len(bearerPrefix):]
	}

	return ""
}

// GetUserInfo extracts user information and roles from Keycloak claims as a map of key-value pairs.
func (t *TokenServiceImpl) GetUserInfo(claims *entities.KeycloakClaims) map[string]interface{} {
	return map[string]interface{}{
		"id":             claims.Subject,
		"username":       claims.PreferredUsername,
		"email":          claims.Email,
		"email_verified": claims.EmailVerified,
		"name":           claims.Name,
		"given_name":     claims.GivenName,
		"family_name":    claims.FamilyName,
		"roles":          claims.RealmAccess.Roles,
		"exp":            claims.ExpiresAt.Unix(),
		"iat":            claims.IssuedAt.Unix(),
	}
}

// HasRole determines if a specific role exists within the realm access roles of the provided Keycloak claims.
func (t *TokenServiceImpl) HasRole(claims *entities.KeycloakClaims, requiredRole string) bool {
	for _, role := range claims.RealmAccess.Roles {
		if role == requiredRole {
			return true
		}
	}
	return false
}

// HasAnyRole checks if the claims contain at least one of the specified roles in the requiredRoles slice.
func (t *TokenServiceImpl) HasAnyRole(claims *entities.KeycloakClaims, requiredRoles []string) bool {
	for _, requiredRole := range requiredRoles {
		if t.HasRole(claims, requiredRole) {
			return true
		}
	}
	return false
}

// loadPublicKey retrieves the public key from the Keycloak JWKS endpoint and caches it for token validation.
// Returns an error if the JWKS fetch, decoding, or key conversion fails, or if no keys are found.
func (t *TokenServiceImpl) loadPublicKey() error {
	if t.publicKey != nil {
		return nil
	}

	jwksURL := fmt.Sprintf("%s/realms/%s/protocol/openid-connect/certs", t.config.KeycloakURL, t.config.Realm)

	resp, err := http.Get(jwksURL)
	if err != nil {
		return fmt.Errorf("failed to fetch JWKS: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			t.logger.Warn("Failed to close JWKS response body", "error", err)
		} else {
			t.logger.Debug("Closed JWKS response body successfully")
		}
	}(resp.Body)

	var jwks entities.KeycloakJWKS
	if err := json.NewDecoder(resp.Body).Decode(&jwks); err != nil {
		return fmt.Errorf("failed to decode JWKS: %w", err)
	}

	if len(jwks.Keys) == 0 {
		return fmt.Errorf("no keys found in JWKS")
	}

	key := jwks.Keys[0]
	publicKey, err := t.jwkToRSAPublicKey(key)
	if err != nil {
		return fmt.Errorf("failed to convert JWK to RSA public key: %w", err)
	}

	t.publicKey = publicKey
	return nil
}

// jwkToRSAPublicKey converts a KeycloakJWK to an *rsa.PublicKey by decoding the modulus and exponent from base64.
func (t *TokenServiceImpl) jwkToRSAPublicKey(jwk entities.KeycloakJWK) (*rsa.PublicKey, error) {
	nBytes, err := base64.RawURLEncoding.DecodeString(jwk.N)
	if err != nil {
		return nil, err
	}

	eBytes, err := base64.RawURLEncoding.DecodeString(jwk.E)
	if err != nil {
		return nil, err
	}

	n := big.NewInt(0).SetBytes(nBytes)
	e := 0
	for _, b := range eBytes {
		e = e*256 + int(b)
	}

	return &rsa.PublicKey{N: n, E: e}, nil
}
