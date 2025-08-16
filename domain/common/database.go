package common

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
)

// Database interface that should match your database operations
type Database interface {
	Get(dest interface{}, query string, args ...interface{}) error
	Select(dest interface{}, query string, args ...interface{}) error
	Exec(query string, args ...interface{}) (sql.Result, error)
	NamedExec(query string, arg interface{}) (sql.Result, error)
	Begin() (*sqlx.Tx, error)
	Ping() error
	Dialect() string
	Stats() sql.DBStats
}
