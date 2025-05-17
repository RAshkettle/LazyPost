// Package components provides UI components for the LazyPost application.
package components

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// SpinnerTickMsg is sent when the spinner animation should advance.
type SpinnerTickMsg time.Time

// Spinner represents a loading spinner that overlays a specific UI component.
// It displays an animation to indicate ongoing processes like HTTP requests.
type Spinner struct {
	Visible  bool     // Whether the spinner is currently visible
	Width    int      // Width of the spinner in characters
	Height   int      // Height of the spinner in characters
	Frames   []string // Animation frames
	FrameIdx int      // Current frame index
	Message  string   // Optional text message to display with the spinner
	X        int      // X position for placing the spinner (default 0)
	Y        int      // Y position for placing the spinner (default 0)
}

// NewSpinner creates a new spinner component with default values.
// The spinner is initially hidden until Show() is called.
func NewSpinner() Spinner {
	return Spinner{
		Visible:  false,
		Width:    0,
		Height:   0,
		Frames:   []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"},
		FrameIdx: 0,
		Message:  "Loading...",
		X:        0,
		Y:        0,
	}
}

// SetWidth sets the width of the spinner component.
func (s *Spinner) SetWidth(width int) {
	s.Width = width
}

// SetHeight sets the height of the spinner component.
func (s *Spinner) SetHeight(height int) {
	s.Height = height
}

// SetPosition sets the position of the spinner.
// x is the horizontal position and y is the vertical position.
func (s *Spinner) SetPosition(x, y int) {
	s.X = x
	s.Y = y
}

// Show displays the spinner with an optional message.
// It returns a command to start the spinner animation.
func (s *Spinner) Show(message string) tea.Cmd {
	s.Visible = true
	if message != "" {
		s.Message = message
	}
	return s.tickCmd()
}

// Hide hides the spinner and stops its animation.
func (s *Spinner) Hide() {
	s.Visible = false
}

// tickCmd returns a command that sends a SpinnerTickMsg after a short delay.
// This creates the animation effect by updating the spinner frame.
func (s *Spinner) tickCmd() tea.Cmd {
	return tea.Tick(time.Millisecond*80, func(t time.Time) tea.Msg {
		return SpinnerTickMsg(t)
	})
}

// Update handles messages and updates the spinner state.
// It advances the animation frame when a SpinnerTickMsg is received.
func (s *Spinner) Update(msg tea.Msg) tea.Cmd {
	if !s.Visible {
		return nil
	}

	switch msg.(type) {
	case SpinnerTickMsg:
		s.FrameIdx = (s.FrameIdx + 1) % len(s.Frames)
		return s.tickCmd()
	}

	return nil
}

// View renders the spinner component.
// If the spinner is not visible, an empty string is returned.
func (s Spinner) View() string {
	if !s.Visible {
		return ""
	}

	// Get the current animation frame
	frame := s.Frames[s.FrameIdx]
	spinnerText := frame + " " + s.Message

	// Create a style for the spinner box
	spinnerStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#5DADE2")). // Light blue border
		Foreground(lipgloss.Color("#FFFFFF")).       // White text
		Background(lipgloss.Color("#2C3E50")).       // Dark blue-gray background
		Padding(1, 1).                              // Add some padding
		Width(s.Width - 4).                         // Adjust for border and padding
		Bold(true)                                  // Make the text bold

	// Render the spinner with its content
	rendered := spinnerStyle.Render(spinnerText)
	
	// Return the rendered spinner (positioning will be handled by the View function)
	return rendered
}
