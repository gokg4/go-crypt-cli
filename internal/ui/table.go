package ui

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"cli-view-crypto-prices/internal/api"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/wordwrap"
)

var ( // styles
	titleStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFDF5")).
		Background(lipgloss.Color("#00B3FF")).
		Padding(0, 1)

	infoStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFDF5")).
		Background(lipgloss.Color("#0087FF")).
		Padding(0, 1)

	viewportStyle = lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		Padding(0, 1)

	errorStyle = lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("#FF0000")). // Red border for errors
		Foreground(lipgloss.Color("#FF0000")).
		Padding(1, 2)
)

type Model struct {
	Table          table.Model
	Spinner        spinner.Model
	Loading        bool
	LoadingMessage string
	SelectedID     string
	Viewport       viewport.Model
	ShowTable      bool
	ShowDetails    bool
	ShowError      bool
	ErrorMessage   string
	ShowForm       bool
	Coin           *api.Crypto
	Currency       string
	Limit          int
	Cryptos        []api.Crypto // Storing the list of cryptos
	Description    string       // Storing the description
	StatusMessage  string
}

// message for when the initial data is loaded
type dataLoadedMsg struct{ cryptos []api.Crypto }

// message for when description is fetched
type descriptionMsg struct{ description string }

// message for error during fetch
type errMsg struct{ err error }

func (m *Model) Init() tea.Cmd {
	return tea.Batch(m.Spinner.Tick, m.fetchData())
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)
	// Clear the status message on any key press
	if _, ok := msg.(tea.KeyMsg); ok {
		m.StatusMessage = ""
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		headerHeight := lipgloss.Height(m.headerView())
		footerHeight := lipgloss.Height(m.footerView())
		m.Viewport.Width = msg.Width
		m.Viewport.Height = msg.Height - headerHeight - footerHeight
		m.Viewport.SetContent(wordwrap.String(m.Description, m.Viewport.Width))
		return m, nil
	case tea.KeyMsg:
		if m.ShowError {
			switch msg.String() {
			case "e":
				m.ShowForm = true
				return m, tea.Quit
			default:
				m.ShowError = false
				m.ErrorMessage = ""
				if len(m.Cryptos) == 0 {
					m.Loading = true
					m.LoadingMessage = "Fetching initial market data"
					return m, m.fetchData()
				}
				m.ShowTable = true
				return m, nil
			}
		}
		if m.ShowDetails {
			switch msg.String() {
			case "m":
				if m.Coin != nil {
					markdown, err := GenerateMarkdown(*m.Coin, m.Currency)
					if err != nil {
						m.StatusMessage = "Error generating markdown"
						return m, nil
					}
					if err := os.MkdirAll("markdown", 0755); err != nil {
						m.StatusMessage = "Error creating markdown directory"
						return m, nil
					}
					filename := fmt.Sprintf("markdown/%s-details.md", m.Coin.ID)
					err = os.WriteFile(filename, []byte(markdown), 0644)
					if err != nil {
						m.StatusMessage = "Error saving markdown file"
					} else {
						m.StatusMessage = fmt.Sprintf("Saved to %s", filename)
					}
				}
				return m, nil
			}

		}

		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "enter":
			if m.ShowTable {
				if len(m.Cryptos) == 0 {
					return m, nil // No cryptos to select
				}
				m.Loading = true
				selectedCrypto := m.Cryptos[m.Table.Cursor()]
				m.SelectedID = selectedCrypto.ID
				m.LoadingMessage = fmt.Sprintf("Fetching details for %s", selectedCrypto.Name)
				return m, fetchDescription(m.SelectedID)
			}
			if m.ShowDetails {
				m.ShowDetails = false
				m.ShowTable = true
			}
		case "esc":
			if m.ShowDetails {
				m.ShowDetails = false
				m.ShowTable = true
			}
		case "e":
			if m.ShowTable {
				m.ShowForm = true
				return m, tea.Quit
			}
		}
	case dataLoadedMsg:
		m.Loading = false
		m.ShowTable = true
		m.Cryptos = msg.cryptos
		m.Table.SetRows(m.createTableRows())
		return m, nil
	case descriptionMsg:
		m.Loading = false
		m.ShowTable = false
		m.ShowDetails = true

		var selectedCoin api.Crypto
		for _, c := range m.Cryptos {
			if c.ID == m.SelectedID {
				selectedCoin = c
				break
			}
		}

		selectedCoin.Description = msg.description
		m.Coin = &selectedCoin

		markdown, err := GenerateMarkdown(*m.Coin, m.Currency)
		if err != nil {
			m.ErrorMessage = "Error generating markdown for view"
			m.ShowError = true
			return m, nil
		}

		out, err := glamour.Render(markdown, "dark")
		if err != nil {
			m.ErrorMessage = "Error rendering markdown for view"
			m.ShowError = true
			return m, nil
		}
		m.Description = out
		m.Viewport.SetContent(m.Description)

		return m, nil

	case errMsg:
		m.Loading = false
		m.ShowDetails = false
		m.ShowTable = false
		m.ShowError = true
		m.ErrorMessage = msg.err.Error()
		return m, nil
	}

	if m.ShowDetails {
		m.Viewport, cmd = m.Viewport.Update(msg)
		cmds = append(cmds, cmd)
	} else {
		m.Table, cmd = m.Table.Update(msg)
		cmds = append(cmds, cmd)
	}

	m.Spinner, cmd = m.Spinner.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

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

func NewModel(currency string, limit int) *Model {
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

	vp := viewport.New(100, 20)
	vp.Style = viewportStyle

	return &Model{
		Table:          t,
		Spinner:        NewSpinner(),
		Viewport:       vp,
		Loading:        true,
		LoadingMessage: "Fetching initial market data",
		ShowTable:      false,
		ShowForm:       false,
		ShowError:      false,
		ErrorMessage:   "",
		Currency:       currency,
		Limit:          limit,
	}
}

func (m *Model) fetchData() tea.Cmd {
	return func() tea.Msg {
		cryptos, err := api.GetMarketData(m.Currency, m.Limit)
		if err != nil {
			return errMsg{err}
		}
		return dataLoadedMsg{cryptos}
	}
}

func fetchDescription(id string) tea.Cmd {
	return func() tea.Msg {
		desc, err := api.GetCoinDescription(id)
		if err != nil {
			return errMsg{err}
		}
		return descriptionMsg{desc}
	}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
