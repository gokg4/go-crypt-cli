package ui

import (
	"cli-view-crypto-prices/internal/api"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/viewport"
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

func NewModel(currency string, limit int) *Model {
	return &Model{
		Table:          newTable(currency),
		Spinner:        NewSpinner(),
		Viewport:       viewport.New(100, 20),
		Loading:        true,
		LoadingMessage: "Fetching initial market data",
		ShowTable:      true,
		ShowForm:       false,
		ShowError:      false,
		ErrorMessage:   "",
		Currency:       currency,
		Limit:          limit,
	}
}
