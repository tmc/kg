package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
	"gopkg.in/yaml.v3"
)

func newConnectCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "connect [concept1] [concept2]",
		Short: "Link two concepts with AI-generated content",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return connectConcepts(args[0], args[1])
		},
	}
}

func connectConcepts(concept1, concept2 string) error {
	// Verify both concepts exist as notes
	notesDir := viper.GetString("notes_directory")
	if notesDir == "" {
		return fmt.Errorf("notes directory not set in config")
	}

	concept1Path := filepath.Join(notesDir, generateFilename(concept1))
	concept2Path := filepath.Join(notesDir, generateFilename(concept2))

	if !fileExists(concept1Path) {
		return fmt.Errorf("concept '%s' does not exist as a note", concept1)
	}
	if !fileExists(concept2Path) {
		return fmt.Errorf("concept '%s' does not exist as a note", concept2)
	}

	// Generate content linking the two concepts
	content, err := generateLinkingContent(concept1, concept2)
	if err != nil {
		return fmt.Errorf("failed to generate linking content: %w", err)
	}

	// Create a new note with the generated content
	newNoteTitle := fmt.Sprintf("%s-%s-connection", concept1, concept2)
	newNotePath := filepath.Join(notesDir, generateFilename(newNoteTitle))
	if err := createNewNote(newNotePath, newNoteTitle, content, []string{concept1, concept2}); err != nil {
		return fmt.Errorf("failed to create new note: %w", err)
	}

	// Update frontmatter of involved notes
	if err := updateFrontmatter(concept1Path, "connected_to", concept2); err != nil {
		return fmt.Errorf("failed to update frontmatter of %s: %w", concept1, err)
	}
	if err := updateFrontmatter(concept2Path, "connected_to", concept1); err != nil {
		return fmt.Errorf("failed to update frontmatter of %s: %w", concept2, err)
	}

	fmt.Printf("Created new note connecting %s and %s: %s\n", concept1, concept2, newNotePath)
	return nil
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func generateLinkingContent(concept1, concept2 string) (string, error) {
	llm, err := openai.New()
	if err != nil {
		return "", fmt.Errorf("failed to create OpenAI client: %w", err)
	}

	prompt := fmt.Sprintf("Generate a short paragraph (3-5 sentences) explaining the connection between %s and %s.", concept1, concept2)
	res, err := llm.GenerateContent(context.TODO(), []llms.MessageContent{
		llms.TextParts(llms.ChatMessageTypeHuman, prompt),
	})
	if err != nil {
		return "", fmt.Errorf("failed to generate content: %w", err)
	}

	return res.Choices[0].Content, err
}

func createNewNote(path, title, content string, tags []string) error {
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	frontmatter := map[string]interface{}{
		"title":    title,
		"tags":     tags,
		"date":     time.Now().Format("2006-01-02"),
		"lastmod":  time.Now().Format("2006-01-02"),
		"draft":    false,
		"connects": tags,
	}

	yamlData, err := yaml.Marshal(frontmatter)
	if err != nil {
		return fmt.Errorf("failed to marshal frontmatter: %w", err)
	}

	_, err = fmt.Fprintf(file, "---\n%s---\n\n# %s\n\n%s\n", string(yamlData), title, content)
	if err != nil {
		return fmt.Errorf("failed to write to file: %w", err)
	}

	return nil
}

func updateFrontmatter(path, key, value string) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	parts := strings.SplitN(string(content), "---", 3)
	if len(parts) != 3 {
		return fmt.Errorf("invalid frontmatter format")
	}

	var frontmatter map[string]interface{}
	if err := yaml.Unmarshal([]byte(parts[1]), &frontmatter); err != nil {
		return fmt.Errorf("failed to unmarshal frontmatter: %w", err)
	}

	if existing, ok := frontmatter[key].([]string); ok {
		frontmatter[key] = append(existing, value)
	} else {
		frontmatter[key] = []string{value}
	}

	updatedYaml, err := yaml.Marshal(frontmatter)
	if err != nil {
		return fmt.Errorf("failed to marshal updated frontmatter: %w", err)
	}

	updatedContent := fmt.Sprintf("---\n%s---\n%s", string(updatedYaml), parts[2])
	if err := os.WriteFile(path, []byte(updatedContent), 0644); err != nil {
		return fmt.Errorf("failed to write updated content: %w", err)
	}

	return nil
}
