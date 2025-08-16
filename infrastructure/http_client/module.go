package http_client

import (
	"github.com/h4rdc0m/aurora-api/domain/common"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"infrastructure.http_client",
	fx.Provide(func() common.HTTPClient { return NewStandardHttpClient(3) }),
)
