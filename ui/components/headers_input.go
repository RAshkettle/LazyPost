// Package components provides UI components for the LazyPost application.
package components

import (
	"strings"

	"github.com/RAshkettle/LazyPost/ui/styles"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// numHeaderRows defines the fixed number of header input rows in the HeadersInputContainer.
const numHeaderRows = 9

// HeaderInput represents a single row in the HeadersInputContainer.
// It consists of a dropdown for selecting a header name and a text input for its value.
type HeaderInput struct {
	HeaderSelect      []string        // HeaderSelect stores the list of available header names for the dropdown.
	SelectedHeader    int             // SelectedHeader is the index of the currently selected header name in HeaderSelect.
	DropdownOpen      bool            // DropdownOpen indicates whether the header name dropdown is currently visible.
	ValueInput        textinput.Model // ValueInput is the text input field for the header value.
	width             int             // width is the total width of this individual header input row (not currently used for individual styling but available).
	headerSelectWidth int             // headerSelectWidth is the allocated width for the header selection part.
	valueInputWidth   int             // valueInputWidth is the allocated width for the value input part.
}

// HeadersInputContainer manages a list of HeaderInput rows, allowing users to input multiple HTTP headers.
// It handles focus navigation between rows and between the header name and value fields within a row.
// It also provides functionality to retrieve all entered headers as a map.
type HeadersInputContainer struct {
	inputs          []HeaderInput  // inputs is the slice of HeaderInput rows.
	focusedRow      int            // focusedRow is the index of the currently focused row.
	focusedInput    int            // focusedInput indicates which part of the focused row has focus (0 for HeaderSelect, 1 for ValueInput).
	Active          bool           // Active indicates if the container itself is focused and interactive.
	width           int            // width is the total width of the container.
	height          int            // height is the total height of the container.
	showHelp        bool           // showHelp determines if the help text is displayed.
	helpText        string         // helpText is the instructional message displayed to the user.
	headerLabel     string         // headerLabel is the text label for the header name column.
	valueLabel      string         // valueLabel is the text label for the header value column.
	baseHeaderStyle lipgloss.Style // baseHeaderStyle is the base style for the header name input area.
	baseValueStyle  lipgloss.Style // baseValueStyle is the base style for the header value input area.
}

// headerOptionsStrings provides a default list of common HTTP header names for the dropdown.
var headerOptionsStrings = []string{
	"Empty", "Accept", "Accept-Charset", "Accept-Encoding", "Accept-Language",
	"Authorization", "Cache-Control", "Connection", "Content-Length",
	"Content-MD5", "Content-Type", "Cookie", "Date", "Expect",
	"Host", "Max-Forwards",
	"Origin", "Pragma", "Proxy-Authorization", "Range", "Referer",
	"TE", "Upgrade", "User-Agent", "Via",
	"X-Csrf-Token", "X-Request-ID", "X-Correlation-ID",
}

// NewHeadersInputContainer creates and initializes a new HeadersInputContainer.
// It pre-populates a fixed number of HeaderInput rows with default values and styles.
func NewHeadersInputContainer() HeadersInputContainer {
	inputs := make([]HeaderInput, numHeaderRows)
	for i := range numHeaderRows {
		valIn := textinput.New()
		valIn.Placeholder = "Value"
		valIn.Prompt = "" // Remove the prompt indicator
		valIn.CharLimit = 256
		valIn.Width = 40 // Default width, will be adjusted

		inputs[i] = HeaderInput{
			HeaderSelect:   make([]string, len(headerOptionsStrings)), // Initialize with a copy
			SelectedHeader: 0,
			DropdownOpen:   false,
			ValueInput:     valIn,
		}
		copy(inputs[i].HeaderSelect, headerOptionsStrings)
	}

	// Define base styles for the input boxes
	baseHeaderStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Padding(0, 1) // Minimal padding

	baseValueStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Padding(0, 1) // Minimal padding

	return HeadersInputContainer{
		inputs:          inputs,
		focusedRow:      0,
		focusedInput:    0,     // Start focus on the first header select
		Active:          false, // Initialize Active state
		showHelp:        true,
		helpText:        "Use ↑/↓/←/→ to navigate, Enter to toggle dropdown/edit.",
		headerLabel:     "Header",
		valueLabel:      "Value",
		baseHeaderStyle: baseHeaderStyle,
		baseValueStyle:  baseValueStyle,
	}
}

// SetActive sets the active state of the HeadersInputContainer.
// When activated, it focuses the currently selected internal input field.
// When deactivated, it blurs all internal input fields.
func (h *HeadersInputContainer) SetActive(active bool) {
	h.Active = active
	if active {
		h.focusCurrentInput() // Focus the internal element when container becomes active
	} else {
		h.blurAllInputs() // Blur internal elements when container becomes inactive
	}
}

// Init is the first command that will be run by Bubble Tea for this component.
// It typically returns textinput.Blink to enable cursor blinking for text inputs.
func (h HeadersInputContainer) Init() tea.Cmd {
	return textinput.Blink
}

// SetWidth sets the rendering width for the HeadersInputContainer and its child HeaderInput rows.
// It distributes the width between the header name and value input fields.
func (h *HeadersInputContainer) SetWidth(width int) {
	h.width = width
	// Distribute width: ~40% for header, ~60% for value, adjust as needed
	// Considering labels and spacing.
	// Let's give Header fixed 30, ValueInput the rest minus some padding/margin
	labelWidth := lipgloss.Width(h.headerLabel + "  ") // Width of "Header  "
	h.inputs[0].headerSelectWidth = 30
	h.inputs[0].valueInputWidth = width - h.inputs[0].headerSelectWidth - labelWidth - lipgloss.Width(h.valueLabel+"  ") - 10 // Adjust 10 for safety/margins

	for i := range h.inputs {
		h.inputs[i].width = width
		h.inputs[i].headerSelectWidth = h.inputs[0].headerSelectWidth
		h.inputs[i].valueInputWidth = h.inputs[0].valueInputWidth
		h.inputs[i].ValueInput.Width = h.inputs[i].valueInputWidth
	}
}

// SetHeight sets the rendering height for the HeadersInputContainer.
func (h *HeadersInputContainer) SetHeight(height int) {
	h.height = height
}

// Update handles messages for the HeadersInputContainer, primarily key presses.
// It manages navigation (up/down rows, left/right between fields), opening/closing dropdowns,
// and delegating character input to the focused ValueInput field.
// It returns the updated HeadersInputContainer and any command to be executed (e.g., focus, blink).
func (h *HeadersInputContainer) Update(msg tea.Msg) (HeadersInputContainer, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	// Use h.inputs[h.focusedRow] directly instead of currentInput for clarity with potential changes to h.focusedRow
	// currentInput := &h.inputs[h.focusedRow] // Keep this for local modifications within a key press

	switch msg := msg.(type) {
	case tea.KeyMsg:
		currentInput := &h.inputs[h.focusedRow] // Get current input for this key event

		keyString := msg.String()
		isNavKey := keyString == "up" || keyString == "down" || keyString == "left" || keyString == "right"
		isEnterKey := keyString == "enter"

		// If ValueInput is the target, is focused for text, and it's NOT a nav or enter key, pass to it.
		if h.focusedInput == 1 && currentInput.ValueInput.Focused() && !isNavKey && !isEnterKey {
			currentInput.ValueInput, cmd = currentInput.ValueInput.Update(msg)
			cmds = append(cmds, cmd)
			return *h, tea.Batch(cmds...) // Character input handled, return.
		}

		// Store previous state for auto-closing dropdown
		prevFocusedRow := h.focusedRow
		prevFocusedInput := h.focusedInput
		prevDropdownOpen := false
		if prevFocusedInput == 0 { // Only a HeaderSelect can have a dropdown open
			prevDropdownOpen = h.inputs[prevFocusedRow].DropdownOpen
		}

		switch keyString {
		case "up":
			if h.focusedInput == 0 && currentInput.DropdownOpen { // Navigating open dropdown
				currentInput.SelectedHeader = (currentInput.SelectedHeader - 1 + len(currentInput.HeaderSelect)) % len(currentInput.HeaderSelect)
			} else { // Navigating rows
				if h.focusedRow > 0 {
					h.focusedRow--
				}
			}
		case "down":
			if h.focusedInput == 0 && currentInput.DropdownOpen { // Navigating open dropdown
				currentInput.SelectedHeader = (currentInput.SelectedHeader + 1) % len(currentInput.HeaderSelect)
			} else { // Navigating rows
				if h.focusedRow < numHeaderRows-1 {
					h.focusedRow++
				}
			}
		case "left":
			if h.focusedInput == 1 { // If on ValueInput
				h.focusedInput = 0 // Move to HeaderSelect
			}
		case "right":
			if h.focusedInput == 0 { // If on HeaderSelect
				h.focusedInput = 1 // Move to ValueInput
			}
		case "enter":
			switch h.focusedInput {
			case 0:
				currentInput.DropdownOpen = !currentInput.DropdownOpen
			case 1:
				if currentInput.ValueInput.Focused() {
					currentInput.ValueInput.Blur()
				} else {
					cmd = currentInput.ValueInput.Focus() // textinput.Focus() returns a command
					cmds = append(cmds, cmd)
				}

			}

		default:
			// Other keys are ignored if not handled by the ValueInput above
			// (e.g. character input when HeaderSelect is the active field)
		}

		// Auto-close dropdown if focus moved away from it
		if prevDropdownOpen {
			// If focus row changed OR focus input changed (from header to value)
			if h.focusedRow != prevFocusedRow || (h.focusedRow == prevFocusedRow && h.focusedInput != prevFocusedInput && prevFocusedInput == 0) {
				h.inputs[prevFocusedRow].DropdownOpen = false
			}
		}
		// currentInput might need to be updated if h.focusedRow changed
		// The final call to focusCurrentInput will use the updated h.focusedRow
	} // end switch msg.(type)

	// Ensure correct input is focused and collect its focus command (e.g., Blink)
	focusCmd := h.focusCurrentInput()
	cmds = append(cmds, focusCmd)

	return *h, tea.Batch(cmds...)
}

// focusCurrentInput ensures that the correct internal input field (HeaderSelect or ValueInput)
// within the currently focused row is appropriately focused or blurred.
// It returns a tea.Cmd, typically textinput.Blink if a ValueInput gains focus.
func (h *HeadersInputContainer) focusCurrentInput() tea.Cmd { // Modified to return tea.Cmd
	var focusCmd tea.Cmd
	for i := range h.inputs {
		if i == h.focusedRow {
			if h.focusedInput == 1 { // Focus ValueInput
				// ValueInput.Focus() returns a command (textinput.Blink).
				focusCmd = h.inputs[i].ValueInput.Focus()
			} else { // Focus HeaderSelect (conceptually, by blurring the value input)
				h.inputs[i].ValueInput.Blur()
			}
		} else { // Not the focused row
			h.inputs[i].ValueInput.Blur()
		}
	}
	return focusCmd
}

// blurAllInputs blurs all ValueInput fields in all HeaderInput rows.
// This is typically called when the HeadersInputContainer itself loses focus.
func (h *HeadersInputContainer) blurAllInputs() {
	for i := range h.inputs {
		h.inputs[i].ValueInput.Blur()
	}
}

// View renders the HeadersInputContainer.
// It displays labels for header and value, followed by each HeaderInput row.
// Each row shows the selected header name (or a dropdown if open) and the value input field.
// Styling is applied to indicate focus and active state. Help text is shown if enabled.
func (h HeadersInputContainer) View() string {
	var rows []string

	headerLabelStyled := lipgloss.NewStyle().Bold(true).Render(h.headerLabel)
	valueLabelStyled := lipgloss.NewStyle().Bold(true).Render(h.valueLabel)

	labelRow := lipgloss.JoinHorizontal(
		lipgloss.Left,
		lipgloss.NewStyle().Width(h.inputs[0].headerSelectWidth+2).Render(headerLabelStyled), // +2 for padding/border
		lipgloss.NewStyle().Width(h.inputs[0].valueInputWidth+2).Render(valueLabelStyled),    // +2 for padding/border
	)
	rows = append(rows, labelRow)

	for i, input := range h.inputs {
		hdrBoxStyle := h.baseHeaderStyle
		valBoxStyle := h.baseValueStyle

		isFocusedRow := i == h.focusedRow

		// --- Header Select Rendering ---
		var headerDisplayContent string
		dropdownIndicator := " ▼"
		if input.DropdownOpen {
			// dropdownIndicator = " ▲"
			var items []string
			for idx, itemStr := range input.HeaderSelect {
				itemStyle := lipgloss.NewStyle()
				prefix := "  "
				if idx == input.SelectedHeader {
					itemStyle = styles.SelectedItemStyle // Assuming styles.SelectedItemStyle is defined
					prefix = "▶ "
				}
				items = append(items, itemStyle.Render(prefix+itemStr))
			}
			headerDisplayContent = strings.Join(items, "\n")
			// Adjust height for open dropdown
			// +1 for border, or consider content height directly
			hdrBoxStyle = hdrBoxStyle.Height(len(input.HeaderSelect) - 1)
		} else {
			if len(input.HeaderSelect) > 0 && input.SelectedHeader >= 0 && input.SelectedHeader < len(input.HeaderSelect) {
				headerDisplayContent = input.HeaderSelect[input.SelectedHeader] + dropdownIndicator
			} else {
				headerDisplayContent = " (empty)" + dropdownIndicator
			}
			hdrBoxStyle = hdrBoxStyle.Height(1) // Standard height for closed dropdown
		}

		if isFocusedRow && h.focusedInput == 0 {
			hdrBoxStyle = hdrBoxStyle.BorderForeground(styles.PrimaryColor)
		} else {
			// Use a default/secondary color for inactive border, similar to MethodSelector's approach
			hdrBoxStyle = hdrBoxStyle.BorderForeground(styles.SecondaryColor) // Or a lipgloss.Color
		}
		headerView := hdrBoxStyle.Width(input.headerSelectWidth).Render(headerDisplayContent)

		// --- End Header Select Rendering ---

		// --- Value Input Rendering ---
		if isFocusedRow && h.focusedInput == 1 {
			valBoxStyle = valBoxStyle.BorderForeground(styles.PrimaryColor)
		} else {
			valBoxStyle = valBoxStyle.BorderForeground(styles.SecondaryColor) // Or a lipgloss.Color
		}
		valueView := valBoxStyle.Width(input.valueInputWidth).Render(input.ValueInput.View())
		// --- End Value Input Rendering ---

		row := lipgloss.JoinHorizontal(lipgloss.Top, headerView, " ", valueView)
		rows = append(rows, row)
	}

	if h.showHelp {
		// Define help style inline, similar to MethodSelector
		helpStyle := lipgloss.NewStyle().
			Foreground(styles.BrightYellow). // Changed to BrightYellow
			Italic(true)
		helpView := helpStyle.Render(h.helpText)
		rows = append(rows, "", helpView)
	}

	return lipgloss.JoinVertical(lipgloss.Left, rows...)
}

// GetHeaders returns a map of all valid headers entered by the user.
// A header is considered valid if its name is not "Empty" and its value is not an empty string.
func (h HeadersInputContainer) GetHeaders() map[string]string {
	headers := make(map[string]string)
	for _, input := range h.inputs {
		if len(input.HeaderSelect) > 0 && input.SelectedHeader < len(input.HeaderSelect) {
			selectedHeaderKey := input.HeaderSelect[input.SelectedHeader]
			value := input.ValueInput.Value()

			// Only add if the selected header is not "Empty" and the value is non-empty
			if selectedHeaderKey != "Empty" && value != "" {
				headers[selectedHeaderKey] = value
			}
		}
	}
	return headers
}

// GetSelectedValues returns the currently selected header name and its corresponding value
// for the currently focused row. This can be useful for context-aware operations.
func (h HeadersInputContainer) GetSelectedValues() (header string, value string) {
	if h.focusedRow < 0 || h.focusedRow >= len(h.inputs) {
		return "", ""
	}
	input := h.inputs[h.focusedRow]
	if len(input.HeaderSelect) > 0 && input.SelectedHeader < len(input.HeaderSelect) {
		header = input.HeaderSelect[input.SelectedHeader]
	}
	value = input.ValueInput.Value()
	return header, value
}

// IsDropdownOpen checks if the header name dropdown for the currently focused row is open.
func (h HeadersInputContainer) IsDropdownOpen() bool {
	if h.focusedInput == 0 && h.focusedRow >= 0 && h.focusedRow < len(h.inputs) {
		return h.inputs[h.focusedRow].DropdownOpen
	}
	return false
}
