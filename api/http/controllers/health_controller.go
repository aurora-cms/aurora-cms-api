package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/h4rdc0m/aurora-api/application/use_cases"
	"github.com/h4rdc0m/aurora-api/domain/common"
	"net/http"
)

// HealthController handles health check endpoints
type HealthController struct {
	healthUseCase *use_cases.HealthUseCase
	logger        common.Logger
}

// NewHealthController creates a new HealthController
func NewHealthController(
	healthUseCase *use_cases.HealthUseCase,
	logger common.Logger,
) *HealthController {
	return &HealthController{
		healthUseCase: healthUseCase,
		logger:        logger,
	}
}

// GetHealth checks the health of all components
func (h *HealthController) GetHealth(c *gin.Context) {
	result := h.healthUseCase.CheckHealth()

	statusCode := http.StatusOK
	if result.Status != "UP" {
		statusCode = http.StatusServiceUnavailable
	}

	c.JSON(statusCode, result)
}

// GetComponentHealth checks the health of a specific component
func (h *HealthController) GetComponentHealth(c *gin.Context) {
	component := c.Param("component")
	result := h.healthUseCase.CheckComponent(component)

	statusCode := http.StatusOK
	if result.Status != "UP" {
		statusCode = http.StatusServiceUnavailable
	}

	c.JSON(statusCode, result)
}
