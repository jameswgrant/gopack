package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/spf13/cobra"
	"gopack/internal"
)

var (
	copy       bool
	estimate   bool
	verbose    bool
	ignorePat  string
	outputFlag string
)

var rootCmd = &cobra.Command{
	Use:   "gopack [path]",
	Short: "Aggregate directory contents into a single formatted string",
	Long: `GoContextPacker is a CLI tool that traverses a directory,
respects .gitignore rules, and aggregates file contents into
a single Markdown-formatted string for easy pasting into LLMs.`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		// Get the target path (default to current directory)
		targetPath := "."
		if len(args) > 0 {
			targetPath = args[0]
		}

		// Create walker
		walker, err := internal.NewWalker(targetPath)
		if err != nil {
			return fmt.Errorf("failed to initialize walker: %w", err)
		}

		// Walk the directory
		files, err := walker.Walk()
		if err != nil {
			return fmt.Errorf("failed to walk directory: %w", err)
		}

		// Show verbose info
		if verbose {
			fmt.Fprintf(os.Stderr, "Found %d files\n", len(files))
			for _, file := range files {
				fmt.Fprintf(os.Stderr, "  %s\n", file.Path)
			}
		}

		// Format the output
		formatter := internal.NewFormatter(files)
		output := formatter.Format()

		// Show token estimate if requested
		if estimate {
			tokenCount := formatter.TokenCount()
			fmt.Fprintln(os.Stderr, formatTokenEstimate(tokenCount))
		}

		// Output the result
		if outputFlag != "" {
			// Write to file
			filePath, err := resolveOutputPath(outputFlag, targetPath)
			if err != nil {
				return err
			}

			if err := os.WriteFile(filePath, []byte(output), 0644); err != nil {
				return fmt.Errorf("failed to write output file: %w", err)
			}
			fmt.Fprintf(os.Stderr, "Done! Context written to %s\n", filePath)
		} else if copy {
			if err := clipboard.WriteAll(output); err != nil {
				fmt.Fprintf(os.Stderr, "⚠ Warning: Failed to copy to clipboard (%v). Printing to terminal instead.\n", err)
				fmt.Print(output)
			} else {
				fmt.Fprintln(os.Stderr, "Done! Context packed to clipboard.")
			}
		} else if !estimate || verbose {
			// Print output unless --estimate was used alone (without --verbose)
			fmt.Print(output)
		}

		return nil
	},
}

// formatWithCommas adds thousand separators to a number
func formatWithCommas(num int) string {
	str := strconv.Itoa(num)
	var result strings.Builder

	for i, ch := range str {
		if i > 0 && (len(str)-i)%3 == 0 {
			result.WriteRune(',')
		}
		result.WriteRune(ch)
	}

	return result.String()
}

// resolveOutputPath determines the final output file path
// If outputPath is empty, returns empty string
// If outputPath is a directory, returns path/context.txt
// Otherwise returns the outputPath as-is
func resolveOutputPath(outputPath string, targetPath string) (string, error) {
	if outputPath == "" {
		return "", nil
	}

	// Check if it's a directory
	info, err := os.Stat(outputPath)
	if err == nil && info.IsDir() {
		return filepath.Join(outputPath, "context.txt"), nil
	}

	// If the path doesn't exist, treat it as a file path
	if os.IsNotExist(err) {
		// Ensure the directory exists
		dir := filepath.Dir(outputPath)
		if dir != "." && dir != "" {
			if err := os.MkdirAll(dir, 0755); err != nil {
				return "", fmt.Errorf("failed to create output directory: %w", err)
			}
		}
		return outputPath, nil
	}

	// If there's another error, return it
	if err != nil {
		return "", fmt.Errorf("failed to check output path: %w", err)
	}

	// Path exists and is not a directory (it's a file)
	return outputPath, nil
}

// formatTokenEstimate returns a professionally formatted token estimate box
func formatTokenEstimate(tokenCount int) string {
	formattedCount := formatWithCommas(tokenCount)
	message := fmt.Sprintf("TOKEN ESTIMATE: ~%s tokens", formattedCount)

	// ANSI color codes
	cyan := "\033[36m"
	bold := "\033[1m"
	reset := "\033[0m"

	// Calculate box width
	boxWidth := len(message) + 4

	// Build the box with ANSI colors
	topLine := strings.Repeat("─", boxWidth)
	bottomLine := strings.Repeat("─", boxWidth)

	box := fmt.Sprintf(
		"%s%s┌%s┐%s\n%s%s│ %s │%s\n%s%s└%s┘%s",
		cyan, bold, topLine, reset,
		cyan, bold, message, reset,
		cyan, bold, bottomLine, reset,
	)

	return box
}

func init() {
	rootCmd.Flags().BoolVarP(&copy, "copy", "c", false, "Copy output to system clipboard")
	rootCmd.Flags().StringVarP(&outputFlag, "output", "o", "", "Write output to a file (defaults to context.txt in the target directory if a directory is provided)")
	rootCmd.Flags().BoolVar(&estimate, "estimate", false, "Calculate token count and display to stderr")
	rootCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Show which files are being packed")
	rootCmd.Flags().StringVar(&ignorePat, "ignore-pattern", "", "Add temporary ignore patterns (e.g., *.test.go)")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
