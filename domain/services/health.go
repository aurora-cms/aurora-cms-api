package services

// HealthService defines methods for evaluating the health state of various components in the system.
type HealthService interface {
	CheckHealth() map[string]HealthStatus

	CheckComponent(component string) HealthStatus
}

// HealthStatus represents the health information of a system or a component.
// Status indicates the health status (e.g., healthy, degraded, etc.).
// Component specifies the system or component name being monitored.
// Details provides additional information about the health status, optional.
// Error contains error details if the component/system is in an unhealthy state, optional.
type HealthStatus struct {
	Status    string                 `json:"status"`
	Component string                 `json:"component"`
	Details   map[string]interface{} `json:"details,omitempty"`
	Error     string                 `json:"error,omitempty"`
}
