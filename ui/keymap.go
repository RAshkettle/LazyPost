package ui

import "github.com/charmbracelet/bubbles/key"

// KeyMap defines the keybindings for the application.
// It maps actions to specific key combinations.
type KeyMap struct {
	FocusMethod key.Binding // Alt+1: Focus the method selector
	FocusURL    key.Binding // Alt+2: Focus the URL input
	FocusSubmit key.Binding // Alt+5: Submit the request
	FocusQuery  key.Binding // Alt+3: Switch to query tab
	FocusResult key.Binding // Alt+4: Switch to result tab
	Next        key.Binding // Tab: Navigate to next inner tab
	Prev        key.Binding // Shift+Tab: Navigate to previous inner tab
	Quit        key.Binding // Ctrl+C/Esc: Quit the application
}

// DefaultKeyMap returns the default keybindings for the application.
// This defines the key mappings that will be used across the app.
var DefaultKeyMap = KeyMap{
	FocusMethod: key.NewBinding(
		key.WithKeys("alt+1"),
		key.WithHelp("alt+1", "focus method"),
	),
	FocusURL: key.NewBinding(
		key.WithKeys("alt+2"),
		key.WithHelp("alt+2", "focus url"),
	),
	FocusQuery: key.NewBinding(
		key.WithKeys("alt+3"),
		key.WithHelp("alt+3", "switch to query tab"),
	),
	FocusResult: key.NewBinding(
		key.WithKeys("alt+4"),
		key.WithHelp("alt+5", "switch to result tab"),
	),
	FocusSubmit: key.NewBinding(
		key.WithKeys("alt+5"),
		key.WithHelp("alt+5", "submit request"),
	),
	Next: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("tab", "next inner tab"),
	),
	Prev: key.NewBinding(
		key.WithKeys("shift+tab"),
		key.WithHelp("shift+tab", "prev inner tab"),
	),
	Quit: key.NewBinding(
		key.WithKeys("ctrl+c", "esc"),
		key.WithHelp("ctrl+c/esc", "quit"),
	),
}
