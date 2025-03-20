package commands

import (
	"github.com/caophuoclong/codifech/lib"
	"github.com/spf13/cobra"
)

type MigrationUp struct{}

func (mu *MigrationUp) Short() string {
	return "Run migrations up"
}

func (mu *MigrationUp) Setup(cmd *cobra.Command) {}
func (mu *MigrationUp) Run() lib.CommandRunner {
	return func() {}
}

func NewMigrationUp() *MigrationUp {
	return &MigrationUp{}
}
