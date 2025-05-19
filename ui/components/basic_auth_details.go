package components

import (
	"fmt"

	"github.com/RAshkettle/LazyPost/ui/styles"
	tea "github.com/charmbracelet/bubbletea"
)

// BasicAuthDetailsComponent placeholder
type BasicAuthDetailsComponent struct {
	width  int
	height int
	active bool
}

func NewBasicAuthDetailsComponent() BasicAuthDetailsComponent {
	return BasicAuthDetailsComponent{}
}
func (c *BasicAuthDetailsComponent) SetActive(active bool) { c.active = active }
func (c *BasicAuthDetailsComponent) SetSize(width, height int) {
	c.width = width
	c.height = height
}
func (c BasicAuthDetailsComponent) Update(msg tea.Msg) tea.Cmd { return nil }
func (c BasicAuthDetailsComponent) View() string {
	if c.width <= 0 || c.height <= 0 {
		return ""
	}
	style := styles.DefaultTheme.BorderStyle.Copy().Width(c.width).Height(c.height)
	if c.active {
		style = styles.DefaultTheme.ActiveBorderStyle.Copy().Width(c.width).Height(c.height)
	}
	return style.Render(fmt.Sprintf("Basic Auth Details\nWidth: %d, Height: %d", c.width, c.height))
}
