package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newAddCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "add [title]",
		Short: "Create new notes with AI-suggested tags",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return addNote(args[0])
		},
	}
}

func addNote(title string) error {
	// TODO: Implement note addition with AI-suggested tags
	fmt.Printf("Adding new note: %s\n", title)
	return nil
}
