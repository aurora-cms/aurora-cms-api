package health

import (
	"context"
	"fmt"
	"github.com/h4rdc0m/aurora-api/domain/common"
	"time"
)

type DatabaseProvider struct {
	db     common.Database
	logger common.Logger
}

func NewDatabaseProvider(db common.Database, logger common.Logger) CheckProvider {
	return &DatabaseProvider{
		db:     db,
		logger: logger,
	}
}

func (d *DatabaseProvider) GetComponentName() string {
	return "database"
}

func (d *DatabaseProvider) Check() (map[string]interface{}, error) {
	details := map[string]interface{}{
		"type": d.db.Dialect(),
	}

	_, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := d.db.Ping()
	if err != nil {
		return details, fmt.Errorf("database ping failed: %w", err)
	}

	stats := d.db.Stats()
	details["connections"] = map[string]interface{}{
		"open":    stats.OpenConnections,
		"idle":    stats.Idle,
		"inUse":   stats.InUse,
		"maxOpen": stats.MaxOpenConnections,
	}

	return details, nil
}
