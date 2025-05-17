// Package components provides UI components for the LazyPost application.
package components

import (
	"github.com/RAshkettle/LazyPost/ui/styles"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ResultTab represents the inner tab component for the Result tab.
// It provides a tabbed interface for viewing different aspects of an HTTP response
// including headers and body content. The component handles tab navigation via Tab/Shift+Tab keys.
type ResultTab struct {
	InnerTabs      []string          // Labels for the inner tabs
	ActiveInnerTab int               // Index of the currently active inner tab
	Width          int               // Width of the component in characters
	Height         int               // Height of the component in characters
	Active         bool              // Whether the component is currently active/focused
	HeadersTab     HeadersContainer  // Container for displaying response headers
	BodyTab        BodyContainer     // Container for displaying response body
}

// NewResultTab creates a new result tab component with predefined inner tabs.
// The component is initialized with the "Headers" tab selected, zero dimensions,
// and inactive state. Each inner tab has default placeholder content.
func NewResultTab() ResultTab {
	headers := NewHeadersContainer()
	body := NewBodyContainer()

	return ResultTab{
		InnerTabs:      []string{"Headers", "Body"},
		ActiveInnerTab: 0,
		Width:          0,
		Height:         0,
		Active:         false,
		HeadersTab:     headers,
		BodyTab:        body,
	}
}

// SetWidth sets the width of the component in characters.
func (r *ResultTab) SetWidth(width int) {
	r.Width = width
	
	// Update sub-components widths
	r.HeadersTab.SetWidth(width - 2) // Adjust for borders
	r.BodyTab.SetWidth(width - 2)    // Adjust for borders
}

// SetHeight sets the height of the component in characters.
func (r *ResultTab) SetHeight(height int) {
	r.Height = height
	
	// Calculate inner container height (95% of available height)
	innerHeight := int(float64(height) * 0.95)
	contentHeight := innerHeight - 4 // Adjust for tabs and borders
	
	// Update sub-components heights
	r.HeadersTab.SetHeight(contentHeight)
	r.BodyTab.SetHeight(contentHeight)
}

// SetActive sets the active state of the component.
// When active, the component has visual styling to indicate focus and responds to key presses.
func (r *ResultTab) SetActive(active bool) {
	r.Active = active
	
	// Set active state on the currently selected tab
	if r.ActiveInnerTab == 0 {
		r.HeadersTab.SetActive(active)
		r.BodyTab.SetActive(false)
	} else {
		r.HeadersTab.SetActive(false)
		r.BodyTab.SetActive(active)
	}
}

// SwitchToInnerTab switches to the specified inner tab by index.
// If the index is out of range, no change is made.
func (r *ResultTab) SwitchToInnerTab(tabIndex int) {
	if tabIndex >= 0 && tabIndex < len(r.InnerTabs) {
		r.ActiveInnerTab = tabIndex
		
		// Update active states of the sub-components
		if r.Active {
			if tabIndex == 0 {
				r.HeadersTab.SetActive(true)
				r.BodyTab.SetActive(false)
			} else {
				r.HeadersTab.SetActive(false)
				r.BodyTab.SetActive(true)
			}
		}
	}
}

// NextTab cycles to the next inner tab.
// It wraps around to the beginning if the end of the tabs is reached.
func (r *ResultTab) NextTab() {
	r.SwitchToInnerTab((r.ActiveInnerTab + 1) % len(r.InnerTabs))
}

// PrevTab cycles to the previous inner tab.
// It wraps around to the end if the beginning of the tabs is reached.
func (r *ResultTab) PrevTab() {
	r.SwitchToInnerTab((r.ActiveInnerTab - 1 + len(r.InnerTabs)) % len(r.InnerTabs))
}

// Update processes input messages and updates the result tab state.
// It handles tab and shift+tab key presses for inner tab navigation.
func (r *ResultTab) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if !r.Active {
			return nil
		}

		switch msg.String() {
		case "tab":
			// Cycle to next inner tab
			r.NextTab()
		case "shift+tab":
			// Cycle to previous inner tab
			r.PrevTab()
		default:
			// Pass key messages to the active inner tab
			if r.ActiveInnerTab == 0 {
				cmd = r.HeadersTab.Update(msg)
			} else {
				cmd = r.BodyTab.Update(msg)
			}
		}
	default:
		// Pass other messages to both containers
		cmd1 := r.HeadersTab.Update(msg)
		cmd2 := r.BodyTab.Update(msg)
		
		// Return the non-nil command if any
		if cmd1 != nil {
			cmd = cmd1
		} else if cmd2 != nil {
			cmd = cmd2
		}
	}
	
	return cmd
}

// SetHeadersContent sets the content for the headers tab.
func (r *ResultTab) SetHeadersContent(content string) {
	r.HeadersTab.SetContent(content)
}

// SetBodyContent sets the content for the body tab.
func (r *ResultTab) SetBodyContent(content string) {
	r.BodyTab.SetContent(content)
}

// SetContent sets the content for a specific inner tab by index.
// This method is for backward compatibility.
func (r *ResultTab) SetContent(tabIndex int, content string) {
	if tabIndex == 0 {
		r.SetHeadersContent(content)
	} else if tabIndex == 1 {
		r.SetBodyContent(content)
	}
}

// View renders the result tab component
func (r ResultTab) View() string {
	if r.Width == 0 || r.Height == 0 {
		return ""
	}

	// Define styles
	borderStyle := styles.BorderStyle
	if r.Active {
		borderStyle = styles.ActiveBorderStyle
	}

	// Create tab styles
	tabStyle := lipgloss.NewStyle().
		Padding(0, 2).
		MarginRight(1).
		Foreground(styles.SecondaryColor)

	activeTabStyle := lipgloss.NewStyle().
		Padding(0, 2).
		MarginRight(1).
		Bold(true)

	if r.Active {
		activeTabStyle = activeTabStyle.Foreground(styles.PrimaryColor)
	} else {
		activeTabStyle = activeTabStyle.Foreground(styles.SecondaryColor)
	}

	// Render inner tabs
	var renderedInnerTabs []string
	for i, tab := range r.InnerTabs {
		if i == r.ActiveInnerTab {
			renderedInnerTabs = append(renderedInnerTabs, activeTabStyle.Render(tab))
		} else {
			renderedInnerTabs = append(renderedInnerTabs, tabStyle.Render(tab))
		}
	}

	// Join inner tabs horizontally
	innerTabBar := lipgloss.JoinHorizontal(lipgloss.Top, renderedInnerTabs...)

	// Calculate inner container height (95% of available height)
	innerHeight := int(float64(r.Height) * 0.95)

	// Get content based on active inner tab
	var content string
	if r.ActiveInnerTab == 0 {
		content = r.HeadersTab.View()
	} else {
		content = r.BodyTab.View()
	}

	// Inner container with border
	innerContainerStyle := borderStyle.
		Width(r.Width).
		Height(innerHeight)
	innerContainer := innerContainerStyle.Render(content)

	// Position the inner tab bar above the inner container with a negative margin
	// to make tabs appear outside the border
	tabBarStyle := lipgloss.NewStyle().
		MarginBottom(-1)

	styledTabBar := tabBarStyle.Render(innerTabBar)

	// Create help text for tab navigation
	helpStyle := lipgloss.NewStyle().
		Foreground(styles.SecondaryColor).
		Align(lipgloss.Right).
		MarginTop(1).
		Width(r.Width).
		Italic(true)
	
	helpText := helpStyle.Render("Press Tab/Shift+Tab to cycle through subitems")

	// Return vertical layout with tab bar, inner container, and help text
	return lipgloss.JoinVertical(
		lipgloss.Left,
		styledTabBar,
		innerContainer,
		helpText,
	)
}
