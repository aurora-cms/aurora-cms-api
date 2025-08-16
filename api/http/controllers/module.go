package controllers

import "go.uber.org/fx"

var Module = fx.Options(
	fx.Provide(NewHealthController),
	fx.Provide(NewAuthController),
	fx.Provide(NewTenantController),
)
