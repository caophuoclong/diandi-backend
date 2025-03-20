package routes

import (
	"diandi-backend/api/controllers"
	"diandi-backend/lib"

	"github.com/gin-gonic/gin"
)

type AuthRoutes struct {
	logger     lib.Logger
	group      *gin.RouterGroup
	controller controllers.AuthController
}

func (s AuthRoutes) SetUp() {
	s.logger.Info("Setting up Auth Routes")
	auth := s.group.Group("/auth")
	{
		auth.POST("/login", s.controller.SignIn)
	}
}

func NewAuthRoutes(
	logger lib.Logger,
	group *gin.RouterGroup,
	controller controllers.AuthController,
) AuthRoutes {
	return AuthRoutes{
		logger:     logger,
		group:      group,
		controller: controller,
	}
}
