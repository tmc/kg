package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "Show all notes with frontmatter info",
		RunE: func(cmd *cobra.Command, args []string) error {
			return listNotes()
		},
	}
}

func listNotes() error {
	// TODO: Implement listing of all notes with frontmatter info
	fmt.Println("Listing all notes")
	return nil
}
