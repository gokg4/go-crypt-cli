package table

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh/spinner"
	"github.com/charmbracelet/lipgloss"

	marketdata "gocryptocli/marketData"
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

var (
	titleStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Right = "├"
		return lipgloss.NewStyle().BorderStyle(b).Padding(0, 1)
	}()

	infoStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Left = "┤"
		return titleStyle.BorderStyle(b)
	}()
)

type model struct {
	table      table.Model
	data       marketdata.GeckoMarketData
	showDetail bool
	selected   int
	currency   string
	quitToMain bool
	vp         viewport.Model // Viewport for the entire detail view
}

func (m *model) Init() tea.Cmd { return nil }

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
		case "m":
			m.quitToMain = true
			return m, tea.Quit
		case "enter":
			if !m.showDetail {
				m.showDetail = true
				m.selected = m.table.Cursor()
				coin := &m.data[m.selected]
				if coin.Description == "" {
					action := func() {
						desc, err := marketdata.FetchCoinDescription(coin.ID)
						if err == nil {
							coin.Description = desc
						} else {
							coin.Description = "Description not available."
						}
						time.Sleep(3 * time.Second)
					}
					_ = spinner.New().Title("Fetching details...").Action(action).Run()
				}
				// Build the full detail string
				detailContent := fmt.Sprintf(
					"Name: %s\nSymbol: %s\nCurrent Price: %.2f %s\nMarket Cap: %d %s\nMarket Cap Rank: %d\n\nDescription:\n%s",
					coin.Name,
					strings.ToUpper(coin.Symbol),
					coin.CurrentPrice, strings.ToUpper(m.currency),
					coin.MarketCap, strings.ToUpper(m.currency),
					coin.MarketCapRank,
					coin.Description,
				)
				// Initialize viewport for the entire detail view
				m.vp = viewport.New(44, 12) // width, height (adjust as needed)
				m.vp.SetContent(detailContent)
				m.vp.SetYOffset(0) // always start at top
			} else {
				m.showDetail = false
			}
			return m, nil
		}
	}

	// Handle scrolling in detail view
	if m.showDetail {
		m.vp, cmd = m.vp.Update(msg)
		return m, cmd
	}

	// Table navigation
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m *model) headerView() string {
	title := titleStyle.Render("Crypto Details " + strings.ToUpper(m.currency))
	line := strings.Repeat("─", max(0, m.vp.Width-lipgloss.Width(title)))
	// titleStyle := lipgloss.NewStyle().
	// 	Bold(true).
	// 	Foreground(lipgloss.Color("205")).
	// 	Padding(0, 1)
	return lipgloss.JoinHorizontal(lipgloss.Center, title, line)
}

func (m *model) footerView() string {
	// infoStyle := lipgloss.NewStyle().
	// 	Foreground(lipgloss.Color("241")).
	// 	PaddingLeft(1).
	// 	PaddingRight(1)
	info := infoStyle.Render(fmt.Sprintf("%3.f%%", m.vp.ScrollPercent()*100))
	line := strings.Repeat("─", max(0, m.vp.Width-lipgloss.Width(info)))
	return lipgloss.JoinHorizontal(lipgloss.Center, line, info)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func (m *model) View() string {
	if m.showDetail {
		detailBox := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("240")).
			Padding(1, 2).
			Width(48)

		content := lipgloss.JoinVertical(lipgloss.Left,
			m.headerView(),
			m.vp.View(),
			m.footerView(),
		)

		return detailBox.Render(content) +
			"\n\nPress Enter to go back. Use ↑/↓, PgUp/PgDn to scroll.\n"
	}
	return baseStyle.Render(m.table.View()) + "\n  " + m.table.HelpView() + "\n"
}

func CreateTable(d marketdata.GeckoMarketData, c string) (quitToMain bool) {
	currency := c
	currentTime := time.Now()
	rows := []table.Row{}
	for _, item := range d {
		rows = append(rows, table.Row{
			fmt.Sprintf("%d", item.MarketCapRank),
			item.Name,
			item.Symbol,
			fmt.Sprintf("%.2f", item.CurrentPrice),
		})
	}

	columns := []table.Column{
		{Title: "Rank", Width: 5},
		{Title: "Name", Width: 20},
		{Title: "Symbol", Width: 7},
		{Title: fmt.Sprintf("Price (%s)", strings.ToUpper(currency)), Width: 14},
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

	fmt.Println("Press 'q' to quit, 'm' to restart, 'e' to focus/unfocus, 'Enter' for details.")
	fmt.Printf("Current market price (%v %v, %v | %v) \n", currentTime.Month(), currentTime.Day(), currentTime.Year(), currentTime.Format(time.Kitchen))

	m := &model{
		table:    t,
		data:     d,
		currency: currency,
	}
	_, err := tea.NewProgram(m).Run()
	if err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
	return m.quitToMain
}
