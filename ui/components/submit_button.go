package components

import (
	"github.com/RAshkettle/LazyPost/ui/styles"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// SubmitButton represents a clickable button component that can be focused and activated.
// It provides a standard button interface with visual feedback for active state and
// handles key press events for interaction.
type SubmitButton struct {
	Label  string // Text displayed on the button
	Width  int    // Width of the button in characters
	Height int    // Height of the button in characters
	Active bool   // Whether the button is currently active/focused
}

// NewButton creates a new button component with the specified label.
// The button is initialized with zero width and height, and inactive state.
func NewButton(label string) SubmitButton {
	return SubmitButton{
		Label:  label,
		Width:  0,
		Height: 0,
		Active: false,
	}
}

// SetWidth sets the width of the button in characters.
func (b *SubmitButton) SetWidth(width int) {
	b.Width = width
}

// SetHeight sets the height of the button in characters.
func (b *SubmitButton) SetHeight(height int) {
	b.Height = height
}

// SetActive sets the active state of the button.
// When a button is active, it has visual styling to indicate focus.
func (b *SubmitButton) SetActive(active bool) {
	b.Active = active
}

// Update processes input messages and updates the button state.
// Returns a tea.Cmd (always nil for Button) and a boolean indicating if the button was activated.
// The boolean is true if Enter was pressed while the button was active.
func (b *SubmitButton) Update(msg tea.Msg) (tea.Cmd, bool) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if !b.Active {
			return nil, false
		}
		
		if msg.String() == "enter" {
			// When Enter is pressed and button is active
			return nil, true
		}
	}
	return nil, false
}

// View renders the button component as a string for terminal display.
// The rendered button includes a border and content, with styling based on the active state.
// When active, the button has a highlighted border and background.
func (b SubmitButton) View() string {
	// Define styles
	borderStyle := styles.BorderStyle
	
	if b.Active {
		borderStyle = styles.ActiveBorderStyle
	}
	
	// Content style
	contentStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FFFFFF")).
		Padding(0, 1)
	
	// If button is active, highlight the text with a different background
	if b.Active {
		contentStyle = contentStyle.Background(lipgloss.Color("#555555"))
	}
	
	// Create content - only show the label text in the button
	content := contentStyle.Render(b.Label)
	
	// Create button title
	titleStyle := lipgloss.NewStyle().Bold(true)
	
	// Change title color based on active state
	if b.Active {
		titleStyle = titleStyle.Foreground(styles.PrimaryColor)
	} else {
		titleStyle = titleStyle.Foreground(styles.SecondaryColor)
	}
	
	// Show hotkey for Submit button, otherwise invisible placeholder
	var title string
	if b.Label == "Submit" {
		title = titleStyle.Render("(Alt+5)")
	} else {
		title = titleStyle.Render(" ")
	}
	
	// Render button with border
	button := borderStyle.
		Width(b.Width).
		Align(lipgloss.Center, lipgloss.Center).
		Render(content)
	
	// Return title plus button for proper vertical alignment
	return title + "\n" + button
}
