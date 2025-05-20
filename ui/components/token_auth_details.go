// Package components defines various UI components for the LazyPost application.
package components

import (
	"github.com/RAshkettle/LazyPost/ui/styles"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// TokenAuthDetailsComponent holds the UI for Bearer Token input.
// It's specifically for Bearer tokens, but named generically as TokenAuthDetailsComponent
// for potential future reuse or extension if other simple token types arise.
type TokenAuthDetailsComponent struct {
	width      int
	height     int
	active     bool // active indicates whether the component is currently focused and accepting input.
	tokenInput textinput.Model // tokenInput is the text input field for the token.
	// No focusedField needed as there's only one input
}

// NewTokenAuthDetailsComponent creates a new instance of TokenAuthDetailsComponent.
// It initializes the text input field for the Bearer token.
func NewTokenAuthDetailsComponent() TokenAuthDetailsComponent {
	ti := textinput.New()
	ti.Placeholder = "Enter Bearer Token"
	ti.Prompt = "Token: "
	ti.Width = 30 // Default width, can be adjusted by SetSize
	// ti.Focus() // Focus will be handled by SetActive or parent Update

	return TokenAuthDetailsComponent{
		tokenInput: ti,
	}
}

// SetActive sets the active state of the component.
// When active, the token input field gains focus. When inactive, it loses focus.
func (c *TokenAuthDetailsComponent) SetActive(active bool) {
	c.active = active
	if active {
		c.tokenInput.Focus() // Focus the input when the component becomes active
	} else {
		c.tokenInput.Blur() // Blur the input when the component becomes inactive
	}
}

// SetSize sets the dimensions for the component's rendering area.
// This influences the overall width and height available for the component to render itself.
func (c *TokenAuthDetailsComponent) SetSize(width, height int) {
	c.width = width
	c.height = height
	// TODO: Consider adjusting c.tokenInput.Width based on c.width,
	// similar to BasicAuthDetailsComponent, if dynamic sizing is desired.
	// Example:
	// inputWidth := width - lipgloss.Width(c.tokenInput.Prompt) - styles.DefaultTheme.ActiveInputStyle.GetHorizontalPadding() - 2 // Rough estimate
	// if inputWidth < 10 { inputWidth = 10 }
	// c.tokenInput.Width = inputWidth
}

// Update handles messages and updates the component's state.
// It only processes messages and updates the token input field if the component is active.
// It returns a tea.Cmd, which might be produced by the text input field's update.
func (c *TokenAuthDetailsComponent) Update(msg tea.Msg) tea.Cmd {
	if !c.active {
		return nil
	}

	var cmd tea.Cmd
	c.tokenInput, cmd = c.tokenInput.Update(msg)
	return cmd
}

// View renders the TokenAuthDetailsComponent.
// It displays the token input field, styled according to its active and focused state,
// within a bordered box. The border style also reflects the component's active state.
// If width or height is zero or negative, it returns an empty string.
func (c TokenAuthDetailsComponent) View() string {
	if c.width <= 0 || c.height <= 0 {
		return ""
	}

	tokenView := c.tokenInput.View()
	var styledTokenView string

	if c.active && c.tokenInput.Focused() { // Check both component active and input focused for active style
		styledTokenView = styles.DefaultTheme.ActiveInputStyle.Render(tokenView)
	} else {
		styledTokenView = styles.DefaultTheme.InactiveInputStyle.Render(tokenView)
	}


	contentWithHelp := lipgloss.JoinVertical(
		lipgloss.Left,
		styledTokenView,

	)

	// Use a general border style, active if the component itself is active.
	// The input field's active/inactive style is handled above.
	componentBorderStyle := styles.DefaultTheme.BorderStyle
	if c.active {
		componentBorderStyle = styles.DefaultTheme.ActiveBorderStyle
	}

	innerWidth := c.width - componentBorderStyle.GetHorizontalFrameSize()
	innerHeight := c.height - componentBorderStyle.GetVerticalFrameSize()
	if innerWidth < 0 {
		innerWidth = 0
	}
	if innerHeight < 0 {
		innerHeight = 0
	}

	return componentBorderStyle.Width(c.width).Height(c.height).Render(
		lipgloss.NewStyle().Width(innerWidth).Height(innerHeight).Render(contentWithHelp),
	)
}

// GetToken returns the current value of the token input field.
func (c *TokenAuthDetailsComponent) GetToken() string {
	return c.tokenInput.Value()
}
