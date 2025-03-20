package routes

import (
	"diandi-backend/lib"

	"github.com/gin-gonic/gin"
)

type ApiVersion struct {
	version string
}

func (api ApiVersion) getVersioning() string {
	return "/api/" + api.version
}

func NewApiV1(handler lib.RequestHandler) *gin.RouterGroup {
	v1 := ApiVersion{version: "v1"}
	return handler.Gin.Group(v1.getVersioning())
}
