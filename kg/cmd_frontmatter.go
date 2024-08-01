package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

func newFrontmatterCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "frontmatter",
		Short: "Bulk update and normalize frontmatter",
		Long:  `Scan all markdown files, normalize frontmatter, and apply bulk updates.`,
	}

	cmd.AddCommand(
		newFrontmatterNormalizeCmd(),
		newFrontmatterUpdateCmd(),
	)

	return cmd
}

func newFrontmatterNormalizeCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "normalize",
		Short: "Normalize frontmatter across all notes",
		RunE: func(cmd *cobra.Command, args []string) error {
			return normalizeFrontmatter()
		},
	}
}

func newFrontmatterUpdateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update [node] [key] [value]",
		Short: "Bulk update a frontmatter field across all notes",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			return updateFrontmatter(args[1], args[2], args[3])
		},
	}

	return cmd
}

func normalizeFrontmatter() error {
	notesDir := viper.GetString("notes_directory")
	if notesDir == "" {
		return fmt.Errorf("notes directory not set in config")
	}

	return filepath.Walk(notesDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(info.Name(), ".md") {
			if err := normalizeFile(path); err != nil {
				return fmt.Errorf("failed to normalize %s: %w", path, err)
			}
		}

		return nil
	})
}

func normalizeFile(path string) error {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	parts := strings.SplitN(string(content), "---", 3)
	if len(parts) != 3 {
		return fmt.Errorf("invalid frontmatter format")
	}

	var frontmatter map[string]interface{}
	if err := yaml.Unmarshal([]byte(parts[1]), &frontmatter); err != nil {
		return fmt.Errorf("failed to parse frontmatter: %w", err)
	}

	// Apply normalization rules
	frontmatter = normalizeFields(frontmatter)

	// Marshal the normalized frontmatter
	normalizedYaml, err := yaml.Marshal(frontmatter)
	if err != nil {
		return fmt.Errorf("failed to marshal normalized frontmatter: %w", err)
	}

	// Reconstruct the file content
	normalizedContent := fmt.Sprintf("---\n%s---\n%s", string(normalizedYaml), parts[2])

	// Write the normalized content back to the file
	if err := ioutil.WriteFile(path, []byte(normalizedContent), 0644); err != nil {
		return fmt.Errorf("failed to write normalized content: %w", err)
	}

	fmt.Printf("Normalized frontmatter in %s\n", path)
	return nil
}

func normalizeFields(frontmatter map[string]interface{}) map[string]interface{} {
	// Normalize date formats
	if date, ok := frontmatter["date"].(string); ok {
		if parsedDate, err := time.Parse("2006-01-02", date); err == nil {
			frontmatter["date"] = parsedDate.Format("2006-01-02")
		}
	}

	// Normalize tag capitalization
	if tags, ok := frontmatter["tags"].([]interface{}); ok {
		normalizedTags := make([]string, len(tags))
		for i, tag := range tags {
			if strTag, ok := tag.(string); ok {
				normalizedTags[i] = strings.ToLower(strTag)
			}
		}
		frontmatter["tags"] = normalizedTags
	}

	// Add more normalization rules as needed

	return frontmatter
}

func updateFileField(path, key, value string) error {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	parts := strings.SplitN(string(content), "---", 3)
	if len(parts) != 3 {
		return fmt.Errorf("invalid frontmatter format")
	}

	var frontmatter map[string]interface{}
	if err := yaml.Unmarshal([]byte(parts[1]), &frontmatter); err != nil {
		return fmt.Errorf("failed to parse frontmatter: %w", err)
	}

	// Update the specified field
	frontmatter[key] = value

	// Marshal the updated frontmatter
	updatedYaml, err := yaml.Marshal(frontmatter)
	if err != nil {
		return fmt.Errorf("failed to marshal updated frontmatter: %w", err)
	}

	// Reconstruct the file content
	updatedContent := fmt.Sprintf("---\n%s---\n%s", string(updatedYaml), parts[2])

	// Write the updated content back to the file
	if err := ioutil.WriteFile(path, []byte(updatedContent), 0644); err != nil {
		return fmt.Errorf("failed to write updated content: %w", err)
	}

	fmt.Printf("Updated %s in %s\n", key, path)
	return nil
}
