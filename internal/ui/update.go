package ui

import (
	"cli-view-crypto-prices/internal/api"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/wordwrap"
)

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
