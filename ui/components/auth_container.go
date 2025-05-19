package components

import (
	"github.com/RAshkettle/LazyPost/ui/styles"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// AuthContainer represents a component for managing authentication details.
// For now, it's a styled empty container.
type AuthContainer struct {
	Width          int    // Width of the component in characters
	Height         int    // Height of the component in characters
	Active         bool   // Whether the component is currently active/focused
	containerStyle lipgloss.Style
	titleStyle     lipgloss.Style
}

// NewAuthContainer creates a new, empty AuthContainer.
func NewAuthContainer() AuthContainer {
	return AuthContainer{
		Width:          0,
		Height:         0,
		Active:         false,
		containerStyle: styles.DefaultTheme.InputContainerStyle.Copy(), // Use a copy to avoid modifying the global theme
		titleStyle:     styles.DefaultTheme.TitleStyle.Copy(),
	}
}

// SetWidth sets the width of the component in characters.
func (ac *AuthContainer) SetWidth(width int) {
	ac.Width = width
}

// SetHeight sets the height of the component in characters.
func (ac *AuthContainer) SetHeight(height int) {
	ac.Height = height
}

// SetActive sets the active state of the component.
func (ac *AuthContainer) SetActive(active bool) {
	ac.Active = active
}

// Update handles any messages to update the component state.
// For now, it does nothing.
func (ac *AuthContainer) Update(msg tea.Msg) tea.Cmd {
	return nil
}

// View renders the AuthContainer.
func (ac AuthContainer) View() string {
	if ac.Width == 0 || ac.Height == 0 {
		return ""
	}

	containerTitle := ac.titleStyle.Render("Authentication")

	// Apply active/inactive border style
	currentContainerStyle := ac.containerStyle
	if ac.Active {
		currentContainerStyle = styles.DefaultTheme.ActiveInputStyle // Or ActiveBorderStyle if preferred
	} else {
		currentContainerStyle = styles.DefaultTheme.InactiveInputStyle // Or BorderStyle
	}

	styledContainer := currentContainerStyle.
		Width(ac.Width).
		Height(ac.Height).
		Padding(0,1) // Match HeadersContainer padding if desired, or keep specific

	// For now, the container is empty except for its title (if rendered inside)
	// To render title inside, we need to manage inner content height.
	// Let's render the title as part of the border for now, similar to other inputs.

	// Calculate inner dimensions for content placement if needed later
	// innerWidth := ac.Width - styledContainer.GetHorizontalFrameSize()
	// innerHeight := ac.Height - styledContainer.GetVerticalFrameSize()



	// If we want the title to be part of the border like text inputs:
	// The lipgloss.Style.Title() method is not standard. Titles are usually rendered
	// as part of the content or drawn separately above the component.
	// For consistency with HeadersInput, let's assume the title is rendered above or as part of the content.
	// For an empty container with a border title, we'd typically do this:

	// Option 1: Title as part of the border (if the style supports it directly, lipgloss doesn't quite do this like bubble textinput)
	// Option 2: Render title, then render the empty box below it.
	// Option 3: Render the box, and place the title string at the top of its content.

	// Let's go with rendering the title string at the top of the content area.
	var finalContent string
	if ac.Height > 0 { // Ensure there's space for the title
		finalContent = containerTitle
		// Fill remaining lines if any
		for i := 1; i < ac.Height - styledContainer.GetVerticalPadding() - styledContainer.GetVerticalFrameSize() ; i++ {
			finalContent = lipgloss.JoinVertical(lipgloss.Left, finalContent, "")
		}
	} else {
		finalContent = "" // Not enough space for title
	}

	return styledContainer.Render(finalContent)
}
