package persistence

import (
	"fmt"
	"github.com/h4rdc0m/aurora-api/infrastructure/config"
	"github.com/h4rdc0m/aurora-api/infrastructure/persistence/mappers"
	"github.com/h4rdc0m/aurora-api/infrastructure/persistence/mysql"
	"github.com/h4rdc0m/aurora-api/infrastructure/persistence/repositories"
	"github.com/redis/go-redis/v9"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"infrastructure.persistence",
	fx.Provide(func(env *config.Env) *redis.Client {
		return redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%s", env.RedisHost, env.RedisPort),
			Username: env.RedisUsername,
			Password: env.RedisPassword,
			DB:       env.RedisDB,
		})
	}),
	mysql.Module,
	mappers.Module,
	repositories.Module,
)
