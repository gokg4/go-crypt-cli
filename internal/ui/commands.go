package ui

import (
	"cli-view-crypto-prices/internal/api"

	tea "github.com/charmbracelet/bubbletea"
)

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
