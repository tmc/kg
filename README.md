# kg

`kg` is a program designed to generate markdown files by connecting concepts within a knowledge graph using the Obsidian markdown format. It utilizes the `langchaingo` library to enable structured content generation via AI.

## Installation

To install `kg`, you will need Go installed on your machine. If you do not have Go installed, you can download and install it from [https://go.dev/dl/](https://go.dev/dl/). Once Go is installed, you can set up `kg` by running:

```bash
go install github.com/tmc/kg@latest
```

## Usage
```shell
kg "Concept 1" "Concept 2"
```

This will process the input concepts, identify a connecting concept, and generate a markdown file with the title named after the connecting concept. The markdown content will explain how the connecting concept relates to the two original concepts and include links using Obsidian's double bracket syntax.

## License

kg is open-source software licensed under the MIT license.
