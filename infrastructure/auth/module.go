package auth

import "go.uber.org/fx"

var Module = fx.Module(
	"infrastructure.auth",
	fx.Provide(NewConfig),
	fx.Provide(NewService),
)
