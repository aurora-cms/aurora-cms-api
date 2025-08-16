package mysql

import "github.com/h4rdc0m/aurora-api/infrastructure/config"

// Config holds the database connection configuration details including host, port, credentials, and database name.
type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	Database string
}

// NewConfig creates a new DBConfig instance by using values from the provided Env configuration.
func NewConfig(env *config.Env) Config {
	return Config{
		Host:     env.DBHost,
		Port:     env.DBPort,
		Username: env.DBUser,
		Password: env.DBPassword,
		Database: env.DBName,
	}
}
