package commands

import (
	"diandi-backend/api/routes"
	"diandi-backend/lib"

	"github.com/spf13/cobra"
)

type ServeCommand struct{}

func (s *ServeCommand) Short() string {
	return "Serve application"
}

func (s *ServeCommand) Setup(cmd *cobra.Command) {}

func (s *ServeCommand) Run() lib.CommandRunner {
	return func(
		logger lib.Logger,
		env lib.Env,
		routes routes.Routes,
		router lib.RequestHandler,
	) {
		routes.SetUp()
		err := router.Gin.SetTrustedProxies([]string{"127.0.0.1"})
		if err != nil {
			return
		}
		if env.ServerPort == "" {
			_ = router.Gin.Run()
		} else {
			_ = router.Gin.Run(":" + env.ServerPort)
		}
	}
}

func NewServeCommand() lib.Command {
	return &ServeCommand{}
}
