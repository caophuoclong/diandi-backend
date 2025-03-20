package routes

import (
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(NewApiV1),
	fx.Provide(NewRoutes),
	fx.Provide(NewAuthRoutes),
)

type Routes []Route

type Route interface {
	SetUp()
}

func NewRoutes(
	authRoutes AuthRoutes,
) Routes {
	return Routes{
		authRoutes,
	}
}

func (r Routes) SetUp() {
	for _, route := range r {
		route.SetUp()
	}
}
