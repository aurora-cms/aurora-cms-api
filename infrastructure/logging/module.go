package logging

import "go.uber.org/fx"

// Module bundles multiple logger providers into a fx.Option for dependency injection in the application.
var Module = fx.Module(
	"infrastructure.logging",
	fx.Provide(NewZapLogger),
	fx.Provide(NewGinLogger),
	fx.Provide(NewLogConfig),
)
