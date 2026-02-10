package ui

import (
	"fmt"
	"os"
)

// PrintSuccess prints a success message
func PrintSuccess(message string) {
	SuccessBox(message)
}

// PrintError prints an error message to stderr
func PrintError(message string) {
	fmt.Fprintf(os.Stderr, "%s", CurrentTheme.Error.Sprintf("%s %s\n", IconError, message))
}

// PrintWarning prints a warning message
func PrintWarning(message string) {
	WarningBox(message)
}

// PrintInfo prints an info message
func PrintInfo(message string) {
	InfoBox(message)
}

// PrintCommand prints a command example
func PrintCommand(command string) {
	CurrentTheme.Muted.Printf("  $ %s\n", command)
}

// PrintResult prints a result key-value pair
func PrintResult(key, value string) {
	CurrentTheme.Label.Printf("  %s: ", key)
	CurrentTheme.Body.Println(value)
}

// PrintExample prints a usage example
func PrintExample(description, command string) {
	fmt.Println()
	CurrentTheme.Muted.Printf("  %s\n", description)
	CurrentTheme.Accent.Printf("  $ %s\n", command)
}
