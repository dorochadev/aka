package ui

import (
	"os"

	"github.com/gookit/color"
)

// Theme defines the color scheme for the CLI
type Theme struct {
	// Core brand colors
	Primary   color.RGBColor
	Secondary color.RGBColor
	Accent    color.RGBColor
	Muted     color.RGBColor

	// Semantic colors
	Success color.RGBColor
	Error   color.RGBColor
	Warning color.RGBColor
	Info    color.RGBColor

	// UI elements
	Border      color.RGBColor
	Highlight   color.RGBColor
	HighlightBg color.RGBColor
	Disabled    color.RGBColor
	Placeholder color.RGBColor

	// Text hierarchy
	Title    color.RGBColor
	Subtitle color.RGBColor
	Body     color.RGBColor
	Label    color.RGBColor
}

// CurrentTheme is the active theme
var CurrentTheme *Theme

func init() {
	// Check for NO_COLOR environment variable
	if os.Getenv("NO_COLOR") != "" || os.Getenv("TERM") == "dumb" {
		CurrentTheme = MonochromeTheme()
	} else {
		CurrentTheme = DefaultTheme()
	}
}

// DefaultTheme creates the default theme with dusk blue and steel blue
func DefaultTheme() *Theme {
	return &Theme{
		// Core brand colors - Dusk Blue palette
		Primary:   color.RGB(48, 77, 109),   // Dusk Blue #304D6D
		Secondary: color.RGB(130, 160, 188), // Steel Blue #82A0BC
		Accent:    color.RGB(100, 140, 170), // Mid-tone blue (between primary and secondary)
		Muted:     color.RGB(120, 130, 140), // Muted gray-blue

		// Semantic colors
		Success: color.RGB(130, 180, 130), // Muted green
		Error:   color.RGB(200, 120, 120), // Muted red
		Warning: color.RGB(200, 180, 130), // Muted yellow
		Info:    color.RGB(130, 160, 188), // Steel Blue

		// UI elements
		Border:      color.RGB(48, 77, 109),   // Dusk Blue
		Highlight:   color.RGB(100, 140, 170), // Mid-tone blue
		HighlightBg: color.RGB(48, 77, 109),   // Dusk Blue background
		Disabled:    color.RGB(120, 130, 140), // Muted
		Placeholder: color.RGB(120, 130, 140), // Muted

		// Text hierarchy
		Title:    color.RGB(48, 77, 109),   // Dusk Blue
		Subtitle: color.RGB(130, 160, 188), // Steel Blue
		Body:     color.RGB(200, 210, 220), // Light gray
		Label:    color.RGB(130, 160, 188), // Steel Blue
	}
}

// MonochromeTheme returns a theme without colors (for NO_COLOR support)
func MonochromeTheme() *Theme {
	noColor := color.RGB(255, 255, 255)
	return &Theme{
		Primary:     noColor,
		Secondary:   noColor,
		Accent:      noColor,
		Muted:       noColor,
		Success:     noColor,
		Error:       noColor,
		Warning:     noColor,
		Info:        noColor,
		Border:      noColor,
		Highlight:   noColor,
		HighlightBg: noColor,
		Disabled:    noColor,
		Placeholder: noColor,
		Title:       noColor,
		Subtitle:    noColor,
		Body:        noColor,
		Label:       noColor,
	}
}
