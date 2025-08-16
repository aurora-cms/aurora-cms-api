package logging

import "github.com/h4rdc0m/aurora-api/infrastructure/config"

type LogConfig struct {
	Environment string
	LogLevel    string
	LogOutput   string
}

// NewLogConfig creates and initializes a new LogConfig instance using the specified environment configuration.
func NewLogConfig(env *config.Env) *LogConfig {
	return &LogConfig{
		Environment: env.Environment,
		LogLevel:    env.LogLevel,
		LogOutput:   env.LogOutput,
	}
}
