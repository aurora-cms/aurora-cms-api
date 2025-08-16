package health

import (
	"context"
	"fmt"
	"github.com/h4rdc0m/aurora-api/domain/common"
	"github.com/h4rdc0m/aurora-api/infrastructure/auth"
	"io"
	"net/http"
	"time"
)

// KeycloakHealthProvider checks the health of the Keycloak connection
type KeycloakHealthProvider struct {
	config auth.Config
	logger common.Logger
	client *http.Client
}

// NewKeycloakProvider creates a new KeycloakHealthProvider
func NewKeycloakProvider(
	config auth.Config,
	logger common.Logger,
) CheckProvider {
	return &KeycloakHealthProvider{
		config: config,
		logger: logger,
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

// GetComponentName returns the name of this component
func (k *KeycloakHealthProvider) GetComponentName() string {
	return "keycloak"
}

// Check performs a health check on the Keycloak connection
func (k *KeycloakHealthProvider) Check() (map[string]interface{}, error) {
	details := map[string]interface{}{
		"url": k.config.KeycloakURL,
	}

	// Check well-known endpoint
	wellKnownURL := fmt.Sprintf(
		"%s/realms/%s/.well-known/openid-configuration",
		k.config.KeycloakURL,
		k.config.KeycloakRealm,
	)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", wellKnownURL, nil)
	if err != nil {
		return details, fmt.Errorf("failed to create request: %w", err)
	}

	start := time.Now()
	resp, err := k.client.Do(req)
	latency := time.Since(start)

	if err != nil {
		return details, fmt.Errorf("failed to connect to Keycloak: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			k.logger.Warn("Failed to close Keycloak response body", "error", err)
		} else {
			k.logger.Debug("Closed Keycloak response body successfully")
		}
	}(resp.Body)

	details["latency_ms"] = latency.Milliseconds()

	if resp.StatusCode != http.StatusOK {
		return details, fmt.Errorf(
			"keycloak returned non-OK status: %d %s",
			resp.StatusCode,
			resp.Status,
		)
	}
	details["type"] = "keycloak"
	details["status"] = resp.Status
	return details, nil
}
