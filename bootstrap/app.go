package bootstrap

import (
	"diandi-backend/commands"
	"diandi-backend/lib"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "codifech",
	Short: "A brief description of your application",
	Long: `
  ____ ___  ____ ___ _____ _____ ____ _   _ 
 / ___/ _ \|  _ \_ _|  ___| ____/ ___| | | |
| |  | | | | | | | || |_  |  _|| |   | |_| |
| |__| |_| | |_| | ||  _| | |__| |___|  _  |
 \____\___/|____/___|_|   |_____\____|_| |_|
This is a command runner or cli for api architecture in golang. 
Using this we can use underlying dependency injection container for running scripts. 
Main advantage is that, we can use same services, repositories, infrastructure present in the application itself`,
}

type App struct {
	*cobra.Command
}

func NewApp() App {
	logger := lib.GetLogger()
	cmd := App{
		Command: rootCmd,
	}
	cmd.AddCommand(
		commands.GetSubCommands(CommonModules)...,
	)
	defer func(logger *lib.Logger) {
		err := logger.Sync()
		if err != nil {

		}
	}(&logger)
	return cmd
}

var RootApp = NewApp()
