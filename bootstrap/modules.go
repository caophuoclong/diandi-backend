package bootstrap

import (
	"diandi-backend/api/controllers"
	"diandi-backend/api/routes"
	"diandi-backend/lib"
	"diandi-backend/services"

	"go.uber.org/fx"
)

var CommonModules = fx.Options(
	lib.Module,
	routes.Module,
	services.Module,
	controllers.Module,
)
