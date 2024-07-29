package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newFrontmatterCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "frontmatter",
		Short: "Bulk update/normalize frontmatter across files",
		RunE: func(cmd *cobra.Command, args []string) error {
			return updateFrontmatter()
		},
	}
}

func updateFrontmatter() error {
	// TODO: Implement frontmatter bulk update/normalization
	fmt.Println("Updating frontmatter across files")
	return nil
}
