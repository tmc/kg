package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

func newListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Show all notes with frontmatter info",
		RunE:  listNotes,
	}

	cmd.Flags().StringP("sort", "s", "title", "Sort by field (title, date, lastmod)")
	cmd.Flags().BoolP("reverse", "r", false, "Reverse sort order")
	cmd.Flags().StringP("filter", "f", "", "Filter by tag")

	return cmd
}

type NoteInfo struct {
	Title    string
	Filename string
	Tags     []string
	Date     time.Time
	LastMod  time.Time
}

func listNotes(cmd *cobra.Command, args []string) error {
	notesDir := viper.GetString("notes_directory")
	if notesDir == "" {
		return fmt.Errorf("notes directory not set in config")
	}

	notes, err := scanNotes(notesDir)
	if err != nil {
		return fmt.Errorf("failed to scan notes: %w", err)
	}

	sortField, _ := cmd.Flags().GetString("sort")
	reverse, _ := cmd.Flags().GetBool("reverse")
	filterTag, _ := cmd.Flags().GetString("filter")

	notes = sortNotes(notes, sortField, reverse)
	notes = filterNotes(notes, filterTag)

	displayNotes(notes)

	return nil
}

func scanNotes(notesDir string) ([]NoteInfo, error) {
	var notes []NoteInfo

	err := filepath.Walk(notesDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(info.Name(), ".md") {
			note, err := parseNoteInfo(path)
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

func parseNoteInfo(path string) (NoteInfo, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return NoteInfo{}, fmt.Errorf("failed to read file: %w", err)
	}

	parts := strings.SplitN(string(content), "---", 3)
	if len(parts) != 3 {
		return NoteInfo{}, fmt.Errorf("invalid note format")
	}

	var frontmatter map[string]interface{}
	if err := yaml.Unmarshal([]byte(parts[1]), &frontmatter); err != nil {
		return NoteInfo{}, fmt.Errorf("failed to parse frontmatter: %w", err)
	}

	note := NoteInfo{
		Title:    frontmatter["title"].(string),
		Filename: filepath.Base(path),
		Tags:     []string{},
	}

	if tags, ok := frontmatter["tags"].([]interface{}); ok {
		for _, tag := range tags {
			note.Tags = append(note.Tags, tag.(string))
		}
	}

	if date, ok := frontmatter["date"].(string); ok {
		note.Date, _ = time.Parse("2006-01-02", date)
	}

	if lastmod, ok := frontmatter["lastmod"].(string); ok {
		note.LastMod, _ = time.Parse("2006-01-02", lastmod)
	}

	return note, nil
}

func sortNotes(notes []NoteInfo, field string, reverse bool) []NoteInfo {
	sort.Slice(notes, func(i, j int) bool {
		var less bool
		switch field {
		case "date":
			less = notes[i].Date.Before(notes[j].Date)
		case "lastmod":
			less = notes[i].LastMod.Before(notes[j].LastMod)
		default:
			less = notes[i].Title < notes[j].Title
		}
		if reverse {
			return !less
		}
		return less
	})
	return notes
}

func filterNotes(notes []NoteInfo, tag string) []NoteInfo {
	if tag == "" {
		return notes
	}

	var filtered []NoteInfo
	for _, note := range notes {
		for _, noteTag := range note.Tags {
			if noteTag == tag {
				filtered = append(filtered, note)
				break
			}
		}
	}
	return filtered
}

func displayNotes(notes []NoteInfo) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "Title\tFilename\tTags\tDate\tLast Modified")
	fmt.Fprintln(w, "-----\t--------\t----\t----\t-------------")

	for _, note := range notes {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
			note.Title,
			note.Filename,
			strings.Join(note.Tags, ", "),
			note.Date.Format("2006-01-02"),
			note.LastMod.Format("2006-01-02"),
		)
	}

	w.Flush()
}
