package components

import (
	"fmt"

	"github.com/RAshkettle/LazyPost/ui/styles"
	tea "github.com/charmbracelet/bubbletea"
)

// TokenAuthDetailsComponent placeholder (for Bearer, JWT)
type TokenAuthDetailsComponent struct {
	width  int
	height int
	active bool
}

func NewTokenAuthDetailsComponent() TokenAuthDetailsComponent {
	return TokenAuthDetailsComponent{}
}
func (c *TokenAuthDetailsComponent) SetActive(active bool) { c.active = active }
func (c *TokenAuthDetailsComponent) SetSize(width, height int) {
	c.width = width
	c.height = height
}
func (c TokenAuthDetailsComponent) Update(msg tea.Msg) tea.Cmd { return nil }
func (c TokenAuthDetailsComponent) View() string {
	if c.width <= 0 || c.height <= 0 {
		return ""
	}
	style := styles.DefaultTheme.BorderStyle.Copy().Width(c.width).Height(c.height)
	if c.active {
		style = styles.DefaultTheme.ActiveBorderStyle.Copy().Width(c.width).Height(c.height)
	}
	return style.Render(fmt.Sprintf("Token Auth Details\n(Bearer/JWT)\nWidth: %d, Height: %d", c.width, c.height))
}
