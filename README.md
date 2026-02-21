# GoContextPacker

A CLI tool that intelligently aggregates your project files into a single, formatted string perfect for sharing with Large Language Models (LLMs). It respects `.gitignore` rules, filters out binary files, and estimates token counts.

## Features

**Smart File Discovery**
- Recursively traverses directories and respects `.gitignore` rules
- Automatically filters out binary files and `.git/` directories
- Includes nested `.gitignore` patterns at any directory level

**Token Estimation**
- Calculate approximate token count (character count / 4) with the `--estimate` flag
- Professional-looking formatted output with thousand separators
- Helps you understand if content fits within LLM context windows

**Clipboard Integration**
- Copy aggregated content directly to system clipboard with `--copy`
- Graceful fallback to terminal output if clipboard is unavailable
- Success confirmation and helpful error messages

**Flexible Output**
- Default: Print to stdout (perfect for piping)
- Markdown-formatted with clear file headers
- All diagnostic output goes to stderr (won't interfere with piped content)

## Installation

### Prerequisites
- Go 1.22 or later

### Build from Source

```bash
git clone <repository-url>
cd gopack
go build -o bin/gopack ./cmd
```

Then add `bin/gopack` to your PATH or use it directly:
```bash
./bin/gopack [path] [flags]
```

## Usage

### Basic Usage

Pack the current directory:
```bash
./bin/gopack
```

Pack a specific directory:
```bash
./bin/gopack ./src
```

### Flags

#### `-c, --copy`
Copy the aggregated content directly to your system clipboard instead of printing to the terminal.

```bash
./bin/gopack ./src --copy
# Output: Done! Context packed to clipboard.
```

If clipboard operations fail, the tool automatically falls back to printing the content to the terminal with a warning:
```
⚠ Warning: Failed to copy to clipboard (...). Printing to terminal instead.
```

#### `--estimate`
Calculate and display the estimated token count using a professional formatted box.

```bash
./bin/gopack ./src --estimate
# Output:
# ┌────────────────────────────┐
# │ TOKEN ESTIMATE: ~1,250 tokens │
# └────────────────────────────┘
# [followed by file contents]
```

#### `-v, --verbose`
Show detailed information about which files are being packed.

```bash
./bin/gopack ./src --verbose
# Output:
# Found 12 files
#   src/main.go
#   src/utils.go
#   ...
# [followed by file contents]
```

#### `--ignore-pattern`
Add temporary ignore patterns (in addition to `.gitignore` rules) using glob syntax.

```bash
./bin/gopack ./src --ignore-pattern "*.test.go"
./bin/gopack ./src --ignore-pattern "*.log"
```

### Combined Examples

```bash
# Get a token estimate and copy to clipboard
./bin/gopack ./src --estimate --copy

# Verbose output with token estimation
./bin/gopack ./src --verbose --estimate

# Copy to clipboard with all diagnostic output
./bin/gopack ./src --copy --verbose --estimate
```

### Output to Files

Since content goes to stdout by default, you can redirect it:

```bash
# Save to a file
./bin/gopack ./src > context.txt

# Send through a pipe
./bin/gopack ./src | head -100

# Copy to clipboard AND save to file
./bin/gopack ./src --copy | tee backup.txt
# (Note: --copy suppresses stdout, use in reverse: tee then pipe to clipboard)
```

## How It Works

### File Selection Process

1. **Directory Traversal** - Recursively walks the target directory
2. **Gitignore Parsing** - Respects `.gitignore` rules at all directory levels
3. **Binary Detection** - Automatically skips binary files (images, executables, etc.)
4. **Content Aggregation** - Combines all text files into a single string

### Output Format

Each file is prefixed with a header for clarity:

```
File: path/to/file.go
[file contents]

File: path/to/another.md
[file contents]
```

### Token Estimation

The token count estimate uses a simple formula: `character count / 4`. This provides a quick approximation useful for understanding context window constraints:

- **GPT-3.5/4**: ~4k-128k tokens
- **Claude**: ~100k-200k tokens
- **Other models**: Check your provider's documentation

## Examples

### Example 1: Quick Copy for LLM Analysis

```bash
./bin/gopack ./myproject --copy --estimate
```

This will:
- Show the token estimate in a formatted box
- Copy the entire project structure and contents to clipboard
- Display a success message
- Leave your clipboard ready to paste into your LLM

### Example 2: Explore Before Copying

```bash
./bin/gopack ./myproject --verbose --estimate
```

This will:
- Show all files being included
- Display token estimate
- Print content to stdout so you can review it first
- Then you can pipe it to `--copy` if satisfied

### Example 3: Custom Filtering

```bash
./bin/gopack ./src --ignore-pattern "*.test.go" --ignore-pattern "*.log"
```

This will:
- Respect `.gitignore` rules
- Additionally exclude `.test.go` and `.log` files
- Aggregate remaining files

## Troubleshooting

### Clipboard Not Working

If you see a warning about clipboard failures:
- **Linux**: Install `xclip` or `xsel` (required for clipboard access)
  ```bash
  # Ubuntu/Debian
  sudo apt-get install xclip

  # Fedora
  sudo dnf install xclip
  ```
- **macOS**: Clipboard should work out of the box
- **Windows**: Clipboard should work out of the box

The tool automatically falls back to printing to the terminal if clipboard is unavailable.

### Too Many/Few Files Included

- Check your `.gitignore` files in the root and subdirectories
- Use `--verbose` to see exactly which files are being included
- Use `--ignore-pattern` to temporarily exclude additional files

### Token Estimate Seems Off

The estimate uses `character count / 4` as a rough approximation. This works well for most use cases but may vary by model:
- Actual token counts depend on the tokenizer used by your specific LLM
- Use this estimate as a general guide, not an absolute measure

## License

GoContextPacker is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

You are free to use, modify, and distribute this software for any purpose, including commercial use.

## Contributing

We welcome contributions! Feel free to submit issues, fork the repository, and create pull requests.
