package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newEditCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "edit [title]",
		Short: "Modify existing notes, preserving frontmatter",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return editNote(args[0])
		},
	}
}

func editNote(title string) error {
	// TODO: Implement note editing functionality
	fmt.Printf("Editing note: %s\n", title)
	return nil
}
