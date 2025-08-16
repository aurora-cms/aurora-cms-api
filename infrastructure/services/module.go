package services

import "go.uber.org/fx"

var Module = fx.Module(
	"infrastructure.services",
	fx.Provide(NewTokenServiceConfig),
	fx.Provide(NewTokenService),
	fx.Provide(NewSessionService),
)
