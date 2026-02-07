# ai-changelog

A CLI tool that generates polished changelogs from your git history using a local LLM via [Ollama](https://ollama.com). It reads your commits, sends them to an AI model running on your machine, and outputs release notes written from the user's perspective.

When Ollama is unavailable, it falls back to structured output by grouping commits using [Conventional Commits](https://www.conventionalcommits.org/) prefixes.

## How It Works

```
git log  ──>  parse commits  ──>  Ollama LLM  ──>  polished changelog
                                      |
                                  (if unavailable)
                                      |
                                      v
                              structured fallback
                          (group by feat/fix/docs/...)
```

1. Reads commits from git history (optionally filtered with `--since`)
2. Sends them to a local Ollama model that collapses related commits into user-facing entries
3. Outputs Markdown or plain text to stdout or a file

The LLM prompt instructs the model to write from the user's perspective, collapse implementation details into high-level entries, skip test/refactor commits, and order by importance.

## Requirements

- **Go** 1.25.6+
- **Git** installed and available in PATH
- **Ollama** running locally (optional, enables AI-powered output)

### Installing Ollama

```bash
# macOS
brew install ollama

# Linux
curl -fsSL https://ollama.com/install.sh | sh
```

Then pull a model and start the server:

```bash
ollama pull llama3.2
ollama serve
```

## Installation

### Homebrew (recommended)

The formula lives in this repo. Tap it and install:

```bash
brew tap brognilucas/ai-changelog https://github.com/brognilucas/ai-changelog
brew install ai-changelog
```

### From source

```bash
git clone https://github.com/brognilucas/ai-changelog.git
cd ai-changelog
go build -o ai-changelog
```

To install the binary into `$GOPATH/bin` (requires the repo to be public or your Git credentials configured):

```bash
go install github.com/brognilucas/ai-changelog@latest
```

### Publishing a new version (maintainers)

1. Tag a release and push: `git tag v1.0.0 && git push origin v1.0.0`.
2. Compute the tarball checksum and update `Formula/ai-changelog.rb` in this repo:
   ```bash
   curl -sL "https://github.com/brognilucas/ai-changelog/archive/refs/tags/v1.0.0.tar.gz" | shasum -a 256
   ```
3. Update the `url` and `sha256` in `Formula/ai-changelog.rb` for the new version and push.

## Usage

Run the command inside any git repository:

```bash
# Generate changelog for all commits (stdout)
ai-changelog

# Since a specific tag
ai-changelog --since v1.0.0

# With a version header, written to a file
ai-changelog -s v1.0.0 -V v1.1.0 -o CHANGELOG.md

# Use a different model
ai-changelog -m mistral -s v2.0.0

# Plain text output
ai-changelog -f plain
```

### Flags

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--since` | `-s` | _(all commits)_ | Generate changelog since a tag or ref (e.g. `v1.0.0`, `HEAD~10`) |
| `--model` | `-m` | `llama3.2` | Ollama model to use for summarization |
| `--format` | `-f` | `markdown` | Output format: `markdown` or `plain` |
| `--output` | `-o` | _(stdout)_ | Write changelog to a file instead of stdout |
| `--version` | `-V` | _(none)_ | Version label for the changelog header |

## Example Output

### AI-powered (with Ollama)

```markdown
# v1.0.0

_First release of ai-changelog, a CLI tool that generates release notes from git history using a local LLM._

## Highlights

- Generate changelogs from git commits using a local Ollama model
- Automatic fallback to structured output when Ollama is unavailable
- Support for Markdown and plain text output formats

## Improvements

- Configurable model selection, version labels, and date filtering
- File output option for CI/CD integration
```

### Structured fallback (without Ollama)

```markdown
# Changelog

## New Features

- add main entry point (7cb2449)
- implement generate command (dece3fd)
- add Ollama fallback behavior (f624c27)

## Bug Fixes

- handle git errors (9987173)
```

## Project Structure

```
ai-changelog/
├── main.go                     # Entry point
├── cmd/
│   ├── root.go                 # CLI flags and command setup
│   └── generate.go             # Core generation logic
├── internal/
│   ├── git/
│   │   └── git.go              # Git log parsing and commit retrieval
│   ├── ollama/
│   │   └── client.go           # Ollama API client and prompt building
│   └── changelog/
│       ├── grouper.go          # Commit categorization and sorting
│       └── renderer.go         # Markdown and plain text renderers
└── tests/                      # Unit tests for all packages
```

## Running Tests

```bash
go test ./...

# With coverage
go test -cover ./...
```

## How It Handles Failures

- **Ollama not running**: Prints a warning to stderr and uses structured fallback
- **LLM returns empty/bad output**: Falls back to structured grouping
- **No commits found**: Prints "No commits found." and exits cleanly
- **Invalid git ref in `--since`**: Returns a descriptive git error

The tool is designed to always produce useful output, even without a running LLM.

## License

MIT