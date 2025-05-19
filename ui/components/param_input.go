// Package components provides UI components for the LazyPost application.
package components

import (
	"strings"

	"github.com/RAshkettle/LazyPost/ui/styles"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const numParamRows = 6

// ParamInput represents a single Name/Value input pair.
type ParamInput struct {
	NameInput  textinput.Model
	ValueInput textinput.Model
}

// ParamsContainer manages a list of parameter inputs (Name/Value pairs).
type ParamsContainer struct {
	Inputs       []ParamInput // Slice of parameter inputs
	Width        int          // Width of the container
	Height       int          // Height of the container
	Active       bool         // Whether the container is currently active/focused
	focusedRow   int          // Index of the currently focused row
	focusedCol   int          // 0 for Name, 1 for Value
	scrollOffset int          // For scrolling if not all rows fit
	contentWidth int          // Calculated width for content area
}

// NewParamsContainer creates a new ParamsContainer with a predefined number of rows.
func NewParamsContainer() ParamsContainer {
	inputs := make([]ParamInput, numParamRows)
	for i := range numParamRows {
		nameInput := textinput.New()
		nameInput.Placeholder = "Name"
		nameInput.Prompt = "" // No prompt, label will be above
		nameInput.CharLimit = 35

		valueInput := textinput.New()
		valueInput.Placeholder = "Value"
		valueInput.Prompt = "" // No prompt
		valueInput.CharLimit = 35

		inputs[i] = ParamInput{NameInput: nameInput, ValueInput: valueInput}
	}

	// Focus the first input by default
	if numParamRows > 0 {
		inputs[0].NameInput.Focus()
	}

	return ParamsContainer{
		Inputs:       inputs,
		Width:        0,
		Height:       0,
		Active:       false,
		focusedRow:   0,
		focusedCol:   0,
		scrollOffset: 0,
		contentWidth: 0,
	}
}

// SetWidth sets the width of the container and its child inputs.
func (pc *ParamsContainer) SetWidth(width int) {
	pc.Width = width

	currentStyle := styles.BorderStyle
	if pc.Active { // Though Active state might change, border/padding are same for current styles
		currentStyle = styles.ActiveBorderStyle
	}
	// Horizontal space taken by container\'s border and padding
	containerChrome := currentStyle.GetHorizontalBorderSize() + currentStyle.GetHorizontalPadding()

	pc.contentWidth = width - containerChrome
	pc.contentWidth = max(pc.contentWidth, 0)

	// Space between name and value inputs
	const spacingBetweenInputs = 1

	// Available width for the two text input columns (outer widths)
	inputsTotalOuterWidth := pc.contentWidth - spacingBetweenInputs
	inputsTotalOuterWidth = max(inputsTotalOuterWidth, 0)

	const textInputHorizontalBorderWidth = 2
	const desiredInputContentWidth = 35
	const idealOuterWidthPerInput = desiredInputContentWidth + textInputHorizontalBorderWidth // 37
	const totalIdealOuterWidth = idealOuterWidthPerInput * 2                                  // 74

	var nameColOuterWidth, valueColOuterWidth int

	if inputsTotalOuterWidth >= totalIdealOuterWidth {
		nameColOuterWidth = idealOuterWidthPerInput
		valueColOuterWidth = idealOuterWidthPerInput
	} else {
		// Distribute available space, try for 50/50 split
		nameColOuterWidth = inputsTotalOuterWidth / 2
		valueColOuterWidth = inputsTotalOuterWidth - nameColOuterWidth // Ensures total is met
	}

	if nameColOuterWidth < 0 {
		nameColOuterWidth = 0
	}
	if valueColOuterWidth < 0 {
		valueColOuterWidth = 0
	}

	for i := range pc.Inputs {
		nameInputContentWidth := nameColOuterWidth - textInputHorizontalBorderWidth
		if nameInputContentWidth < 0 {
			nameInputContentWidth = 0
		}
		pc.Inputs[i].NameInput.Width = nameInputContentWidth

		valueInputContentWidth := valueColOuterWidth - textInputHorizontalBorderWidth
		if valueInputContentWidth < 0 {
			valueInputContentWidth = 0
		}
		pc.Inputs[i].ValueInput.Width = valueInputContentWidth
	}
}

// SetHeight sets the height of the container.
func (pc *ParamsContainer) SetHeight(height int) {
	pc.Height = height
	// Height calculation might be needed if scrolling is implemented for many rows.
	// For now, assume all 6 rows are visible or handled by overall layout.
}

// SetActive sets the active state of the container.
func (pc *ParamsContainer) SetActive(active bool) {
	pc.Active = active
	if active {
		pc.ensureFocusedInputVisible()
		pc.focusCurrentInput()
	} else {
		pc.blurAllInputs()
	}
}

func (pc *ParamsContainer) blurAllInputs() {
	for i := range pc.Inputs {
		pc.Inputs[i].NameInput.Blur()
		pc.Inputs[i].ValueInput.Blur()
	}
}

func (pc *ParamsContainer) focusCurrentInput() {
	pc.blurAllInputs()
	if pc.focusedRow >= 0 && pc.focusedRow < len(pc.Inputs) {
		if pc.focusedCol == 0 {
			pc.Inputs[pc.focusedRow].NameInput.Focus()
		} else {
			pc.Inputs[pc.focusedRow].ValueInput.Focus()
		}
	}
}

func (pc *ParamsContainer) getNumDisplayableInputRows() int {
	if pc.Height <= 0 {
		return numParamRows // If height not set, assume all are displayable
	}

	currentStyle := styles.BorderStyle
	if pc.Active {
		currentStyle = styles.ActiveBorderStyle
	}
	borderSize := currentStyle.GetVerticalBorderSize()

	// Header=1, Separator=1. Total 2 fixed lines for these.
	// displayable is lines available for input rows AND scroll indicator
	displayable := pc.Height - borderSize - 2

	// If scrolling will be active (not all rows fit *before* accounting for scrollbar line)
	// and there's space for at least one row + scrollbar
	if numParamRows > displayable && displayable > 0 {
		displayable-- // Reserve one line for the scroll indicator
	}

	if displayable < 0 {
		displayable = 0
	}
	if displayable > numParamRows {
		displayable = numParamRows
	}
	return displayable
}

func (pc *ParamsContainer) ensureFocusedInputVisible() {
	numDisplayable := pc.getNumDisplayableInputRows()

	if numDisplayable <= 0 || numDisplayable >= numParamRows { // No scrolling needed or possible
		pc.scrollOffset = 0
		return
	}

	// If focused row is above the visible window, scroll up
	if pc.focusedRow < pc.scrollOffset {
		pc.scrollOffset = pc.focusedRow
	} else if pc.focusedRow >= pc.scrollOffset+numDisplayable {
		// If focused row is below the visible window, scroll down
		pc.scrollOffset = pc.focusedRow - numDisplayable + 1
	}

	// Clamp scrollOffset to valid range
	if pc.scrollOffset < 0 {
		pc.scrollOffset = 0
	}
	maxScrollOffset := numParamRows - numDisplayable

	maxScrollOffset = max(maxScrollOffset, 0)
	if pc.scrollOffset > maxScrollOffset {
		pc.scrollOffset = maxScrollOffset
	}
}

// Focus sets the focus to the first input field in the container.
func (pc *ParamsContainer) Focus() {
	pc.focusedRow = 0
	pc.focusedCol = 0
	pc.focusCurrentInput()
	pc.ensureFocusedInputVisible() // Ensure the newly focused input is visible
}

// Blur removes focus from all input fields in the container.
func (pc *ParamsContainer) Blur() {
	pc.blurAllInputs()
}

// Update handles messages for the ParamsContainer.
func (pc *ParamsContainer) Update(msg tea.Msg) tea.Cmd {
	if !pc.Active {
		return nil
	}

	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Ensure that if an input is focused, it gets the key press first,
		// unless it\'s a navigation key we want to intercept.
		// Intercept navigation keys regardless of input focus.
		switch msg.String() {
		case "up":
			if pc.focusedRow > 0 {
				pc.focusedRow--
				pc.focusCurrentInput()
				pc.ensureFocusedInputVisible()
			}
			return nil
		case "down":
			if pc.focusedRow < numParamRows-1 {
				pc.focusedRow++
				pc.focusCurrentInput()
				pc.ensureFocusedInputVisible()
			}
			return nil
		case "left":
			if pc.focusedCol == 1 { // If on Value, move to Name
				pc.focusedCol = 0
				pc.focusCurrentInput()
			} else if pc.focusedCol == 0 && pc.focusedRow > 0 { // If on Name and not first row, move to Value of prev row
				pc.focusedRow--
				pc.focusedCol = 1 // Move to Value column of previous row
				pc.focusCurrentInput()
				pc.ensureFocusedInputVisible() // Row changed
			}
			return nil
		case "right": // Treat Tab as right
			if pc.focusedCol == 0 { // If on Name, move to Value
				pc.focusedCol = 1
				pc.focusCurrentInput()
			} else if pc.focusedCol == 1 && pc.focusedRow < numParamRows-1 { // If on Value and not last row, move to Name of next row
				pc.focusedRow++
				pc.focusedCol = 0 // Move to Name column of next row
				pc.focusCurrentInput()
				pc.ensureFocusedInputVisible() // Row changed
			}
			return nil
		case "shift+tab": // Treat Shift+Tab as left
			if pc.focusedCol == 1 { // If on Value, move to Name of current row
				pc.focusedCol = 0
				pc.focusCurrentInput()
			} else if pc.focusedCol == 0 && pc.focusedRow > 0 { // If on Name and not first row, move to Value of prev row
				pc.focusedRow--
				pc.focusedCol = 1 // Move to Value column of previous row
				pc.focusCurrentInput()
				pc.ensureFocusedInputVisible() // Row changed
			}
			return nil
		default:
			// If not a navigation key, pass to the focused input
			if pc.focusedRow >= 0 && pc.focusedRow < len(pc.Inputs) {
				if pc.focusedCol == 0 {
					pc.Inputs[pc.focusedRow].NameInput, cmd = pc.Inputs[pc.focusedRow].NameInput.Update(msg)
					cmds = append(cmds, cmd)
				} else {
					pc.Inputs[pc.focusedRow].ValueInput, cmd = pc.Inputs[pc.focusedRow].ValueInput.Update(msg)
					cmds = append(cmds, cmd)
				}
			}
		}
	}
	return tea.Batch(cmds...)
}

// View renders the ParamsContainer.
func (pc *ParamsContainer) View() string {
	var rows []string

	// Labels
	nameLabel := "Name"
	valueLabel := "Value"
	labelStyle := lipgloss.NewStyle().Bold(true).Foreground(styles.SecondaryColor)

	const spacingBetweenInputs = 1
	const textInputHorizontalBorderWidth = 2 // Should match SetWidth
	const desiredInputContentWidth = 35      // Should match SetWidth
	const idealOuterWidthPerInput = desiredInputContentWidth + textInputHorizontalBorderWidth
	const totalIdealOuterWidth = idealOuterWidthPerInput * 2

	// Available width for the two text input columns (outer widths)
	// This is pc.contentWidth (container content area) minus spacing between inputs.
	inputsTotalOuterWidth := pc.contentWidth - spacingBetweenInputs

	inputsTotalOuterWidth = max(inputsTotalOuterWidth, 0)

	var nameInputRenderWidth, valueInputRenderWidth int

	if inputsTotalOuterWidth >= totalIdealOuterWidth {
		nameInputRenderWidth = idealOuterWidthPerInput
		valueInputRenderWidth = idealOuterWidthPerInput
	} else {
		nameInputRenderWidth = inputsTotalOuterWidth / 2
		valueInputRenderWidth = inputsTotalOuterWidth - nameInputRenderWidth
	}

	if nameInputRenderWidth < 0 {
		nameInputRenderWidth = 0
	}
	if valueInputRenderWidth < 0 {
		valueInputRenderWidth = 0
	}

	actualContentWidth := pc.contentWidth // Used for separator, scroll indicator and help text
	if actualContentWidth < 0 {
		actualContentWidth = 0
	}

	header := lipgloss.JoinHorizontal(lipgloss.Top,
		lipgloss.NewStyle().Width(nameInputRenderWidth).Render(labelStyle.Render(nameLabel)),
		lipgloss.NewStyle().Width(spacingBetweenInputs).Render(""), // Spacer cell
		lipgloss.NewStyle().Width(valueInputRenderWidth).Render(labelStyle.Render(valueLabel)),
	)
	rows = append(rows, header)
	// Use a simple line character for the separator
	rows = append(rows, strings.Repeat("─", actualContentWidth)) // Separator

	numDisplayable := pc.getNumDisplayableInputRows()

	startRow := pc.scrollOffset
	endRow := pc.scrollOffset + numDisplayable
	if endRow > numParamRows {
		endRow = numParamRows
	}
	if startRow > endRow {
		startRow = endRow
	}

	for i := startRow; i < endRow; i++ {
		nameView := pc.Inputs[i].NameInput.View()   // Content for name input
		valueView := pc.Inputs[i].ValueInput.View() // Content for value input

		// rowStyle := lipgloss.NewStyle() // Removed background highlight
		// if pc.Active && pc.focusedRow == i {
		// 	rowStyle = rowStyle.Background(lipgloss.Color("240")) // A light gray for focused row
		// }

		// Base style for input boxes
		inputBoxBaseStyle := lipgloss.NewStyle()

		// Style for the name input box
		nameBoxStyle := inputBoxBaseStyle.
			Width(nameInputRenderWidth).
			BorderForeground(styles.SecondaryColor)
		if nameInputRenderWidth >= 3 {
			nameBoxStyle = nameBoxStyle.Border(lipgloss.RoundedBorder())
		} else if nameInputRenderWidth > 0 { // Use normal border if not enough space for rounded
			nameBoxStyle = nameBoxStyle.Border(lipgloss.NormalBorder())
		} // else no border if width is 0

		// Style for the value input box
		valueBoxStyle := inputBoxBaseStyle.
			Width(valueInputRenderWidth).
			BorderForeground(styles.SecondaryColor)
		if valueInputRenderWidth >= 3 {
			valueBoxStyle = valueBoxStyle.Border(lipgloss.RoundedBorder())
		} else if valueInputRenderWidth > 0 { // Use normal border if not enough space for rounded
			valueBoxStyle = valueBoxStyle.Border(lipgloss.NormalBorder())
		} // else no border if width is 0

		// Highlight focused input by changing its border color
		if pc.Active && pc.focusedRow == i {
			if pc.focusedCol == 0 { // Name input is focused
				nameBoxStyle = nameBoxStyle.BorderForeground(styles.PrimaryColor)
			} else { // Value input is focused
				valueBoxStyle = valueBoxStyle.BorderForeground(styles.PrimaryColor)
			}
		}

		styledNameView := nameBoxStyle.Render(nameView)
		styledValueView := valueBoxStyle.Render(valueView)

		rowRender := lipgloss.JoinHorizontal(lipgloss.Top,
			styledNameView,
			lipgloss.NewStyle().Width(spacingBetweenInputs).Render(""), // Spacer cell
			styledValueView,
		)
		// rows = append(rows, rowStyle.Render(rowRender)) // Render without the rowStyle background
		rows = append(rows, rowRender)
	}

	if numParamRows > numDisplayable && numDisplayable > 0 {
		scrollIndicator := ""
		if pc.scrollOffset > 0 {
			scrollIndicator += "↑ "
		} else {
			scrollIndicator += "  "
		}
		if pc.scrollOffset+numDisplayable < numParamRows {
			scrollIndicator += "↓"
		} else {
			scrollIndicator += " "
		}
		if strings.TrimSpace(scrollIndicator) != "" {
			rows = append(rows, lipgloss.NewStyle().Width(actualContentWidth).Align(lipgloss.Center).Render(scrollIndicator))
		}
	}

	// Add help text
	helpTextStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("226")) // Yellow
	helpText := "Use ↑/↓/←/→ to navigate."
	// Ensure help text doesn't exceed container width if it's very narrow
	// It might be better to let it wrap or truncate based on lipgloss behavior if Width is set.
	// For now, just render it. If actualContentWidth is too small, it will be truncated by the container.
	rows = append(rows, helpTextStyle.Width(actualContentWidth).Render(helpText))

	containerContent := lipgloss.JoinVertical(lipgloss.Left, rows...)

	currentContainerStyle := styles.BorderStyle
	if pc.Active {
		currentContainerStyle = styles.ActiveBorderStyle
	}

	return currentContainerStyle.Width(pc.Width).Height(pc.Height).Render(containerContent)
}

// GetParams returns the current parameters as a map.
func (pc *ParamsContainer) GetParams() map[string]string {
	params := make(map[string]string)
	for _, p := range pc.Inputs {
		name := strings.TrimSpace(p.NameInput.Value())
		value := strings.TrimSpace(p.ValueInput.Value())
		if name != "" { // Only include if name is not empty
			params[name] = value
		}
	}
	return params
}

// ClearParams clears all input fields.
func (pc *ParamsContainer) ClearParams() {
	for i := range pc.Inputs {
		pc.Inputs[i].NameInput.Reset()
		pc.Inputs[i].ValueInput.Reset()
	}
	pc.focusedRow = 0
	pc.focusedCol = 0
	if numParamRows > 0 {
		// pc.Inputs[0].NameInput.Focus() // Focus is handled by SetActive or Focus()
		pc.focusCurrentInput() // Ensure the correct input is focused after clearing
	}
}

// IsAnyInputFocused checks if any text input within the ParamsContainer is currently focused.
func (pc *ParamsContainer) IsAnyInputFocused() bool {
	if pc.focusedRow < 0 || pc.focusedRow >= len(pc.Inputs) {
		return false
	}
	if pc.focusedCol == 0 {
		return pc.Inputs[pc.focusedRow].NameInput.Focused()
	}
	return pc.Inputs[pc.focusedRow].ValueInput.Focused()
}
