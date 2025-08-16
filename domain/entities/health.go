package entities

// HealthCheckResult represents the result of a health check for an application, including overall status and components.
// Status indicates the overall health status of the application, typically "healthy" or "unhealthy".
// Components provides a map of individual component names to their corresponding health details.
// Timestamp represents the time when the health check was performed or reported.
type HealthCheckResult struct {
	Status     string                     `json:"status"`
	Components map[string]ComponentHealth `json:"components"`
	Timestamp  string                     `json:"timestamp"`
}

// ComponentHealth represents the health status of a system component.
// Status indicates the overall health as a string value.
// Details provide additional information about the component in key-value pairs.
// Error contains error information if there is a health issue, and is omitted if no error exists.
type ComponentHealth struct {
	Status  string                 `json:"status"`
	Details map[string]interface{} `json:"details"`
	Error   string                 `json:"error,omitempty"`
}

// HealthStatusUp represents a healthy state of the system or component.
// HealthStatusDown represents an unhealthy state of the system or component.
// HealthStatusDegraded represents a partially functional or degraded state of the system or component.
const (
	HealthStatusUp       = "UP"
	HealthStatusDown     = "DOWN"
	HealthStatusDegraded = "DEGRADED"
)
