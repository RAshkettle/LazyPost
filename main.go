package main

import (
	"fmt"
	"os"

	"github.com/RAshkettle/LazyPost/ui"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	app := ui.NewApp()
	p := tea.NewProgram(app, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v", err)
		os.Exit(1)
	}
}
