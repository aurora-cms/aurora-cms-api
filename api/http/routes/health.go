package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/h4rdc0m/aurora-api/api/http/controllers"
	"github.com/h4rdc0m/aurora-api/domain/common"
)

// HealthRoutes configures the health check routes
type HealthRoutes struct {
	router           common.Router
	healthController *controllers.HealthController
}

// NewHealthRoutes creates a new HealthRoutes instance
func NewHealthRoutes(
	router *gin.Engine,
	healthController *controllers.HealthController,
) *HealthRoutes {
	return &HealthRoutes{
		router:           router,
		healthController: healthController,
	}
}

// Setup sets up the health check routes
func (r *HealthRoutes) Setup() {
	health := r.router.Group("/health")
	{
		health.GET("", r.healthController.GetHealth)
		health.GET("/:component", r.healthController.GetComponentHealth)
	}
}
