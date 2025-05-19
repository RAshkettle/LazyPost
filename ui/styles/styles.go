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
	ErrorColor     = lipgloss.Color("#FF0000") // Red for error messages

	// Border Styles
	// Standard border style for inactive components
	BorderStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(SecondaryColor)

	// Border style for active/focused components
	ActiveBorderStyle = BorderStyle.Copy(). // Use Copy() to avoid modifying the original
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

	// Style for general input fields (active state)
	ActiveInputStyle = ActiveBorderStyle.Copy().
		Padding(0, 1) // Add some horizontal padding for text inside input

	// Style for general input fields (inactive state)
	InactiveInputStyle = BorderStyle.Copy().
		Padding(0, 1) // Add some horizontal padding for text inside input

	// Style for the items in an open dropdown
	DropdownItemStyle = lipgloss.NewStyle().
		Padding(0, 1) // Add some horizontal padding

	// Style for the currently highlighted item in an open dropdown
	DropdownSelectedItemStyle = DropdownItemStyle.Copy().
		Background(PrimaryColor).
		Foreground(SecondaryColor)

	// Style for containers holding inputs or other components
	InputContainerStyle = BorderStyle.Copy()

	// Style for text within a dropdown
	DropdownTextStyle = lipgloss.NewStyle().
		Foreground(SecondaryColor)

	// Style for the dropdown arrow
	DropdownArrowStyle = lipgloss.NewStyle().
		Foreground(SecondaryColor)

	// Create warning style
	ToastStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#FFD700")). // Gold border
		Foreground(lipgloss.Color("#FFFFFF")).       // White text
		Background(lipgloss.Color("#A52A2A")).       // Brown-red background
		Padding(0, 1).                               // Add some padding
		Align(lipgloss.Center, lipgloss.Center).     // Center content
		Bold(true)                                   // Make the text bold

)

// Theme struct to hold all application styles
type Theme struct {
	PrimaryColor        lipgloss.Color
	SecondaryColor      lipgloss.Color
	URLColor            lipgloss.Color
	MethodColor         lipgloss.Color
	ErrorColor          lipgloss.Color
	BrightYellow        lipgloss.Color
	BorderStyle         lipgloss.Style
	ActiveBorderStyle   lipgloss.Style
	TitleStyle          lipgloss.Style
	URLTitleStyle       lipgloss.Style
	MethodTitleStyle    lipgloss.Style
	SelectedItemStyle   lipgloss.Style
	ActiveInputStyle    lipgloss.Style
	InactiveInputStyle  lipgloss.Style
	DropdownItemStyle lipgloss.Style // New style for dropdown items
	DropdownSelectedItemStyle lipgloss.Style // New style for selected dropdown items
	InputContainerStyle lipgloss.Style
	DropdownTextStyle   lipgloss.Style
	DropdownArrowStyle  lipgloss.Style
	ToastStyle          lipgloss.Style

	// New fields for additional colors and styles
	HelpTextColor          lipgloss.Color // Color for help text
	ErrorStyle          lipgloss.Style
	SuccessStyle        lipgloss.Style
	SpinnerStyle        lipgloss.Style
	HelpTextStyle       lipgloss.Style // New style for help text
}

// DefaultTheme is the instance of Theme with default styles
var DefaultTheme = Theme{
	PrimaryColor:        PrimaryColor,
	SecondaryColor:      SecondaryColor,
	URLColor:            URLColor,
	MethodColor:         MethodColor,
	ErrorColor:          ErrorColor,
	BrightYellow:        BrightYellow,
	BorderStyle:         BorderStyle,
	ActiveBorderStyle:   ActiveBorderStyle,
	TitleStyle:          TitleStyle,
	URLTitleStyle:       URLTitleStyle,
	MethodTitleStyle:    MethodTitleStyle,
	SelectedItemStyle:   SelectedItemStyle,
	ActiveInputStyle:    ActiveInputStyle,
	InactiveInputStyle:  InactiveInputStyle,
	DropdownItemStyle: DropdownItemStyle, // Initialize new style
	DropdownSelectedItemStyle: DropdownSelectedItemStyle, // Initialize new style
	InputContainerStyle: InputContainerStyle,
	DropdownTextStyle:   DropdownTextStyle,
	DropdownArrowStyle:  DropdownArrowStyle,
	ToastStyle:          ToastStyle,

	// Initialize new fields
	HelpTextColor:          lipgloss.Color("#E5C07B"), // Yellow for help text
	ErrorStyle:          lipgloss.NewStyle().Foreground(ErrorColor),
	SuccessStyle:        lipgloss.NewStyle().Foreground(BrightYellow),
	SpinnerStyle:        lipgloss.NewStyle().Foreground(PrimaryColor),
	HelpTextStyle:       lipgloss.NewStyle().Foreground(lipgloss.Color("#E5C07B")), // Yellow for help text
}
