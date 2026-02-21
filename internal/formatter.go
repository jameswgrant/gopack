package internal

import (
	"bytes"
	"fmt"
)

// Formatter handles converting files to output format.
type Formatter struct {
	files []File
}

// NewFormatter creates a new Formatter with the given files.
func NewFormatter(files []File) *Formatter {
	return &Formatter{files: files}
}

// Format returns the formatted output as a string.
func (f *Formatter) Format() string {
	var buf bytes.Buffer

	for i, file := range f.files {
		// Write file header
		fmt.Fprintf(&buf, "File: %s\n", file.Path)
		// Write file content
		buf.Write(file.Content)
		// Add blank line between files (except after the last one)
		if i < len(f.files)-1 {
			buf.WriteString("\n\n")
		}
	}

	return buf.String()
}

// TokenCount returns an estimated token count (character count / 4).
func (f *Formatter) TokenCount() int {
	totalChars := 0
	for _, file := range f.files {
		totalChars += len(file.Path) + len("File: \n")
		totalChars += len(file.Content)
	}
	return totalChars / 4
}
