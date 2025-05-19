package components

import (
	"fmt"
	"strings"

	"github.com/RAshkettle/LazyPost/ui/styles"
	"github.com/atotto/clipboard" // Added for clipboard functionality
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// BodyContainer represents a scrollable component for displaying HTTP response bodies.
// It uses a viewport for scrolling through large content.
type BodyContainer struct {
	Viewport   viewport.Model // Viewport for scrollable content
	rawContent string         // Store raw content for copying
	Width      int            // Width of the component in characters
	Height     int            // Height of the component in characters
	Active     bool           // Whether the component is currently active/focused
}

// NewBodyContainer creates a new body container with a scrollable viewport.
func NewBodyContainer() BodyContainer {
	vp := viewport.New(0, 0)
	vp.SetContent("Response body will be displayed here.")

	// Configure viewport keybindings
	vp.KeyMap = viewport.KeyMap{
		Up:           key.NewBinding(key.WithKeys("up", "k")),
		Down:         key.NewBinding(key.WithKeys("down", "j")),
		Left:         key.NewBinding(key.WithKeys("left", "h")),
		Right:        key.NewBinding(key.WithKeys("right", "l")),
		PageUp:       key.NewBinding(key.WithKeys("pgup")),
		PageDown:     key.NewBinding(key.WithKeys("pgdown")),
		HalfPageUp:   key.NewBinding(key.WithKeys("ctrl+u")),
		HalfPageDown: key.NewBinding(key.WithKeys("ctrl+d")),
	}

	return BodyContainer{
		Viewport:   vp,
		rawContent: "Response body will be displayed here.", // Initialize rawContent
		Width:      0,
		Height:     0,
		Active:     false,
	}
}

// SetContent updates the body content to display and resets scroll position.
func (b *BodyContainer) SetContent(content string) {
	b.rawContent = content // Store raw content
	// Make sure we have valid dimensions before setting content
	if b.Width > 0 && b.Height > 0 {
		// Store the content and ensure the viewport is properly sized
		effectiveWidth := b.Width - 4  // Account for 2 chars padding on both sides plus border
		b.Viewport.Width = b.Width - 2 // Account for border padding
		b.Viewport.Height = b.Height - 2

		// Apply text wrapping to ensure content fits within the viewport width
		wrappedContent := wrapText(content, effectiveWidth)

		// Set the wrapped content and reset the scroll position

		b.Viewport.SetContent(wrappedContent)
		b.Viewport.GotoTop()
	} else {
		// Just store the content for now, the viewport will be updated when dimensions are set
		b.Viewport.SetContent(content) // Keep this for initial placeholder
	}
}

// wrapText wraps the text to ensure it fits within the specified width.
// This ensures all content is visible and properly formatted within the viewport.
func wrapText(content string, width int) string {
	if width <= 0 {
		return content
	}

	var result strings.Builder
	lines := strings.Split(content, "\n")

	for i, line := range lines {
		if len(line) <= width {
			result.WriteString(line)
		} else {
			// Wrap lines longer than width
			for j := 0; j < len(line); j += width {
				end := j + width
				end = min(end, len(line))

				result.WriteString(line[j:end])
				if end < len(line) {
					result.WriteString("\n")
				}
			}
		}
		// Add newline after each original line except the last one
		if i < len(lines)-1 {
			result.WriteString("\n")
		}
	}

	return result.String()
}

// SetWidth sets the width of the component in characters.
func (b *BodyContainer) SetWidth(width int) {
	b.Width = width
	if width > 2 { // Only set reasonable dimensions
		b.Viewport.Width = width - 2 // Account for border padding

		// Re-wrap content when width changes if we have content
		content := b.Viewport.View()
		if content != "" && content != "Response body will be displayed here." {
			effectiveWidth := width - 6 // Account for 2 chars padding on both sides plus border
			wrappedContent := wrapText(content, effectiveWidth)
			b.Viewport.SetContent(wrappedContent)
		}
	}
}

// SetHeight sets the height of the component in characters.
// Also adds 3 extra rows to extend the container's height.
func (b *BodyContainer) SetHeight(height int) {
	// Add 3 extra rows to extend the height
	b.Height = height + 3
	if height > 2 { // Only set reasonable dimensions
		b.Viewport.Height = height + 1 // Account for border padding and extension
	}
}

// SetActive sets the active state of the component.
func (b *BodyContainer) SetActive(active bool) {
	b.Active = active
}

// Update handles viewport navigation and other messages.
func (b *BodyContainer) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	// Always update the viewport on window resize regardless of active state
	if _, ok := msg.(tea.WindowSizeMsg); ok {
		if b.Width > 2 && b.Height > 2 {
			// Update dimensions
			b.Viewport.Width = b.Width - 2
			b.Viewport.Height = b.Height - 2

			// Re-wrap content based on new width
			origContent := b.Viewport.View()
			if origContent != "" && origContent != "Response body will be displayed here." {
				// Save current scroll position
				currentPosition := b.Viewport.YOffset

				// Re-wrap text for new dimensions
				effectiveWidth := b.Width - 6 // Account for 2 chars padding on both sides plus border
				wrappedContent := wrapText(origContent, effectiveWidth)
				b.Viewport.SetContent(wrappedContent)

				// Try to restore scroll position (approximately)
				if currentPosition > 0 && currentPosition < b.Viewport.TotalLineCount() {
					b.Viewport.YOffset = currentPosition
				}
			}
		}
	}

	// Only handle key navigation when active
	if !b.Active {
		return nil
	}

	switch msgType := msg.(type) {
	case tea.KeyMsg:
		switch msgType.String() {
		case "y":
			if b.Active {
				err := clipboard.WriteAll(b.rawContent)
				if err != nil {
					// Optionally, you could send a message back to the app to show a toast
					// For now, just print to stderr or log
					fmt.Println("Error copying to clipboard:", err)
				}
				// We might want to provide feedback to the user, e.g., a short message
				// This could be a new tea.Msg that the main app handles.
				// For simplicity, returning nil for now.
				return nil
			}
		case "home":
			// Jump to the top of the content
			b.Viewport.GotoTop()
			return nil
		case "end":
			// Jump to the bottom of the content
			b.Viewport.GotoBottom()
			return nil
		case "up", "k", "down", "j", "pgup", "pgdn", "ctrl+u", "ctrl+d":
			// Let viewport handle other navigation keys
			b.Viewport, cmd = b.Viewport.Update(msg)
			cmds = append(cmds, cmd)
			return tea.Batch(cmds...)
		}
	}

	return tea.Batch(cmds...)
}

// addPadding adds the specified number of spaces to the left and right of each line.
func addPadding(content string, paddingSize int) string {
	if paddingSize <= 0 {
		return content
	}

	padding := strings.Repeat(" ", paddingSize)
	var result strings.Builder

	lines := strings.Split(content, "\n")
	for i, line := range lines {
		result.WriteString(padding + line + padding)
		if i < len(lines)-1 {
			result.WriteString("\n")
		}
	}

	return result.String()
}

// View renders the body container with scrolling.
func (b BodyContainer) View() string {
	if b.Width == 0 || b.Height == 0 {
		return ""
	}

	// Get viewport content and add padding
	content := addPadding(b.Viewport.View(), 2)

	// Show scrolling help text when body is active
	if b.Active {
		helpStyle := styles.HelpStyle
		helpStyle.Width(b.Width - 2)
		// Show helpful scrolling indicators
		var helpParts []string

		// Check if content needs scrolling
		atBottom := b.Viewport.AtBottom()

		// If we're not at the top or not at the bottom, content is scrollable
		if !atBottom || b.Viewport.YOffset > 0 {
			currLine := fmt.Sprintf("Line %d", b.Viewport.YOffset+1)
			helpParts = append(helpParts, "↑/↓ to scroll • PgUp/PgDn for faster scrolling • "+currLine)

			// Add indicator if we're not at the bottom
			if !atBottom {
				helpParts[len(helpParts)-1] += " (more ↓)"
			}
		}

		helpParts = append(helpParts, "'y' to copy")

		helpText := strings.Join(helpParts, " • ")

		if helpText != "" {
			content = lipgloss.JoinVertical(lipgloss.Left, content, helpStyle.Render(helpText))
		}
	}

	return content
}
