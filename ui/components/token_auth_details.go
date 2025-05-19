package components

import (
	"fmt"

	"github.com/RAshkettle/LazyPost/ui/styles"
	tea "github.com/charmbracelet/bubbletea"
)

// TokenAuthDetailsComponent placeholder (for Bearer)
// TokenAuthDetailsComponent ...
type TokenAuthDetailsComponent struct {
	width  int
	height int
	active bool
}

// NewTokenAuthDetailsComponent ...
func NewTokenAuthDetailsComponent() TokenAuthDetailsComponent {
	return TokenAuthDetailsComponent{}
}

// SetActive ...
func (c *TokenAuthDetailsComponent) SetActive(active bool) { c.active = active }

// SetSize ...
func (c *TokenAuthDetailsComponent) SetSize(width, height int) {
	c.width = width
	c.height = height
}

// Update ...
func (c TokenAuthDetailsComponent) Update(msg tea.Msg) tea.Cmd { return nil }

// View ...
func (c TokenAuthDetailsComponent) View() string {
	if c.width <= 0 || c.height <= 0 {
		return ""
	}
	style := styles.DefaultTheme.BorderStyle.Width(c.width).Height(c.height)
	if c.active {
		style = styles.DefaultTheme.ActiveBorderStyle.Width(c.width).Height(c.height)
	}
	return style.Render(fmt.Sprintf("Bearer Auth Details"))
}
