package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

func newImportCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "import [file]",
		Short: "Add external markdown files, parsing frontmatter",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return importFile(args[0])
		},
	}
}

func importFile(filePath string) error {
	// Read the file content
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	// Parse the markdown file and extract frontmatter
	frontmatter, body, err := parseFrontmatter(string(content))
	if err != nil {
		return fmt.Errorf("failed to parse frontmatter: %w", err)
	}

	// Validate frontmatter
	if err := validateFrontmatter(frontmatter); err != nil {
		return fmt.Errorf("invalid frontmatter: %w", err)
	}

	// Generate filename from title
	title, ok := frontmatter["title"].(string)
	if !ok {
		return fmt.Errorf("title not found in frontmatter")
	}
	filename := generateFilename(title)

	// Check for conflicts with existing notes
	notesDir := viper.GetString("notes_directory")
	if notesDir == "" {
		return fmt.Errorf("notes directory not set in config")
	}
	destPath := filepath.Join(notesDir, filename)

	if _, err := os.Stat(destPath); !os.IsNotExist(err) {
		return fmt.Errorf("a note with the title '%s' already exists", title)
	}

	// Copy the file to the appropriate location
	if err := copyFile(filePath, destPath); err != nil {
		return fmt.Errorf("failed to copy file: %w", err)
	}

	// Update internal links
	updatedBody, err := updateInternalLinks(body, notesDir)
	if err != nil {
		return fmt.Errorf("failed to update internal links: %w", err)
	}

	// Write the updated content back to the file
	updatedContent := fmt.Sprintf("---\n%s\n---\n%s", marshalFrontmatter(frontmatter), updatedBody)
	if err := ioutil.WriteFile(destPath, []byte(updatedContent), 0644); err != nil {
		return fmt.Errorf("failed to write updated content: %w", err)
	}

	fmt.Printf("Successfully imported '%s' to %s\n", title, destPath)
	return nil
}

func parseFrontmatter(content string) (map[string]interface{}, string, error) {
	parts := strings.SplitN(content, "---", 3)
	if len(parts) != 3 {
		return nil, "", fmt.Errorf("invalid frontmatter format")
	}

	var frontmatter map[string]interface{}
	if err := yaml.Unmarshal([]byte(parts[1]), &frontmatter); err != nil {
		return nil, "", fmt.Errorf("failed to parse frontmatter: %w", err)
	}

	return frontmatter, parts[2], nil
}

func validateFrontmatter(frontmatter map[string]interface{}) error {
	requiredFields := []string{"title", "date"}
	for _, field := range requiredFields {
		if _, ok := frontmatter[field]; !ok {
			return fmt.Errorf("missing required field: %s", field)
		}
	}
	return nil
}

func copyFile(src, dst string) error {
	input, err := ioutil.ReadFile(src)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(dst, input, 0644)
}

func updateInternalLinks(content, notesDir string) (string, error) {
	// A more robust implementation would use a proper Markdown parser
	return strings.ReplaceAll(content, "[[", "["+notesDir+"/"), fmt.Errorf("not implemented")
}

func marshalFrontmatter(frontmatter map[string]interface{}) string {
	yamlData, err := yaml.Marshal(frontmatter)
	if err != nil {
		// If marshaling fails, return an empty string
		return ""
	}
	return string(yamlData)
}
