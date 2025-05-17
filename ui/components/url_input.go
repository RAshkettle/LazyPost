// Package components provides UI components for the LazyPost application.
package components

import (
	"github.com/RAshkettle/LazyPost/ui/styles"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// URLInput represents the URL input component where users can enter
// the target URL for HTTP requests. It wraps the textinput.Model from
// the Bubble Tea framework to provide specialized URL input functionality.
type URLInput struct {
	TextInput textinput.Model // The underlying text input model
	Width     int             // Width of the component in characters
	Active    bool            // Whether the component is currently active/focused
}

// NewURLInput creates a new URL input component with default configuration.
// The input is initially focused and has a placeholder text.
func NewURLInput() URLInput {
	input := textinput.New()
	input.Placeholder = "Enter URL"
	input.Focus()
	input.CharLimit = 256
	input.Width = 80

	return URLInput{
		TextInput: input,
		Width:     0,
		Active:    true,
	}
}

// SetWidth sets the width of the URL input component.
// It adjusts the internal TextInput width to account for border and padding.
func (u *URLInput) SetWidth(width int) {
	u.Width = width
	u.TextInput.Width = width - 4 // Adjust for border and padding
}

// SetActive sets the active state of the URL input.
// When active, the input is focused and can receive keyboard input.
// When inactive, the input is blurred and displays with different styling.
func (u *URLInput) SetActive(active bool) {
	u.Active = active
	if active {
		u.TextInput.Focus()
	} else {
		u.TextInput.Blur()
	}
}

// GetText returns the current URL text entered by the user.
func (u *URLInput) GetText() string {
	return u.TextInput.Value()
}

// SelectAllText selects all text in the input field.
// This is used when focusing the input to allow quick replacement of the URL.
func (u *URLInput) SelectAllText() {
	// The TextInput model doesn't have a built-in select all function,
	// but we can simulate it by positioning the cursor at the end and
	// setting the text to be selected
	textLen := len(u.TextInput.Value())
	u.TextInput.SetCursor(textLen)
	// Select all text by positioning the selection at the start
	u.TextInput.CursorStart()
	// Select all text
	for i := 0; i < textLen; i++ {
		u.TextInput.SetCursor(i + 1)
	}
}

// Update processes input messages and updates the URLInput component.
// It only processes messages when the component is active.
// Returns any commands that need to be executed.
func (u *URLInput) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	if u.Active {
		u.TextInput, cmd = u.TextInput.Update(msg)
	}
	return cmd
}

// View renders the URLInput component with the appropriate styling.
// It displays a title with hotkey and the text input field with border.
func (u URLInput) View() string {
	// Define styles
	borderStyle := styles.BorderStyle

	if u.Active {
		borderStyle = styles.ActiveBorderStyle
	}

	// Use minimal padding to make it just one line tall
	borderStyle = borderStyle.Padding(0, 1)
	
	// Create simple title with number hotkey
	titleStyle := lipgloss.NewStyle().
		Bold(true)
	
	// Change title color based on active state
	if u.Active {
		titleStyle = titleStyle.Foreground(styles.PrimaryColor)
	} else {
		titleStyle = titleStyle.Foreground(styles.SecondaryColor)
	}
	
	title := titleStyle.Render("(Alt+2) URL")
	
	// Render the URL box with the title directly above it
	inputBox := borderStyle.Width(u.Width).Render(u.TextInput.View())
	
	// Position the title at the top-left of the input box
	return title + "\n" + inputBox 
}
