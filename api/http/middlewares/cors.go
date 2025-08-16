package middlewares

import (
	"github.com/h4rdc0m/aurora-api/domain/common"
	"github.com/h4rdc0m/aurora-api/infrastructure/config"
	cors "github.com/rs/cors/wrapper/gin"
)

// CorsMiddleware is a middleware that handles CORS requests.
type CorsMiddleware struct {
	handler common.Router
	logger  common.Logger
	env     *config.Env
}

// NewCorsMiddleware creates a new instance of CorsMiddleware.
func NewCorsMiddleware(handler common.Router, logger common.Logger, env *config.Env) *CorsMiddleware {
	return &CorsMiddleware{
		handler: handler,
		logger:  logger,
		env:     env,
	}
}

// Setup initializes the CORS middleware with the appropriate settings.
func (m *CorsMiddleware) Setup() {
	m.logger.Info("Setting up CORS middleware")
	debug := m.env.Environment == "development"
	m.handler.Use(cors.New(cors.Options{
		Logger:           m.logger,
		AllowCredentials: true,
		AllowOriginFunc:  func(origin string) bool { return true },
		AllowedHeaders:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		Debug:            debug,
	}))
}
