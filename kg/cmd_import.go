package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newImportCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "import [file]",
		Short: "Add external markdown files, parsing frontmatter",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return importFile(args[0])
		},
	}
}

func importFile(file string) error {
	// TODO: Implement file import functionality
	fmt.Printf("Importing file: %s\n", file)
	return nil
}
