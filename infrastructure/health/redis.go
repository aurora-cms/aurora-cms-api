package health

import (
	"context"
	"fmt"
	"github.com/h4rdc0m/aurora-api/domain/common"
	"github.com/redis/go-redis/v9"
	"strings"
	"time"
)

// RedisProvider encapsulates functionality for interacting with a Redis instance and provides health checks.
type RedisProvider struct {
	client *redis.Client
	logger common.Logger
}

// NewRedisProvider initializes and returns a RedisProvider instance with the provided redis client and logger.
func NewRedisProvider(client *redis.Client, logger common.Logger) CheckProvider {
	return &RedisProvider{
		client: client,
		logger: logger,
	}
}

// GetComponentName returns the name of the component ("redis") associated with the RedisProvider.
func (r *RedisProvider) GetComponentName() string {
	return "redis"
}

// Check validates the Redis connection by pinging the server and gathering basic info, returning details and possible errors.
func (r *RedisProvider) Check() (map[string]interface{}, error) {
	details := map[string]interface{}{
		"type": "redis",
	}

	// Set a timeout for the ping operation
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Try to ping Redis
	start := time.Now()
	pong, err := r.client.Ping(ctx).Result()
	latency := time.Since(start)

	if err != nil {
		return details, fmt.Errorf("redis ping failed: %w", err)
	}

	details["ping"] = pong
	details["latency_ms"] = latency.Milliseconds()

	// Get some Redis info
	info, err := r.client.Info(ctx).Result()
	if err == nil {
		parsedInfo := make(map[string]string)
		for _, line := range strings.Split(info, "\n") {
			if strings.Contains(line, ":") {
				parts := strings.SplitN(line, ":", 2)
				if len(parts) == 2 {
					parsedInfo[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
				}
			}
		}
		details["info"] = parsedInfo
	}

	return details, nil
}
