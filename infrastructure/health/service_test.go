package health

import (
	"errors"
	"github.com/h4rdc0m/aurora-api/domain/entities"
	"github.com/h4rdc0m/aurora-api/domain/services"
	"testing"
)

type mockLogger struct{}

func (m *mockLogger) Debug(_ string, _ ...interface{})  {}
func (m *mockLogger) Info(_ string, _ ...interface{})   {}
func (m *mockLogger) Warn(_ string, _ ...interface{})   {}
func (m *mockLogger) Error(_ string, _ ...interface{})  {}
func (m *mockLogger) Fatal(_ string, _ ...interface{})  {}
func (m *mockLogger) Panic(_ string, _ ...interface{})  {}
func (m *mockLogger) Print(_ string, _ ...interface{})  {}
func (m *mockLogger) Debugf(_ string, _ ...interface{}) {}
func (m *mockLogger) Infof(_ string, _ ...interface{})  {}
func (m *mockLogger) Warnf(_ string, _ ...interface{})  {}
func (m *mockLogger) Errorf(_ string, _ ...interface{}) {}
func (m *mockLogger) Fatalf(_ string, _ ...interface{}) {}
func (m *mockLogger) Panicf(_ string, _ ...interface{}) {}
func (m *mockLogger) Printf(_ string, _ ...interface{}) {}

type mockCheckProvider struct {
	componentName string
	checkResult   map[string]interface{}
	checkError    error
}

func (m *mockCheckProvider) Check() (map[string]interface{}, error) {
	return m.checkResult, m.checkError
}

func (m *mockCheckProvider) GetComponentName() string {
	return m.componentName
}

func TestService_CheckHealth(t *testing.T) {
	tests := []struct {
		name      string
		providers []CheckProvider
		expected  map[string]services.HealthStatus
	}{
		{
			name: "All components healthy",
			providers: []CheckProvider{
				&mockCheckProvider{
					componentName: "component1",
					checkResult:   map[string]interface{}{"key": "value"},
					checkError:    nil,
				},
				&mockCheckProvider{
					componentName: "component2",
					checkResult:   map[string]interface{}{"key2": "value2"},
					checkError:    nil,
				},
			},
			expected: map[string]services.HealthStatus{
				"component1": {
					Status:    entities.HealthStatusUp,
					Component: "component1",
					Details:   map[string]interface{}{"key": "value"},
				},
				"component2": {
					Status:    entities.HealthStatusUp,
					Component: "component2",
					Details:   map[string]interface{}{"key2": "value2"},
				},
			},
		},
		{
			name: "Some components down",
			providers: []CheckProvider{
				&mockCheckProvider{
					componentName: "component1",
					checkResult:   nil,
					checkError:    errors.New("connection failed"),
				},
				&mockCheckProvider{
					componentName: "component2",
					checkResult:   map[string]interface{}{"key2": "value2"},
					checkError:    nil,
				},
			},
			expected: map[string]services.HealthStatus{
				"component1": {
					Status:    entities.HealthStatusDown,
					Component: "component1",
					Error:     "connection failed",
				},
				"component2": {
					Status:    entities.HealthStatusUp,
					Component: "component2",
					Details:   map[string]interface{}{"key2": "value2"},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewService(&mockLogger{}, tt.providers)
			result := service.CheckHealth()
			for component, expectedStatus := range tt.expected {
				if result[component].Status != expectedStatus.Status {
					t.Errorf("unexpected health status for component %s. Got %+v, expected %+v", component, result[component], expectedStatus)
				}
			}
		})
	}
}

func TestService_CheckComponent(t *testing.T) {
	tests := []struct {
		name      string
		providers []CheckProvider
		component string
		expected  services.HealthStatus
	}{
		{
			name: "Component exists and is healthy",
			providers: []CheckProvider{
				&mockCheckProvider{
					componentName: "component1",
					checkResult:   map[string]interface{}{"key": "value"},
					checkError:    nil,
				},
			},
			component: "component1",
			expected: services.HealthStatus{
				Status:    entities.HealthStatusUp,
				Component: "component1",
				Details:   map[string]interface{}{"key": "value"},
			},
		},
		{
			name: "Component exists but is down",
			providers: []CheckProvider{
				&mockCheckProvider{
					componentName: "component1",
					checkResult:   nil,
					checkError:    errors.New("connection failed"),
				},
			},
			component: "component1",
			expected: services.HealthStatus{
				Status:    entities.HealthStatusDown,
				Component: "component1",
				Error:     "connection failed",
			},
		},
		{
			name:      "Component does not exist",
			providers: []CheckProvider{},
			component: "unknown",
			expected: services.HealthStatus{
				Status:    entities.HealthStatusDown,
				Component: "unknown",
				Error:     "Component not found",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewService(&mockLogger{}, tt.providers)
			result := service.CheckComponent(tt.component)
			if result.Status != tt.expected.Status {
				t.Errorf("unexpected health status. Got %+v, expected %+v", result, tt.expected)
			}
		})
	}
}
