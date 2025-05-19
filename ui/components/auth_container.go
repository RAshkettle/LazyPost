package components

import (
	"fmt"

	"github.com/RAshkettle/LazyPost/ui/styles"
	"github.com/charmbracelet/bubbles/key" // For placeholder keymap
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var authTypeOptions = []string{"None", "Basic", "Bearer", "JWT", "OAuth2", "API Key"} // Added "None"

// AuthSelectorKeyMap defines keybindings for the AuthSelector.
// (Placeholder for future interactivity)
type AuthSelectorKeyMap struct {
	Open   key.Binding
	Close  key.Binding
	Next   key.Binding
	Prev   key.Binding
	Select key.Binding
}

// DefaultAuthSelectorKeyMap provides default keybindings.
// (Placeholder for future interactivity)
var DefaultAuthSelectorKeyMap = AuthSelectorKeyMap{
	Open:   key.NewBinding(key.WithKeys("enter", " ")),
	Close:  key.NewBinding(key.WithKeys("esc")),
	Next:   key.NewBinding(key.WithKeys("down", "j")),
	Prev:   key.NewBinding(key.WithKeys("up", "k")),
	Select: key.NewBinding(key.WithKeys("enter")),
}

// AuthSelector manages the dropdown for authentication types.
type AuthSelector struct {
	options            []string
	selectedIndex      int
	highlightedIndex   int // Used when isOpen is true
	isOpen             bool
	active             bool
	width              int
	activeStyle        lipgloss.Style
	inactiveStyle      lipgloss.Style
	dropdownTextStyle  lipgloss.Style
	dropdownArrowStyle lipgloss.Style
	dropdownItemStyle         lipgloss.Style // Style for items when dropdown is open
	keymap             AuthSelectorKeyMap
}

// NewAuthSelector creates a new AuthSelector.
func NewAuthSelector() AuthSelector {
	return AuthSelector{
		options:            authTypeOptions,
		selectedIndex:      0,
		highlightedIndex:   0, // Initialize highlightedIndex
		isOpen:             false,
		active:             false,
		activeStyle:        styles.DefaultTheme.ActiveInputStyle.Copy(),
		inactiveStyle:      styles.DefaultTheme.InactiveInputStyle.Copy(),
		dropdownTextStyle:  styles.DefaultTheme.DropdownTextStyle.Copy(),
		dropdownArrowStyle: styles.DefaultTheme.DropdownArrowStyle.Copy(),
		dropdownItemStyle: styles.DefaultTheme.DropdownItemStyle.Copy(),             // Initialize new style
		keymap:             DefaultAuthSelectorKeyMap,
	}
}

// View renders the AuthSelector.
// It now handles both closed and open states.
func (as AuthSelector) View() string {
	var currentStyle lipgloss.Style
	if as.active {
		currentStyle = as.activeStyle
	} else {
		currentStyle = as.inactiveStyle
	}

	// Calculate effective width for content inside the style's padding
	effectiveContentWidth := as.width - currentStyle.GetHorizontalPadding()
	if effectiveContentWidth < 0 {
		effectiveContentWidth = 0
	}

	if !as.isOpen {
		selectedOptionText := as.dropdownTextStyle.Render(as.options[as.selectedIndex])
		arrow := as.dropdownArrowStyle.Render(" ▼")
		optionStrPaddedWidth := effectiveContentWidth - lipgloss.Width(arrow) - 2 // -2 for spaces around text and arrow
		if optionStrPaddedWidth < 0 {
			optionStrPaddedWidth = 0
		}
		viewString := fmt.Sprintf(" %-*s%s ", optionStrPaddedWidth, selectedOptionText, arrow)
		return currentStyle.Copy().Width(as.width).Render(viewString)
	}

	// Render open state
	var items []string
	for i, optionText := range as.options {
		var renderedText string
		if i == as.highlightedIndex {
			displayText := "▶ " + optionText
			renderedText = styles.DefaultTheme.SelectedItemStyle.Render(displayText)
		} else {
			displayText := "  " + optionText
			renderedText = as.dropdownTextStyle.Render(displayText)
		}
		// Each item line uses dropdownItemStyle for padding and width, then renders the specific text style inside.
		line := as.dropdownItemStyle.Copy().Width(effectiveContentWidth).Render(renderedText)
		items = append(items, line)
	}

	dropdownContent := lipgloss.JoinVertical(lipgloss.Left, items...)
	return currentStyle.Copy().Width(as.width).Render(dropdownContent)
}

// SetWidth sets the width of the AuthSelector.
func (as *AuthSelector) SetWidth(width int) {
	as.width = width
}

// SetActive sets the active state of the AuthSelector.
func (as *AuthSelector) SetActive(active bool) {
	as.active = active
}

// Update handles messages for the AuthSelector.
func (as *AuthSelector) Update(msg tea.Msg) tea.Cmd {
	if !as.active { // Only process messages if active and the component is supposed to be interactive
		return nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if as.isOpen { // Handle keys when dropdown is open
			switch {
			case key.Matches(msg, as.keymap.Close):
				as.isOpen = false
				as.highlightedIndex = as.selectedIndex // Reset highlight to current selection
				return nil
			case key.Matches(msg, as.keymap.Next):
				as.highlightedIndex = (as.highlightedIndex + 1) % len(as.options)
				return nil
			case key.Matches(msg, as.keymap.Prev):
				as.highlightedIndex = (as.highlightedIndex - 1 + len(as.options)) % len(as.options)
				return nil
			case key.Matches(msg, as.keymap.Select):
				as.selectedIndex = as.highlightedIndex
				as.isOpen = false
				// TODO: Potentially send a message indicating selection changed, if other components need to react immediately.
				// For now, the change will be reflected in the next View() call.
				return nil
			}
		} else { // Handle keys when dropdown is closed
			switch {
			case key.Matches(msg, as.keymap.Open):
				as.isOpen = true
				as.highlightedIndex = as.selectedIndex // Start highlighting from current selection
				return nil
			}
		}
	}
	return nil
}

// AuthContainer encapsulates the AuthSelector and manages the "Authentication" section.
type AuthContainer struct {
	Width          int    // Width of the component in characters
	Height         int    // Height of the component in characters
	Active         bool   // If the container or its children are focused
	authSelector   AuthSelector
	titleStyle     lipgloss.Style
	// Removed ac.containerStyle, will use BorderStyle/ActiveBorderStyle directly
}

// NewAuthContainer creates a new AuthContainer.
func NewAuthContainer() AuthContainer {
	selector := NewAuthSelector()
	return AuthContainer{
		Width:          0,
		Height:         0,
		Active:         false,
		authSelector:   selector,
		titleStyle:     styles.DefaultTheme.TitleStyle.Copy(),
	}
}

// SetWidth sets the width of the AuthContainer.
func (ac *AuthContainer) SetWidth(width int) {
	ac.Width = width
}

// SetHeight sets the height of the AuthContainer.
func (ac *AuthContainer) SetHeight(height int) {
	ac.Height = height
}

// SetActive sets the active state of the AuthContainer and its focusable children.
func (ac *AuthContainer) SetActive(active bool) {
	ac.Active = active
	// For now, if the container is active, the selector is also made active.
	// More complex focus logic might be needed if more elements are added.
	ac.authSelector.SetActive(active)
}

// Update handles messages for the AuthContainer.
func (ac *AuthContainer) Update(msg tea.Msg) tea.Cmd {
	if ac.Active { // Only delegate to authSelector if the container is active
		// If there were multiple focusable elements in AuthContainer,
		// we'd need to check which one has focus before delegating.
		// For now, if AuthContainer is active, AuthSelector is the active element.
		return ac.authSelector.Update(msg)
	}
	return nil
}

// View renders the AuthContainer.
func (ac AuthContainer) View() string {
	if ac.Width == 0 || ac.Height == 0 {
		return ""
	}

	var currentFrameStyle lipgloss.Style
	if ac.Active {
		currentFrameStyle = styles.DefaultTheme.ActiveBorderStyle.Copy()
	} else {
		currentFrameStyle = styles.DefaultTheme.BorderStyle.Copy()
	}

	outerFrame := currentFrameStyle.
		Width(ac.Width).
		Height(ac.Height).
		Padding(0, 1) // Padding between frame and inner content

	trueInnerWidth := ac.Width - outerFrame.GetHorizontalFrameSize()
	trueInnerHeight := ac.Height - outerFrame.GetVerticalFrameSize()
	if trueInnerWidth < 0 { trueInnerWidth = 0 }
	if trueInnerHeight < 0 { trueInnerHeight = 0 }

	var contentLines []string

	// Line 1: Title "Authentication" - REMOVED
	// if trueInnerHeight > 0 {
	// 	titleRendered := ac.titleStyle.Render("Authentication")
	// 	titleLine := lipgloss.NewStyle().Width(trueInnerWidth).Render(titleRendered)
	// 	contentLines = append(contentLines, titleLine)
	// }

	// Line 2 (now Line 1): "Type: " label + AuthSelector
	// Ensure there's enough height for at least one line of content.
	if trueInnerHeight > 0 { // Adjusted condition to check if any space is available
		// label := "Type: " // REMOVED
		// labelWidth := lipgloss.Width(label) // REMOVED
		// labelPart := lipgloss.NewStyle().Width(labelWidth).Render(label) // REMOVED
		
		selectorWidth := 30 // FIXED WIDTH for AuthSelector
		// selectorWidth := trueInnerWidth // OLD: AuthSelector now takes full inner width

		// Use a mutable copy of the selector to set width for rendering
		tempSelector := ac.authSelector 
		tempSelector.SetWidth(selectorWidth) // Set the fixed width
		// The active state of ac.authSelector is managed by AuthContainer.SetActive
		selectorView := tempSelector.View() // This can be a multi-line block if dropdown is open

		// Render the selectorView. If trueInnerWidth is greater than selectorWidth (30),
		// the selectorView will be left-aligned within the available trueInnerWidth.
		// If trueInnerWidth is less than 30, selectorView will be clipped by this rendering step.
		contentLines = append(contentLines, lipgloss.NewStyle().Width(trueInnerWidth).Render(selectorView))
	}
	
	// Build the main content block from what we have so far.
	innerContentBlock := lipgloss.JoinVertical(lipgloss.Left, contentLines...)

	// Calculate padding needed to fill the rest of the trueInnerHeight.
	paddingHeight := trueInnerHeight - lipgloss.Height(innerContentBlock)
	if paddingHeight < 0 { 
		paddingHeight = 0 // Content is already taller than the container.
	}

	var finalInnerContent string
	if paddingHeight > 0 {
		paddingBlock := lipgloss.NewStyle().Width(trueInnerWidth).Height(paddingHeight).Render("")
		finalInnerContent = lipgloss.JoinVertical(lipgloss.Left, innerContentBlock, paddingBlock)
	} else {
		finalInnerContent = innerContentBlock // No padding needed, or content overflows.
	}
	
	return outerFrame.Render(finalInnerContent)
}

// GetAuthHeaders returns the authentication headers based on the selected type and inputs. (Placeholder)
func (ac AuthContainer) GetAuthHeaders() map[string]string {
	// TODO: Implement logic based on authSelector.selectedIndex and future input fields
	return nil
}

// IsFocused checks if the AuthContainer or its components are focused. (Placeholder)
func (ac AuthContainer) IsFocused() bool {
	return ac.Active // Simple check for now, might need to check authSelector.active
}
