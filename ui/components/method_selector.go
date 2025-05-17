package components

import (
	"github.com/RAshkettle/LazyPost/ui/styles"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// MethodSelector represents the HTTP method selection component.
// It provides a dropdown-like interface for selecting common HTTP methods
// (GET, POST, PUT, DELETE, PATCH) with visual feedback for the current selection.
type MethodSelector struct {
	Methods        []string // Available HTTP methods to select from
	SelectedMethod int      // Index of the currently selected method
	Width          int      // Width of the component in characters
	Active         bool     // Whether the component is currently active/focused
	DropdownOpen   bool     // Whether the dropdown list is currently open
}

// NewMethodSelector creates a new method selector component with predefined HTTP methods.
// The component is initialized with GET method selected, zero width, and inactive state.
func NewMethodSelector() MethodSelector {
	return MethodSelector{
		Methods:        []string{"GET", "POST", "PUT", "DELETE", "PATCH"},
		SelectedMethod: 0,
		Width:          0,
		Active:         false,
		DropdownOpen:   false,
	}
}

// SetWidth sets the width of the method selector in characters.
func (m *MethodSelector) SetWidth(width int) {
	m.Width = width
}

// SetActive sets the active state of the method selector.
// When active, the component has visual styling to indicate focus and responds to key presses.
func (m *MethodSelector) SetActive(active bool) {
	m.Active = active
}

// GetSelectedMethod returns the currently selected HTTP method as a string.
// If there are no methods available (which shouldn't happen under normal circumstances),
// it returns an empty string.
func (m *MethodSelector) GetSelectedMethod() string {
	if len(m.Methods) == 0 {
		return ""
	}
	return m.Methods[m.SelectedMethod]
}

// Next selects the next HTTP method in the list.
// It wraps around to the beginning if the end of the list is reached.
func (m *MethodSelector) Next() {
	m.SelectedMethod = (m.SelectedMethod + 1) % len(m.Methods)
}

// Prev selects the previous HTTP method in the list.
// It wraps around to the end if the beginning of the list is reached.
func (m *MethodSelector) Prev() {
	m.SelectedMethod = (m.SelectedMethod - 1 + len(m.Methods)) % len(m.Methods)
}

// Update processes input messages and updates the method selector state.
// It handles key presses for dropdown navigation, selection, and toggling.
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

// View renders the method selector component as a string for terminal display.
// When the dropdown is closed, it shows only the selected method.
// When the dropdown is open, it shows all available methods with the selected one highlighted.
// The component includes a border and title, with styling based on the active state.
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
