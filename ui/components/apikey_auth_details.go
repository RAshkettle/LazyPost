package components

import (
	"fmt"

	"github.com/RAshkettle/LazyPost/ui/styles"
	tea "github.com/charmbracelet/bubbletea"
)

// APIKeyAuthDetailsComponent placeholder
type APIKeyAuthDetailsComponent struct {
	width  int
	height int
	active bool
}

func NewAPIKeyAuthDetailsComponent() APIKeyAuthDetailsComponent {
	return APIKeyAuthDetailsComponent{}
}
func (c *APIKeyAuthDetailsComponent) SetActive(active bool) { c.active = active }
func (c *APIKeyAuthDetailsComponent) SetSize(width, height int) {
	c.width = width
	c.height = height
}
func (c APIKeyAuthDetailsComponent) Update(msg tea.Msg) tea.Cmd { return nil }
func (c APIKeyAuthDetailsComponent) View() string {
	if c.width <= 0 || c.height <= 0 {
		return ""
	}
	style := styles.DefaultTheme.BorderStyle.Copy().Width(c.width).Height(c.height)
	if c.active {
		style = styles.DefaultTheme.ActiveBorderStyle.Copy().Width(c.width).Height(c.height)
	}
	return style.Render(fmt.Sprintf("API Key Auth Details\nWidth: %d, Height: %d", c.width, c.height))
}
