package components

import (
	"fmt"

	"github.com/RAshkettle/LazyPost/ui/styles"
	tea "github.com/charmbracelet/bubbletea"
)

// JWTAuthDetailsComponent placeholder
type JWTAuthDetailsComponent struct {
	width  int
	height int
	active bool
}

func NewJWTAuthDetailsComponent() JWTAuthDetailsComponent {
	return JWTAuthDetailsComponent{}
}
func (c *JWTAuthDetailsComponent) SetActive(active bool) { c.active = active }
func (c *JWTAuthDetailsComponent) SetSize(width, height int) {
	c.width = width
	c.height = height
}
func (c JWTAuthDetailsComponent) Update(msg tea.Msg) tea.Cmd { return nil }
func (c JWTAuthDetailsComponent) View() string {
	if c.width <= 0 || c.height <= 0 {
		return ""
	}
	style := styles.DefaultTheme.BorderStyle.Width(c.width).Height(c.height)
	if c.active {
		style = styles.DefaultTheme.ActiveBorderStyle.Width(c.width).Height(c.height)
	}
	return style.Render(fmt.Sprintf("JWT Auth Details"))
}
