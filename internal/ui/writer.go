package ui

import (
	"regexp"
	"strings"

	"github.com/rivo/tview"
)

// WriterWrapper is used to capture stdout and stderr and print them to the logView.
type WriterWrapper struct {
	logView     *tview.TextView
	app         *tview.Application
	accumulated string // Accumulate the formatted text with color tags
}

// Write method to implement io.Writer, appending content to logView.
func (w *WriterWrapper) Write(p []byte) (n int, err error) {
	// Convert ANSI escape codes to tview-compatible color tags
	convertedText := ansiToTview(string(p))

	// Append the new converted text to accumulated text
	w.accumulated += convertedText

	// Set the accumulated text to the logView
	w.logView.SetText(w.accumulated)
	w.app.Draw()
	return len(p), nil
}

// ansiToTview converts ANSI escape sequences into tview-compatible color tags.
func ansiToTview(text string) string {
	// Define the ANSI escape sequence regex pattern
	ansiRegex := regexp.MustCompile(`\033\[(\d+)(;\d+)?m`)

	// Map ANSI color codes to tview color tags
	ansiToTviewMap := map[string]string{
		"30": "[black]",
		"31": "[red]",
		"32": "[green]",
		"33": "[yellow]",
		"34": "[blue]",
		"35": "[magenta]",
		"36": "[cyan]",
		"37": "[white]",
		"0":  "[-]", // Reset/Default color
	}

	// Replace ANSI codes with tview-compatible color tags
	return ansiRegex.ReplaceAllStringFunc(text, func(match string) string {
		// Extract the color code (e.g., "34" for blue)
		code := strings.Trim(match, "\033[m")
		// Convert to tview tag if available
		if colorTag, ok := ansiToTviewMap[code]; ok {
			return colorTag
		}
		// If no match, return the original sequence (could be an unsupported code)
		return match
	})
}
