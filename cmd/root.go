package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/spf13/cobra"
	"gopack/internal"
)

var (
	copy      bool
	estimate  bool
	verbose   bool
	ignorePat string
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
		if copy {
			if err := clipboard.WriteAll(output); err != nil {
				fmt.Fprintf(os.Stderr, "⚠ Warning: Failed to copy to clipboard (%v). Printing to terminal instead.\n", err)
				fmt.Print(output)
			} else {
				fmt.Fprintln(os.Stderr, "Done! Context packed to clipboard.")
			}
		} else {
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
