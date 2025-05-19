// Package styles provides consistent styling for the LazyPost application.
// It defines colors, border styles, and text formats used throughout the UI.
package styles

import (
	"github.com/charmbracelet/lipgloss"
)

// Common styling constants used throughout the application
var (
	// Colors
	PrimaryColor   = lipgloss.Color("#00FF00") // Green for active borders
	BrightYellow   = lipgloss.Color("#FFFF00") // Bright yellow for selected method
	SecondaryColor = lipgloss.Color("#FFFFFF") // White for general text and inactive borders
	URLColor       = lipgloss.Color("#00BFFF") // Bright blue color for URL elements
	MethodColor    = lipgloss.Color("#00BFFF") // Blue color for Method elements

	// Border Styles
	// Standard border style for inactive components
	BorderStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(SecondaryColor)

	// Border style for active/focused components
	ActiveBorderStyle = BorderStyle.
				BorderForeground(PrimaryColor)

	// Text Styles
	// General title style for components
	TitleStyle = lipgloss.NewStyle().
			Foreground(SecondaryColor).
			Bold(true)

	// Title style specific for URL components
	URLTitleStyle = lipgloss.NewStyle().
			Foreground(URLColor).
			Bold(true)

	// Title style specific for Method components
	MethodTitleStyle = lipgloss.NewStyle().
				Foreground(MethodColor).
				Bold(true)

	// Style for selected items in lists or dropdowns
	SelectedItemStyle = lipgloss.NewStyle().
				Foreground(BrightYellow).
				Bold(true)

		// Create warning style
	ToastStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#FFD700")). // Gold border
			Foreground(lipgloss.Color("#FFFFFF")).       // White text
			Background(lipgloss.Color("#A52A2A")).       // Brown-red background
			Padding(0, 1).                               // Add some padding
			Align(lipgloss.Center, lipgloss.Center).     // Center content                             // Use the specified width
			Bold(true)                                   // Make the text bold

	HelpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFF00")). // Yellow color
			Align(lipgloss.Right).
			Bold(true) // Make it bold

)
