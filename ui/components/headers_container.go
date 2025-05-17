// Package components provides UI components for the LazyPost application.
package components

import (
	"fmt"

	"github.com/atotto/clipboard" // Added for clipboard functionality
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// HeadersContainer represents a component for displaying HTTP response headers.
// It formats and displays header information in a visually appealing way.
type HeadersContainer struct {
	Content    string // The content to display (formatted header text)
	rawContent string // Store raw content for copying
	Width      int    // Width of the component in characters
	Height     int    // Height of the component in characters
	Active     bool   // Whether the component is currently active/focused
}

// NewHeadersContainer creates a new headers container.
func NewHeadersContainer() HeadersContainer {
	return HeadersContainer{
		Content:    "Response headers will be displayed here.",
		rawContent: "Response headers will be displayed here.", // Initialize rawContent
		Width:      0,
		Height:     0,
		Active:     false,
	}
}

// SetContent updates the header content to display.
func (h *HeadersContainer) SetContent(content string) {
	h.Content = content
	h.rawContent = content // Store raw content
}

// SetWidth sets the width of the component in characters.
func (h *HeadersContainer) SetWidth(width int) {
	h.Width = width
}

// SetHeight sets the height of the component in characters.
func (h *HeadersContainer) SetHeight(height int) {
	h.Height = height
}

// SetActive sets the active state of the component.
func (h *HeadersContainer) SetActive(active bool) {
	h.Active = active
}

// Update handles any messages to update the component state.
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

// View renders the headers container.
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
