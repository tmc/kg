package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newSearchCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "search [query]",
		Short: "Find keywords in content and frontmatter",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return searchNotes(args[0])
		},
	}
}

func searchNotes(query string) error {
	// TODO: Implement search functionality
	fmt.Printf("Searching for: %s\n", query)
	return nil
}
