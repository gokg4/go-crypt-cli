package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m *Model) View() string {
	if m.ShowError {
		return errorStyle.Render(fmt.Sprintf("Error: %s\n\nPress 'e' to edit preferences, or any other key to retry.", m.ErrorMessage))
	}
	if m.Loading {
		return fmt.Sprintf("\n   %s %s... \n\n", m.Spinner.View(), m.LoadingMessage)
	}
	if m.ShowDetails {
		return fmt.Sprintf("%s\n%s\n%s", m.headerView(), m.Viewport.View(), m.footerView()+"\nPress 'm' to save as markdown, 'enter' or 'esc' to go back.")
	}
	return "\n" + m.Table.View() + "\nPress 'q' to quit, 'enter' for details, 'e' to edit preferences.\n"
}

func (m *Model) headerView() string {
	title := titleStyle.Render("Crypto Details " + strings.ToUpper(m.Currency))
	line := strings.Repeat("─", max(0, m.Viewport.Width-lipgloss.Width(title)))
	return lipgloss.JoinHorizontal(lipgloss.Center, title, line)
}

func (m *Model) footerView() string {
	info := infoStyle.Render(fmt.Sprintf("%3.f%%", m.Viewport.ScrollPercent()*100))
	status := ""
	if m.StatusMessage != "" {
		status = infoStyle.Render(m.StatusMessage)
	}
	line := strings.Repeat("─", max(0, m.Viewport.Width-lipgloss.Width(info)-lipgloss.Width(status)))
	return lipgloss.JoinHorizontal(lipgloss.Center, line, status, info)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
