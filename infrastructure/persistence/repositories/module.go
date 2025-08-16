package repositories

import "go.uber.org/fx"

var Module = fx.Options(
	fx.Provide(NewTenantRepository),
	fx.Provide(NewSiteRepository),
	fx.Provide(NewTemplateRepository),
	fx.Provide(NewUserRepository),
	fx.Provide(NewPageRepository),
	fx.Provide(NewPageVersionRepository),
	fx.Provide(NewPageBlockRepository),
)
