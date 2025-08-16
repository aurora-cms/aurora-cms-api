package infrastructure

import (
	"github.com/h4rdc0m/aurora-api/infrastructure/auth"
	"github.com/h4rdc0m/aurora-api/infrastructure/config"
	"github.com/h4rdc0m/aurora-api/infrastructure/health"
	"github.com/h4rdc0m/aurora-api/infrastructure/http"
	"github.com/h4rdc0m/aurora-api/infrastructure/http_client"
	"github.com/h4rdc0m/aurora-api/infrastructure/logging"
	"github.com/h4rdc0m/aurora-api/infrastructure/persistence"
	"github.com/h4rdc0m/aurora-api/infrastructure/services"
	"github.com/h4rdc0m/aurora-api/infrastructure/time"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"infrastructure",
	config.Module,
	logging.Module,
	time_provider.Module,
	http.Module,
	http_client.Module,
	auth.Module,
	persistence.Module,
	health.Module,
	services.Module,
)
