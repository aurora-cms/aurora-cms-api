package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/h4rdc0m/aurora-api/domain/common"
	"github.com/h4rdc0m/aurora-api/domain/entities"
	"github.com/h4rdc0m/aurora-api/domain/services"
	"net/http"
)

// KeycloakMiddleware provides middleware for handling Keycloak authentication and authorization in a Gin application.
type KeycloakMiddleware struct {
	logger       common.Logger
	router       common.Router
	tokenService services.TokenService
}

// NewKeycloakMiddleware creates a new instance of KeycloakMiddleware with the provided logger, router, and token service.
func NewKeycloakMiddleware(
	logger common.Logger,
	router common.Router,
	tokenService services.TokenService,
) *KeycloakMiddleware {
	return &KeycloakMiddleware{
		logger:       logger,
		router:       router,
		tokenService: tokenService,
	}
}

// Setup initializes the Keycloak middleware. NOOP for now, as no specific setup is required.
func (k *KeycloakMiddleware) Setup() {}

// AuthRequired is a Gin middleware that checks for a valid Keycloak token in the request header.
func (k *KeycloakMiddleware) AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		token := k.tokenService.ExtractTokenFromHeader(authHeader)
		if token == "" {
			k.logger.Error("Authorization token not found in header")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization token not found"})
			return
		}

		claims, err := k.tokenService.ValidateToken(token)
		if err != nil {
			k.logger.Error("Failed to validate token", "error", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		// Set user context
		c.Set("user_claims", claims)
		c.Set("user_id", claims.Subject)
		c.Set("user_email", claims.Email)
		c.Set("user_username", claims.PreferredUsername)
		c.Set("user_name", claims.Name)
		c.Set("user_roles", claims.RealmAccess.Roles)
		c.Next()
	}
}

// RequireRoles is a Gin middleware that checks if the user has any of the required roles.
func (k *KeycloakMiddleware) RequireRoles(requiredRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, exists := c.Get("user_claims")
		if !exists {
			k.logger.Error("User claims not found in context")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "User claims not found"})
			return
		}

		keycloakClaims := claims.(*entities.KeycloakClaims)
		if k.tokenService.HasAnyRole(keycloakClaims, requiredRoles) {
			c.Next()
			return
		}

		k.logger.Error("User does not have required roles", "required_roles", requiredRoles)
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Forbidden: insufficient roles"})
	}
}
