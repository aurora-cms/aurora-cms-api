package services

import "github.com/h4rdc0m/aurora-api/domain/entities"

type TokenService interface {
	ValidateToken(tokenString string) (*entities.KeycloakClaims, error)
	ExtractTokenFromHeader(authHeader string) string
	GetUserInfo(claims *entities.KeycloakClaims) map[string]interface{}
	HasRole(claims *entities.KeycloakClaims, role string) bool
	HasAnyRole(claims *entities.KeycloakClaims, roles []string) bool
}
