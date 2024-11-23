package table

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	marketdata "gocryptocli/marketData"
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

type model struct {
	table table.Model
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "e":
			if m.table.Focused() {
				m.table.Blur()
			} else {
				m.table.Focus()
			}
		case "q", "ctrl+c":
			return m, tea.Quit
		case "enter":
			return m, tea.Batch(
				tea.Printf("Let's go to %s!", m.table.SelectedRow()[1]),
			)
		}
	}
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return baseStyle.Render(m.table.View()) + "\n  " + m.table.HelpView() + "\n"
}

func CreateTable(d marketdata.GeckoMarketData, c string) {
	currency := c
	currentTime := time.Now()
	rows := []table.Row{}
	for _, item := range d {
		rows = append(rows, table.Row{
			fmt.Sprintf("%d", item.MarketCapRank),
			item.Name,
			item.Symbol,
			fmt.Sprintf("%.2f", item.CurrentPrice), // Format price
		})
	}

	// Create and render table
	columns := []table.Column{
		{Title: "Rank", Width: 5},
		{Title: "Name", Width: 20},
		{Title: "Symbol", Width: 7},
		{Title: fmt.Sprintf("Price (%s)", strings.ToUpper(currency)), Width: 12},
	}
	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(12),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)

	fmt.Println("Press 'q' to quit")
	fmt.Printf("Current market price (%v %v, %v | %v) \n", currentTime.Month(), currentTime.Day(), currentTime.Year(), currentTime.Format(time.Kitchen))

	m := model{t}
	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
