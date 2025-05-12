# SHOUT 📢
A CLI tool to generate project dumps for LLM analysis in Markdown format.

## Features

- Generates structured markdown of your project
- Supports file filtering by extension
- Excludes specified directories and patterns
- Interactive TUI mode available
- Preserves project structure in output
## Installation

```sh
git clone https://github.com/CinematicCow/shout.git
cd shout
go build -o shout cmd/main.go
```
    
## Usage

```sh
shout [flags]
```

| Flag              | Description                | Example               |
|-------------------|----------------------------|-----------------------|
| -e, --extensions  | File extensions to include | -e go,md              |
| -d, --directories | Directories to scan        | -d cmd,internal       |
| -s, --skip        | Patterns to skip           | -s node_modules,*.tmp |
| -o, --output      | Output file path           | -o docs/project.md    |
| -i, --interactive | Use interactive TUI mode   | -i                    |

## Examples
- Scan a go project excluding tests:
```sh
shout -e go -s *_test.go
```
- Document multiple directories:
```
shout -d src,lib -o documentation.md
```

### Output Format
The generated markdown includes:

- Project structure tree
- All source files with syntax highlighting
- Organized by directory structure
## FAQ

#### How do I exclude multiple patterns?

Use commas: `-s node_modules,.git,*.tmp

#### Can I scan multiple root directories?

Yes, `shout -d dir1,dir2,dir3`

#### How do I include all file types?

Simply omit the `-e` flag to include all files

#### Can I use wildcards in extensions?

No, specify exact extensions: `-e go,js` not `*.go`

## Authors

- [@CinematicCow](https://www.github.com/cinematiccow)
