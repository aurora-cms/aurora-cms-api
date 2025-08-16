package config

import (
	"go.uber.org/fx"
)

var Module = fx.Module(
	"infrastructure.config",
	fx.Provide(NewEnv),
)
