package use_cases

import (
	"github.com/h4rdc0m/aurora-api/domain/common"
	"github.com/h4rdc0m/aurora-api/domain/entities"
	"github.com/h4rdc0m/aurora-api/domain/services"
	"time"
)

type HealthUseCase struct {
	healthService services.HealthService
	logger        common.Logger
}

func NewHealthUseCase(healthService services.HealthService, logger common.Logger) *HealthUseCase {
	return &HealthUseCase{
		healthService: healthService,
		logger:        logger,
	}
}

func (h *HealthUseCase) CheckHealth() *entities.HealthCheckResult {
	healthResults := h.healthService.CheckHealth()

	result := &entities.HealthCheckResult{
		Status:     entities.HealthStatusUp,
		Components: make(map[string]entities.ComponentHealth),
		Timestamp:  time.Now().UTC().Format(time.RFC3339),
	}

	for name, status := range healthResults {
		componentHealth := entities.ComponentHealth{
			Status:  status.Status,
			Details: status.Details,
			Error:   status.Error,
		}

		result.Components[name] = componentHealth

		if status.Status != entities.HealthStatusUp {
			result.Status = entities.HealthStatusDown
		}
	}

	return result
}

// CheckComponent checks the health of a specific component
func (h *HealthUseCase) CheckComponent(component string) entities.ComponentHealth {
	healthStatus := h.healthService.CheckComponent(component)

	return entities.ComponentHealth{
		Status:  healthStatus.Status,
		Details: healthStatus.Details,
		Error:   healthStatus.Error,
	}
}
