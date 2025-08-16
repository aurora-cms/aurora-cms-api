package commands

import (
	"github.com/gin-gonic/gin"
	"github.com/h4rdc0m/aurora-api/api/http/middlewares"
	"github.com/h4rdc0m/aurora-api/api/http/routes"
	"github.com/h4rdc0m/aurora-api/domain/common"
	"github.com/h4rdc0m/aurora-api/infrastructure/config"
	"github.com/spf13/cobra"
)

// ServeCommand test command
type ServeCommand struct{}

func (s *ServeCommand) Short() string {
	return "Serve the Aurora API"
}

func (s *ServeCommand) Setup(_ *cobra.Command) {}

func (s *ServeCommand) Run() common.CommandRunner {
	return func(
		middleware middlewares.Middlewares,
		env *config.Env,
		router *gin.Engine,
		route routes.Routes,
		logger common.Logger,
		database common.Database,
	) {
		middleware.Setup()
		route.Setup()
		logger.Info("Starting Aurora API server")
		if env.ServerPort == "" {
			_ = router.Run()
		} else {
			_ = router.Run(":" + env.ServerPort)
		}
	}
}

// NewServeCommand creates a new instance of ServeCommand.
func NewServeCommand() *ServeCommand {
	return &ServeCommand{}
}
