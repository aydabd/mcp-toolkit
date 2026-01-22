// Package prompt provides cross-platform terminal input utilities.
package prompt

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

var (
	reader           = bufio.NewReader(os.Stdin)
	output io.Writer = os.Stdout
)

// SetReader sets the input reader (for testing).
func SetReader(r io.Reader) {
	reader = bufio.NewReader(r)
}

// SetOutput sets the output writer (for testing).
func SetOutput(w io.Writer) {
	output = w
}

// Reset restores default stdin/stdout.
func Reset() {
	reader = bufio.NewReader(os.Stdin)
	output = os.Stdout
}

// String prompts for a string value.
func String(label string) string {
	fmt.Fprintf(output, "%s: ", label)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

// Default prompts with a default value.
func Default(label, defaultVal string) string {
	fmt.Fprintf(output, "%s [%s]: ", label, defaultVal)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	if input == "" {
		return defaultVal
	}
	return input
}

// Secret prompts for sensitive input (shows input on terminal, no masking for compatibility).
func Secret(label string) string {
	fmt.Fprintf(output, "%s: ", label)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

// Confirm prompts for yes/no confirmation (default: no).
func Confirm(label string) bool {
	fmt.Fprintf(output, "%s [y/N]: ", label)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(strings.ToLower(input))
	return input == "y" || input == "yes"
}

// ConfirmDefault prompts for yes/no with specified default.
func ConfirmDefault(label string, defaultYes bool) bool {
	prompt := "[y/N]"
	if defaultYes {
		prompt = "[Y/n]"
	}
	fmt.Fprintf(output, "%s %s: ", label, prompt)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(strings.ToLower(input))
	if input == "" {
		return defaultYes
	}
	return input == "y" || input == "yes"
}
