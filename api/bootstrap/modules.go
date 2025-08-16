package bootstrap

import (
	"github.com/h4rdc0m/aurora-api/api"
	"github.com/h4rdc0m/aurora-api/application"
	"github.com/h4rdc0m/aurora-api/infrastructure"
	"go.uber.org/fx"
)

var CommonModules = fx.Options(
	infrastructure.Module,
	application.Module,
	api.Module,
)
