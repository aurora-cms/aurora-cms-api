package use_cases

import "go.uber.org/fx"

var Module = fx.Options(
	fx.Provide(NewAuthUseCase),
	fx.Provide(NewHealthUseCase),
	fx.Provide(NewSiteUseCase),
	fx.Provide(NewTenantUseCase),
)
