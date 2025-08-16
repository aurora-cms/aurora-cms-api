package mysql

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

// Simple mock logger for testing
type mockLogger struct{}

func (m *mockLogger) Panic(_ string, _ ...interface{}) {

}

func (m *mockLogger) Print(_ string, _ ...interface{}) {

}

func (m *mockLogger) Debugf(_ string, _ ...interface{}) {

}

func (m *mockLogger) Infof(_ string, _ ...interface{}) {

}

func (m *mockLogger) Warnf(_ string, _ ...interface{}) {

}

func (m *mockLogger) Errorf(_ string, _ ...interface{}) {

}

func (m *mockLogger) Fatalf(_ string, _ ...interface{}) {

}

func (m *mockLogger) Panicf(_ string, _ ...interface{}) {

}

func (m *mockLogger) Printf(_ string, _ ...interface{}) {

}

func (m *mockLogger) Debug(_ string, _ ...interface{}) {}
func (m *mockLogger) Info(_ string, _ ...interface{})  {}
func (m *mockLogger) Warn(_ string, _ ...interface{})  {}
func (m *mockLogger) Error(_ string, _ ...interface{}) {}
func (m *mockLogger) Fatal(_ string, _ ...interface{}) {}

func TestDatabase_Get(t *testing.T) {
	logger := &mockLogger{}

	testCases := []struct {
		name       string
		query      string
		args       []interface{}
		dbResponse error
		expectErr  bool
	}{
		{"success_no_args", "SELECT * FROM users", nil, nil, false},
		{"success_with_args", "SELECT * FROM users WHERE id = ?", []interface{}{1}, nil, false},
		{"db_error", "SELECT * FROM users WHERE id = ?", []interface{}{1}, sql.ErrNoRows, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db, mockDb, err := sqlmock.New()
			assert.NoError(t, err)
			defer func(db *sql.DB) {
				err := db.Close()
				if err != nil {
					return
				}
			}(db)

			database := &Database{db: sqlx.NewDb(db, "mysql"), logger: logger}

			dest := struct {
				ID int `db:"id"`
			}{}

			if tc.expectErr {
				if len(tc.args) > 0 {
					mockDb.ExpectQuery("SELECT \\* FROM users WHERE id = \\?").WithArgs(tc.args[0]).WillReturnError(tc.dbResponse)
				} else {
					mockDb.ExpectQuery("SELECT \\* FROM users").WillReturnError(tc.dbResponse)
				}
			} else {
				rows := sqlmock.NewRows([]string{"id"}).AddRow(1)
				if len(tc.args) > 0 {
					mockDb.ExpectQuery("SELECT \\* FROM users WHERE id = \\?").WithArgs(tc.args[0]).WillReturnRows(rows)
				} else {
					mockDb.ExpectQuery("SELECT \\* FROM users").WillReturnRows(rows)
				}
			}

			err = database.Get(&dest, tc.query, tc.args...)

			if tc.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.NoError(t, mockDb.ExpectationsWereMet())
		})
	}
}

func TestDatabase_Select(t *testing.T) {
	logger := new(mockLogger)

	testCases := []struct {
		name       string
		query      string
		args       []interface{}
		dbResponse error
		expectErr  bool
	}{
		{"success_no_args", "SELECT * FROM users", nil, nil, false},
		{"success_with_args", "SELECT * FROM users WHERE id = ?", []interface{}{1}, nil, false},
		{"db_error", "SELECT * FROM users WHERE id = ?", []interface{}{1}, sql.ErrNoRows, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db, mockDb, err := sqlmock.New()
			assert.NoError(t, err)
			defer func(db *sql.DB) {
				err := db.Close()
				if err != nil {
					return
				}
			}(db)

			database := &Database{db: sqlx.NewDb(db, "mysql"), logger: logger}

			var dest []struct {
				ID int `db:"id"`
			}

			if tc.expectErr {
				if len(tc.args) > 0 {
					mockDb.ExpectQuery("SELECT \\* FROM users WHERE id = \\?").WithArgs(tc.args[0]).WillReturnError(tc.dbResponse)
				} else {
					mockDb.ExpectQuery("SELECT \\* FROM users").WillReturnError(tc.dbResponse)
				}
			} else {
				rows := sqlmock.NewRows([]string{"id"}).AddRow(1)
				if len(tc.args) > 0 {
					mockDb.ExpectQuery("SELECT \\* FROM users WHERE id = \\?").WithArgs(tc.args[0]).WillReturnRows(rows)
				} else {
					mockDb.ExpectQuery("SELECT \\* FROM users").WillReturnRows(rows)
				}
			}

			err = database.Select(&dest, tc.query, tc.args...)

			if tc.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.NoError(t, mockDb.ExpectationsWereMet())
		})
	}
}

func TestDatabase_NamedExec(t *testing.T) {
	logger := &mockLogger{}

	testCases := []struct {
		name       string
		query      string
		arg        interface{}
		dbResponse error
		expectErr  bool
	}{
		{"success", "INSERT INTO users (name) VALUES (:name)", map[string]interface{}{"name": "John"}, nil, false},
		{"db_error", "INSERT INTO users (name) VALUES (:name)", map[string]interface{}{"name": "John"}, sql.ErrConnDone, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db, mockDb, err := sqlmock.New()
			assert.NoError(t, err)
			defer func(db *sql.DB) {
				err := db.Close()
				if err != nil {
					return
				}
			}(db)

			database := &Database{db: sqlx.NewDb(db, "mysql"), logger: logger}

			// NamedExec converts named parameters to positional parameters
			// The query becomes: "INSERT INTO users (name) VALUES (?)"
			// And the argument becomes just the value: "John"
			if tc.expectErr {
				mockDb.ExpectExec("INSERT INTO users \\(name\\) VALUES \\(\\?\\)").WithArgs("John").WillReturnError(tc.dbResponse)
			} else {
				mockDb.ExpectExec("INSERT INTO users \\(name\\) VALUES \\(\\?\\)").WithArgs("John").WillReturnResult(sqlmock.NewResult(1, 1))
			}

			_, err = database.NamedExec(tc.query, tc.arg)

			if tc.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.NoError(t, mockDb.ExpectationsWereMet())
		})
	}
}

func TestDatabase_Exec(t *testing.T) {
	logger := &mockLogger{}

	testCases := []struct {
		name       string
		query      string
		args       []interface{}
		dbResponse error
		expectErr  bool
	}{
		{"success_no_args", "DELETE FROM users", nil, nil, false},
		{"success_with_args", "DELETE FROM users WHERE id = ?", []interface{}{1}, nil, false},
		{"db_error", "DELETE FROM users WHERE id = ?", []interface{}{1}, sql.ErrConnDone, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db, mockDb, err := sqlmock.New()
			assert.NoError(t, err)
			defer func(db *sql.DB) {
				err := db.Close()
				if err != nil {
					return
				}
			}(db)

			database := &Database{db: sqlx.NewDb(db, "mysql"), logger: logger}

			if tc.expectErr {
				if len(tc.args) > 0 {
					mockDb.ExpectExec("DELETE FROM users WHERE id = \\?").WithArgs(tc.args[0]).WillReturnError(tc.dbResponse)
				} else {
					mockDb.ExpectExec("DELETE FROM users").WillReturnError(tc.dbResponse)
				}
			} else {
				if len(tc.args) > 0 {
					mockDb.ExpectExec("DELETE FROM users WHERE id = \\?").WithArgs(tc.args[0]).WillReturnResult(sqlmock.NewResult(1, 1))
				} else {
					mockDb.ExpectExec("DELETE FROM users").WillReturnResult(sqlmock.NewResult(1, 1))
				}
			}

			_, err = database.Exec(tc.query, tc.args...)

			if tc.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.NoError(t, mockDb.ExpectationsWereMet())
		})
	}
}

func TestDatabase_Begin(t *testing.T) {
	logger := &mockLogger{}
	db, mockDb, err := sqlmock.New()
	assert.NoError(t, err)
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			return
		}
	}(db)

	database := &Database{db: sqlx.NewDb(db, "mysql"), logger: logger}

	mockDb.ExpectBegin()

	tx, err := database.Begin()
	assert.NoError(t, err)
	assert.NotNil(t, tx)

	assert.NoError(t, mockDb.ExpectationsWereMet())
}

func TestDatabase_Ping(t *testing.T) {
	logger := &mockLogger{}

	testCases := []struct {
		name       string
		dbResponse error
		expectErr  bool
	}{
		{"success", nil, false},
		{"db_error", sql.ErrConnDone, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Enable ping monitoring for this test
			db, mockDb, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
			assert.NoError(t, err)
			defer func(db *sql.DB) {
				err := db.Close()
				if err != nil {
					return
				}
			}(db)

			database := &Database{db: sqlx.NewDb(db, "mysql"), logger: logger}

			if tc.expectErr {
				mockDb.ExpectPing().WillReturnError(tc.dbResponse)
			} else {
				mockDb.ExpectPing()
			}

			err = database.Ping()

			if tc.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.NoError(t, mockDb.ExpectationsWereMet())
		})
	}
}

func TestDatabase_Dialect(t *testing.T) {
	logger := &mockLogger{}
	db, _, err := sqlmock.New()
	assert.NoError(t, err)
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			return
		}
	}(db)

	database := &Database{db: sqlx.NewDb(db, "mysql"), logger: logger}

	assert.Equal(t, "mysql", database.Dialect())
}

func TestDatabase_Stats(t *testing.T) {
	logger := &mockLogger{}
	db, _, err := sqlmock.New()
	assert.NoError(t, err)
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			return
		}
	}(db)

	database := &Database{db: sqlx.NewDb(db, "mysql"), logger: logger}

	stats := database.Stats()
	assert.NotNil(t, stats)
}
