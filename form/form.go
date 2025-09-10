package form

import (
	"github.com/charmbracelet/huh"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	Form *huh.Form
	Currency string
	TopList  string
}

func New() Model {
	var currency, topList string
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Choose your preferred currency").
				Options(
					huh.NewOption("USD", "usd"),
					huh.NewOption("EUR", "eur"),
					huh.NewOption("GBP", "gbp"),
					huh.NewOption("INR", "inr"),
					huh.NewOption("AUD", "aud"),
				).
				Value(&currency),

			huh.NewSelect[string]().
				Title("Top List (eg: top 10 cryptocurrencies)").
				Options(
					huh.NewOption("Top 10", "10").Selected(true),
					huh.NewOption("Top 25", "25"),
					huh.NewOption("Top 50", "50"),
				).
				Value(&topList),
		),
	)
	return Model{Form: form, Currency: currency, TopList: topList}
}

func (m Model) Init() tea.Cmd {
	return m.Form.Init()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	form, cmd := m.Form.Update(msg)
	m.Form = form.(*huh.Form)
	return m, cmd
}

func (m Model) View() string {
	return m.Form.View()
}

func (m Model) GetFormData() (string, string) {
	return m.Currency, m.TopList
}
