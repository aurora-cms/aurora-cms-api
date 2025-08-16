package http

import (
	"github.com/gin-gonic/gin"
	"github.com/h4rdc0m/aurora-api/domain/common"
)

// Router is a struct that holds the Gin engine instance for routing HTTP requests.
type Router struct {
	*gin.Engine
}

// NewRouter initializes a new Router with a Gin engine and sets the default logger.
func NewRouter(logger common.GinLogger) common.Router {
	gin.DefaultWriter = logger

	engine := gin.New()

	return &Router{Engine: engine}
}
