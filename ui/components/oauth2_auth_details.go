package components

import (
	"fmt"

	"github.com/RAshkettle/LazyPost/ui/styles"
	tea "github.com/charmbracelet/bubbletea"
)

// OAuth2AuthDetailsComponent placeholder
type OAuth2AuthDetailsComponent struct {
	width  int
	height int
	active bool
}

func NewOAuth2AuthDetailsComponent() OAuth2AuthDetailsComponent {
	return OAuth2AuthDetailsComponent{}
}
func (c *OAuth2AuthDetailsComponent) SetActive(active bool) { c.active = active }
func (c *OAuth2AuthDetailsComponent) SetSize(width, height int) {
	c.width = width
	c.height = height
}
func (c OAuth2AuthDetailsComponent) Update(msg tea.Msg) tea.Cmd { return nil }
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
