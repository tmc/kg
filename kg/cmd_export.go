package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newExportCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "export [format]",
		Short: "Output to JSON/CSV, including frontmatter",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return exportGraph(args[0])
		},
	}
}

func exportGraph(format string) error {
	// TODO: Implement graph export functionality
	fmt.Printf("Exporting graph to %s format\n", format)
	return nil
}
