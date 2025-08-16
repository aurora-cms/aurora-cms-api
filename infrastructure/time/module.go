package time_provider

import "go.uber.org/fx"

var Module = fx.Module(
	"infrastructure.time_provider",
	fx.Provide(NewStandardTimeProvider),
)
