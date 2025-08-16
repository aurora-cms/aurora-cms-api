package services

import (
	"testing"

	"github.com/h4rdc0m/aurora-api/infrastructure/config"
)

func TestNewTokenServiceConfig(t *testing.T) {
	tests := []struct {
		name       string
		env        config.Env
		wantConfig TokenServiceConfig
	}{
		{
			name: "ValidInput",
			env: config.Env{
				KeycloakURL:   "https://keycloak.example.com",
				KeycloakRealm: "example-realm",
			},
			wantConfig: TokenServiceConfig{
				KeycloakURL: "https://keycloak.example.com",
				Realm:       "example-realm",
			},
		},
		{
			name: "EmptyValues",
			env: config.Env{
				KeycloakURL:   "",
				KeycloakRealm: "",
			},
			wantConfig: TokenServiceConfig{
				KeycloakURL: "",
				Realm:       "",
			},
		},
		{
			name: "PartialValues",
			env: config.Env{
				KeycloakURL:   "https://keycloak.partial.com",
				KeycloakRealm: "",
			},
			wantConfig: TokenServiceConfig{
				KeycloakURL: "https://keycloak.partial.com",
				Realm:       "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotConfig := NewTokenServiceConfig(&tt.env)

			if gotConfig.KeycloakURL != tt.wantConfig.KeycloakURL {
				t.Errorf("NewTokenServiceConfig().KeycloakURL = %v; want %v", gotConfig.KeycloakURL, tt.wantConfig.KeycloakURL)
			}

			if gotConfig.Realm != tt.wantConfig.Realm {
				t.Errorf("NewTokenServiceConfig().Realm = %v; want %v", gotConfig.Realm, tt.wantConfig.Realm)
			}
		})
	}
}
