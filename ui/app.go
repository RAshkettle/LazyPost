// Package ui provides the user interface for the LazyPost application.
// It implements a TUI (Text User Interface) using the Bubble Tea framework.
package ui

import (
	"fmt"
	"os"
	"strings"

	"github.com/RAshkettle/LazyPost/ui/components"
	"github.com/RAshkettle/LazyPost/ui/styles"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// App represents the main application model.
// It embeds all UI components and manages the application state and logic.
type App struct {
	methodSelector components.MethodSelector // Component for selecting HTTP method.
	urlInput       components.URLInput       // Component for URL input.
	submitButton   components.SubmitButton   // Component for the submit button.
	tabContainer   components.TabsContainer  // Component for managing query and result tabs.
	toast          components.Toast          // Component for displaying toast notifications.
	spinner        components.Spinner        // Component for displaying a loading spinner.
	width          int                       // Current width of the terminal window.
	height         int                       // Current height of the terminal window.
	urlInputWidth  int                       // Cached width of the URL input, used for spinner positioning.
	urlInputX      int                       // Cached X coordinate of the URL input, used for spinner positioning.
	keymap         KeyMap                    // Defines keybindings for the application.
	banner         string                    // Stores the application banner text, loaded from banner.txt.
}

// NewApp initializes and returns a new App model.
// It sets up all the necessary UI components, loads the banner, and prepares the initial state.
func NewApp() App {
	methodSelector := components.NewMethodSelector()
	urlInput := components.NewURLInput()
	submitButton := components.NewButton("Submit")
	tabContainer := components.NewTabsContainer()
	toast := components.NewToast()
	spinner := components.NewSpinner()

	// Load the banner from file
	bannerData, err := os.ReadFile("banner.txt")
	bannerText := ""
	if err == nil {
		bannerText = string(bannerData)
	}

	return App{
		methodSelector: methodSelector,
		urlInput:       urlInput,
		submitButton:   submitButton,
		tabContainer:   tabContainer,
		toast:          toast,
		spinner:        spinner,
		width:          0,
		height:         0,
		keymap:         DefaultKeyMap,
		banner:         bannerText,
	}
}

// Init is the first command that is run when the application starts.
// It satisfies the tea.Model interface.
func (a App) Init() tea.Cmd {
	return tea.Batch(
		a.urlInput.TextInput.Focus(),
	)
}

// Update handles incoming messages and updates the App model accordingly.
// It is a central part of the Bubble Tea event loop and satisfies the tea.Model interface.
func (a App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case RequestCompleteMsg:
		a.handleRequestCompleteMsg(msg)
		return a, nil

	case components.SpinnerTickMsg:
		// Update spinner animation and continue ticking if visible
		if cmd := a.spinner.Update(msg); cmd != nil {
			cmds = append(cmds, cmd)
		}
		return a, tea.Batch(cmds...)

	case tea.KeyMsg:
		// First check if there's a toast visible - it should capture all key presses
		var shouldReturn bool

		var c tea.Cmd
		cmds, shouldReturn, c = a.handleKeyMsg(msg, cmds)
		if shouldReturn {
			return a, c
		}

	case tea.WindowSizeMsg:
		a.handleWindowSizeMsg(msg) // Position at the URL input
	}

	return a, tea.Batch(cmds...)
}

func (a *App) handleKeyMsg(msg tea.KeyMsg, cmds []tea.Cmd) ([]tea.Cmd, bool, tea.Cmd) {
	if a.toast.Visible && msg.String() == "enter" {
		// Dismiss the toast and focus the URL input
		a.toast.Hide()
		a.methodSelector.SetActive(false)
		a.urlInput.SetActive(true)
		a.submitButton.SetActive(false)
		a.tabContainer.SetActive(false)

		// Select all text in URL input
		a.urlInput.SelectAllText()
		return nil, true, nil
	}

	switch {
	case key.Matches(msg, a.keymap.Quit):
		return nil, true, tea.Quit

	case key.Matches(msg, a.keymap.FocusMethod):
		// Focus method selector
		a.methodSelector.SetActive(true)
		a.urlInput.SetActive(false)
		a.submitButton.SetActive(false)
		a.tabContainer.SetActive(false)
		return nil, true, nil

	case key.Matches(msg, a.keymap.FocusURL):
		// Focus URL input
		a.methodSelector.SetActive(false)
		a.urlInput.SetActive(true)
		a.submitButton.SetActive(false)
		a.tabContainer.SetActive(false)
		return nil, true, nil

	case key.Matches(msg, a.keymap.FocusSubmit):
		// Directly execute the submit action (not just focus)
		cmd := a.handleSubmit()
		return nil, true, cmd

	case key.Matches(msg, a.keymap.FocusQuery):
		// Switch to Query tab
		a.tabContainer.SwitchToTab(0)
		a.methodSelector.SetActive(false)
		a.urlInput.SetActive(false)
		a.submitButton.SetActive(false)
		a.tabContainer.SetActive(true)
		return nil, true, nil

	case key.Matches(msg, a.keymap.FocusResult):
		// Switch to Result tab
		a.tabContainer.SwitchToTab(1)
		a.methodSelector.SetActive(false)
		a.urlInput.SetActive(false)
		a.submitButton.SetActive(false)
		a.tabContainer.SetActive(true)
		return nil, true, nil

	case key.Matches(msg, a.keymap.Next), key.Matches(msg, a.keymap.Prev):
		// Tab and Shift+Tab only work in tab containers
		if a.tabContainer.Active {
			a.tabContainer.Update(msg)
			return nil, true, nil
		}
		// Otherwise, ignore tab/shift+tab
		return nil, true, nil

	// Let the active component handle other key presses
	default:
		// Special handling for arrow keys
		switch msg.String() {
		case "up", "down", "left", "right":
			// If method selector is active, let it handle arrow keys
			if a.methodSelector.Active {
				a.methodSelector.Update(msg)
				return nil, true, nil
			} else if a.urlInput.Active {
				// URL input handles arrow keys internally
				if cmd := a.urlInput.Update(msg); cmd != nil {
					cmds = append(cmds, cmd)
				}
				return nil, true, tea.Batch(cmds...)
			} else if a.tabContainer.Active {
				// Tab container might handle arrow keys
				a.tabContainer.Update(msg)
				return nil, true, nil
			}
		}

		// Handle other keys
		if a.methodSelector.Active {
			a.methodSelector.Update(msg)
		} else if a.urlInput.Active {
			if cmd := a.urlInput.Update(msg); cmd != nil {
				cmds = append(cmds, cmd)
			}

			// Special handling for Enter in URL field (submit the form)
			if msg.String() == "enter" {
				cmd := a.handleSubmit()
				return nil, true, cmd
			}
		} else if a.submitButton.Active {
			if _, submitted := a.submitButton.Update(msg); submitted {
				cmd := a.handleSubmit()
				return nil, true, cmd
			}
		} else if a.tabContainer.Active {
			a.tabContainer.Update(msg)
		}

	}
	return cmds, false, nil
}

func (a *App) handleWindowSizeMsg(msg tea.WindowSizeMsg) {
	a.width = msg.Width
	a.height = msg.Height

	// Calculate the available width after accounting for 10% padding (5% on each side)
	availableWidth := int(float64(a.width) * 0.9)
	paddingWidth := int(float64(a.width) * 0.05) // 5% padding on each side

	// Update component widths
	methodBoxWidth := int(float64(availableWidth) * 0.2)

	// Set button width to reasonable size (about 15% of available space)
	buttonWidth := int(float64(availableWidth) * 0.15)

	// URL gets the remaining space after method and button
	urlBoxWidth := availableWidth - methodBoxWidth - buttonWidth - 4 // -4 for spacing

	// Set tab container size - full width and most of the height
	tabContainerWidth := availableWidth
	// Reduce height by 15% from the previous calculation and accommodate for banner (7 lines)
	tabContainerHeight := int(float64(a.height-15) * 0.85) // Reduced to account for banner

	// Store URL input position and dimensions for the spinner
	a.urlInputWidth = urlBoxWidth
	a.urlInputX = methodBoxWidth + paddingWidth + 1 // Add paddingWidth (5%) and 1 for spacing

	a.methodSelector.SetWidth(methodBoxWidth)
	a.urlInput.SetWidth(urlBoxWidth)
	a.submitButton.SetWidth(buttonWidth)
	// Mirror button height to match URL container (no fixed height)
	a.tabContainer.SetWidth(tabContainerWidth)
	a.tabContainer.SetHeight(tabContainerHeight)

	// Set toast dimensions
	toastWidth := int(float64(availableWidth) * 0.5) // Half the available width
	a.toast.SetWidth(toastWidth)
	a.toast.SetHeight(5) // Fixed height

	// Set spinner dimensions to match the URL input
	a.spinner.SetWidth(urlBoxWidth)
	a.spinner.SetHeight(3) // URL input height (1 for title + 2 for input)
	a.spinner.SetPosition(a.urlInputX, 3)
}

func (a *App) handleRequestCompleteMsg(msg RequestCompleteMsg) {
	a.spinner.Hide()

	if msg.Error != nil {
		// Show error toast and allow user to try again
		a.toast.Show(fmt.Sprintf("Error: %s", msg.Error.Error()))
		// Move focus back to URL input
		a.methodSelector.SetActive(false)
		a.urlInput.SetActive(true)
		a.submitButton.SetActive(false)
		a.tabContainer.SetActive(false)
	}

	// Update the result tabs with response data
	resultTab := a.tabContainer.GetResultTab()
	resultTab.SetHeadersContent(msg.Headers) // Headers tab
	resultTab.SetBodyContent(msg.Body)       // Body tab

	// Activate the result tab and set it to show headers first
	a.tabContainer.SetActive(true)
	a.tabContainer.SwitchToTab(1) // Switch to Result tab (index 1)
	resultTab.SwitchToInnerTab(0) // Ensure Headers tab is active (index 0)
	resultTab.SetActive(true)     // Make sure the result tab is active
}

// View renders the current state of the application as a string.
// It satisfies the tea.Model interface.
func (a App) View() string {
	if a.width == 0 {
		return "Initializing..."
	}

	// Create the main view
	centeredView := a.renderMainView()

	// Check if toast should be shown
	if a.toast.Visible {
		return a.renderToastOverlay()
	}

	// Check if spinner should be shown
	if a.spinner.Visible {
		return a.renderSpinnerOverlay(centeredView)
	}

	return centeredView
}

// renderMainView creates the main UI layout with banner, inputs, and tabs
func (a App) renderMainView() string {
	// Render the banner with green color to match the primary color theme
	bannerBox := a.renderBanner()

	// Render the components
	methodBox := a.methodSelector.View()
	urlBox := a.urlInput.View()
	submitBox := a.submitButton.View()
	tabBox := a.tabContainer.View()

	// Arrange the top boxes side by side
	topRow := lipgloss.JoinHorizontal(lipgloss.Top, methodBox, urlBox, submitBox)

	// Add vertical arrangement with the banner at top, then input row, then tab container
	// Add a 2-line gap between the components for better spacing
	fullView := lipgloss.JoinVertical(lipgloss.Left, bannerBox, "", topRow, "", tabBox)

	// Add 5% padding on each side for centering
	paddingWidth := int(float64(a.width) * 0.05)

	// Create a centered style
	centeredStyle := lipgloss.NewStyle().
		PaddingLeft(paddingWidth).
		PaddingRight(paddingWidth)

	// Apply the centered style
	return centeredStyle.Render(fullView)
}

// renderBanner styles and returns the application banner
func (a App) renderBanner() string {
	bannerStyle := lipgloss.NewStyle().
		Foreground(styles.PrimaryColor). // Green (matching active elements)
		Bold(true).
		Align(lipgloss.Center)

	return bannerStyle.Render(a.banner)
}

// renderToastOverlay creates an overlay with a toast notification centered on the screen
func (a App) renderToastOverlay() string {
	toastView := a.toast.View()

	// Position the toast in the center of the screen
	toastStyle := lipgloss.NewStyle().
		Align(lipgloss.Center, lipgloss.Center).
		Padding((a.height / 2) - 6) // Truly center with padding

	toastView = toastStyle.Render(toastView)

	// Create an overlay that covers the entire screen with the toast in the center
	return lipgloss.Place(a.width, a.height, lipgloss.Center, lipgloss.Center, toastView)
}

// renderSpinnerOverlay creates an overlay with a spinner positioned over the URL input
func (a App) renderSpinnerOverlay(baseView string) string {
	spinnerView := a.spinner.View()

	// Calculate the line position of the URL input (3 lines from top: banner + empty line + title)
	urlLinePosition := 3

	// Now position the spinner directly on top of the URL input
	lines := strings.Split(baseView, "\n")

	// Replace the URL input lines with the spinner
	spinnerLines := strings.Split(spinnerView, "\n")

	// Insert the spinner at the URL position
	for i, spinnerLine := range spinnerLines {
		lineIndex := urlLinePosition + i
		if lineIndex < len(lines) {
			// Pad the spinner line to align it with the URL input
			paddedSpinnerLine := strings.Repeat(" ", a.urlInputX) + spinnerLine
			lines[lineIndex] = paddedSpinnerLine
		}
	}

	return strings.Join(lines, "\n")
}
