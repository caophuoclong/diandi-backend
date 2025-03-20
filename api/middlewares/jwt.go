package middlewares

import (
	"diandi-backend/lib"
	"diandi-backend/services"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type JWTMiddleware struct {
	logger  lib.Logger
	service services.AuthService
}

func NewAuthMiddleware(
	logger lib.Logger,
	service services.AuthService,
) JWTMiddleware {
	return JWTMiddleware{
		logger:  logger,
		service: service,
	}
}

func (m JWTMiddleware) SetUp() {}

func (m JWTMiddleware) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		t := strings.Split(authHeader, " ")
		if len(t) == 2 {
			authToken := t[1]
			authorized, err := m.service.Authorize(authToken)
			if authorized {
				c.Next()
				return
			}
			c.JSON(
				http.StatusInternalServerError, gin.H{
					"error": err.Error(),
				})
		}
		c.JSON(
			http.StatusUnauthorized,
			gin.H{
				"error": "You are not authorized",
			},
		)
		c.Abort()
	}
}
