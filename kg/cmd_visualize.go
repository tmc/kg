package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tmc/dot"
)

func newVisualizeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "visualize",
		Short: "Generate graph representation",
		RunE: func(cmd *cobra.Command, args []string) error {
			return visualizeGraph(cmd, args)
		},
	}

	cmd.Flags().StringP("output", "o", "knowledge_graph.dot", "Output file name")
	cmd.Flags().StringP("format", "f", "dot", "Output format (dot or html)")
	cmd.Flags().StringP("layout", "l", "dot", "Graph layout algorithm (dot, neato, fdp, sfdp, twopi, circo)")
	cmd.Flags().StringSliceP("filter", "t", []string{}, "Filter nodes by tags")

	return cmd
}

func visualizeGraph(cmd *cobra.Command, args []string) error {
	notesDir := viper.GetString("notes_directory")
	if notesDir == "" {
		return fmt.Errorf("notes directory not set in config")
	}

	outputFile, _ := cmd.Flags().GetString("output")
	format, _ := cmd.Flags().GetString("format")
	layout, _ := cmd.Flags().GetString("layout")
	filterTags, _ := cmd.Flags().GetStringSlice("filter")

	graph := dot.NewGraph("KnowledgeGraph")
	// attach args:
	for _, arg := range args {
		graph.Set("args", arg)
	}

	err := filepath.Walk(notesDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(info.Name(), ".md") {
			note, err := parseNote(path)
			if err != nil {
				return fmt.Errorf("failed to parse note %s: %w", path, err)
			}

			if shouldIncludeNote(note, filterTags) {
				addNodeToGraph(graph, note)
				addEdgesToGraph(graph, note)
			}
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to build graph: %w", err)
	}

	if format == "html" {
		return generateInteractiveHTML(graph, outputFile)
	}

	return generateDOTFile(graph, outputFile, layout)
}

func shouldIncludeNote(note Note, filterTags []string) bool {
	if len(filterTags) == 0 {
		return true
	}

	for _, tag := range note.Tags {
		for _, filterTag := range filterTags {
			if tag == filterTag {
				return true
			}
		}
	}

	return false
}

func addNodeToGraph(graph *dot.Graph, note Note) {
	n := dot.NewNode("KnowledgeGraph")
	n.Set("shape", "box")
	n.Set("label", note.Title)
	graph.AddNode(n)
}

func addEdgesToGraph(graph *dot.Graph, note Note) {
	for _, connection := range note.Connections {
		_ = connection
		// TODO:
		// e := dot.NewEdge(conn..
		// graph.AddEdge(note.Filename, connection, true, nil)
	}
}

func generateDOTFile(graph *dot.Graph, outputFile, layout string) error {
	dotContent := graph.String()
	dotContent = fmt.Sprintf("digraph KnowledgeGraph {\n  layout=%s;\n%s\n}", layout, dotContent[23:])

	err := os.WriteFile(outputFile, []byte(dotContent), 0644)
	if err != nil {
		return fmt.Errorf("failed to write DOT file: %w", err)
	}

	fmt.Printf("Graph visualization saved to %s\n", outputFile)
	return nil
}

func generateInteractiveHTML(graph *dot.Graph, outputFile string) error {
	// This is a simplified version. You might want to use a proper HTML template.
	htmlContent := `
<!DOCTYPE html>
<html>
<head>
    <title>Knowledge Graph Visualization</title>
    <script src="https://d3js.org/d3.v5.min.js"></script>
    <script src="https://unpkg.com/@hpcc-js/wasm@0.3.11/dist/index.min.js"></script>
    <script src="https://unpkg.com/d3-graphviz@3.0.5/build/d3-graphviz.js"></script>
</head>
<body>
    <div id="graph" style="text-align: center;"></div>
    <script>
        d3.select("#graph").graphviz()
            .renderDot(` + "`" + graph.String() + "`" + `);
    </script>
</body>
</html>`

	err := os.WriteFile(outputFile, []byte(htmlContent), 0644)
	if err != nil {
		return fmt.Errorf("failed to write HTML file: %w", err)
	}

	fmt.Printf("Interactive graph visualization saved to %s\n", outputFile)
	return nil
}
