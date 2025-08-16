package http

import (
	"github.com/gin-gonic/gin"
	"github.com/h4rdc0m/aurora-api/domain/common"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"infrastructure.http",
	fx.Provide(NewRouter),
	// For convenience let other layers request *gin.Engine directly.
	fx.Provide(
		func(r common.Router) *gin.Engine { return r.(*Router).Engine },
	),
)
