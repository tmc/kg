package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newStatsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "stats",
		Short: "Display knowledge graph metrics",
		RunE: func(cmd *cobra.Command, args []string) error {
			return displayStats()
		},
	}
}

func displayStats() error {
	// TODO: Implement statistics display
	fmt.Println("Displaying knowledge graph statistics")
	return nil
}
