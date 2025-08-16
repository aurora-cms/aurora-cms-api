package mocks

import (
	"database/sql"
	"github.com/h4rdc0m/aurora-api/domain/common"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/mock"
)

// Database is a mock implementation of the Database interface
type Database struct {
	mock.Mock
}

var _ common.Database = (*Database)(nil)

func (d *Database) Get(dest interface{}, query string, args ...interface{}) error {
	args = append([]interface{}{dest, query}, args...)
	call := d.Called(args...)
	if call.Error(0) != nil {
		return call.Error(0)
	}
	return nil
}

func (d *Database) Select(dest interface{}, query string, args ...interface{}) error {
	args = append([]interface{}{dest, query}, args...)
	call := d.Called(args...)
	if call.Error(0) != nil {
		return call.Error(0)
	}
	return nil
}

func (d *Database) NamedExec(query string, arg interface{}) (sql.Result, error) {
	args := d.Called(query, arg)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(sql.Result), args.Error(1)
}

func (d *Database) Exec(query string, args ...interface{}) (sql.Result, error) {
	args = append([]interface{}{query}, args...)
	call := d.Called(args...)
	if call.Get(0) == nil {
		return nil, call.Error(1)
	}
	return call.Get(0).(sql.Result), call.Error(1)
}

func (d *Database) Begin() (*sqlx.Tx, error) {
	args := d.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*sqlx.Tx), args.Error(1)
}

func (d *Database) Dialect() string {
	args := d.Called()
	return args.String(0)
}

func (d *Database) Ping() error {
	args := d.Called()
	return args.Error(0)
}

func (d *Database) Stats() sql.DBStats {
	args := d.Called()
	if args.Get(0) == nil {
		return sql.DBStats{}
	}
	return args.Get(0).(sql.DBStats)
}

type SqlResult struct {
	mock.Mock
}

func (m *SqlResult) LastInsertId() (int64, error) {
	args := m.Called()
	return args.Get(0).(int64), args.Error(1)
}

func (m *SqlResult) RowsAffected() (int64, error) {
	args := m.Called()
	return args.Get(0).(int64), args.Error(1)
}
