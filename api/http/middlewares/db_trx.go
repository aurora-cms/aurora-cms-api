package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/h4rdc0m/aurora-api/constants"
	"github.com/h4rdc0m/aurora-api/domain/common"
	"go.uber.org/zap"
	"net/http"
)

// DatabaseTrx is a middleware that handles database transactions.
type DatabaseTrx struct {
	handler common.Router
	logger  common.Logger
	db      common.Database
}

// NewDatabaseTrx creates a new instance of DatabaseTrx middleware.
func NewDatabaseTrx(handler common.Router, logger common.Logger, db common.Database) *DatabaseTrx {
	return &DatabaseTrx{
		handler: handler,
		logger:  logger,
		db:      db,
	}
}

// Setup initializes the database transaction middleware.
func (m *DatabaseTrx) Setup() {
	m.logger.Info("Setting up database transaction middleware")

	m.handler.Use(func(c *gin.Context) {
		db, ok := m.db.(common.Database)
		if !ok || db == nil {
			m.logger.Error("Failed to get common.Database for transaction")
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		txHandle, err := db.Begin()
		if err != nil {
			m.logger.Error("Failed to begin database transaction", zap.Error(err))
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		m.logger.Info("Starting database transaction")

		defer func() {
			if r := recover(); r != nil {
				err := txHandle.Rollback()
				if err != nil {
					m.logger.Error("Failed to rollback database transaction after panic", zap.Error(err))
				} else {
					m.logger.Info("Rolled back database transaction after panic")
				}
			}
		}()

		c.Set(constants.DBTransaction, txHandle)
		c.Next()

		// Commit transaction on success status
		if statusInList(c.Writer.Status(), []int{http.StatusOK, http.StatusCreated, http.StatusNoContent}) {
			m.logger.Info("Committing database transaction")
			if err := txHandle.Commit(); err != nil {
				m.logger.Error("Failed to commit database transaction", zap.Error(err))
			}
		} else {
			m.logger.Info("rolling back database transaction")
			err := txHandle.Rollback()
			if err != nil {
				m.logger.Error("Failed to rollback database transaction", zap.Error(err))
			} else {
				m.logger.Info("Rolled back database transaction")
			}
		}
	})
}

// statusInList checks if context writer status is in the provided list of statuses.
func statusInList(status int, statusList []int) bool {
	for _, s := range statusList {
		if status == s {
			return true
		}
	}
	return false
}
