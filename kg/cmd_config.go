package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newConfigCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "config [key] [value]",
		Short: "Manage .kgrc settings",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return updateConfig(args[0], args[1])
		},
	}
}

func updateConfig(key, value string) error {
	// TODO: Implement configuration management
	fmt.Printf("Updating config: %s = %s\n", key, value)
	return nil
}
