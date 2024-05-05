package main

import (
	"context"
	_ "embed"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/anthropic"
	"github.com/tmc/langchaingo/schema"
)

//go:embed system-prompt.txt
var systemPrompt string

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}

// extract "title: (title)" from a string:
var titleExtractRe = regexp.MustCompile(`title: (.*)\n`)

type fileWriter struct {
	parts []string
}

func (fw *fileWriter) addPart(chunk []byte) {
	fmt.Print(string(chunk))
	fw.parts = append(fw.parts, string(chunk))
	// extract title from the first part:
	combined := strings.Join(fw.parts, "")
	if !titleExtractRe.MatchString(combined) {
		return
	}
	title := titleExtractRe.FindStringSubmatch(combined)[1]
	if err := os.WriteFile(title+".md", []byte(combined), 0644); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
}

func run() error {
	ctx := context.Background()
	llm, err := anthropic.New(
		anthropic.WithModel("claude-3-opus-20240229"),
	)
	if err != nil {
		return err
	}
	fw := &fileWriter{}
	messages := []llms.MessageContent{
		llms.TextParts(schema.ChatMessageTypeSystem, systemPrompt),
		llms.TextParts(schema.ChatMessageTypeHuman, strings.Join(os.Args[1:], " ")),
	}
	_, err = llm.GenerateContent(ctx,
		messages,
		llms.WithTemperature(0.9),
		llms.WithMaxTokens(4096),
		llms.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {
			fw.addPart(chunk)
			return nil
		}),
	)
	fmt.Println()
	if err != nil {
		return err
	}
	return nil
}
