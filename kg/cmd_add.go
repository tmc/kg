package main

import (
	"bufio"
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

func newAddCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "add [title]",
		Short: "Create new notes with AI-suggested tags",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return addNote(args[0])
		},
	}
}

func addNote(title string) error {
	// Generate filename
	filename := generateFilename(title)

	// Get AI-suggested tags
	suggestedTags, err := getSuggestedTags(title)
	if err != nil {
		return fmt.Errorf("failed to get AI-suggested tags: %w", err)
	}

	// Allow user to edit/confirm suggested tags
	confirmedTags := confirmTags(suggestedTags)

	// Create frontmatter
	frontmatter := generateFrontmatter(title, confirmedTags)

	// Create the file
	notesDir := viper.GetString("notes_directory")
	filePath := filepath.Join(notesDir, filename)
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	// Write frontmatter and initial content
	_, err = file.WriteString(frontmatter + "\n\n# " + title + "\n")
	if err != nil {
		return fmt.Errorf("failed to write to file: %w", err)
	}

	fmt.Printf("Note created: %s\n", filePath)
	return nil
}

func generateFilename(title string) string {
	// Convert title to kebab-case
	kebabTitle := strings.ToLower(strings.ReplaceAll(title, " ", "-"))
	return kebabTitle + ".md"
}

func getSuggestedTags(title string) ([]string, error) {
	llm, err := openai.New()
	if err != nil {
		return nil, fmt.Errorf("failed to create OpenAI client: %w", err)
	}

	prompt := fmt.Sprintf("Suggest 3-5 relevant tags for a note titled '%s'. Respond with only the tags, separated by commas.", title)
	completion, err := llm.Call(prompt, llms.WithMaxTokens(50))
	if err != nil {
		return nil, fmt.Errorf("failed to get AI response: %w", err)
	}

	tags := strings.Split(strings.TrimSpace(completion), ",")
	for i, tag := range tags {
		tags[i] = strings.TrimSpace(tag)
	}

	return tags, nil
}

func confirmTags(suggestedTags []string) []string {
	fmt.Println("Suggested tags:")
	for i, tag := range suggestedTags {
		fmt.Printf("%d. %s\n", i+1, tag)
	}

	fmt.Println("Enter the numbers of the tags you want to keep, separated by spaces.")
	fmt.Println("To add a new tag, type '+' followed by the tag.")
	fmt.Print("Your selection: ")

	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	confirmedTags := []string{}
	for _, selection := range strings.Fields(input) {
		if strings.HasPrefix(selection, "+") {
			confirmedTags = append(confirmedTags, strings.TrimPrefix(selection, "+"))
		} else if index, err := strconv.Atoi(selection); err == nil && index > 0 && index <= len(suggestedTags) {
			confirmedTags = append(confirmedTags, suggestedTags[index-1])
		}
	}

	return confirmedTags
}

func generateFrontmatter(title string, tags []string) string {
	frontmatter := map[string]interface{}{
		"title":     title,
		"tags":      tags,
		"date":      time.Now().Format("2006-01-02"),
		"lastmod":   time.Now().Format("2006-01-02"),
		"draft":     false,
	}

	yamlData, err := yaml.Marshal(frontmatter)
	if err != nil {
		// If marshaling fails, return a basic frontmatter
		return fmt.Sprintf("---\ntitle: %s\ntags: %v\ndate: %s\n---\n", title, tags, time.Now().Format("2006-01-02"))
	}

	return fmt.Sprintf("---\n%s---\n", string(yamlData))
}