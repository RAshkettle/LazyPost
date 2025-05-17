package components

import (
	"github.com/RAshkettle/LazyPost/ui/styles"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// QueryTab represents the inner tab component for the Query tab.
// It provides a tabbed interface for configuring different aspects of an HTTP request
// including parameters, authentication, headers, and body.
// The component handles tab navigation via Tab/Shift+Tab keys.
type QueryTab struct {
	InnerTabs      []string        // Labels for the inner tabs
	ActiveInnerTab int             // Index of the currently active inner tab
	Width          int             // Width of the component in characters
	Height         int             // Height of the component in characters
	Active         bool            // Whether the component is currently active/focused
	ParamsInput    ParamsContainer // Container for parameter inputs
	HeadersInput   HeadersInputContainer // Container for header inputs
	QueryBodyInput textarea.Model  // Textarea for request body input 

	// Placeholder content for other tabs
	authContent    string
	headersContent string
}

// NewQueryTab creates a new query tab component with predefined inner tabs.
// The component is initialized with the "Params" tab selected, zero dimensions,
// and inactive state. Each inner tab has default placeholder content.
func NewQueryTab() QueryTab {
	authContent := "Configure authentication settings here."
	headersContent := "Configure request headers here."

	paramsInput := NewParamsContainer()
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
		HeadersInput:   headersInput,
		QueryBodyInput: bodyInput, 
		authContent:    authContent,
		headersContent: headersContent,
	}
}

// SetWidth sets the width of the component in characters.
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
	q.HeadersInput.SetWidth(actualContentDisplayWidth)

	queryBodyInputWidth := actualContentDisplayWidth - 2
	if queryBodyInputWidth < 0 {
		queryBodyInputWidth = 0
	}
	q.QueryBodyInput.SetWidth(queryBodyInputWidth)
}

// SetHeight sets the height of the component in characters.
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
	q.HeadersInput.SetHeight(actualContentDisplayHeight)

	queryBodyInputHeight := actualContentDisplayHeight - 2
	if queryBodyInputHeight < 0 {
		queryBodyInputHeight = 0
	}
	q.QueryBodyInput.SetHeight(queryBodyInputHeight)
}

// SetActive sets the active state of the component.
func (q *QueryTab) SetActive(active bool) {
	q.Active = active
	q.updateFocus()
}

// updateFocus manages focus for internal components based on active state and active inner tab.
func (q *QueryTab) updateFocus() {
	isParamsActive := q.Active && q.InnerTabs[q.ActiveInnerTab] == "Params"
	isBodyActive := q.Active && q.InnerTabs[q.ActiveInnerTab] == "Body"
	isHeadersActive := q.Active && q.InnerTabs[q.ActiveInnerTab] == "Headers"

	if isParamsActive {
		q.ParamsInput.SetActive(true) 
		q.QueryBodyInput.Blur()
		q.HeadersInput.SetActive(false)
	} else if isBodyActive {
		q.ParamsInput.SetActive(false)
		q.QueryBodyInput.Focus() 
		q.HeadersInput.SetActive(false)
	} else if isHeadersActive {
		q.ParamsInput.SetActive(false)
		q.QueryBodyInput.Blur()
		q.HeadersInput.SetActive(true)
	} else {
		q.ParamsInput.SetActive(false)
		q.QueryBodyInput.Blur()
		q.HeadersInput.SetActive(false)
	}
}

// SwitchToInnerTab switches to the specified inner tab by index.
func (q *QueryTab) SwitchToInnerTab(tabIndex int) {
	if tabIndex >= 0 && tabIndex < len(q.InnerTabs) {
		currentActiveTabName := q.InnerTabs[q.ActiveInnerTab]
		if currentActiveTabName == "Params" {
			q.ParamsInput.Blur() 
			q.ParamsInput.SetActive(false) // Also explicitly deactivate
		} else if currentActiveTabName == "Body" {
			q.QueryBodyInput.Blur()
		}

		q.ActiveInnerTab = tabIndex
		q.updateFocus() 
	}
}

// NextTab cycles to the next inner tab.
func (q *QueryTab) NextTab() {
	newTabIndex := (q.ActiveInnerTab + 1) % len(q.InnerTabs)
	q.SwitchToInnerTab(newTabIndex)
}

// PrevTab cycles to the previous inner tab.
func (q *QueryTab) PrevTab() {
	newTabIndex := (q.ActiveInnerTab - 1 + len(q.InnerTabs)) % len(q.InnerTabs)
	q.SwitchToInnerTab(newTabIndex)
}

// Update processes input messages and updates the query tab state.
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
					cmd = q.ParamsInput.Update(msg) // ParamsInput handles its own internal nav keys
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

// View renders the query tab component as a string for terminal display.
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
		case "Auth":
			placeholderText = q.authContent
		case "Headers":
			placeholderText = q.headersContent
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
		currentContent = placeholderStyle.Render(placeholderText)
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

// GetBodyContent returns the content of the query body input.
func (q *QueryTab) GetBodyContent() string {
	return q.QueryBodyInput.Value()
}

// IsAnyInputFocused checks if any input within the QueryTab is focused.
func (q *QueryTab) IsAnyInputFocused() bool {
	if q.InnerTabs[q.ActiveInnerTab] == "Params" && q.ParamsInput.IsAnyInputFocused() {
		return true
	}
	if q.InnerTabs[q.ActiveInnerTab] == "Body" && q.QueryBodyInput.Focused() {
		return true
	}
	return false
}
