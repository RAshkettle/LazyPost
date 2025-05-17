// Package components provides UI components for the LazyPost application.
package components

import (
	"github.com/RAshkettle/LazyPost/ui/styles"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/charmbracelet/lipgloss"
)

const numHeaderRows = 9

// HeaderInput represents a single Header/Value input pair with a select for header names.
type HeaderInput struct {
	HeaderSelect list.Model    // Select input for header names
	ValueInput   textinput.Model // Text input for header values
}

// HeadersInputContainer manages a list of header inputs (Header/Value pairs).
type HeadersInputContainer struct {
	Inputs       []HeaderInput // Slice of header inputs
	Width        int           // Width of the container
	Height       int           // Height of the container
	Active       bool          // Whether the container is currently active/focused
	focusedRow   int           // Index of the currently focused row
	focusedCol   int           // 0 for HeaderSelect, 1 for ValueInput
	scrollOffset int           // For scrolling if not all rows fit
	contentWidth int           // Calculated width for content area
}

// stringItem is a simple implementation of list.Item for header options.
type stringItem string

// FilterValue returns the string value for filtering and display.
func (i stringItem) FilterValue() string { return string(i) }

// headerOptions defines the selectable header names.
var headerOptions = []list.Item{
	stringItem("Accept"),
	stringItem("Accept-Charset"),
	stringItem("Accept-Encoding"),
	stringItem("Accept-Language"),
	stringItem("Accept-Datetime"),
	stringItem("Authorization"),
	stringItem("Cookie"),
	stringItem("Expect"),
	stringItem("From"),
	stringItem("Host"),
	stringItem("If-Match"),
	stringItem("If-Modified-Since"),
	stringItem("If-None-Match"),
	stringItem("If-Range"),
	stringItem("If-Unmodified-Since"),
	stringItem("Max-Forwards"),
	stringItem("Proxy-Authorization"),
	stringItem("Range"),
	stringItem("Referer"),
	stringItem("TE"),
	stringItem("User-Agent"),
	stringItem("X-Requested-With"),
	stringItem("X-Forwarded-For"),
	stringItem("X-Forwarded-Host"),
	stringItem("X-Forwarded-Proto"),
	stringItem("X-HTTP-Method-Override"),
	stringItem("X-Csrf-Token"),
	stringItem("X-XSS-Protection"),
	stringItem("X-Content-Type-Options"),
	stringItem("X-Real-IP"),
	stringItem("X-Powered-By"),
	stringItem("DNT"),
	stringItem("X-Api-Key"),
	stringItem("X-Auth-Token"),
}

// NewHeadersInputContainer creates a new HeadersInputContainer with a predefined number of rows.
func NewHeadersInputContainer() HeadersInputContainer {
	inputs := make([]HeaderInput, numHeaderRows)
	// Use headerOptions for each select widget
	for i := 0; i < numHeaderRows; i++ {
		headerSel := list.New(headerOptions, list.NewDefaultDelegate(), 0, 1)
		headerSel.Title = "Header"
		// Remove extra UI chrome: no title bar, status bar, help text, or filter input
		headerSel.SetShowTitle(false)
		headerSel.SetShowStatusBar(false)
		headerSel.SetShowHelp(false)
		headerSel.SetFilteringEnabled(false)

		valueInput := textinput.New()
		valueInput.Placeholder = "Value"
		valueInput.Prompt = ""
		valueInput.CharLimit = 35

		inputs[i] = HeaderInput{HeaderSelect: headerSel, ValueInput: valueInput}
	}

	// Focus first row's header select by default
	if numHeaderRows > 0 {
		inputs[0].HeaderSelect.Select(0)
	}

	return HeadersInputContainer{
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
func (hc *HeadersInputContainer) SetWidth(width int) {
	hc.Width = width
	currentStyle := styles.BorderStyle
	if hc.Active {
		currentStyle = styles.ActiveBorderStyle
	}
	// Compute inner content width
	containerChrome := currentStyle.GetHorizontalBorderSize() + currentStyle.GetHorizontalPadding()
	hc.contentWidth = width - containerChrome
	if hc.contentWidth < 0 {
		hc.contentWidth = 0
	}

	// Layout: two columns (HeaderSelect and ValueInput) separated by 1 space
	const spacing = 1
	inputsTotal := hc.contentWidth - spacing
	if inputsTotal < 0 {
		inputsTotal = 0
	}

	// Ideal widths: same as ParamsContainer's 35 chars + 2 border
	const idealContent = 35
	const borderWidth = 2
	const idealOuter = idealContent + borderWidth
	const totalIdeal = idealOuter * 2

	var hdrOuter, valOuter int
	if inputsTotal >= totalIdeal {
		hdrOuter = idealOuter
		valOuter = idealOuter
	} else {
		hdrOuter = inputsTotal / 2
		valOuter = inputsTotal - hdrOuter
	}
	if hdrOuter < 0 {
		hdrOuter = 0
	}
	if valOuter < 0 {
		valOuter = 0
	}

	// Fixed column widths: Header select width=30, Value input width=40
	for i := range hc.Inputs {
		hc.Inputs[i].HeaderSelect.SetWidth(30)
		hc.Inputs[i].ValueInput.Width = 40
	}
}

// SetHeight sets the height of the container.
func (hc *HeadersInputContainer) SetHeight(height int) {
	hc.Height = height
}

// SetActive sets the active state of the container.
func (hc *HeadersInputContainer) SetActive(active bool) {
	hc.Active = active
}

// Update handles incoming Bubble Tea messages for the headers container (stub for layout).
func (hc *HeadersInputContainer) Update(msg tea.Msg) tea.Cmd {
	// TODO: implement navigation and input handling
	return nil
}

// View renders the HeadersInputContainer.
func (hc HeadersInputContainer) View() string {
	var rows []string

	// Header row labels
	headLabel := "Header"
	valLabel := "Value"
	labelStyle := lipgloss.NewStyle().Bold(true).Foreground(styles.SecondaryColor)

	// Render header labels with fixed widths
	head := lipgloss.JoinHorizontal(lipgloss.Top,
		lipgloss.NewStyle().Width(30).Render(labelStyle.Render(headLabel)),
		lipgloss.NewStyle().Width(1).Render(""),
		lipgloss.NewStyle().Width(40).Render(labelStyle.Render(valLabel)),
	)
	rows = append(rows, head)

	// Render each input row
	for i := 0; i < numHeaderRows; i++ {
		hdrView := hc.Inputs[i].HeaderSelect.View()
		valView := hc.Inputs[i].ValueInput.View()
		// Wrap the value input in a border like Params
		valBox := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(styles.SecondaryColor).
			Width(40).
			Render(valView)
		row := lipgloss.JoinHorizontal(lipgloss.Top,
			lipgloss.NewStyle().Width(30).Render(hdrView),
			lipgloss.NewStyle().Width(1).Render(""),
			valBox,
		)
		rows = append(rows, row)
	}

	// Combine rows
	content := lipgloss.JoinVertical(lipgloss.Left, rows...)
	style := styles.BorderStyle
	if hc.Active {
		style = styles.ActiveBorderStyle
	}
	return style.Width(hc.Width).Height(hc.Height).Render(content)
}
