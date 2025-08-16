package routes

import (
	"github.com/h4rdc0m/aurora-api/api/http/controllers"
	"github.com/h4rdc0m/aurora-api/domain/common"
)

type TenantRoutes struct {
	logger           common.Logger
	handler          common.Router
	tenantController *controllers.TenantController
}
