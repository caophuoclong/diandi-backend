package commands

import (
	"github.com/caophuoclong/codifech/lib"
	"github.com/spf13/cobra"
)

type MigrationDown struct{}

func (m *MigrationDown) Run() lib.CommandRunner {
	return func() {}
}

func (m *MigrationDown) Short() string {
	return "Run migrations down"
}

func (m *MigrationDown) Setup(cmd *cobra.Command) {
}

func NewMigrationDown() *MigrationDown {
	return &MigrationDown{}
}
