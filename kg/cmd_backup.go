package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newBackupCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "backup",
		Short: "Create compressed, timestamped backups",
		RunE: func(cmd *cobra.Command, args []string) error {
			return createBackup()
		},
	}
}

func createBackup() error {
	// TODO: Implement backup creation
	fmt.Println("Creating backup")
	return nil
}
