package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
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
	Title      string                 `json:"title"`
	Filename   string                 `json:"filename"`
	Frontmatter map[string]interface{} `json:"frontmatter"`
	Content    string                 `json:"content"`
	Connections []string               `json:"connections"`
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

	if connections, ok := frontmatter["connected_to"].([]interface{}); ok {
		for _, conn := range connections {
			note.Connections = append(note.Connections, conn.(string))
		}
	}

	return note, nil
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