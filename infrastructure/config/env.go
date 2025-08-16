package config

import (
	"github.com/spf13/viper"
	"log"
)

// Env holds the environment configuration for the Aurora API.
type Env struct {
	ServerPort                 string `mapstructure:"AURORA_SERVER_PORT"`
	Environment                string `mapstructure:"AURORA_ENVIRONMENT"`
	BaseURL                    string `mapstructure:"AURORA_BASE_URL"`
	LogOutput                  string `mapstructure:"AURORA_LOG_OUTPUT"`
	LogLevel                   string `mapstructure:"AURORA_LOG_LEVEL"`
	DBUser                     string `mapstructure:"AURORA_DB_USER"`
	DBPassword                 string `mapstructure:"AURORA_DB_PASSWORD"`
	DBHost                     string `mapstructure:"AURORA_DB_HOST"`
	DBPort                     string `mapstructure:"AURORA_DB_PORT"`
	DBName                     string `mapstructure:"AURORA_DB_NAME"`
	RedisHost                  string `mapstructure:"AURORA_REDIS_HOST"`
	RedisPort                  string `mapstructure:"AURORA_REDIS_PORT"`
	RedisDB                    int    `mapstructure:"AURORA_REDIS_DB"`
	RedisUsername              string `mapstructure:"AURORA_REDIS_USERNAME"`
	RedisPassword              string `mapstructure:"AURORA_REDIS_PASSWORD"`
	KeycloakURL                string `mapstructure:"AURORA_KEYCLOAK_URL"`
	KeycloakRealm              string `mapstructure:"AURORA_KEYCLOAK_REALM"`
	KeycloakClientID           string `mapstructure:"AURORA_KEYCLOAK_CLIENT_ID"`
	KeycloakClientSecret       string `mapstructure:"AURORA_KEYCLOAK_CLIENT_SECRET"`
	KeycloakDefaultRedirectURI string `mapstructure:"AURORA_KEYCLOAK_DEFAULT_REDIRECT_URI"`
}

// NewEnv initializes and returns an Env struct by reading and unmarshaling the configuration from a .env file.
func NewEnv() *Env {
	env := Env{}
	viper.SetConfigFile(".env")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}

	err = viper.Unmarshal(&env)
	if err != nil {
		log.Fatalf("unable to decode into struct, %v", err)
	}

	return &env
}
