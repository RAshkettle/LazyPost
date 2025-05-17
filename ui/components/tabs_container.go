// Package components provides UI components for the LazyPost application.
package components

import (
	"fmt"

	"github.com/RAshkettle/LazyPost/ui/styles"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// TabsContainer represents a tabbed container with multiple tabs.
// It manages a main set of tabs (Query and Result) and renders the appropriate
// inner tab component based on the active tab selection.
type TabsContainer struct {
	Tabs        []string    // Labels for the main tabs
	ActiveTab   int         // Index of the currently active main tab
	Width       int         // Width of the container in characters
	Height      int         // Height of the container in characters
	Active      bool        // Whether the component is currently active/focused
	TabContents []string    // Default content for each tab (used as fallback)
	QueryTab    QueryTab    // The query tab component with its inner tabs
	ResultTab   ResultTab   // The result tab component with its inner tabs
}

// NewTabsContainer creates a new tab container with Query and Result tabs.
// It initializes both tabs with default content and proper configuration.
func NewTabsContainer() TabsContainer {
	queryContent := "Enter request parameters here.\n\n" +
		"Headers:\n" +
		"Content-Type: application/json\n\n" +
		"Body:\n" +
		"{\n  \"key\": \"value\"\n}"
	
	resultContent := "Response will be displayed here after request is sent."
	
	return TabsContainer{
		Tabs:        []string{"Query", "Result"},
		ActiveTab:   0,
		Width:       0,
		Height:      0,
		Active:      false,
		TabContents: []string{queryContent, resultContent},
		QueryTab:    NewQueryTab(),
		ResultTab:   NewResultTab(),
	}
}

// SetWidth sets the width of the tab container and propagates
// the appropriate width to the inner tab components, with reduced right margin.
func (t *TabsContainer) SetWidth(width int) {
	t.Width = width
	// Reduced margin on the right by 50%
	contentWidth := width - 2 // Reduced from width - 4
	t.QueryTab.SetWidth(contentWidth)
	t.ResultTab.SetWidth(contentWidth)
}

// SetHeight sets the height of the tab container and propagates
// the height to the inner tab components, giving the QueryTab more vertical space.
func (t *TabsContainer) SetHeight(height int) {
	t.Height = height
	// Give the QueryTab more height (we're adding an extra 10%)
	queryTabHeight := height - 4 + int(float64(height-4)*0.1)
	t.QueryTab.SetHeight(queryTabHeight) 
	t.ResultTab.SetHeight(queryTabHeight)
}

// SetActive sets the active state of the tab container and propagates
// the active state to the inner tab components.
func (t *TabsContainer) SetActive(active bool) {
	t.Active = active
	t.QueryTab.SetActive(active)
	t.ResultTab.SetActive(active)
}

// SwitchToTab switches to the specified tab by index.
// If the index is out of range, no change is made.
func (t *TabsContainer) SwitchToTab(tabIndex int) {
	if tabIndex >= 0 && tabIndex < len(t.Tabs) {
		t.ActiveTab = tabIndex
	}
}

// Update processes input messages and updates the container state.
// It handles alt+key combinations for tab switching and delegates
// tab/shift+tab navigation to the appropriate inner tab component.
func (t *TabsContainer) Update(msg tea.Msg) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if !t.Active {
			return
		}
		
		switch msg.String() {
		case "alt+4":
			// Switch to Query tab
			t.SwitchToTab(0)
		case "alt+5":
			// Switch to Result tab
			t.SwitchToTab(1)
		case "tab", "shift+tab":
			// Handle tab/shift+tab events in the active tab
			if t.ActiveTab == 0 {
				t.QueryTab.Update(msg)
			} else if t.ActiveTab == 1 {
				t.ResultTab.Update(msg)
			}
		default:
			// Pass other messages to the active tab
			if t.ActiveTab == 0 {
				t.QueryTab.Update(msg)
			} else if t.ActiveTab == 1 {
				t.ResultTab.Update(msg)
			}
		}
	}
}

// View renders the tab container component with the active tab's content.
// It creates a tabbed interface with hotkey indicators and renders the appropriate
// inner tab component (QueryTab or ResultTab) based on which main tab is active.
func (t TabsContainer) View() string {
	if t.Width == 0 || t.Height == 0 {
		return ""
	}
	
	// Define styles
	borderStyle := styles.BorderStyle
	
	if t.Active {
		borderStyle = styles.ActiveBorderStyle
	}
	
	// Create tab styles
	tabStyle := lipgloss.NewStyle().
		Padding(0, 2).
		MarginRight(1).
		Foreground(styles.SecondaryColor)
	
	// Base active tab style - green if tab container is active, white with bold if not
	activeTabStyle := lipgloss.NewStyle().
		Padding(0, 2).
		MarginRight(1).
		Bold(true)
	
	// Apply appropriate color based on active state
	if t.Active {
		activeTabStyle = activeTabStyle.Foreground(styles.PrimaryColor)
	} else {
		activeTabStyle = activeTabStyle.Foreground(styles.SecondaryColor)
	}
	
	// Create holistic tab rendering function
	renderTab := func(text string, index int, isActive bool) string {
		// Choose the appropriate style for the tab
		var baseStyle lipgloss.Style
		if isActive {
			baseStyle = activeTabStyle
		} else {
			baseStyle = tabStyle
		}
		
		// Create tab text with Alt+number hotkey
		tabText := fmt.Sprintf("(Alt+%d) %s", index+3, text)
		return baseStyle.Render(tabText)
	}
	
	// Render tabs
	var renderedTabs []string
	for i, tab := range t.Tabs {
		renderedTabs = append(renderedTabs, renderTab(tab, i, i == t.ActiveTab))
	}
	
	// Join tabs horizontally
	tabBar := lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)
	
	// Create content area
	contentStyle := lipgloss.NewStyle().
		Width(t.Width - 2). // Reduced from width - 4
		Height(t.Height - 4).
		Padding(1, 2)
	
	var content string
	if t.ActiveTab == 0 {
		// Render QueryTab component
		content = t.QueryTab.View()
	} else if t.ActiveTab == 1 {
		// Render ResultTab component
		content = t.ResultTab.View()
	} else {
		// Render other tabs normally
		content = contentStyle.Render(t.TabContents[t.ActiveTab])
	}
	
	// Put it all together with a border
	mainStyle := borderStyle.
		Width(t.Width).
		Height(t.Height)
	
	// Create content area with border
	contentBox := mainStyle.Render(content)
	
	// Position the tab bar above the content box
	return lipgloss.JoinVertical(lipgloss.Left, tabBar, contentBox)
}

// GetResultTab returns a pointer to the result tab component.
func (t *TabsContainer) GetResultTab() *ResultTab {
	return &t.ResultTab
}

// GetQueryTab returns a pointer to the query tab component.
func (t *TabsContainer) GetQueryTab() *QueryTab {
	return &t.QueryTab
}
