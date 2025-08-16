package health

import (
	"github.com/h4rdc0m/aurora-api/domain/common"
	"github.com/h4rdc0m/aurora-api/domain/services"
	"github.com/h4rdc0m/aurora-api/infrastructure/auth"
	"github.com/redis/go-redis/v9"
	"go.uber.org/fx"
)

func ProvideCheckProviders(
	db common.Database,
	redisClient *redis.Client,
	authConfig auth.Config,
	logger common.Logger,
) []CheckProvider {
	return []CheckProvider{
		NewDatabaseProvider(db, logger),
		NewRedisProvider(redisClient, logger),
		NewKeycloakProvider(authConfig, logger),
	}
}

func ProvideHealthService(
	checkProviders []CheckProvider,
	logger common.Logger,
) services.HealthService {
	return NewService(logger, checkProviders)
}

var Module = fx.Module(
	"infrastructure.health",
	fx.Provide(ProvideCheckProviders),
	fx.Provide(ProvideHealthService),
)
