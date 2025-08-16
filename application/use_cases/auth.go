package use_cases

import (
	"github.com/h4rdc0m/aurora-api/domain/common"
	"github.com/h4rdc0m/aurora-api/domain/entities"
	"github.com/h4rdc0m/aurora-api/domain/services"
)

type AuthUseCase struct {
	authService services.AuthService
	logger      common.Logger
}

func NewAuthUseCase(authService services.AuthService, logger common.Logger) *AuthUseCase {
	return &AuthUseCase{
		authService: authService,
		logger:      logger,
	}
}

func (a *AuthUseCase) InitiateLogin(redirectURI string) (*entities.LoginResponse, *entities.AuthSession, error) {
	req := entities.LoginRequest{
		RedirectURI: redirectURI,
	}

	res, session, err := a.authService.InitiateLogin(req)
	if err != nil {
		a.logger.Error("Failed to initiate login", "error", err)
		return nil, nil, err
	}
	return res, session, nil
}

func (a *AuthUseCase) HandleCallback(code string, session *entities.AuthSession) (*entities.TokenResponse, error) {
	tokens, err := a.authService.ExchangeCodeForTokens(code, session)
	if err != nil {
		a.logger.Error("Failed to exchange code for tokens", "error", err)
		return nil, err
	}
	return tokens, nil
}

func (a *AuthUseCase) RefreshToken(refreshToken string) (*entities.TokenResponse, error) {
	tokens, err := a.authService.RefreshToken(refreshToken)
	if err != nil {
		a.logger.Error("Failed to refresh tokens", "error", err)
		return nil, err
	}
	return tokens, nil
}

func (a *AuthUseCase) Logout(refreshToken string, redirectURI string) (*entities.LogoutResponse, error) {
	req := entities.LogoutRequest{
		RefreshToken: refreshToken,
		RedirectURI:  redirectURI,
	}

	res, err := a.authService.Logout(req)
	if err != nil {
		a.logger.Error("Failed to logout", "error", err)
		return nil, err
	}

	return res, nil
}

// GetAuthConfig returns the authentication configuration
func (a *AuthUseCase) GetAuthConfig() map[string]interface{} {
	return a.authService.GetAuthConfig()
}

// GetDiscoveryConfig returns the OpenID Connect discovery configuration
func (a *AuthUseCase) GetDiscoveryConfig() map[string]interface{} {
	return a.authService.GetDiscoveryConfig()
}
