package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/search/query"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

var (
	index bleve.Index
	once  sync.Once
)

func newSearchCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "search [query]",
		Short: "Find keywords in content and frontmatter",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return searchNotes(args[0])
		},
	}

	cmd.Flags().BoolP("fuzzy", "f", false, "Enable fuzzy matching")
	cmd.Flags().IntP("context", "c", 50, "Number of characters to show as context")

	return cmd
}

func searchNotes(queryString string) error {
	// Initialize the index if it hasn't been created yet
	once.Do(func() {
		var err error
		index, err = createIndex()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to create search index: %v\n", err)
			os.Exit(1)
		}
	})

	// Parse the query
	q := parseQuery(queryString)

	// Perform the search
	searchRequest := bleve.NewSearchRequest(q)
	searchRequest.Fields = []string{"title", "tags", "content"}
	searchRequest.Highlight = bleve.NewHighlight()

	searchResults, err := index.Search(searchRequest)
	if err != nil {
		return fmt.Errorf("search failed: %w", err)
	}

	// Display results
	displayResults(searchResults)

	return nil
}

func createIndex() (bleve.Index, error) {
	notesDir := viper.GetString("notes_directory")
	if notesDir == "" {
		return nil, fmt.Errorf("notes directory not set in config")
	}

	indexPath := filepath.Join(notesDir, ".kg_search_index")

	// Open existing index or create a new one
	index, err := bleve.Open(indexPath)
	if err == bleve.ErrorIndexPathDoesNotExist {
		mapping := bleve.NewIndexMapping()
		index, err = bleve.New(indexPath, mapping)
		if err != nil {
			return nil, fmt.Errorf("failed to create index: %w", err)
		}
	} else if err != nil {
		return nil, fmt.Errorf("failed to open index: %w", err)
	}

	// Index all notes
	err = filepath.Walk(notesDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(info.Name(), ".md") {
			if err := indexNote(index, path); err != nil {
				return fmt.Errorf("failed to index note %s: %w", path, err)
			}
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to index notes: %w", err)
	}

	return index, nil
}

func indexNote(index bleve.Index, path string) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	parts := strings.SplitN(string(content), "---", 3)
	if len(parts) != 3 {
		return fmt.Errorf("invalid note format")
	}

	var frontmatter map[string]interface{}
	if err := yaml.Unmarshal([]byte(parts[1]), &frontmatter); err != nil {
		return fmt.Errorf("failed to parse frontmatter: %w", err)
	}

	doc := struct {
		Title   string   `json:"title"`
		Tags    []string `json:"tags"`
		Content string   `json:"content"`
	}{
		Title:   frontmatter["title"].(string),
		Tags:    frontmatter["tags"].([]string),
		Content: strings.TrimSpace(parts[2]),
	}

	return index.Index(path, doc)
}

func parseQuery(queryString string) query.Query {
	// Implement advanced query parsing (AND, OR, NOT)
	// This is a simple implementation; you may want to use a proper query parser
	terms := strings.Fields(queryString)
	queries := make([]query.Query, len(terms))

	for i, term := range terms {
		switch {
		case strings.HasPrefix(term, "+"):
			queries[i] = query.NewMatchQuery(strings.TrimPrefix(term, "+"))
		case strings.HasPrefix(term, "-"):
			queries[i] = query.NewBooleanQuery(nil, nil, []query.Query{query.NewMatchQuery(strings.TrimPrefix(term, "-"))})
		default:
			queries[i] = query.NewMatchQuery(term)
		}
	}

	return query.NewBooleanQuery(queries, nil, nil)
}

func displayResults(results *bleve.SearchResult) {
	fmt.Printf("Found %d results\n\n", results.Total)

	highlighter := color.New(color.FgYellow).Add(color.Bold).SprintFunc()

	for _, hit := range results.Hits {
		fmt.Printf("Title: %s\n", hit.Fields["title"])
		fmt.Printf("Tags: %v\n", hit.Fields["tags"])

		if content, ok := hit.Fields["content"].(string); ok {
			// Display content with context and highlighted matches
			for _, fragments := range hit.Fragments {
				for _, fragment := range fragments {
					fmt.Printf("... %s ...\n", highlightMatches(fragment, highlighter))
				}
			}
			_ = content
		}

		fmt.Println(strings.Repeat("-", 40))
	}
}

func highlightMatches(text string, highlighter func(a ...interface{}) string) string {
	parts := strings.Split(text, "<mark>")
	for i := 1; i < len(parts); i++ {
		subParts := strings.SplitN(parts[i], "</mark>", 2)
		if len(subParts) == 2 {
			parts[i] = highlighter(subParts[0]) + subParts[1]
		}
	}
	return strings.Join(parts, "")
}
