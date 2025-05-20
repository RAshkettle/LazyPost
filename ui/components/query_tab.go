// Package components defines various UI components for the LazyPost application.
package components

import (
	"github.com/RAshkettle/LazyPost/ui/styles"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// QueryTab represents the main interactive area for constructing an HTTP request.
// It contains several inner tabs (Params, Auth, Headers, Body) allowing the user
// to configure different parts of the request. It manages focus between these inner tabs
// and delegates interactions to the active inner component.
type QueryTab struct {
	InnerTabs      []string              // InnerTabs stores the labels for the switchable inner sections (e.g., "Params", "Auth").
	ActiveInnerTab int                   // ActiveInnerTab is the index of the currently visible and interactive inner tab.
	Width          int                   // Width is the rendering width of the entire QueryTab component.
	Height         int                   // Height is the rendering height of the entire QueryTab component.
	Active         bool                  // Active indicates if the QueryTab itself (and thus its active inner tab) is focused.
	ParamsInput    ParamsContainer       // ParamsInput is the component for managing URL query parameters.
	AuthInput      AuthContainer         // AuthInput is the component for managing authentication settings.
	HeadersInput   HeadersInputContainer // HeadersInput is the component for managing request headers.
	QueryBodyInput textarea.Model        // QueryBodyInput is the text area for inputting the request body.

	// headersContent was a placeholder, now HeadersInput component is used.
	headersContent string // This might still be used if Headers tab is not fully componentized
}

// NewQueryTab creates and initializes a new QueryTab component.
// It sets up the inner tabs and their corresponding child components (ParamsContainer, AuthContainer, etc.).
func NewQueryTab() QueryTab {
	// authContent := "Configure authentication settings here." // No longer needed
	headersContent := "Configure request headers here."

	paramsInput := NewParamsContainer()
	authInput := NewAuthContainer() // Initialize AuthContainer
	headersInput := NewHeadersInputContainer()

	bodyInput := textarea.New()
	bodyInput.Placeholder = "Enter request body here in JSON..."
	bodyInput.ShowLineNumbers = false 

	return QueryTab{
		InnerTabs:      []string{"Params", "Auth", "Headers", "Body"},
		ActiveInnerTab: 0,
		Width:          0,
		Height:         0,
		Active:         false,
		ParamsInput:    paramsInput,
		AuthInput:      authInput, // Add AuthContainer to initialization
		HeadersInput:   headersInput,
		QueryBodyInput: bodyInput,
		// authContent:    authContent, // No longer needed
		headersContent: headersContent,
	}
}

// SetWidth sets the rendering width for the QueryTab and propagates it to its child components.
// The width is adjusted for borders and padding before being passed to children.
func (q *QueryTab) SetWidth(width int) {
	q.Width = width
	innerContainerBorderStyle := styles.BorderStyle
	if q.Active {
		innerContainerBorderStyle = styles.ActiveBorderStyle
	}

	actualContentDisplayWidth := q.Width - innerContainerBorderStyle.GetHorizontalBorderSize() - innerContainerBorderStyle.GetHorizontalPadding()
	if actualContentDisplayWidth < 0 {
		actualContentDisplayWidth = 0
	}
	q.ParamsInput.SetWidth(actualContentDisplayWidth)
	q.AuthInput.SetWidth(actualContentDisplayWidth) // Set width for AuthContainer
	q.HeadersInput.SetWidth(actualContentDisplayWidth)

	queryBodyInputWidth := actualContentDisplayWidth - 2
	if queryBodyInputWidth < 0 {
		queryBodyInputWidth = 0
	}
	q.QueryBodyInput.SetWidth(queryBodyInputWidth)
}

// SetHeight sets the rendering height for the QueryTab and propagates it to its child components.
// The height is adjusted for the tab bar, borders, and padding before being passed to children.
func (q *QueryTab) SetHeight(height int) {
	q.Height = height
	innerContainerHeight := q.Height - 2
	if innerContainerHeight < 0 {
		innerContainerHeight = 0
	}

	innerContainerBorderStyle := styles.BorderStyle
	if q.Active { 
		innerContainerBorderStyle = styles.ActiveBorderStyle
	}

	actualContentDisplayHeight := innerContainerHeight - innerContainerBorderStyle.GetVerticalBorderSize() - innerContainerBorderStyle.GetVerticalPadding()
	if actualContentDisplayHeight < 0 {
		actualContentDisplayHeight = 0
	}
	q.ParamsInput.SetHeight(actualContentDisplayHeight)
	q.AuthInput.SetHeight(actualContentDisplayHeight) // Set height for AuthContainer
	q.HeadersInput.SetHeight(actualContentDisplayHeight)

	queryBodyInputHeight := actualContentDisplayHeight - 2
	if queryBodyInputHeight < 0 {
		queryBodyInputHeight = 0
	}
	q.QueryBodyInput.SetHeight(queryBodyInputHeight)
}

// SetActive sets the active state of the QueryTab.
// This also triggers an update to the focus state of its internal components.
func (q *QueryTab) SetActive(active bool) {
	q.Active = active
	q.updateFocus()
}

// updateFocus manages which internal component (Params, Auth, Headers, Body)
// should be active and focused based on the QueryTab's overall active state
// and the currently selected ActiveInnerTab.
func (q *QueryTab) updateFocus() {
	isParamsActive := q.Active && q.InnerTabs[q.ActiveInnerTab] == "Params"
	isAuthActive := q.Active && q.InnerTabs[q.ActiveInnerTab] == "Auth" // Check for Auth tab
	isBodyActive := q.Active && q.InnerTabs[q.ActiveInnerTab] == "Body"
	isHeadersActive := q.Active && q.InnerTabs[q.ActiveInnerTab] == "Headers"

	if isParamsActive {
		q.ParamsInput.SetActive(true)
		q.AuthInput.SetActive(false) // Deactivate AuthContainer
		q.QueryBodyInput.Blur()
		q.HeadersInput.SetActive(false)
	} else if isAuthActive { // Handle Auth tab focus
		q.ParamsInput.SetActive(false)
		q.AuthInput.SetActive(true) // Activate AuthContainer
		q.QueryBodyInput.Blur()
		q.HeadersInput.SetActive(false)
	} else if isBodyActive {
		q.ParamsInput.SetActive(false)
		q.AuthInput.SetActive(false) // Deactivate AuthContainer
		q.QueryBodyInput.Focus()
		q.HeadersInput.SetActive(false)
	} else if isHeadersActive {
		q.ParamsInput.SetActive(false)
		q.AuthInput.SetActive(false) // Deactivate AuthContainer
		q.QueryBodyInput.Blur()
		q.HeadersInput.SetActive(true)
	} else {
		q.ParamsInput.SetActive(false)
		q.AuthInput.SetActive(false) // Deactivate AuthContainer
		q.QueryBodyInput.Blur()
		q.HeadersInput.SetActive(false)
	}
}

// SwitchToInnerTab changes the active inner tab to the one specified by tabIndex.
// It deactivates the previously active inner component and activates the new one.
func (q *QueryTab) SwitchToInnerTab(tabIndex int) {
	if tabIndex >= 0 && tabIndex < len(q.InnerTabs) {
		currentActiveTabName := q.InnerTabs[q.ActiveInnerTab]
		if currentActiveTabName == "Params" {
			q.ParamsInput.Blur() 
			q.ParamsInput.SetActive(false) // Also explicitly deactivate
		} else if currentActiveTabName == "Auth" { // Handle Auth tab deactivation
			q.AuthInput.SetActive(false)
		} else if currentActiveTabName == "Body" {
			q.QueryBodyInput.Blur()
		} else if currentActiveTabName == "Headers" {
			q.HeadersInput.SetActive(false)
		}

		q.ActiveInnerTab = tabIndex
		q.updateFocus() 
	}
}

// NextTab cycles to the next inner tab in the sequence.
func (q *QueryTab) NextTab() {
	newTabIndex := (q.ActiveInnerTab + 1) % len(q.InnerTabs)
	q.SwitchToInnerTab(newTabIndex)
}

// PrevTab cycles to the previous inner tab in the sequence.
func (q *QueryTab) PrevTab() {
	newTabIndex := (q.ActiveInnerTab - 1 + len(q.InnerTabs)) % len(q.InnerTabs)
	q.SwitchToInnerTab(newTabIndex)
}

// Update handles messages for the QueryTab.
// It manages Tab/Shift+Tab navigation between inner tabs.
// For other messages, it delegates to the Update method of the currently active inner component.
// It ensures that components like the textarea receive necessary updates for cursor blinking even if not fully active.
func (q *QueryTab) Update(msg tea.Msg) tea.Cmd {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	currentInnerTab := q.InnerTabs[q.ActiveInnerTab]

	if q.Active {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			// Handle Tab and Shift+Tab for QueryTab navigation first.
			// These should take precedence over component-internal Tab/Shift+Tab.
			switch msg.String() {
			case "tab":
				q.NextTab()
				return nil // Absorb Tab, prevent further processing by children
			case "shift+tab":
				q.PrevTab()
				return nil // Absorb Shift+Tab
			default:
				// If not Tab/Shift+Tab, pass to the active component if it's focused/active
				if currentInnerTab == "Params" && q.ParamsInput.Active {
					cmd = q.ParamsInput.Update(msg)
					cmds = append(cmds, cmd)
				} else if currentInnerTab == "Auth" && q.AuthInput.Active { // Delegate to AuthInput
					cmd = q.AuthInput.Update(msg)
					cmds = append(cmds, cmd)
				} else if currentInnerTab == "Headers" && q.HeadersInput.Active { // Check Active field
					// Update returns (HeadersInputContainer, tea.Cmd)
					newHeadersInput, headerCmd := q.HeadersInput.Update(msg)
					q.HeadersInput = newHeadersInput
					cmds = append(cmds, headerCmd)
				} else if currentInnerTab == "Body" && q.QueryBodyInput.Focused() {
					q.QueryBodyInput, cmd = q.QueryBodyInput.Update(msg)
					cmds = append(cmds, cmd)
				}
			}
		default:
			// Pass non-key messages to relevant components
			// This is important for components to process their own command results (like focus/blur)
			if currentInnerTab == "Params" {
				cmd = q.ParamsInput.Update(msg)
				cmds = append(cmds, cmd)
			}
			if currentInnerTab == "Auth" { // Pass non-key messages to AuthInput
				cmd = q.AuthInput.Update(msg)
				cmds = append(cmds, cmd)
			}
			if currentInnerTab == "Headers" {
				// Update returns (HeadersInputContainer, tea.Cmd)
				newHeadersInput, headerCmd := q.HeadersInput.Update(msg)
				q.HeadersInput = newHeadersInput
				cmds = append(cmds, headerCmd)
			}
			// QueryBodyInput also needs updates for its state (e.g., cursor blink)
			// even if it's not the active tab, but especially if it is.
			// The textarea.Update method is generally safe to call.
			q.QueryBodyInput, cmd = q.QueryBodyInput.Update(msg)
			cmds = append(cmds, cmd)
		}
	} else {
		// If QueryTab itself is not active, we still need to pass some messages
		// for components to update their internal state (e.g., cursor blinking for textarea).
		if _, ok := msg.(tea.KeyMsg); !ok { // Don't pass key messages if not active
			q.QueryBodyInput, cmd = q.QueryBodyInput.Update(msg)
			cmds = append(cmds, cmd)
			
			// ParamsInput might also need non-key messages if it has ongoing operations
			// For now, let's assume its SetActive(false) handles its state sufficiently.
			// If ParamsInput needs updates when QueryTab is inactive, add its update here too.
			// cmd = q.ParamsInput.Update(msg)
			// cmds = append(cmds, cmd)
		}
	}
	return tea.Batch(cmds...)
}

// View renders the QueryTab component.
// It displays a bar with inner tab labels, with the active tab highlighted.
// Below the tab bar, it renders the View of the currently active inner component.
// Help text is displayed at the bottom, contextual to the active inner tab and its state.
func (q QueryTab) View() string {
	if q.Width == 0 || q.Height == 0 {
		return ""
	}

	tabStyle := lipgloss.NewStyle().
		Padding(0, 2).
		MarginRight(1).
		Foreground(styles.SecondaryColor)

	activeTabStyle := lipgloss.NewStyle().
		Padding(0, 2).
		MarginRight(1).
		Bold(true)

	if q.Active {
		activeTabStyle = activeTabStyle.Foreground(styles.PrimaryColor)
	} else {
		activeTabStyle = activeTabStyle.Foreground(styles.SecondaryColor) 
	}

	var renderedInnerTabs []string
	for i, tab := range q.InnerTabs {
		if i == q.ActiveInnerTab {
			renderedInnerTabs = append(renderedInnerTabs, activeTabStyle.Render(tab))
		} else {
			renderedInnerTabs = append(renderedInnerTabs, tabStyle.Render(tab))
		}
	}

	innerTabBar := lipgloss.JoinHorizontal(lipgloss.Top, renderedInnerTabs...)

	innerContentBoxHeight := q.Height - 2 
	if innerContentBoxHeight < 0 {
		innerContentBoxHeight = 0
	}

	currentContentBoxBorderStyle := styles.BorderStyle
	if q.Active { 
		currentContentBoxBorderStyle = styles.ActiveBorderStyle
	}

	actualContentDisplayWidth := q.Width - currentContentBoxBorderStyle.GetHorizontalBorderSize() - currentContentBoxBorderStyle.GetHorizontalPadding()
	actualContentDisplayHeight := innerContentBoxHeight - currentContentBoxBorderStyle.GetVerticalBorderSize() - currentContentBoxBorderStyle.GetVerticalPadding()
	if actualContentDisplayWidth < 0 {
		actualContentDisplayWidth = 0
	}
	if actualContentDisplayHeight < 0 {
		actualContentDisplayHeight = 0
	}

	var currentContent string
	activeInnerTabName := q.InnerTabs[q.ActiveInnerTab]

	switch activeInnerTabName {
	case "Params":
		currentContent = q.ParamsInput.View()
	case "Auth": // Render AuthContainer
		currentContent = q.AuthInput.View()
	case "Headers":
		currentContent = q.HeadersInput.View()
	case "Body":
		activeQueryTabBorderColor := styles.PrimaryColor
		inactiveQueryTabBorderColor := styles.SecondaryColor

		focusedTAStyle := q.QueryBodyInput.FocusedStyle
		blurredTAStyle := q.QueryBodyInput.BlurredStyle

		if q.Active {
			focusedTAStyle.Base = focusedTAStyle.Base.BorderForeground(activeQueryTabBorderColor)
			blurredTAStyle.Base = blurredTAStyle.Base.BorderForeground(activeQueryTabBorderColor) 
		} else {
			focusedTAStyle.Base = focusedTAStyle.Base.BorderForeground(inactiveQueryTabBorderColor)
			blurredTAStyle.Base = blurredTAStyle.Base.BorderForeground(inactiveQueryTabBorderColor)
		}
		q.QueryBodyInput.FocusedStyle = focusedTAStyle
		q.QueryBodyInput.BlurredStyle = blurredTAStyle
		
		bodyView := q.QueryBodyInput.View()
		
		currentContent = lipgloss.NewStyle().
			Width(actualContentDisplayWidth).
			Height(actualContentDisplayHeight).
			Align(lipgloss.Center, lipgloss.Top). 
			Render(bodyView)

	default:
		var placeholderText string
		switch activeInnerTabName {
		// case "Auth": // No longer needed, AuthInput.View() is used
		// 	placeholderText = q.authContent
		case "Headers":
			// If Headers is not yet a full component, this might still be used.
			// For now, assuming HeadersInput.View() handles it.
			// If HeadersInput.View() can be empty or not fully cover the area,
			// a placeholder might still be relevant under certain conditions.
			// Let's remove headersContent for now, assuming HeadersInput.View() is sufficient.
			// placeholderText = q.headersContent 
			placeholderText = "Headers content via HeadersInput.View()"
		default:
			// This case should ideally not be reached if ActiveInnerTab is always valid
			// and corresponds to one of the defined InnerTabs ("Params", "Auth", "Headers", "Body").
			// If it is reached, it implies a state inconsistency.
			placeholderText = "Unknown tab content."
		}
		placeholderStyle := lipgloss.NewStyle().
			Width(actualContentDisplayWidth).
			Height(actualContentDisplayHeight).
			Padding(1, 2).
			Align(lipgloss.Center, lipgloss.Center)

		// Only render placeholder if not handled by a specific component view
		if activeInnerTabName != "Params" && activeInnerTabName != "Auth" && activeInnerTabName != "Body" && activeInnerTabName != "Headers" {
		    currentContent = placeholderStyle.Render(placeholderText)
		} else if activeInnerTabName == "Headers" && q.HeadersInput.View() == "" { // Example: if HeadersInput can be empty
			 // currentContent = placeholderStyle.Render("Configure request headers here.")
             // This is now handled by HeadersInput.View(), if it's empty, it's empty.
		}

	}

	innerContainer := currentContentBoxBorderStyle.
		Width(q.Width). 
		Height(innerContentBoxHeight).
		Render(currentContent)

	tabBarStyle := lipgloss.NewStyle().
		MarginBottom(-1) 

	styledTabBar := tabBarStyle.Render(innerTabBar)

	helpStyle := lipgloss.NewStyle().
		Foreground(styles.SecondaryColor).
		Align(lipgloss.Right).
		MarginTop(1). 
		Width(q.Width).
		Italic(true)
	
	helpTextString := "Press Tab/Shift+Tab to cycle items"
	if q.Active && activeInnerTabName == "Body" && q.QueryBodyInput.Focused() {
		helpTextString = "Esc to release focus; Tab/Shift+Tab to cycle tabs"
	} else if q.Active && activeInnerTabName == "Params" && q.ParamsInput.IsAnyInputFocused() {
		helpTextString = "Use Arrows/Tab to navigate fields; Tab/Shift+Tab to cycle tabs"
	}
	
	helpText := helpStyle.Render(helpTextString)
	
	return lipgloss.JoinVertical(
		lipgloss.Left,
		styledTabBar,
		innerContainer,
		helpText,
	)
}

// GetBodyContent returns the current content of the QueryBodyInput (request body text area).
func (q *QueryTab) GetBodyContent() string {
	return q.QueryBodyInput.Value()
}

// IsAnyInputFocused checks if any interactive element within the currently active inner tab is focused.
// This is used to determine context for keybindings or help text.
func (q *QueryTab) IsAnyInputFocused() bool {
	if q.InnerTabs[q.ActiveInnerTab] == "Params" && q.ParamsInput.IsAnyInputFocused() {
		return true
	}
	if q.InnerTabs[q.ActiveInnerTab] == "Auth" && q.AuthInput.IsFocused() { // Check AuthInput focus
		return true
	}
	if q.InnerTabs[q.ActiveInnerTab] == "Body" && q.QueryBodyInput.Focused() {
		return true
	}
	return false
}
