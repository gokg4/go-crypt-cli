package main

import (
	"fmt"
	"os"
	"time"

	"gocryptocli/form"
	marketdata "gocryptocli/marketData"
	"gocryptocli/table"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
)

type tickMsg time.Time

type (
	errMsg struct{ err error }
	loadedDataMsg struct {
		data     marketdata.GeckoMarketData
		currency string
	}
)

func (e errMsg) Error() string { return e.err.Error() }

type model struct {
	state          string
	form           form.Model
	table          table.Model
	progress       progress.Model
	loadingMessage string
	quitting       bool
	err            error
}

func initialModel() model {
	return model{
		state: "form",
		form:  form.New(),
	}
}

func (m model) Init() tea.Cmd {
	return m.form.Init()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit
		default:
			if m.state == "table" {
				var cmd tea.Cmd
				tableModel, cmd := m.table.Update(msg)
				m.table = tableModel.(table.Model)
				return m, cmd
			}
		}

	case errMsg:
		m.err = msg
		return m, nil

	case loadedDataMsg:
		m.state = "table"
		m.table = table.New(msg.data, msg.currency)
		return m, nil

	case tickMsg:
		if m.state == "loading" {
			newProgress := m.progress
			if m.progress.Percent() < 1.0 {
				newProgress.SetPercent(m.progress.Percent() + 0.2)
			}
			m.progress = newProgress
		}
		return m, tea.Tick(time.Millisecond*100, func(t time.Time) tea.Msg {
			return tickMsg(t)
		})

	}

	switch m.state {
	case "form":
		formModel, cmd := m.form.Update(msg)
		m.form = formModel.(form.Model)

		if m.form.Form.State == huh.StateCompleted {
			m.state = "loading"
			m.loadingMessage = "Fetching data..."
			m.progress = progress.New()
			return m, tea.Batch(m.fetchData(), tea.Tick(time.Millisecond*100, func(t time.Time) tea.Msg {
				return tickMsg(t)
			}))
		}

		return m, cmd

	case "loading":
		var cmd tea.Cmd
		progressModel, cmd := m.progress.Update(msg)
		if newProgress, ok := progressModel.(progress.Model); ok {
			m.progress = newProgress
		}
		return m, cmd

	case "table":
		var cmd tea.Cmd
		tableModel, cmd := m.table.Update(msg)
		m.table = tableModel.(table.Model)
		return m, cmd
	}

	return m, nil
}

func (m model) View() string {
	if m.quitting {
		return ""
	}
	if m.err != nil {
		return fmt.Sprintf("Error: %s\n", m.err)
	}

	switch m.state {
	case "form":
		return m.form.View()
	case "loading":
		return m.loadingMessage + "\n" + m.progress.View()
	case "table":
		return m.table.View()
	default:
		return ""
	}
}

func (m model) fetchData() tea.Cmd {
	return func() tea.Msg {
		currency, topList := m.form.GetFormData()
		data, err := marketdata.GetMarketData(currency, topList)
		if err != nil {
			return errMsg{err}
		}
		return loadedDataMsg{data: data, currency: currency}
	}
}

func main() {
	m := initialModel()
	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
