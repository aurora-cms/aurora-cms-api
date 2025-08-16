package services

import "github.com/h4rdc0m/aurora-api/infrastructure/config"

type TokenServiceConfig struct {
	KeycloakURL string
	Realm       string
}

func NewTokenServiceConfig(env *config.Env) TokenServiceConfig {
	return TokenServiceConfig{
		KeycloakURL: env.KeycloakURL,
		Realm:       env.KeycloakRealm,
	}
}
