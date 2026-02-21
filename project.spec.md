# Specification: GoContextPacker (CLI Tool)

## 1. Overview
A CLI tool that traverses a directory, respects `.gitignore` rules, and aggregates file contents into a single Markdown-formatted string. It copies the output to the system clipboard for easy pasting into LLMs.

## 2. Technical Stack
- **Language:** Go
- **CLI Framework:** `github.com/spf13/cobra`
- **Gitignore Parsing:** `github.com/go-git/go-git/v5/plumbing/format/gitignore` (or native implementation if simpler)
- **Clipboard:** `github.com/atotto/clipboard`

## 3. Core Features
### A. File Discovery
- Recursively walk the current directory (default) or a specified path.
- **CRITICAL:** Must read and respect `.gitignore` files at the root and subdirectories.
- Ignore `.git/` folder and binary files (images, executables) by default.

### B. Output Formatting
The output must look like this:

File: src/main.go
[File Content Here]

File: README.md
[File Content Here]

### C. Token Estimation
- specific command or flag `--estimate` that calculates a rough token count (Character Count / 4) and displays it in the terminal (stderr) so the user knows if it fits in the context window.

### D. Clipboard
- Default behavior: Print to stdout.
- Flag `-c` or `--copy`: Pipe output directly to system clipboard.

## 4. CLI Interface
- `gopack [path] [flags]`
- Flags:
  - `-c, --copy`: Copy to clipboard.
  - `-v, --verbose`: Show which files are being packed.
  - `--ignore-pattern`: Add temporary ignore patterns (e.g., `*.test.go`).

## 5. Implementation Plan
1. **Setup:** Initialize module, set up Cobra skeleton.
2. **Walker:** Implement the recursive file walker with `.gitignore` logic.
3. **Filter:** Add binary file detection (read first 512 bytes -> http.DetectContentType).
4. **Output:** format the stream to buffer.
5. **Clipboard:** Integrate clipboard library.
