// Package components defines various UI components for the LazyPost application.
package components

import (
	"fmt"

	"github.com/RAshkettle/LazyPost/ui/styles"
	tea "github.com/charmbracelet/bubbletea"
)

// OAuth2AuthDetailsComponent is a placeholder for OAuth2 authentication details UI.
// It currently displays a simple message and will be implemented with actual
// input fields and logic for OAuth2 flows in the future.
type OAuth2AuthDetailsComponent struct {
	width  int  // width is the width of the component.
	height int  // height is the height of the component.
	active bool // active indicates whether the component is currently focused.
}

// NewOAuth2AuthDetailsComponent creates a new instance of OAuth2AuthDetailsComponent.
func NewOAuth2AuthDetailsComponent() OAuth2AuthDetailsComponent {
	return OAuth2AuthDetailsComponent{}
}

// SetActive sets the active state of the component.
func (c *OAuth2AuthDetailsComponent) SetActive(active bool) { c.active = active }

// SetSize sets the dimensions for the component's rendering area.
func (c *OAuth2AuthDetailsComponent) SetSize(width, height int) {
	c.width = width
	c.height = height
}

// Update handles messages and updates the component's state.
// Currently, it's a no-op as the component is a placeholder.
func (c OAuth2AuthDetailsComponent) Update(msg tea.Msg) tea.Cmd { return nil }

// View renders the OAuth2AuthDetailsComponent.
// It displays a placeholder message within a styled border.
// If width or height is zero or negative, it returns an empty string.
func (c OAuth2AuthDetailsComponent) View() string {
	if c.width <= 0 || c.height <= 0 {
		return ""
	}
	style := styles.DefaultTheme.BorderStyle.Width(c.width).Height(c.height)
	if c.active {
		style = styles.DefaultTheme.ActiveBorderStyle.Width(c.width).Height(c.height)
	}
	return style.Render(fmt.Sprintf("OAuth2 Auth Details"))
}
