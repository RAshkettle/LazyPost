package components

import (
	"github.com/RAshkettle/LazyPost/ui/styles"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// TokenAuthDetailsComponent holds the UI for Bearer Token input.
// It's specifically for Bearer tokens, but named generically as TokenAuthDetailsComponent.
type TokenAuthDetailsComponent struct {
	width      int
	height     int
	active     bool // Is the component itself active
	tokenInput textinput.Model
	// No focusedField needed as there's only one input
}

// NewTokenAuthDetailsComponent creates a new instance of TokenAuthDetailsComponent.
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
func (c *TokenAuthDetailsComponent) SetActive(active bool) {
	c.active = active
	if active {
		c.tokenInput.Focus() // Focus the input when the component becomes active
	} else {
		c.tokenInput.Blur() // Blur the input when the component becomes inactive
	}
}

// SetSize sets the dimensions for the component's rendering area.
func (c *TokenAuthDetailsComponent) SetSize(width, height int) {
	c.width = width
	c.height = height
	// Adjust input width based on component width, similar to BasicAuthDetailsComponent if needed
	// For now, let's assume a fixed prompt and the input takes available space or a max width.
	// If we want it to be dynamic like BasicAuth, we'd do:
	// inputWidth := width - lipgloss.Width(c.tokenInput.Prompt) - styles.DefaultTheme.ActiveInputStyle.GetHorizontalPadding() - 2 // Rough estimate
	// if inputWidth < 10 { inputWidth = 10 }
	// c.tokenInput.Width = inputWidth
}

// Update handles messages and updates the component's state.
func (c *TokenAuthDetailsComponent) Update(msg tea.Msg) tea.Cmd {
	if !c.active {
		return nil
	}

	var cmd tea.Cmd
	c.tokenInput, cmd = c.tokenInput.Update(msg)
	return cmd
}

// View renders the TokenAuthDetailsComponent.
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

	helpTextView := styles.DefaultTheme.HelpTextStyle.Foreground(styles.BrightYellow).Render("Enter your Bearer token.")

	contentWithHelp := lipgloss.JoinVertical(
		lipgloss.Left,
		styledTokenView,
		helpTextView,
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
