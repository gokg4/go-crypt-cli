package core

import (
	"fmt"
	"time"

	"gocryptocli/form"
	marketdata "gocryptocli/marketData"

	"github.com/charmbracelet/huh/spinner"
)

const VsCurrencyDefault = "inr"

// Handles both form input and data fetching
func RestartApp() (marketdata.GeckoMarketData, string) {
	currency, topList := form.GetFormData()
	if currency == "" {
		currency = VsCurrencyDefault
	}

	var data marketdata.GeckoMarketData
	action := func() {
		data = marketdata.GetMarketData(currency, topList)
		time.Sleep(2 * time.Second)
	}

	fmt.Println()
	_ = spinner.New().Title("Getting data...").Action(action).Run()

	return data, currency
}
