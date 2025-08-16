package routes

import "go.uber.org/fx"

// Module provides the routes module for the application.
var Module = fx.Options(
	fx.Provide(NewHealthRoutes),
	fx.Provide(NewAuthRoutes),
	fx.Provide(NewRoutes),
)

// Routes is a collection of route definitions for the application.
type Routes []Route

// Route defines the interface for setting up routes in the application.
type Route interface {
	Setup()
}

// NewRoutes creates a new instance of Routes with the provided health routes.
func NewRoutes(
	healthRoutes *HealthRoutes,
	authRoutes *AuthRoutes,
) Routes {
	return Routes{
		healthRoutes,
		authRoutes,
	}
}

// Setup initializes all routes in the Routes collection.
func (r Routes) Setup() {
	for _, route := range r {
		route.Setup()
	}
}
