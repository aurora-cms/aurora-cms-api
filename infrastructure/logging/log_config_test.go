package logging

import (
	"reflect"
	"testing"

	"github.com/h4rdc0m/aurora-api/infrastructure/config"
)

func TestNewLogConfig(t *testing.T) {
	tests := []struct {
		name       string
		input      *config.Env
		wantConfig *LogConfig
	}{
		{
			name: "Valid environment with all fields",
			input: &config.Env{
				Environment: "production",
				LogLevel:    "info",
				LogOutput:   "stdout",
			},
			wantConfig: &LogConfig{
				Environment: "production",
				LogLevel:    "info",
				LogOutput:   "stdout",
			},
		},
		{
			name: "Empty environment fields",
			input: &config.Env{
				Environment: "",
				LogLevel:    "",
				LogOutput:   "",
			},
			wantConfig: &LogConfig{
				Environment: "",
				LogLevel:    "",
				LogOutput:   "",
			},
		},
		{
			name: "Custom environment fields",
			input: &config.Env{
				Environment: "development",
				LogLevel:    "debug",
				LogOutput:   "/var/log/app.log",
			},
			wantConfig: &LogConfig{
				Environment: "development",
				LogLevel:    "debug",
				LogOutput:   "/var/log/app.log",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewLogConfig(tt.input)
			if !reflect.DeepEqual(got, tt.wantConfig) {
				t.Errorf("NewLogConfig() = %+v, want %+v", got, tt.wantConfig)
			}
		})
	}
}
