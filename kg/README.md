# Knowledge Graph Manager (kg)

Knowledge Graph Manager (kg) is a CLI tool designed to manage a knowledge graph using markdown files with YAML frontmatter. It provides various commands to create, edit, search, and visualize your knowledge graph.

## Features

- Create new notes with AI-suggested tags
- Edit existing notes while preserving frontmatter
- Connect concepts with AI-generated content
- Search for keywords in content and frontmatter
- Visualize the knowledge graph
- Export the graph to JSON or CSV formats
- Import external markdown files
- Create backups of your knowledge graph
- Manage configuration settings
- Bulk update and normalize frontmatter across files
- Display statistics about your knowledge graph

## Installation

To install kg, make sure you have Go installed on your system, then run:

```
go get github.com/yourusername/kg
```

## Configuration

kg uses a configuration file named `.kgrc`. By default, it looks for this file in the following locations:

1. Current directory
2. Parent directory
3. User's home directory

You can also specify a custom configuration file using the `--config` flag.

## Usage

Here are some example commands:

```
kg add "New Note Title"
kg edit "Existing Note Title"
kg connect "Concept A" "Concept B"
kg search "keyword"
kg visualize
kg export json
kg import /path/to/file.md
kg backup
kg config key value
kg frontmatter
kg stats
```

For more detailed information on each command, use the `--help` flag:

```
kg --help
kg <command> --help
```

## Development Status

This project is currently under development. The following features are planned or in progress:

- Implement note addition with AI-suggested tags
- Implement backup creation
- Implement configuration management
- Implement AI-assisted content generation to link concepts
- Implement note editing functionality
- Implement graph export functionality
- Implement frontmatter bulk update/normalization
- Implement file import functionality
- Implement listing of all notes with frontmatter info
- Implement search functionality
- Implement statistics display
- Implement graph visualization
- Implement YAML frontmatter parsing and generation
- Integrate langchaingo for interacting with the Anthropic Claude API
- Implement error handling and logging throughout the application
- Create a default .kgrc file with sensible defaults
- Implement the hierarchical config system that merges multiple .kgrc files
- Add comprehensive usage instructions for each command in the help text
- Implement file naming conventions and management
- Add unit tests for all implemented functionality

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the [MIT License](LICENSE).