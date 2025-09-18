package ui

import (
	"fmt"
	"strconv"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
)

func newTable(currency string) table.Model {
	columns := []table.Column{
		{Title: "Rank", Width: 4},
		{Title: "Name", Width: 20},
		{Title: "Symbol", Width: 10},
		{Title: fmt.Sprintf("Price (%s)", currency), Width: 20},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows([]table.Row{}),
		table.WithFocused(true),
		table.WithHeight(10),
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

	return t
}

func (m *Model) createTableRows() []table.Row {
	var rows []table.Row
	for _, crypto := range m.Cryptos {
		rows = append(rows, table.Row{
			strconv.Itoa(crypto.MarketCapRank),
			crypto.Name,
			crypto.Symbol,
			fmt.Sprintf("%.2f", crypto.CurrentPrice),
		})
	}
	return rows
}
