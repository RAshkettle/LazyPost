// Package components provides UI components for the LazyPost application.
package components

import (
	"fmt"

	"github.com/atotto/clipboard" // Added for clipboard functionality
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// HeadersContainer represents a component for displaying HTTP response headers.
// It formats and displays header information. If active, it also shows a hint
// for copying the content to the clipboard using the 'y' key.
type HeadersContainer struct {
	Content    string // Content is the formatted header text to be displayed.
	rawContent string // rawContent stores the unformatted content for clipboard copying.
	Width      int    // Width is the width of the component in characters.
	Height     int    // Height is the height of thecomponent in characters.
	Active     bool   // Active indicates whether the component is currently focused and can respond to key presses like 'y'.
}

// NewHeadersContainer creates and initializes a new HeadersContainer.
// It starts with placeholder content and default dimensions.
func NewHeadersContainer() HeadersContainer {
	return HeadersContainer{
		Content:    "Response headers will be displayed here.",
		rawContent: "Response headers will be displayed here.", // Initialize rawContent
		Width:      0,
		Height:     0,
		Active:     false,
	}
}

// SetContent updates the header content to be displayed and the raw content for copying.
func (h *HeadersContainer) SetContent(content string) {
	h.Content = content
	h.rawContent = content // Store raw content
}

// SetWidth sets the rendering width for the HeadersContainer.
func (h *HeadersContainer) SetWidth(width int) {
	h.Width = width
}

// SetHeight sets the rendering height for the HeadersContainer.
func (h *HeadersContainer) SetHeight(height int) {
	h.Height = height
}

// SetActive sets the active state of the HeadersContainer.
// When active, it may display additional help text or respond to keys.
func (h *HeadersContainer) SetActive(active bool) {
	h.Active = active
}

// Update handles messages for the HeadersContainer.
// If the container is active and the 'y' key is pressed, it attempts to copy the raw content to the clipboard.
func (h *HeadersContainer) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if h.Active && msg.String() == "y" {
			err := clipboard.WriteAll(h.rawContent)
			if err != nil {
				// Optionally, send a message back to the app to show a toast
				fmt.Println("Error copying to clipboard:", err)
			}
			// Optionally, provide user feedback (e.g., via a toast message)
			return nil
		}
	}
	return nil
}

// View renders the HeadersContainer.
// It displays the formatted header content. If active, it appends a help message for copying.
// The content is rendered within a styled box, respecting the component's width and height.
// If width or height is zero or negative, it returns an empty string.
func (h HeadersContainer) View() string {
	if h.Width == 0 || h.Height == 0 {
		return ""
	}

	baseContent := h.Content

	if h.Active {
		helpStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFF00")). // Yellow color
			Align(lipgloss.Right).
			Bold(true).
			Width(h.Width - 4) // Account for padding of contentStyle and this style

		helpText := "'y' to copy"
		baseContent = lipgloss.JoinVertical(lipgloss.Left, baseContent, helpStyle.Render(helpText))
	}

	contentStyle := lipgloss.NewStyle().
		Width(h.Width).
		Height(h.Height).
		Padding(1, 2)

	return contentStyle.Render(baseContent)
}
