package auth

import "github.com/h4rdc0m/aurora-api/infrastructure/config"

type Config struct {
	KeycloakURL                string
	KeycloakRealm              string
	KeycloakClientID           string
	KeycloakClientSecret       string
	KeycloakDefaultRedirectURI string
	BaseURL                    string
}

// NewConfig creates a new AuthConfig instance with the provided environment variables.
func NewConfig(env *config.Env) Config {
	return Config{
		KeycloakURL:                env.KeycloakURL,
		KeycloakRealm:              env.KeycloakRealm,
		KeycloakClientID:           env.KeycloakClientID,
		KeycloakClientSecret:       env.KeycloakClientSecret,
		KeycloakDefaultRedirectURI: env.KeycloakDefaultRedirectURI,
		BaseURL:                    env.BaseURL,
	}
}
