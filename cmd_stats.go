package main

import (
	"fmt"
	"sort"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
	notesDir := viper.GetString("notes_directory")
	if notesDir == "" {
		return fmt.Errorf("notes directory not set in config")
	}

	notes, err := loadNotes(notesDir)
	if err != nil {
		return fmt.Errorf("failed to load notes: %w", err)
	}

	fmt.Printf("Total number of notes: %d\n\n", len(notes))

	displayTagDistribution(notes)
	displayMostConnectedNotes(notes)
	displayDateBasedStats(notes)
	displayAverageNoteLength(notes)

	return nil
}

func displayTagDistribution(notes []Note) {
	tagCount := make(map[string]int)
	for _, note := range notes {
		for _, tag := range note.Tags {
			tagCount[tag]++
		}
	}

	fmt.Println("Tag distribution:")
	for tag, count := range tagCount {
		fmt.Printf("  %s: %d\n", tag, count)
	}
	fmt.Println()
}

func displayMostConnectedNotes(notes []Note) {
	sort.Slice(notes, func(i, j int) bool {
		return len(notes[i].Connections) > len(notes[j].Connections)
	})

	fmt.Println("Most connected notes:")
	for i := 0; i < 5 && i < len(notes); i++ {
		fmt.Printf("  %s: %d connections\n", notes[i].Title, len(notes[i].Connections))
	}
	fmt.Println()
}

func displayDateBasedStats(notes []Note) {
	notesByMonth := make(map[string]int)
	for _, note := range notes {
		monthKey := note.Date.Format("2006-01")
		notesByMonth[monthKey]++
	}

	fmt.Println("Notes per month:")
	for month, count := range notesByMonth {
		fmt.Printf("  %s: %d\n", month, count)
	}
	fmt.Println()
}

func displayAverageNoteLength(notes []Note) {
	var totalLength int
	for _, note := range notes {
		totalLength += len(note.Content)
	}
	averageLength := float64(totalLength) / float64(len(notes))

	fmt.Printf("Average note length: %.2f characters\n", averageLength)
}
