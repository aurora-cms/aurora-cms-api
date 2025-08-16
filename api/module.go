package api

import (
	"github.com/h4rdc0m/aurora-api/api/http/controllers"
	"github.com/h4rdc0m/aurora-api/api/http/middlewares"
	"github.com/h4rdc0m/aurora-api/api/http/routes"
	"go.uber.org/fx"
)

var Module = fx.Options(
	controllers.Module,
	routes.Module,
	middlewares.Module,
)
