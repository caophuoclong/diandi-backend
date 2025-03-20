package commands

import (
	"context"
	"diandi-backend/lib"

	"github.com/spf13/cobra"
	"go.uber.org/fx"
)

var cmds = map[string]lib.Command{
	"app:serve":      NewServeCommand(),
	"migration:up":   NewMigrationUp(),
	"migration:down": NewMigrationDown(),
}

func GetSubCommands(opt fx.Option) []*cobra.Command {
	var subCommands []*cobra.Command

	for name, cmd := range cmds {
		subCommands = append(subCommands, WrapSubCommand(name, cmd, opt))
	}

	return subCommands
}

func WrapSubCommand(name string, cmd lib.Command, opt fx.Option) *cobra.Command {
	wrappedCmd := &cobra.Command{
		Use:   name,
		Short: cmd.Short(),
		Run: func(c *cobra.Command, args []string) {
			logger := lib.GetLogger()
			opts := fx.Options(
				fx.Invoke(cmd.Run()))
			ctx := context.Background()
			app := fx.New(opt, opts)
			err := app.Start(ctx)
			defer func(app *fx.App, ctx context.Context) {
				err := app.Stop(ctx)
				if err != nil {
					logger.Fatal(err)
				}
			}(app, ctx)
			if err != nil {
				logger.Fatal(err)
			}

		},
	}
	cmd.Setup(wrappedCmd)
	return wrappedCmd
}
