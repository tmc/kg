package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
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
	// Find the note file
	notesDir := viper.GetString("notes_directory")
	if notesDir == "" {
		return fmt.Errorf("notes directory not set in config")
	}

	filename := generateFilename(title)
	filePath := filepath.Join(notesDir, filename)

	// Check if the file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fmt.Errorf("note '%s' does not exist", title)
	}

	// Load existing note content
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read note: %w", err)
	}

	// Split content into frontmatter and body
	parts := strings.SplitN(string(content), "---", 3)
	if len(parts) != 3 {
		return fmt.Errorf("invalid note format")
	}

	frontmatter := parts[1]
	body := parts[2]

	// Create a temporary file for editing
	tempFile, err := ioutil.TempFile("", "kg-edit-*.md")
	if err != nil {
		return fmt.Errorf("failed to create temporary file: %w", err)
	}
	defer os.Remove(tempFile.Name())

	// Write the content to the temporary file
	_, err = tempFile.WriteString(fmt.Sprintf("---\n%s---\n%s", frontmatter, body))
	if err != nil {
		return fmt.Errorf("failed to write to temporary file: %w", err)
	}
	tempFile.Close()

	// Open the note in the user's preferred editor
	editor := viper.GetString("editor")
	if editor == "" {
		editor = os.Getenv("EDITOR")
	}
	if editor == "" {
		editor = "nano" // Default to nano if no editor is specified
	}

	cmd := exec.Command(editor, tempFile.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to run editor: %w", err)
	}

	// Read the edited content
	editedContent, err := ioutil.ReadFile(tempFile.Name())
	if err != nil {
		return fmt.Errorf("failed to read edited content: %w", err)
	}

	// Parse and validate the updated frontmatter
	editedParts := strings.SplitN(string(editedContent), "---", 3)
	if len(editedParts) != 3 {
		return fmt.Errorf("invalid edited note format")
	}

	var updatedFrontmatter map[string]interface{}
	err = yaml.Unmarshal([]byte(editedParts[1]), &updatedFrontmatter)
	if err != nil {
		return fmt.Errorf("invalid frontmatter: %w", err)
	}

	// Validate required fields
	requiredFields := []string{"title", "date"}
	for _, field := range requiredFields {
		if _, ok := updatedFrontmatter[field]; !ok {
			return fmt.Errorf("missing required frontmatter field: %s", field)
		}
	}

	// Save the changes
	err = ioutil.WriteFile(filePath, editedContent, 0644)
	if err != nil {
		return fmt.Errorf("failed to save changes: %w", err)
	}

	fmt.Printf("Note '%s' updated successfully\n", title)
	return nil
}
