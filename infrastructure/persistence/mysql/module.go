package mysql

import (
	"github.com/h4rdc0m/aurora-api/domain/common"
	"github.com/h4rdc0m/aurora-api/infrastructure/config"
	"go.uber.org/fx"
)

func ProvideDatabase(cfg Config, logger common.Logger) (common.Database, error) {
	return NewDatabase(cfg, logger)
}

func ProvideConfig(env *config.Env) Config {
	return NewConfig(env)
}

var Module = fx.Options(
	fx.Provide(ProvideConfig),
	fx.Provide(ProvideDatabase),
)
