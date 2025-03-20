package lib

import "go.uber.org/fx"

var Module = fx.Options(
	fx.Provide(GetLogger),
	fx.Provide(NewEnv),
	fx.Provide(NewRequestHandler),

)
