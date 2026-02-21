package internal

import (
	"bufio"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// File represents a file to be included in the output.
type File struct {
	Path    string
	Content []byte
}

// Walker traverses a directory and filters files based on .gitignore rules.
type Walker struct {
	rootPath string
	patterns map[string][]string // dir -> patterns
}

// NewWalker creates a new Walker for the given root path.
func NewWalker(rootPath string) (*Walker, error) {
	if rootPath == "" {
		rootPath = "."
	}

	// Resolve to absolute path
	absPath, err := filepath.Abs(rootPath)
	if err != nil {
		return nil, err
	}

	w := &Walker{
		rootPath: absPath,
		patterns: make(map[string][]string),
	}

	// Load root .gitignore
	w.loadGitignore(absPath)

	return w, nil
}

// Walk traverses the directory and returns a slice of File structs.
func (w *Walker) Walk() ([]File, error) {
	var files []File

	err := filepath.Walk(w.rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, _ := filepath.Rel(w.rootPath, path)

		// Skip .git directory
		if info.IsDir() && info.Name() == ".git" {
			return filepath.SkipDir
		}

		// For directories, try to load .gitignore
		if info.IsDir() {
			w.loadGitignore(path)
		}

		// Check if path is ignored
		if w.isIgnored(relPath) {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		// Only process regular files
		if !info.IsDir() && info.Mode().IsRegular() {
			// Check if file is binary
			if isBinary(path) {
				return nil
			}

			// Read file content
			content, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			files = append(files, File{
				Path:    relPath,
				Content: content,
			})
		}

		return nil
	})

	return files, err
}

// loadGitignore loads patterns from a .gitignore file in the directory.
func (w *Walker) loadGitignore(dirPath string) {
	gitignorePath := filepath.Join(dirPath, ".gitignore")
	file, err := os.Open(gitignorePath)
	if err != nil {
		return // .gitignore doesn't exist or can't be read
	}
	defer file.Close()

	var patterns []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		patterns = append(patterns, line)
	}

	if len(patterns) > 0 {
		w.patterns[dirPath] = patterns
	}
}

// isIgnored checks if a path matches any gitignore patterns.
func (w *Walker) isIgnored(relPath string) bool {
	// Normalize path separators
	relPath = strings.ReplaceAll(relPath, string(filepath.Separator), "/")
	parts := strings.Split(relPath, "/")

	// Check patterns from root directory
	for _, pattern := range w.patterns[w.rootPath] {
		if matchPattern(relPath, parts, pattern) {
			return true
		}
	}

	return false
}

// matchPattern checks if a path matches a gitignore pattern.
func matchPattern(fullPath string, parts []string, pattern string) bool {
	// Remove trailing slash from pattern
	pattern = strings.TrimSuffix(pattern, "/")

	// If pattern starts with /, it's relative to root
	if strings.HasPrefix(pattern, "/") {
		pattern = strings.TrimPrefix(pattern, "/")
		return simpleMatch(fullPath, pattern)
	}

	// Pattern can match any part of the path
	if strings.Contains(pattern, "/") {
		return simpleMatch(fullPath, pattern)
	}

	// Pattern matches any path component
	for _, part := range parts {
		if simpleMatch(part, pattern) {
			return true
		}
	}

	return false
}

// simpleMatch performs a simple glob-style match.
// Supports * (any chars) and ? (single char).
func simpleMatch(name, pattern string) bool {
	matched, _ := filepath.Match(pattern, name)
	return matched
}

// isBinary detects if a file is binary by reading its first 512 bytes.
func isBinary(filePath string) bool {
	file, err := os.Open(filePath)
	if err != nil {
		return true // Assume binary on error
	}
	defer file.Close()

	// Read first 512 bytes
	buffer := make([]byte, 512)
	n, err := io.ReadFull(file, buffer)
	if err != nil && err != io.ErrUnexpectedEOF {
		return true // Assume binary on error
	}

	// Use http.DetectContentType to check if it's a text file
	contentType := http.DetectContentType(buffer[:n])
	return !strings.HasPrefix(contentType, "text/")
}
