package table

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	marketdata "gocryptocli/marketData"
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

type (
	DescriptionFetchedMsg struct {
		Index       int
		Description string
	}
	ErrMsg struct{ err error }
)

type Model struct {
	table      table.Model
	data       marketdata.GeckoMarketData
	showDetail bool
	selected   int
	currency   string
	vp         viewport.Model
}

func New(d marketdata.GeckoMarketData, c string) Model {
	rows := make([]table.Row, len(d))
	for i, item := range d {
		rows[i] = table.Row{
			fmt.Sprintf("%d", item.MarketCapRank),
			item.Name,
			item.Symbol,
			fmt.Sprintf("%.2f", item.CurrentPrice),
		}
	}

	columns := []table.Column{
		{Title: "Rank", Width: 5},
		{Title: "Name", Width: 20},
		{Title: "Symbol", Width: 7},
		{Title: fmt.Sprintf("Price (%s)", strings.ToUpper(c)), Width: 14},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(12),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.BorderStyle(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color("240")).BorderBottom(true).Bold(false)
	s.Selected = s.Selected.Foreground(lipgloss.Color("229")).Background(lipgloss.Color("57")).Bold(false)
	t.SetStyles(s)

	return Model{
		table:    t,
		data:     d,
		currency: c,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			if !m.showDetail {
				m.showDetail = true
				m.selected = m.table.Cursor()
				coin := m.data[m.selected]
				if coin.Description == "" {
					return m, func() tea.Msg {
						desc, err := marketdata.FetchCoinDescription(coin.ID)
						if err != nil {
							return ErrMsg{err}
						}
						return DescriptionFetchedMsg{Index: m.selected, Description: desc}
					}
				}
				m.updateDetailView()
			} else {
				m.showDetail = false
			}
			return m, nil
		}

	case DescriptionFetchedMsg:
		m.data[msg.Index].Description = msg.Description
		if msg.Index == m.selected {
			m.updateDetailView()
		}
		return m, nil

	case ErrMsg:
		fmt.Println(msg.err)
		return m, nil
	}

	if m.showDetail {
		m.vp, cmd = m.vp.Update(msg)
	} else {
		m.table, cmd = m.table.Update(msg)
	}
	return m, cmd
}

func (m *Model) updateDetailView() {
	coin := m.data[m.selected]
	content := fmt.Sprintf(`Name: %s
Symbol: %s
Current Price: %.2f %s
Market Cap: %d %s
Market Cap Rank: %d

Description:
%s`, 
		coin.Name, 
		strings.ToUpper(coin.Symbol), 
		coin.CurrentPrice, 
		strings.ToUpper(m.currency), 
		coin.MarketCap, 
		strings.ToUpper(m.currency), 
		coin.MarketCapRank, 
		coin.Description,
	)
	m.vp = viewport.New(60, 15)
	m.vp.SetContent(content)
}

func (m Model) View() string {
	currentTime := time.Now()
	if m.showDetail {
		return lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("240")).
			Padding(1, 2).
			Render(m.vp.View()) + "\n\nPress Enter to go back."
	}
	return fmt.Sprintf("Current market price (%v %v, %v | %v) \n", currentTime.Month(), currentTime.Day(), currentTime.Year(), currentTime.Format(time.Kitchen)) +
		baseStyle.Render(m.table.View()) + "\n" + m.table.HelpView()
}
