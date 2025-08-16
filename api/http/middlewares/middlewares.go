package middlewares

import "go.uber.org/fx"

// Module provides a collection of middleware interfaces and their setup functionality.
var Module = fx.Options(
	fx.Provide(NewCorsMiddleware),
	fx.Provide(NewDatabaseTrx),
	fx.Provide(NewKeycloakMiddleware),
	fx.Provide(NewMiddlewares),
)

// IMiddleware defines the interface for middleware setup.
type IMiddleware interface {
	Setup()
}

// Middlewares is a collection of middleware interfaces that can be used in the application.
type Middlewares []IMiddleware

// NewMiddlewares creates a new instance of Middlewares.
func NewMiddlewares(
	corsMiddleware *CorsMiddleware,
	dbTrxMiddleware *DatabaseTrx,
	keycloakMiddleware *KeycloakMiddleware,
) Middlewares {
	return Middlewares{
		corsMiddleware,
		dbTrxMiddleware,
		keycloakMiddleware,
	}
}

// Setup initializes all middlewares in the Middlewares collection.
func (m Middlewares) Setup() {
	for _, middleware := range m {
		middleware.Setup()
	}
}
