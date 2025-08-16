package mysql

import (
	"database/sql"
	"fmt"
	"github.com/h4rdc0m/aurora-api/domain/common"
	"github.com/jmoiron/sqlx"
)

// Database is a wrapper around gorm.Database to provide a consistent interface for database operations
type Database struct {
	db     *sqlx.DB
	logger common.Logger // Assuming common.Logger is defined in your project
}

var _ common.Database = (*Database)(nil)

// NewDatabase creates a new instance of Database with the provided Database connection
func NewDatabase(config Config, logger common.Logger) (common.Database, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", config.Username, config.Password, config.Host, config.Port, config.Database)
	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		return nil, err // Handle error appropriately in production code
	}

	return &Database{
		db:     db,
		logger: logger,
	}, nil
}

func (d *Database) Get(dest interface{}, query string, args ...interface{}) error {
	if len(args) == 0 {
		return d.db.Get(dest, query)
	}
	err := d.db.Get(dest, query, args...)
	if err != nil {
		d.logger.Error("Failed to execute Get query", "query", query, "args", args, "error", err)
		return err
	}
	return nil
}

func (d *Database) Select(dest interface{}, query string, args ...interface{}) error {
	if len(args) == 0 {
		return d.db.Select(dest, query)
	}
	err := d.db.Select(dest, query, args...)
	if err != nil {
		d.logger.Error("Failed to execute Select query", "query", query, "args", args, "error", err)
		return err
	}
	return nil
}

func (d *Database) NamedExec(query string, arg interface{}) (sql.Result, error) {
	result, err := d.db.NamedExec(query, arg)
	if err != nil {
		d.logger.Error("Failed to execute NamedExec query", "query", query, "arg", arg, "error", err)
		return nil, err
	}
	return result, nil
}

func (d *Database) Exec(query string, args ...interface{}) (sql.Result, error) {
	if len(args) == 0 {
		return d.db.Exec(query)
	}
	result, err := d.db.Exec(query, args...)
	if err != nil {
		d.logger.Error("Failed to execute Exec query", "query", query, "args", args, "error", err)
		return nil, err
	}
	return result, nil
}

func (d *Database) Begin() (*sqlx.Tx, error) {
	tx, err := d.db.Beginx()
	if err != nil {
		d.logger.Error("Failed to begin transaction", "error", err)
		return nil, err
	}
	return tx, nil
}

func (d *Database) Ping() error {
	err := d.db.Ping()
	if err != nil {
		d.logger.Error("Failed to ping database", "error", err)
		return err
	}
	return nil
}

func (d *Database) Dialect() string {
	return d.db.DriverName()
}

func (d *Database) Stats() sql.DBStats {
	stats := d.db.Stats()
	if stats.MaxOpenConnections == 0 {
		d.logger.Warn("Database stats indicate no open connections", "stats", stats)
	}
	return stats
}
