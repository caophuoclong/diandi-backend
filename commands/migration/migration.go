package migration

import (
	"diandi-backend/lib"

	"github.com/spf13/cobra"
)

type MigrationCommand struct{}

func (m *MigrationCommand) Short() string {
	return "Run migrations"
}

func (m *MigrationCommand) Setup(cmd *cobra.Command) {

}

func (m *MigrationCommand) Run() lib.CommandRunner {
	logger := lib.GetLogger()
	logger.Info("Running migrations")
	return func() {

	}
}

func NewMigrationCommand() *MigrationCommand {
	return &MigrationCommand{}
}
