package health

import (
	"github.com/h4rdc0m/aurora-api/domain/common"
	"github.com/h4rdc0m/aurora-api/domain/entities"
	"github.com/h4rdc0m/aurora-api/domain/services"
)

type CheckProvider interface {
	Check() (map[string]interface{}, error)
	GetComponentName() string
}

type Service struct {
	logger     common.Logger
	components map[string]CheckProvider
}

func NewService(logger common.Logger, providers []CheckProvider) services.HealthService {
	components := make(map[string]CheckProvider)
	for _, provider := range providers {
		components[provider.GetComponentName()] = provider
	}

	return &Service{
		logger:     logger,
		components: components,
	}
}

func (s *Service) CheckHealth() map[string]services.HealthStatus {
	res := make(map[string]services.HealthStatus)

	for _, provider := range s.components {
		componentName := provider.GetComponentName()
		status := s.checkProvider(provider)
		res[componentName] = status
	}

	return res
}

func (s *Service) CheckComponent(component string) services.HealthStatus {
	provider, exists := s.components[component]
	if !exists {
		return services.HealthStatus{
			Status:    entities.HealthStatusDown,
			Component: component,
			Error:     "Component not found",
		}
	}

	return s.checkProvider(provider)
}

func (s *Service) checkProvider(provider CheckProvider) services.HealthStatus {
	componentName := provider.GetComponentName()
	details, err := provider.Check()

	status := services.HealthStatus{
		Component: componentName,
		Details:   details,
	}

	if err != nil {
		status.Status = entities.HealthStatusDown
		status.Error = err.Error()
	} else {
		status.Status = entities.HealthStatusUp
	}

	return status
}
