// Package components defines various UI components for the LazyPost application.
package components

import (
	"github.com/RAshkettle/LazyPost/ui/styles"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// MethodSelector represents the HTTP method selection component.
// It allows the user to choose an HTTP method (e.g., GET, POST) from a predefined list.
// The component can display as a simple selection or an open dropdown list.
type MethodSelector struct {
	Methods        []string // Methods is the list of available HTTP method strings.
	SelectedMethod int      // SelectedMethod is the index of the currently selected method in the Methods slice.
	Width          int      // Width is the rendering width of the component.
	Active         bool     // Active indicates whether the component is currently focused and interactive.
	DropdownOpen   bool     // DropdownOpen indicates whether the list of methods is currently displayed as a dropdown.
}

// NewMethodSelector creates and initializes a new MethodSelector component.
// It populates the list of HTTP methods and sets initial default values.
func NewMethodSelector() MethodSelector {
	return MethodSelector{
		Methods:        []string{"GET", "POST", "PUT", "DELETE", "PATCH"},
		SelectedMethod: 0,
		Width:          0,
		Active:         false,
		DropdownOpen:   false,
	}
}

// SetWidth sets the rendering width for the MethodSelector component.
func (m *MethodSelector) SetWidth(width int) {
	m.Width = width
}

// SetActive sets the active state of the MethodSelector.
// An active selector responds to key presses and has distinct visual styling.
func (m *MethodSelector) SetActive(active bool) {
	m.Active = active
}

// GetSelectedMethod returns the string representation of the currently selected HTTP method.
// If no methods are available or selected (which is unlikely in normal operation),
// it might return an empty string or the default method.
func (m *MethodSelector) GetSelectedMethod() string {
	if len(m.Methods) == 0 {
		return ""
	}
	return m.Methods[m.SelectedMethod]
}

// Next selects the next HTTP method in the list, wrapping around to the beginning if necessary.
func (m *MethodSelector) Next() {
	m.SelectedMethod = (m.SelectedMethod + 1) % len(m.Methods)
}

// Prev selects the previous HTTP method in the list, wrapping around to the end if necessary.
func (m *MethodSelector) Prev() {
	m.SelectedMethod = (m.SelectedMethod - 1 + len(m.Methods)) % len(m.Methods)
}

// Update handles messages for the MethodSelector, primarily key presses when it's active.
// It allows toggling the dropdown with Enter and navigating with Up/Down arrows when the dropdown is open.
func (m *MethodSelector) Update(msg tea.Msg) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if !m.Active {
			return
		}
		
		switch msg.String() {
		case "enter":
			// Toggle dropdown open/closed
			m.DropdownOpen = !m.DropdownOpen
		
		case "down":
			// Navigate down in the dropdown
			if m.DropdownOpen {
				m.Next()
			}
		
		case "up":
			// Navigate up in the dropdown
			if m.DropdownOpen {
				m.Prev()
			}
		}
	}
}

// View renders the MethodSelector component.
// If the dropdown is open, it lists all methods with the current selection highlighted.
// If closed, it shows only the selected method with a dropdown indicator.
// The component includes a title and is bordered, with styles changing based on the active state.
func (m MethodSelector) View() string {
	// Define styles
	borderStyle := styles.BorderStyle

	if m.Active {
		borderStyle = styles.ActiveBorderStyle
	}
	
	// Use minimal padding for consistency with URL component
	borderStyle = borderStyle.Padding(0, 1)

	// Create simple title with number hotkey
	titleStyle := lipgloss.NewStyle().
		Bold(true)
	
	// Change title color based on active state
	if m.Active {
		titleStyle = titleStyle.Foreground(styles.PrimaryColor)
	} else {
		titleStyle = titleStyle.Foreground(styles.SecondaryColor)
	}
	
	title := titleStyle.Render("(Alt+1) Method")
	
	// Build method content based on dropdown state
	var methodContent string
	
	// Create dropdown indicator
	dropdownIndicator := "▼" // Unicode down arrow
	if m.DropdownOpen {
		dropdownIndicator = "▲" // Unicode up arrow
	}
	
	selectedMethod := m.Methods[m.SelectedMethod]
	
	if m.DropdownOpen {
		// When dropdown is open, show all options
		methodContent = ""
		for i, method := range m.Methods {
			methodStyle := lipgloss.NewStyle()
			prefix := "  " // Space for indentation
			
			if i == m.SelectedMethod {
				methodStyle = styles.SelectedItemStyle
				prefix = "▶ " // Unicode right pointer
			}
			
			methodContent += methodStyle.Render(prefix + method) + "\n"
		}
		
		// Add instruction at the bottom
		helpStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#888888")).
			Italic(true)
		methodContent += helpStyle.Render("Press Enter to select")
	} else {
		// When dropdown is closed, show only selected method
		selectedStyle := lipgloss.NewStyle().
			Foreground(styles.BrightYellow).
			Bold(true)
		
		// Create a dropdown-like display with the currently selected method
		methodContent = selectedStyle.Render(selectedMethod) + " " + dropdownIndicator
		
		// Add hint if component is active
		if m.Active {
			hintStyle := lipgloss.NewStyle().
				Foreground(lipgloss.Color("#888888")).
				Italic(true)
			methodContent += "\n" + hintStyle.Render("Press Enter to open")
		}
	}
	
	// Render the method box
	methodBox := borderStyle.Width(m.Width).Render(methodContent)
	
	// Position the title above the method box
	return title + "\n" + methodBox
}
