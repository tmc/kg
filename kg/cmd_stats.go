package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
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

func loadNotes(notesDir string) ([]Note, error) {
	var notes []Note

	err := filepath.Walk(notesDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(info.Name(), ".md") {
			note, err := parseNote(path)
			if err != nil {
				return fmt.Errorf("failed to parse note %s: %w", path, err)
			}
			notes = append(notes, note)
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to walk notes directory: %w", err)
	}

	return notes, nil
}

func parseNote(path string) (Note, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return Note{}, fmt.Errorf("failed to read file: %w", err)
	}

	parts := strings.SplitN(string(content), "---", 3)
	if len(parts) != 3 {
		return Note{}, fmt.Errorf("invalid note format")
	}

	var frontmatter map[string]interface{}
	if err := yaml.Unmarshal([]byte(parts[1]), &frontmatter); err != nil {
		return Note{}, fmt.Errorf("failed to parse frontmatter: %w", err)
	}

	note := Note{
		Title:       frontmatter["title"].(string),
		Filename:    filepath.Base(path),
		Frontmatter: frontmatter,
		Content:     strings.TrimSpace(parts[2]),
		Connections: []string{},
	}

	if tags, ok := frontmatter["tags"].([]interface{}); ok {
		for _, tag := range tags {
			note.Tags = append(note.Tags, tag.(string))
		}
	}

	if connections, ok := frontmatter["connected_to"].([]interface{}); ok {
		for _, conn := range connections {
			note.Connections = append(note.Connections, conn.(string))
		}
	}

	if date, ok := frontmatter["date"].(string); ok {
		note.Date, _ = time.Parse("2006-01-02", date)
	}

	return note, nil
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