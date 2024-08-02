package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

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
