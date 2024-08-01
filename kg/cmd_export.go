package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func newExportCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "export [format]",
		Short: "Output to JSON/CSV, including frontmatter",
		Long:  `Export the knowledge graph to JSON or CSV format, including all frontmatter data.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			format := strings.ToLower(args[0])
			if format != "json" && format != "csv" {
				return fmt.Errorf("unsupported format: %s. Use 'json' or 'csv'", format)
			}
			return exportGraph(format)
		},
	}
}

type Note struct {
	Title       string                 `json:"title"`
	Filename    string                 `json:"filename"`
	Frontmatter map[string]interface{} `json:"frontmatter"`
	Content     string                 `json:"content"`
	Connections []string               `json:"connections"`

	Tags []string  `json:"tags"`
	Date time.Time `json:"date"`
}

func exportGraph(format string) error {
	notesDir := viper.GetString("notes_directory")
	if notesDir == "" {
		return fmt.Errorf("notes directory not set in config")
	}

	notes, err := loadNotes(notesDir)
	if err != nil {
		return fmt.Errorf("failed to load notes: %w", err)
	}

	if format == "json" {
		return exportJSON(notes)
	} else {
		return exportCSV(notes)
	}
}

func exportJSON(notes []Note) error {
	jsonData, err := json.MarshalIndent(notes, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal notes to JSON: %w", err)
	}

	if err := os.WriteFile("knowledge_graph_export.json", jsonData, 0644); err != nil {
		return fmt.Errorf("failed to write JSON file: %w", err)
	}

	fmt.Println("Exported knowledge graph to knowledge_graph_export.json")
	return nil
}

func exportCSV(notes []Note) error {
	nodesFile, err := os.Create("knowledge_graph_nodes.csv")
	if err != nil {
		return fmt.Errorf("failed to create nodes CSV file: %w", err)
	}
	defer nodesFile.Close()

	edgesFile, err := os.Create("knowledge_graph_edges.csv")
	if err != nil {
		return fmt.Errorf("failed to create edges CSV file: %w", err)
	}
	defer edgesFile.Close()

	nodesWriter := csv.NewWriter(nodesFile)
	edgesWriter := csv.NewWriter(edgesFile)

	// Write headers
	nodesWriter.Write([]string{"ID", "Title", "Filename", "Tags", "Date", "LastMod"})
	edgesWriter.Write([]string{"Source", "Target"})

	for _, note := range notes {
		// Write node
		nodesWriter.Write([]string{
			note.Filename,
			note.Title,
			note.Filename,
			strings.Join(note.Frontmatter["tags"].([]string), "|"),
			note.Frontmatter["date"].(string),
			note.Frontmatter["lastmod"].(string),
		})

		// Write edges
		for _, connection := range note.Connections {
			edgesWriter.Write([]string{note.Filename, connection})
		}
	}

	nodesWriter.Flush()
	edgesWriter.Flush()

	if err := nodesWriter.Error(); err != nil {
		return fmt.Errorf("error writing nodes CSV: %w", err)
	}
	if err := edgesWriter.Error(); err != nil {
		return fmt.Errorf("error writing edges CSV: %w", err)
	}

	fmt.Println("Exported knowledge graph to knowledge_graph_nodes.csv and knowledge_graph_edges.csv")
	return nil
}
