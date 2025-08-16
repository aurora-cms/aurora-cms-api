package application

import (
	"github.com/h4rdc0m/aurora-api/application/use_cases"
	"go.uber.org/fx"
)

var Module = fx.Options(
	use_cases.Module,
)
