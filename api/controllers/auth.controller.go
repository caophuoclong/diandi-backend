package controllers

import (
	"diandi-backend/domains"
	"diandi-backend/lib"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	logger  lib.Logger
	service domains.AuthService
}

func NewAuthController(
	logger lib.Logger,
	service domains.AuthService,
) AuthController {
	return AuthController{
		logger:  logger,
		service: service,
	}
}

func (ac AuthController) SignIn(c *gin.Context) {
	ac.logger.Info(c.ClientIP())
	c.JSON(
		http.StatusAccepted,
		gin.H{
			"accessToken": "aaaa",
		},
	)
}
