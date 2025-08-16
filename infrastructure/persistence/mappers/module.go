package mappers

import "go.uber.org/fx"

var Module = fx.Options(
	fx.Provide(NewTenantMapper),
	fx.Provide(NewUserMapper),
	fx.Provide(NewSiteMapper),
	fx.Provide(NewTemplateMapper),
	fx.Provide(NewPageMapper),
	fx.Provide(NewPageVersionMapper),
	fx.Provide(NewPageBlockMapper),
)
