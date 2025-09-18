package ui

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func Run(currency string, limit int) {
	m := NewModel(currency, limit)
	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
