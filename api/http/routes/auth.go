package routes

import (
	"github.com/h4rdc0m/aurora-api/api/http/controllers"
	"github.com/h4rdc0m/aurora-api/api/http/middlewares"
	"github.com/h4rdc0m/aurora-api/domain/common"
)

type AuthRoutes struct {
	logger     common.Logger
	handler    common.Router
	controller *controllers.AuthController
	middleware *middlewares.KeycloakMiddleware
}

func NewAuthRoutes(
	logger common.Logger,
	handler common.Router,
	controller *controllers.AuthController,
	middleware *middlewares.KeycloakMiddleware,
) *AuthRoutes {
	return &AuthRoutes{
		logger:     logger,
		handler:    handler,
		controller: controller,
		middleware: middleware,
	}
}

func (r *AuthRoutes) Setup() {
	r.logger.Info("Setting up auth routes")

	auth := r.handler.Group("/auth")
	{
		auth.GET("/login", r.controller.InitiateLogin)
		auth.GET("/callback", r.controller.HandleCallback)
		auth.POST("/refresh", r.controller.RefreshToken)
		auth.GET("/config", r.controller.GetAuthConfig)

		// Protected routes
		auth.POST("/logout", r.middleware.AuthRequired(), r.controller.Logout)
		auth.GET("/userinfo", r.middleware.AuthRequired(), r.controller.GetUserInfo)

		// OIDC Discovery endpoint
		auth.GET("/.well-known/openid-configuration", r.controller.GetDiscoveryConfig)

		// Health check for auth service
		auth.GET("/health", r.controller.HealthCheck)
	}
}
