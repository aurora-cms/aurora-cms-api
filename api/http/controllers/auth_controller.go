package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/h4rdc0m/aurora-api/domain/common"
	"github.com/h4rdc0m/aurora-api/domain/entities"
	"github.com/h4rdc0m/aurora-api/domain/services"
	"net/http"
)

// AuthController handles authentication-related requests.
type AuthController struct {
	BaseController
	logger         common.Logger
	authService    services.AuthService
	sessionService services.SessionService
	tokenService   services.TokenService
}

// NewAuthController creates a new instance of AuthController with the provided services and environment.
func NewAuthController(
	logger common.Logger,
	authService services.AuthService,
	sessionService services.SessionService,
	tokenService services.TokenService,
) *AuthController {
	return &AuthController{
		BaseController: BaseController{},
		logger:         logger,
		authService:    authService,
		sessionService: sessionService,
		tokenService:   tokenService,
	}
}

func (a *AuthController) InitiateLogin(c *gin.Context) {
	var request entities.LoginRequest
	request.RedirectURI = c.Query("redirect_uri")

	loginResponse, session, err := a.authService.InitiateLogin(request)
	if err != nil {
		a.logger.Error("Failed to initiate login", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	// Store the session for later use
	if err := a.sessionService.StoreSession(session); err != nil {
		a.logger.Error("Failed to store session", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, loginResponse)
}

func (a *AuthController) HandleCallback(c *gin.Context) {
	code := c.Query("code")
	state := c.Query("state")
	errorParam := c.Query("error")

	if errorParam != "" {
		errorDescription := c.Query("error_description")
		a.logger.Info("Authentication error", "error", errorParam, "description", errorDescription)
		c.JSON(http.StatusBadRequest, gin.H{"error": errorParam, "description": errorDescription})
		return
	}

	if code == "" || state == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing code or state"})
		return
	}

	// Validate and consume the session
	session, err := a.sessionService.ValidateAndConsumeSession(state)
	if err != nil {
		a.logger.Info("Failed to validate session", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session"})
		return
	}

	// Exchange the code for tokens
	tokenResponse, err := a.authService.ExchangeCodeForTokens(code, session)
	if err != nil {
		a.logger.Info("Failed to exchange code for tokens", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	// Validate the access token
	claims, err := a.tokenService.ValidateToken(tokenResponse.AccessToken)
	if err != nil {
		a.logger.Error("Failed to validate access token", "error", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid access token"})
		return
	}

	response := gin.H{
		"access_token":       tokenResponse.AccessToken,
		"refresh_token":      tokenResponse.RefreshToken,
		"expires_in":         tokenResponse.ExpiresIn,
		"refresh_expires_in": tokenResponse.RefreshExpiresIn,
		"token_type":         tokenResponse.TokenType,
		"id_token":           tokenResponse.IDToken,
		"scope":              tokenResponse.Scope,
		"user":               a.tokenService.GetUserInfo(claims),
	}
	c.JSON(http.StatusOK, response)
}

func (a *AuthController) RefreshToken(c *gin.Context) {
	var request entities.RefreshTokenRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		a.logger.Error("Invalid request body", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	tokenResponse, err := a.authService.RefreshToken(request.RefreshToken)
	if err != nil {
		a.logger.Error("Failed to refresh token", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, tokenResponse)
}

func (a *AuthController) Logout(c *gin.Context) {
	var request entities.LogoutRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		a.logger.Error("Invalid request body", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	logoutResponse, err := a.authService.Logout(request)
	if err != nil {
		a.logger.Error("Failed to logout", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, logoutResponse)
}

func (a *AuthController) GetUserInfo(c *gin.Context) {
	claims, exists := c.Get("user_claims")
	if !exists {
		a.logger.Error("User claims not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	keycloakClaims := claims.(*entities.KeycloakClaims)
	c.JSON(http.StatusOK, a.tokenService.GetUserInfo(keycloakClaims))
}

func (a *AuthController) GetAuthConfig(c *gin.Context) {
	c.JSON(http.StatusOK, a.authService.GetAuthConfig())
}

func (a *AuthController) GetDiscoveryConfig(c *gin.Context) {
	c.JSON(http.StatusOK, a.authService.GetDiscoveryConfig())
}

func (a *AuthController) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "ok",
		"timestamp": gin.H{},
	})
}
