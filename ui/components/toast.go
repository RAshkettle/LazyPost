// Package components provides UI components for the LazyPost application.
package components

import (
	"time"

	"github.com/RAshkettle/LazyPost/ui/styles"
	tea "github.com/charmbracelet/bubbletea"
)

// TickMsg is sent when the timer ticks.
// It is used for automatic dismissal timing of toast notifications.
type TickMsg time.Time

// Toast represents a temporary notification that displays messages to the user.
// It can show success, warning, or error messages with a dismissal option.
type Toast struct {
	Message   string // The text message to display in the toast
	Visible   bool   // Whether the toast is currently visible
	Width     int    // Width of the toast in characters
	Height    int    // Height of the toast in characters
	Dismissed bool   // Whether the toast has been dismissed by the user
}

// NewToast creates a new toast notification with default values.
// The toast is initially hidden until Show() is called.
func NewToast() Toast {
	return Toast{
		Message:   "",
		Visible:   false,
		Width:     0,
		Height:    0,
		Dismissed: false,
	}
}

// SetWidth sets the width of the toast notification in characters.
func (t *Toast) SetWidth(width int) {
	t.Width = width
}

// SetHeight sets the height of the toast notification in characters.
func (t *Toast) SetHeight(height int) {
	t.Height = height
}

// Show displays a toast message with the provided text.
// This makes the toast visible and updates its message content.
func (t *Toast) Show(message string) {
	t.Message = message
	t.Visible = true
}

// Hide hides the toast notification and resets its state.
// This clears the message and sets the dismissed flag to false.
func (t *Toast) Hide() {
	t.Visible = false
	t.Message = ""
	t.Dismissed = false
}

// Update processes input messages and updates the toast state.
// Returns a boolean indicating whether the update resulted in any state change.
// Note: Enter keypresses are now handled by the App's Update method.
func (t *Toast) Update(msg tea.Msg) bool {
	if !t.Visible {
		return false
	}

	// Note: We no longer handle Enter keypresses here, as they're now 
	// handled directly in the App's Update method
	return false
}

// View renders the toast component as a styled notification box.
// The toast has a gold border, white text, and a brown-red background.
// If the toast is not visible, an empty string is returned.
func (t Toast) View() string {
	if !t.Visible {
		return ""
	}



	// Add a dismiss hint
	content := t.Message + "\n\nPress Enter to dismiss"
	

	return styles.ToastStyle.Render(content)
}
