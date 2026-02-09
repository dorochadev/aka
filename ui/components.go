package ui

import (
	"fmt"
	"strings"
)

// Box drawing characters
const (
	BoxTopLeft     = "╭"
	BoxTopRight    = "╮"
	BoxBottomLeft  = "╰"
	BoxBottomRight = "╯"
	BoxHorizontal  = "─"
	BoxVertical    = "│"
)

// Standard icons
const (
	IconSuccess = "✓"
	IconError   = "✗"
	IconWarning = "⚠"
	IconInfo    = "ℹ"
	IconPointer = "▸"
	IconBullet  = "•"
	IconArrow   = "→"
)

// Header prints a styled header with optional subtitle
func Header(title string, subtitle string) {
	CurrentTheme.Title.Println(title)
	if subtitle != "" {
		CurrentTheme.Subtitle.Println(subtitle)
	}
	Divider()
	fmt.Println()
}

// Section prints a section header
func Section(title string) {
	fmt.Println()
	CurrentTheme.Primary.Printf("%s %s\n", IconPointer, title)
}

// Divider prints a horizontal line separator
func Divider() {
	CurrentTheme.Border.Println(strings.Repeat(BoxHorizontal, 50))
}

// InfoBox prints an info message with icon
func InfoBox(message string) {
	CurrentTheme.Info.Printf("%s %s\n", IconInfo, message)
}

// SuccessBox prints a success message with icon
func SuccessBox(message string) {
	CurrentTheme.Success.Printf("%s %s\n", IconSuccess, message)
}

// WarningBox prints a warning message with icon
func WarningBox(message string) {
	CurrentTheme.Warning.Printf("%s %s\n", IconWarning, message)
}

// ErrorBox prints an error message with icon
func ErrorBox(message string) {
	CurrentTheme.Error.Printf("%s %s\n", IconError, message)
}

// List prints a bulleted list
func List(items []string) {
	for _, item := range items {
		CurrentTheme.Accent.Printf("  %s ", IconBullet)
		CurrentTheme.Body.Println(item)
	}
}

// KeyValue prints a key-value pair
func KeyValue(key, value string) {
	CurrentTheme.Label.Printf("  %-20s ", key+":")
	CurrentTheme.Body.Println(value)
}

// Table prints a formatted table with headers and rows
func Table(headers []string, rows [][]string) {
	if len(headers) == 0 || len(rows) == 0 {
		return
	}

	// Calculate column widths
	colWidths := make([]int, len(headers))
	for i, header := range headers {
		colWidths[i] = len(header)
	}
	for _, row := range rows {
		for i, cell := range row {
			if i < len(colWidths) && len(cell) > colWidths[i] {
				colWidths[i] = len(cell)
			}
		}
	}

	// Print headers
	fmt.Print("  ")
	for i, header := range headers {
		CurrentTheme.Primary.Printf("%-*s", colWidths[i]+3, header)
	}
	fmt.Println()

	// Print separator
	fmt.Print("  ")
	CurrentTheme.Border.Print(strings.Repeat(BoxHorizontal, sum(colWidths)+len(headers)*3))
	fmt.Println()

	// Print rows
	for _, row := range rows {
		fmt.Print("  ")
		for i, cell := range row {
			if i < len(colWidths) {
				CurrentTheme.Body.Printf("%-*s", colWidths[i]+3, cell)
			}
		}
		fmt.Println()
	}
}

// sum calculates the sum of integers in a slice
func sum(nums []int) int {
	total := 0
	for _, n := range nums {
		total += n
	}
	return total
}
