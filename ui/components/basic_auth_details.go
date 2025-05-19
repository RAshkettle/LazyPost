package components

import (
	"github.com/RAshkettle/LazyPost/ui/styles"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	basicAuthUsernameField = 0
	basicAuthPasswordField = 1
)

// BasicAuthDetailsComponent holds the UI for Basic Auth input fields.
type BasicAuthDetailsComponent struct {
	width  int
	height int
	active bool // Is the component itself active (e.g., visible and interactable in a tab)

	usernameInput textinput.Model
	passwordInput textinput.Model
	focusedField  int // 0 for username, 1 for password
}

// NewBasicAuthDetailsComponent creates a new instance of BasicAuthDetailsComponent.
func NewBasicAuthDetailsComponent() BasicAuthDetailsComponent {
	username := textinput.New()
	username.Placeholder = "Enter username"
	username.Prompt = "Username: "
	username.Width = 30 // Width of the text area
	// username.Focus() // Initial focus will be handled by Update or SetActive

	password := textinput.New()
	password.Placeholder = "Enter password"
	password.Prompt = "Password: "
	password.EchoMode = textinput.EchoPassword
	password.EchoCharacter = '*'
	password.Width = 30 // Width of the text area

	return BasicAuthDetailsComponent{
		usernameInput: username,
		passwordInput: password,
		focusedField:  basicAuthUsernameField, // Default to username when component is first created/activated
	}
}

// SetActive sets the active state of the component.
// When active, it can receive focus and process updates.
func (c *BasicAuthDetailsComponent) SetActive(active bool) {
	c.active = active
	if !active {
		c.usernameInput.Blur()
		c.passwordInput.Blur()
		// c.focusedField = -1 // Or keep last focused to restore later if desired
	} else {
		// When the component becomes active, ensure the default field is ready for focus.
		// Actual Focus() command is usually sent via Update to manage cursor blinking.
		if c.focusedField == basicAuthUsernameField {
			c.usernameInput.Focus()
			c.passwordInput.Blur()
		} else if c.focusedField == basicAuthPasswordField {
			c.passwordInput.Focus()
			c.usernameInput.Blur()
		} else {
			// Default to focusing username if no specific field was pre-focused
			c.focusedField = basicAuthUsernameField
			c.usernameInput.Focus()
			c.passwordInput.Blur()
		}
	}
}

// SetSize sets the dimensions for the component's rendering area.
func (c *BasicAuthDetailsComponent) SetSize(width, height int) {
	c.width = width
	c.height = height
}

// Update handles messages and updates the component's state.
func (c *BasicAuthDetailsComponent) Update(msg tea.Msg) tea.Cmd {
	if !c.active {
		return nil
	}

	var cmds []tea.Cmd
	var cmd tea.Cmd

	// Ensure the correct input is focused if component is active but inputs aren't (e.g. first activation)
	// This logic is now partly in SetActive, but Update is where commands are returned.
	if c.active {
		if c.focusedField == basicAuthUsernameField && !c.usernameInput.Focused() {
			cmds = append(cmds, c.usernameInput.Focus()) // Ensure blink command is initiated
		}
		// No need to explicitly focus password here as SetActive or key-handling will manage it.
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab", "down":
			if c.focusedField == basicAuthUsernameField {
				c.usernameInput.Blur()
				c.focusedField = basicAuthPasswordField
				cmds = append(cmds, c.passwordInput.Focus())
			} else {
				c.passwordInput.Blur()
				c.focusedField = basicAuthUsernameField
				cmds = append(cmds, c.usernameInput.Focus())
			}
			return tea.Batch(cmds...)

		case "shift+tab", "up":
			if c.focusedField == basicAuthPasswordField {
				c.passwordInput.Blur()
				c.focusedField = basicAuthUsernameField
				cmds = append(cmds, c.usernameInput.Focus())
			} else {
				c.usernameInput.Blur()
				c.focusedField = basicAuthPasswordField
				cmds = append(cmds, c.passwordInput.Focus())
			}
			return tea.Batch(cmds...)
		}
	}

	// Delegate message to the currently focused input field
	if c.focusedField == basicAuthUsernameField {
		c.usernameInput, cmd = c.usernameInput.Update(msg)
		cmds = append(cmds, cmd)
	} else if c.focusedField == basicAuthPasswordField {
		c.passwordInput, cmd = c.passwordInput.Update(msg)
		cmds = append(cmds, cmd)
	}

	return tea.Batch(cmds...)
}

// View renders the BasicAuthDetailsComponent.
func (c BasicAuthDetailsComponent) View() string {
	if c.width <= 0 || c.height <= 0 {
		return "" // Not enough space to render
	}

	// Get the view strings from the text input components
	usernameView := c.usernameInput.View() // This includes the prompt
	passwordView := c.passwordInput.View() // This includes the prompt

	var styledUsernameView, styledPasswordView string

	// Apply active/inactive styling based on which field is focused
	if c.focusedField == basicAuthUsernameField {
		styledUsernameView = styles.DefaultTheme.ActiveInputStyle.Render(usernameView)
		styledPasswordView = styles.DefaultTheme.InactiveInputStyle.Render(passwordView)
	} else if c.focusedField == basicAuthPasswordField {
		styledUsernameView = styles.DefaultTheme.InactiveInputStyle.Render(usernameView)
		styledPasswordView = styles.DefaultTheme.ActiveInputStyle.Render(passwordView)
	} else {
		// Neither is focused (e.g., component itself is not active, though View might be called)
		styledUsernameView = styles.DefaultTheme.InactiveInputStyle.Render(usernameView)
		styledPasswordView = styles.DefaultTheme.InactiveInputStyle.Render(passwordView)
	}

	// Join the styled input fields vertically
	inputsView := lipgloss.JoinVertical(lipgloss.Left, styledUsernameView, styledPasswordView)

	// Help text
	helpTextView := styles.DefaultTheme.HelpTextStyle.Foreground(styles.BrightYellow).Render("Tab/Shift+Tab or Up/Down to navigate fields.")

	// Combine inputs and help text
	contentWithHelp := lipgloss.JoinVertical(
		lipgloss.Left,
		inputsView,
		helpTextView,
	)

	// Determine the overall border style for the component
	componentBorderStyle := styles.DefaultTheme.BorderStyle
	if c.active { // c.active refers to the component's overall active state
		componentBorderStyle = styles.DefaultTheme.ActiveBorderStyle
	}

	// Ensure the component respects its given width and height, applying padding from the border style.
	// The content (inputsView) will be placed inside this border.
	// We need to account for the border's frame size.
	innerWidth := c.width - componentBorderStyle.GetHorizontalFrameSize()
	innerHeight := c.height - componentBorderStyle.GetVerticalFrameSize()

	if innerWidth < 0 {
		innerWidth = 0
	}
	if innerHeight < 0 {
		innerHeight = 0
	}

	// If the content is taller than the available inner height, it might be truncated by Render.
	// Or, we could choose to make the component scrollable if needed in the future.
	// For now, we just render it into the available space.
	finalView := componentBorderStyle.Width(c.width).Height(c.height).Render(
		lipgloss.NewStyle().Width(innerWidth).Height(innerHeight).Render(contentWithHelp),
	)

	return finalView
}

// GetValues returns the current username and password. (Placeholder)
func (c *BasicAuthDetailsComponent) GetValues() (username string, password string) {
	return c.usernameInput.Value(), c.passwordInput.Value()
}
